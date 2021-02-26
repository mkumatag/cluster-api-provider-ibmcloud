package scope

import (
	"github.com/IBM-Cloud/power-go-client/clients/instance"
	"github.com/IBM-Cloud/power-go-client/ibmpisession"
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
