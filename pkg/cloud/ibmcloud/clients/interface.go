package clients

import (
	"fmt"
	"github.com/softlayer/softlayer-go/datatypes"
	"io/ioutil"
	"k8s.io/klog"
	"os"
	//"github.com/davecgh/go-spew/spew"
	"github.com/softlayer/softlayer-go/session"
	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
	clusterv1 "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	ibmcloudv1 "sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/apis/ibmcloud/v1alpha1"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud"
	"sigs.k8s.io/cluster-api-provider-ibmcloud/pkg/cloud/ibmcloud/clients/v1alpha1"
)

const CloudsYamlFile = "/etc/ibmcloud/clouds.yaml"

type GuestService interface {
	CreateGuest(string, string, *ibmcloudv1.IBMCloudMachineProviderSpec, string)
	DeleteGuest(int) error
	GetGuest(string, string) (*datatypes.Virtual_Guest, error)
}

// func NewGuestService(sess *session.Session) *GuestService {
// 	return v1alpha1.NewGuestServiceV1alpha1(sess)
// }

func NewInstanceServiceFromMachine(kubeClient kubernetes.Interface, machine *clusterv1.Machine) (GuestService, error) {
	klog.Info("NewInstanceServiceFromMachine: entry")
	// AuthConfig is mounted into controller pod for clouds authentication
	fileName := CloudsYamlFile
	if _, err := os.Stat(fileName); err != nil {
		return nil, fmt.Errorf("Cannot stat %q: %v", fileName, err)
	}
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("Cannot read %q: %v", fileName, err)
	}

	config := cloud.Config{}
	yaml.Unmarshal(bytes, &config)
	authConfig := config.Clouds.IBMCloud.Auth

	if authConfig.APIUserName == "" || authConfig.AuthenticationKey == "" {
		return nil, fmt.Errorf("Failed getting IBM Cloud config API Username %q, Authentication Key %q", authConfig.APIUserName, authConfig.AuthenticationKey)
	}
	var g GuestService
/*
	fmt.Println("-------------------------------")
	spew.Dump(machine)
	fmt.Println("-------------------------------")
	out := machine.Spec.ProviderSpec.Value.Object.DeepCopyObject()
	switch out.(type) {
	case *ibmcloudv1.IBMCloudMachineProviderSpec:
		sess := session.New(authConfig.APIUserName, authConfig.AuthenticationKey)
		g = v1alpha1.NewGuestServiceV1alpha1(sess)
	default:
		return nil, fmt.Errorf("Invalid type: %v", out)
	}
*/
	sess := session.New(authConfig.APIUserName, authConfig.AuthenticationKey)
	g  = v1alpha1.NewGuestServiceV1alpha1(sess)
	return g, nil
}
