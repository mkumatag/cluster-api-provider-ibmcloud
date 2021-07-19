package scope

import (
	"context"
	"fmt"
	"github.com/IBM/go-sdk-core/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/go-logr/logr"
	infrav1 "github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/api/v1alpha3"
	"github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/pkg"
	"github.com/pkg/errors"
	utils "github.com/ppc64le-cloud/powervs-utils"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/klog/klogr"
	"k8s.io/klog/v2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type PowerVSClusterScopeParams struct {
	Client            client.Client
	Logger            logr.Logger
	Cluster           *clusterv1.Cluster
	IBMPowerVSCluster *infrav1.IBMPowerVSCluster
}

type PowerVSClusterScope struct {
	*IBMVPCClients

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

	vpcClients, err := NewIBMVPCClients(params.IBMPowerVSCluster)
	if err != nil {
		return nil, fmt.Errorf("failed to create NewIBMVPCClients")
	}

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
		IBMVPCClients:     vpcClients,
		IBMPowerVSClient:  c,
		Cluster:           params.Cluster,
		IBMPowerVSCluster: params.IBMPowerVSCluster,
		patchHelper:       helper,
	}, nil
}

func (s *PowerVSClusterScope) xx() error {
	return nil
}

// APIServerPort returns the APIServerPort to use when creating the load balancer.
func (s *PowerVSClusterScope) APIServerPort() int32 {
	if s.Cluster.Spec.ClusterNetwork != nil && s.Cluster.Spec.ClusterNetwork.APIServerPort != nil {
		return *s.Cluster.Spec.ClusterNetwork.APIServerPort
	}
	return 6443
}

// ListSubnets lists all the subnets in a given vpcID
func (s *PowerVSClusterScope) ListSubnets(vpcID string) (subnets []vpcv1.SubnetIdentityIntf, err error) {
	subnetCollection, _, err := s.IBMVPCClients.VPCService.ListSubnets(&vpcv1.ListSubnetsOptions{})
	if err != nil {
		return
	}
	for _, subnet := range subnetCollection.Subnets {
		if subnet.VPC != nil && *subnet.VPC.ID == vpcID {
			id, _ := s.IBMVPCClients.VPCService.NewSubnetIdentityByID(*subnet.ID)
			subnets = append(subnets, id)
		}
	}
	return
}

// CreateSecurityGroup creates a security group, by default allows all the outboud traffic and mentioned traffic for
// inbound port for an vpcID
func (s *PowerVSClusterScope) CreateSecurityGroup(vpcID string, port int64) (secGroup *vpcv1.SecurityGroup, err error) {
	secOpt := s.IBMVPCClients.VPCService.NewCreateSecurityGroupOptions(
		&vpcv1.VPCIdentityByID{
			ID: &vpcID,
		})
	secOpt.SetRules(
		[]vpcv1.SecurityGroupRulePrototypeIntf{
			&vpcv1.SecurityGroupRulePrototype{
				Direction: core.StringPtr(vpcv1.SecurityGroupRulePrototypeDirectionInboundConst),
				IPVersion: core.StringPtr(vpcv1.SecurityGroupRulePrototypeIPVersionIpv4Const),
				Protocol:  core.StringPtr("tcp"),
				PortMin:   core.Int64Ptr(port),
				PortMax:   core.Int64Ptr(port),
			},
			&vpcv1.SecurityGroupRulePrototypeSecurityGroupRuleProtocolAll{
				Direction: core.StringPtr(vpcv1.SecurityGroupRulePrototypeSecurityGroupRuleProtocolAllDirectionOutboundConst),
				IPVersion: core.StringPtr(vpcv1.SecurityGroupRulePrototypeSecurityGroupRuleProtocolAllIPVersionIpv4Const),
				Protocol:  core.StringPtr(vpcv1.SecurityGroupRulePrototypeSecurityGroupRuleProtocolAllProtocolAllConst),
			},
		})
	secGroup, _, err = s.IBMVPCClients.VPCService.CreateSecurityGroup(secOpt)
	return
}

// CreateLoadBalancer is a function to create an Application Load Balancer with given subnet and security groups
func (s *PowerVSClusterScope) CreateLoadBalancer(
	subnets []vpcv1.SubnetIdentityIntf, securityGroups []vpcv1.SecurityGroupIdentityIntf) (
	lb *vpcv1.LoadBalancer, err error) {
	opt := s.IBMVPCClients.VPCService.NewCreateLoadBalancerOptions(
		true,
		subnets,
	)
	// TODO: set the security groups
	//	opt.SetSecurityGroups(securityGroups)
	lb, _, err = s.IBMVPCClients.VPCService.CreateLoadBalancer(opt)
	return
}

// WaitForUpdateLB is a function to wait till load balancer becomes active and online
func (s *PowerVSClusterScope) WaitForUpdateLB(id *string) error {
	return wait.PollImmediate(1*time.Minute, 10*time.Minute, func() (bool, error) {
		lb, _, err := s.IBMVPCClients.VPCService.GetLoadBalancer(
			&vpcv1.GetLoadBalancerOptions{
				ID: id,
			})
		if err != nil {
			return false, err
		}
		if *lb.ProvisioningStatus == vpcv1.LoadBalancerProvisioningStatusActiveConst &&
			*lb.OperatingStatus == vpcv1.LoadBalancerOperatingStatusOnlineConst {
			return true, nil
		}
		klog.Infof("LB create in-progress, current ProvisioningStatus: %s, OperatingStatus: %s",
			*lb.ProvisioningStatus, *lb.OperatingStatus)
		return false, nil
	})
}

// CreateLoadBalancerPool is a function to create a load balancer pool in the given pool
// with port mentioned for the health check
func (s *PowerVSClusterScope) CreateLoadBalancerPool(lbID *string, port int64) (pool *vpcv1.LoadBalancerPool, err error) {
	pool, _, err = s.IBMVPCClients.VPCService.CreateLoadBalancerPool(
		&vpcv1.CreateLoadBalancerPoolOptions{
			LoadBalancerID: lbID,
			Name:           core.StringPtr("kube-apiserver"),
			Protocol:       core.StringPtr(vpcv1.LoadBalancerPoolPrototypeProtocolTCPConst),
			Algorithm:      core.StringPtr(vpcv1.LoadBalancerPoolPrototypeAlgorithmRoundRobinConst),
			HealthMonitor: &vpcv1.LoadBalancerPoolHealthMonitorPrototype{
				Delay:      core.Int64Ptr(5),
				MaxRetries: core.Int64Ptr(2),
				Timeout:    core.Int64Ptr(2),
				URLPath:    core.StringPtr("/"),
				Port:       core.Int64Ptr(port),
				Type:       core.StringPtr(vpcv1.LoadBalancerPoolHealthMonitorPrototypeTypeTCPConst),
			},
			SessionPersistence: &vpcv1.LoadBalancerPoolSessionPersistencePrototype{
				Type: core.StringPtr(vpcv1.LoadBalancerPoolSessionPersistencePrototypeTypeSourceIPConst),
			},
			//TODO: enable this when updated to latest client
			//ProxyProtocol: core.StringPtr(vpcv1.LoadBalancerPoolPrototypeProxyProtocolDisabledConst),
		})
	return
}

func (s *PowerVSClusterScope) CreateLoadBalancerPoolMember(
	lbID, poolID, ipaddress string, port int64) (member *vpcv1.LoadBalancerPoolMember, err error) {
	member, _, err = s.IBMVPCClients.VPCService.CreateLoadBalancerPoolMember(
		&vpcv1.CreateLoadBalancerPoolMemberOptions{
			LoadBalancerID: &lbID,
			PoolID:         &poolID,
			Port:           core.Int64Ptr(port),
			Target: &vpcv1.LoadBalancerPoolMemberTargetPrototypeIP{
				Address: core.StringPtr(ipaddress),
			},
			// TODO: if missed understand the impact
			//Weight: core.Int64Ptr(100),
		})
	return
}

func (s *PowerVSClusterScope) CreateLoadBalancerListener(
	lbID, defaultPoolID *string, port int64) (listener *vpcv1.LoadBalancerListener, err error) {
	opt := s.IBMVPCClients.VPCService.NewCreateLoadBalancerListenerOptions(
		*lbID,
		port,
		vpcv1.CreateLoadBalancerListenerOptionsProtocolTCPConst)
	opt.SetDefaultPool(&vpcv1.LoadBalancerPoolIdentityByID{ID: defaultPoolID})
	listener, _, err = s.IBMVPCClients.VPCService.CreateLoadBalancerListener(opt)
	return
}

// PatchObject persists the cluster configuration and status.
func (s *PowerVSClusterScope) PatchObject() error {
	return s.patchHelper.Patch(context.TODO(), s.IBMPowerVSCluster)
}

// Close closes the current scope persisting the cluster configuration and status.
func (s *PowerVSClusterScope) Close() error {
	return s.PatchObject()
}

//func (s *PowerVSClusterScope) CreatePort() (*models.NetworkPort, error) {
//	params := &p_cloud_networks.PcloudNetworksPortsPostParams{
//		CloudInstanceID: s.IBMPowerVSCluster.Spec.CloudInstanceID,
//		Body: &models.NetworkPortCreate{
//			Description: s.IBMPowerVSCluster.Name + " Network Port",
//		},
//	}
//	return s.IBMPowerVSClient.NetworkClient.CreatePort(s.IBMPowerVSCluster.Spec.Network, s.IBMPowerVSCluster.Spec.CloudInstanceID, params, TIMEOUT)
//}

//func (s *PowerVSClusterScope) DeletePort() (err error) {
//	if s.IBMPowerVSCluster.Status.APIEndpoint.PortID != nil {
//		_, err = s.IBMPowerVSClient.NetworkClient.DeletePort(s.IBMPowerVSCluster.Spec.Network, s.IBMPowerVSCluster.Spec.CloudInstanceID, *s.IBMPowerVSCluster.Status.APIEndpoint.PortID, TIMEOUT)
//	}
//	return
//}

//func (s *PowerVSClusterScope) GetPort() (*models.NetworkPort, error) {
//	return s.IBMPowerVSClient.NetworkClient.GetPort(s.IBMPowerVSCluster.Spec.Network, s.IBMPowerVSCluster.Spec.CloudInstanceID, *s.IBMPowerVSCluster.Status.APIEndpoint.PortID, TIMEOUT)
//}
