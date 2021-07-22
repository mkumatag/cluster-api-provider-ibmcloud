package scope

import (
	"context"
	"fmt"
	"github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/pkg"
	utils "github.com/ppc64le-cloud/powervs-utils"

	"github.com/go-logr/logr"
	infrav1 "github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/api/v1alpha3"
	"github.com/pkg/errors"
	"k8s.io/klog/v2/klogr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PowerVSClusterScopeParams struct {
	Client            client.Client
	Logger            logr.Logger
	Cluster           *clusterv1.Cluster
	IBMPowerVSCluster *infrav1.IBMPowerVSCluster
}

type PowerVSClusterScope struct {
	logr.Logger
	client      client.Client
	patchHelper *patch.Helper

	IBMPowerVSClient  *IBMPowerVSClient
	Cluster           *clusterv1.Cluster
	IBMPowerVSCluster *infrav1.IBMPowerVSCluster
}

func NewPowerVSClusterScope(params PowerVSClusterScopeParams) (*PowerVSClusterScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}
	if params.IBMPowerVSCluster == nil {
		return nil, errors.New("failed to generate new scope from nil IBMVPCCluster")
	}

	if params.Logger == nil {
		params.Logger = klogr.New()
	}

	spec := params.IBMPowerVSCluster.Spec
	resource, err := pkg.IBMCloud.ResourceClient.GetInstance(spec.CloudInstanceID)
	if err != nil {
		return nil, err
	}
	region, err := utils.GetRegion(resource.RegionID)
	if err != nil {
		return nil, err
	}
	zone := resource.RegionID

	c, err := NewIBMPowerVSClient(pkg.IBMCloud.Config.IAMAccessToken, pkg.IBMCloud.User.Account, spec.CloudInstanceID, region, zone, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create NewIBMPowerVSClient")
	}

	helper, err := patch.NewHelper(params.IBMPowerVSCluster, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}

	return &PowerVSClusterScope{
		Logger:            params.Logger,
		client:            params.Client,
		IBMPowerVSClient:  c,
		Cluster:           params.Cluster,
		IBMPowerVSCluster: params.IBMPowerVSCluster,
		patchHelper:       helper,
	}, nil
}

// PatchObject persists the cluster configuration and status.
func (s *PowerVSClusterScope) PatchObject() error {
	return s.patchHelper.Patch(context.TODO(), s.IBMPowerVSCluster)
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *PowerVSClusterScope) Close() error {
	return s.PatchObject()
}
