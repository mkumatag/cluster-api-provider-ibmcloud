package scope

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/IBM-Cloud/power-go-client/power/client/p_cloud_p_vm_instances"
	"github.com/IBM-Cloud/power-go-client/power/models"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"strconv"
	"time"

	"github.com/go-logr/logr"
	infrav1 "github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/api/v1alpha3"
	"github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/pkg"
	"github.com/pkg/errors"
	"github.com/ppc64le-cloud/powervs-utils"
	"k8s.io/klog/v2/klogr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha4"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type PowerVSMachineScopeParams struct {
	Logger            logr.Logger
	Client            client.Client
	Cluster           *clusterv1.Cluster
	Machine           *clusterv1.Machine
	IBMPowerVSCluster *infrav1.IBMPowerVSCluster
	IBMPowerVSMachine *infrav1.IBMPowerVSMachine
}

type PowerVSMachineScope struct {
	logr.Logger
	client      client.Client
	patchHelper *patch.Helper

	IBMPowerVSClient  *IBMPowerVSClient
	Cluster           *clusterv1.Cluster
	Machine           *clusterv1.Machine
	IBMPowerVSCluster *infrav1.IBMPowerVSCluster
	IBMPowerVSMachine *infrav1.IBMPowerVSMachine
}

func NewPowerVSMachineScope(params PowerVSMachineScopeParams) (*PowerVSMachineScope, error) {
	if params.Client == nil {
		return nil, errors.New("client is required when creating a MachineScope")
	}
	if params.Machine == nil {
		return nil, errors.New("machine is required when creating a MachineScope")
	}
	if params.Cluster == nil {
		return nil, errors.New("cluster is required when creating a MachineScope")
	}
	if params.IBMPowerVSMachine == nil {
		return nil, errors.New("aws machine is required when creating a MachineScope")
	}

	if params.Logger == nil {
		params.Logger = klogr.New()
	}

	m := params.IBMPowerVSMachine

	resource, err := pkg.IBMCloud.ResourceClient.GetInstance(m.Spec.CloudInstanceID)
	if err != nil {
		return nil, err
	}
	region, err := utils.GetRegion(resource.RegionID)
	if err != nil {
		return nil, err
	}
	zone := resource.RegionID

	c, err := NewIBMPowerVSClient(pkg.IBMCloud.Config.IAMAccessToken, pkg.IBMCloud.User.Account, m.Spec.CloudInstanceID, region, zone, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create NewIBMPowerVSClient")
	}

	helper, err := patch.NewHelper(params.IBMPowerVSMachine, params.Client)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init patch helper")
	}
	return &PowerVSMachineScope{
		Logger:      params.Logger,
		client:      params.Client,
		patchHelper: helper,

		IBMPowerVSClient:  c,
		Cluster:           params.Cluster,
		Machine:           params.Machine,
		IBMPowerVSMachine: params.IBMPowerVSMachine,
		IBMPowerVSCluster: params.IBMPowerVSCluster,
	}, nil
}

func (m *PowerVSMachineScope) ensureInstanceUnique(instanceName string) (*models.PVMInstanceReference, error) {
	//resource, err := pkg.IBMCloud.ResourceClient.GetInstance(m.IBMPowerVSMachine.Spec.CloudInstanceID)
	//if err != nil {
	//	return nil, err
	//}
	//region, err := utils.GetRegion(resource.RegionID)
	//if err != nil {
	//	return nil, err
	//}
	//session, err := ibmpisession.New(pkg.IBMCloud.Config.IAMAccessToken, region, true, 60*time.Minute, pkg.IBMCloud.User.Account, resource.RegionID)
	//if err != nil {
	//	return nil, err
	//}
	//
	//InstanceClient := instance.NewIBMPIInstanceClient(session, m.IBMPowerVSMachine.Spec.CloudInstanceID)
	instances, err := m.IBMPowerVSClient.InstanceClient.GetAll(m.IBMPowerVSMachine.Spec.CloudInstanceID, 60*time.Minute)
	if err != nil {
		return nil, err
	}
	for _, ins := range instances.PvmInstances {
		fmt.Printf("ServerName: %s, compare with: %s", *ins.ServerName, instanceName)
		if *ins.ServerName == instanceName {
			return ins, nil
		}
	}
	return nil, nil
}

func (m *PowerVSMachineScope) CreateMachine() (*models.PVMInstanceReference, error) {
	s := m.IBMPowerVSMachine.Spec

	//address := ""
	//// Get the IP address from the cluster status APIEndpoint if machine is a controller-plane
	//_, ok := m.IBMPowerVSMachine.Labels[clusterv1.MachineControlPlaneLabelName]
	//if ok && m.IBMPowerVSCluster.Status.APIEndpoint.InternalAddress != nil {
	//	address = *m.IBMPowerVSCluster.Status.APIEndpoint.InternalAddress
	//}

	instanceReply, err := m.ensureInstanceUnique(m.IBMPowerVSMachine.Name)
	if err != nil {
		return nil, err
	} else {
		if instanceReply != nil {
			//TODO need a resonable wraped error
			return instanceReply, nil
		}
	}
	cloudInitData, err := m.GetBootstrapData()
	if err != nil {
		return nil, err
	}

	memory, err := strconv.ParseFloat(s.Memory, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert memory(%s) to float64", s.Memory)
	}
	cores, err := strconv.ParseFloat(s.Cores, 64)
	if err != nil {
		return nil, fmt.Errorf("failed to convert Cores(%s) to float64", s.Cores)
	}

	params := &p_cloud_p_vm_instances.PcloudPvminstancesPostParams{
		Body: &models.PVMInstanceCreate{
			ImageID:     &s.Image,
			KeyPairName: s.SSHKey,
			Networks: []*models.PVMInstanceAddNetwork{
				{
					NetworkID: &s.Network,
					//IPAddress: address,
				},
			},
			ServerName: &s.Name,
			Memory:     &memory,
			Processors: &cores,
			ProcType:   &s.Processor,
			SysType:    s.MachineType,
			UserData:   cloudInitData,
		},
	}
	_, err = m.IBMPowerVSClient.InstanceClient.Create(params, s.CloudInstanceID, time.Hour)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// Close closes the current scope persisting the cluster configuration and status.
func (m *PowerVSMachineScope) Close() error {
	return m.PatchObject()
}

// PatchObject persists the cluster configuration and status.
func (m *PowerVSMachineScope) PatchObject() error {
	return m.patchHelper.Patch(context.TODO(), m.IBMPowerVSMachine)
}

func (m *PowerVSMachineScope) DeleteMachine() error {
	return nil
}

//var hackCloudData = `#cloud-config
//runcmd:
//  - 'kubeadm init'
//`

// GetBootstrapData returns the base64 encoded bootstrap data from the secret in the Machine's bootstrap.dataSecretName
func (m *PowerVSMachineScope) GetBootstrapData() (string, error) {
	if m.Machine.Spec.Bootstrap.DataSecretName == nil {
		return "", errors.New("error retrieving bootstrap data: linked Machine's bootstrap.dataSecretName is nil")
	}

	secret := &corev1.Secret{}
	key := types.NamespacedName{Namespace: m.Machine.Namespace, Name: *m.Machine.Spec.Bootstrap.DataSecretName}
	if err := m.client.Get(context.TODO(), key, secret); err != nil {
		return "", errors.Wrapf(err, "failed to retrieve bootstrap data secret for IBMVPCMachine %s/%s", m.Machine.Namespace, m.Machine.Name)
	}

	value, ok := secret.Data["value"]
	if !ok {
		return "", errors.New("error retrieving bootstrap data: secret value key is missing")
	}

	return base64.StdEncoding.EncodeToString(value), nil
	//return base64.StdEncoding.EncodeToString([]byte(hackCloudData)), nil
}
