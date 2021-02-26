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
	"fmt"
	"github.com/IBM-Cloud/power-go-client/power/models"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-logr/logr"
	"github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/cloud/scope"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	infrastructurev1alpha3 "github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/api/v1alpha3"
)

// IBMPowerVSMachineReconciler reconciles a IBMPowerVSMachine object
type IBMPowerVSMachineReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsmachines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsmachines/status,verbs=get;update;patch

func (r *IBMPowerVSMachineReconciler) Reconcile(req ctrl.Request) (_ ctrl.Result, reterr error) {
	ctx := context.Background()
	log := r.Log.WithValues("ibmpowervsmachine", req.NamespacedName)

	ibmPowerVSMachine := &infrastructurev1alpha3.IBMPowerVSMachine{}
	err := r.Get(ctx, req.NamespacedName, ibmPowerVSMachine)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	fmt.Printf("ibmPowerVSMachine spec: %+v", ibmPowerVSMachine.Spec)

	// Fetch the Machine.
	machine, err := util.GetOwnerMachine(ctx, r.Client, ibmPowerVSMachine.ObjectMeta)
	if err != nil {
		return ctrl.Result{}, err
	}
	if machine == nil {
		log.Info("Machine Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	fmt.Printf("machine spec: %+v", machine.Spec)

	// Fetch the Cluster.
	cluster, err := util.GetClusterFromMetadata(ctx, r.Client, ibmPowerVSMachine.ObjectMeta)
	if err != nil {
		log.Info("Machine is missing cluster label or cluster does not exist")
		return ctrl.Result{}, nil
	}

	log = log.WithValues("cluster", cluster.Name)

	ibmCluster := &infrastructurev1alpha3.IBMPowerVSCluster{}
	ibmPowerVSClusterName := client.ObjectKey{
		Namespace: ibmPowerVSMachine.Namespace,
		Name:      cluster.Spec.InfrastructureRef.Name,
	}
	if err := r.Client.Get(ctx, ibmPowerVSClusterName, ibmCluster); err != nil {
		log.Info("IBMPowerVSCluster is not available yet")
		return ctrl.Result{}, nil
	}

	spew.Dump(ibmCluster)

	// Create the machine scope
	machineScope, err := scope.NewPowerVSMachineScope(scope.PowerVSMachineScopeParams{
		Client:            r.Client,
		Logger:            log,
		Cluster:           cluster,
		IBMPowerVSCluster: ibmCluster,
		Machine:           machine,
		IBMPowerVSMachine: ibmPowerVSMachine,
	})
	if err != nil {
		return ctrl.Result{}, errors.Errorf("failed to create scope: %+v", err)
	}
	// Always close the scope when exiting this function so we can persist any GCPMachine changes.

	defer func() {
		if err := machineScope.Close(); err != nil && reterr == nil {
			reterr = err
		}
	}()

	// Handle deleted machines
	if !ibmPowerVSMachine.ObjectMeta.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(machineScope)
	}

	// Handle non-deleted machines
	return r.reconcileNormal(ctx, machineScope)
}

func (r *IBMPowerVSMachineReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha3.IBMPowerVSMachine{}).
		Complete(r)
}

func (r *IBMPowerVSMachineReconciler) reconcileDelete(scope *scope.PowerVSMachineScope) (_ ctrl.Result, reterr error) {
	scope.Info("Handling deleted IBMPowerVSMachine")

	if err := scope.DeleteMachine(); err != nil {
		scope.Info("error deleting IBMPowerVSMachine")
		return ctrl.Result{}, errors.Wrapf(err, "error deleting IBMPowerVSMachine %s/%s", scope.IBMPowerVSMachine.Namespace, scope.IBMPowerVSMachine.Spec.Name)
	}

	defer func() {
		if reterr == nil {
			// VSI is deleted so remove the finalizer.
			controllerutil.RemoveFinalizer(scope.IBMPowerVSMachine, infrastructurev1alpha3.MachineFinalizer)
		}
	}()

	return ctrl.Result{}, nil
}

func (r *IBMPowerVSMachineReconciler) getOrCreate(scope *scope.PowerVSMachineScope) (*models.PVMInstanceReference, error) {
	instance, err := scope.CreateMachine()
	return instance, err
}

func (r *IBMPowerVSMachineReconciler) reconcileNormal(ctx context.Context, machineScope *scope.PowerVSMachineScope) (ctrl.Result, error) {
	controllerutil.AddFinalizer(machineScope.IBMPowerVSMachine, infrastructurev1alpha3.MachineFinalizer)

	// Make sure bootstrap data is available and populated.
	if machineScope.Machine.Spec.Bootstrap.DataSecretName == nil {
		machineScope.Info("Bootstrap data secret reference is not yet available")
		return ctrl.Result{}, nil
	}

	ins, err := r.getOrCreate(machineScope)
	if err != nil {
		return ctrl.Result{}, errors.Wrapf(err, "failed to reconcile VSI for IBMVPCMachine %s/%s", machineScope.IBMPowerVSMachine.Namespace, machineScope.IBMPowerVSMachine.Name)
	}

	if ins != nil {
		instance, err := machineScope.IBMPowerVSClient.InstanceClient.Get(*ins.PvmInstanceID, machineScope.IBMPowerVSMachine.Spec.CloudInstanceID, 60*time.Minute)
		if err != nil {
			return ctrl.Result{}, err
		}
		machineScope.IBMPowerVSMachine.Status.InstanceID = *instance.PvmInstanceID
		var addresses []v1.NodeAddress
		for _, network := range instance.Networks {
			addresses = append(addresses, v1.NodeAddress{
				Type:    v1.NodeInternalIP,
				Address: network.IPAddress,
			})
			if network.ExternalIP != "" {
				addresses = append(addresses, v1.NodeAddress{
					Type:    v1.NodeExternalIP,
					Address: network.ExternalIP,
				})
			}
		}
		machineScope.IBMPowerVSMachine.Status.Addresses = addresses
		if instance.Health != nil {
			machineScope.IBMPowerVSMachine.Status.Health = instance.Health.Status
		}
		machineScope.IBMPowerVSMachine.Status.InstanceState = *instance.Status
		if machineScope.IBMPowerVSMachine.Status.InstanceState == "ACTIVE" {
			machineScope.IBMPowerVSMachine.Status.Ready = true
		}
		machineScope.Info(*ins.PvmInstanceID)
	}
	_, ok := machineScope.IBMPowerVSMachine.Labels[clusterv1.MachineControlPlaneLabelName]
	machineScope.IBMPowerVSMachine.Spec.ProviderID = pointer.StringPtr(fmt.Sprintf("ibmpowervs://%s/%s", machineScope.Machine.Spec.ClusterName, machineScope.IBMPowerVSMachine.Name))
	if ok {
		//machineScope.IBMPowerVSMachine.Status.Ready = true
		for _, address := range machineScope.IBMPowerVSMachine.Status.Addresses {
			//spew.Dump(machineScope)
			//if address.Type == v1.NodeExternalIP {
			if address.Type == v1.NodeInternalIP {
				machineScope.IBMPowerVSCluster.Status.APIEndpoint = infrastructurev1alpha3.PowerVSAPIEndpoint{
					Address: &address.Address,
				}
				machineScope.IBMPowerVSCluster.Spec.ControlPlaneEndpoint = clusterv1.APIEndpoint{
					Host: address.Address,
					Port: 6443,
				}
			}
		}
		if err := r.Client.Update(ctx, machineScope.IBMPowerVSCluster); err != nil {
			return ctrl.Result{}, err
		}
		if err := r.Client.Status().Update(ctx, machineScope.IBMPowerVSCluster); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}
