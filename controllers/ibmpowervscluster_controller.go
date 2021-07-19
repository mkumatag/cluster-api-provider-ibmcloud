/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/cloud/scope"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrastructurev1alpha3 "github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/api/v1alpha3"
)

// IBMPowerVSClusterReconciler reconciles a IBMPowerVSCluster object
type IBMPowerVSClusterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclusters/status,verbs=get;update;patch

func (r *IBMPowerVSClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (_ ctrl.Result, reterr error) {
	log := r.Log.WithValues("ibmpowervscluster", req.NamespacedName)

	// your logic here

	// Fetch the IBMPowerVSCluster instance
	ibmCluster := &infrastructurev1alpha3.IBMPowerVSCluster{}
	err := r.Get(ctx, req.NamespacedName, ibmCluster)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Fetch the Cluster.
	cluster, err := util.GetOwnerCluster(ctx, r.Client, ibmCluster.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	clusterScope, err := scope.NewPowerVSClusterScope(scope.PowerVSClusterScopeParams{
		Client:            r.Client,
		Logger:            log,
		Cluster:           cluster,
		IBMPowerVSCluster: ibmCluster,
	})

	// Always close the scope when exiting this function so we can persist any GCPMachine changes.
	defer func() {
		if err := clusterScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// Handle deleted clusters
	if !ibmCluster.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(clusterScope)
	}

	if err != nil {
		return reconcile.Result{}, errors.Errorf("failed to create scope: %+v", err)
	} else {
		return r.reconcile(ctx, clusterScope)
	}
}

func (r *IBMPowerVSClusterReconciler) reconcile(ctx context.Context, clusterScope *scope.PowerVSClusterScope) (ctrl.Result, error) {
	if !controllerutil.ContainsFinalizer(clusterScope.IBMPowerVSCluster, infrastructurev1alpha3.IBMPowerVSClusterFinalizer) {
		controllerutil.AddFinalizer(clusterScope.IBMPowerVSCluster, infrastructurev1alpha3.IBMPowerVSClusterFinalizer)
		return ctrl.Result{}, nil
	}

	clusterScope.IBMPowerVSCluster.Status.Ready = true

	//if clusterScope.IBMPowerVSCluster.Status.APIEndpoint.PortID == nil {
	//	port, err := clusterScope.CreatePort()
	//	if err != nil {
	//		return ctrl.Result{}, errors.Wrap(err, "failed to create a port for APIEndpoint")
	//	}
	//
	//	clusterScope.IBMPowerVSCluster.Status.APIEndpoint = infrastructurev1alpha3.PowerVSAPIEndpoint{
	//		PortID:          port.PortID,
	//		InternalAddress: port.IPAddress,
	//	}
	//}
	//
	//if clusterScope.IBMPowerVSCluster.Status.APIEndpoint.PortID != nil && clusterScope.IBMPowerVSCluster.Status.APIEndpoint.Address == nil {
	//	portInfo, err := clusterScope.GetPort()
	//	if err != nil {
	//		return ctrl.Result{}, errors.Wrap(err, "failed to GetPort")
	//	}
	//	if portInfo.ExternalIP != "" {
	//		clusterScope.IBMPowerVSCluster.Status.APIEndpoint.Address = &portInfo.ExternalIP
	//	} else {
	//		return ctrl.Result{}, errors.Wrap(err, "failed to the external IP address for the port")
	//	}
	//	clusterScope.IBMPowerVSCluster.Status.Ready = true
	//}

	return ctrl.Result{}, nil
}

func (r *IBMPowerVSClusterReconciler) reconcileDelete(clusterScope *scope.PowerVSClusterScope) (ctrl.Result, error) {
	//if err := clusterScope.DeletePort(); err != nil {
	//	return ctrl.Result{}, errors.Wrap(err, "failed to delete a port for APIEndpoint")
	//} else {
	controllerutil.RemoveFinalizer(clusterScope.IBMPowerVSCluster, infrastructurev1alpha3.IBMPowerVSClusterFinalizer)
	return ctrl.Result{}, nil
	//}
}

func (r *IBMPowerVSClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha3.IBMPowerVSCluster{}).
		Complete(r)
}
