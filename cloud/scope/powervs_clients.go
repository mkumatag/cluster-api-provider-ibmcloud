package scope

import (
	"fmt"
	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
	infrav1 "github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/api/v1alpha3"
	"github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/pkg"
	"os"
	"time"
)

const TIMEOUT = 1 * time.Hour

type IBMPowerVSClient struct {
	region string
	zone   string

	session        *ibmpisession.IBMPISession
	InstanceClient *instance.IBMPIInstanceClient
	NetworkClient  *instance.IBMPINetworkClient
}

func NewIBMPowerVSClient(token, account, cloudInstanceID, region, zone string, debug bool) (_ *IBMPowerVSClient, err error) {
	client := &IBMPowerVSClient{}
	client.session, err = ibmpisession.New(token, region, debug, TIMEOUT, account, zone)
	if err != nil {
		return nil, err
	}

	client.InstanceClient = instance.NewIBMPIInstanceClient(client.session, cloudInstanceID)
	client.NetworkClient = instance.NewIBMPINetworkClient(client.session, cloudInstanceID)
	return client, nil
}

func NewIBMVPCClients(cluster *infrav1.IBMPowerVSCluster) (*IBMVPCClients, error) {
	iamEndpoint := os.Getenv("IAM_ENDPOINT")
	if iamEndpoint == "" {
		iamEndpoint = "https://iam.cloud.ibm.com/identity/token"
	}
	svcEndpoint := os.Getenv("SERVICE_ENDPOINT")
	if svcEndpoint == "" {
		svcEndpoint = fmt.Sprintf("https://%s.iaas.cloud.ibm.com/v1", cluster.Spec.VPCRegion)
	}

	client := &IBMVPCClients{}

	vpcErr := client.setIBMVPCService(iamEndpoint, svcEndpoint, pkg.IBMCloud.Config.BluemixAPIKey)
	if vpcErr != nil {
		return nil, vpcErr
	}
	return client, nil
}
