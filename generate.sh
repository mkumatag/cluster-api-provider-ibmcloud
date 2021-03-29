#!/bin/bash

IBMPOWERVS_SERVICE_INSTANCE_ID=7845d372-d4e1-46b8-91fc-41051c984601 \
IBMPOWERVS_NETWORK_ID=f3895fde-e53e-4243-afa7-e11a1fb33726 \
IBMPOWERVS_VIP=192.168.151.101 \
IBMPOWERVS_VIP_EXTERNAL=158.175.162.101 \
IBMPOWERVS_SSHKEY_NAME=mkumatag-pub-key \
IBMPOWERVS_IMAGE_ID=71a00ad7-58e6-4903-90f3-82437bc1f145 \
clusterctl config cluster ibm-powervs-1 --kubernetes-version v1.19.8 \
--target-namespace default \
--control-plane-machine-count=3 \
--worker-machine-count=1 \
--from ./templates/cluster-template-powervs.yaml > deploy.yaml
