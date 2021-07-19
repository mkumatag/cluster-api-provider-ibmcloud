#!/usr/bin/env bash

set -x

IBMPOWERVS_SSHKEY_NAME="mkumatag-pub-key" \
IBMPOWERVS_VIP="" \
IBMPOWERVS_VIP_EXTERNAL="" \
IBMPOWERVS_VPC_REGION="eu-gb" \
IBMPOWERVS_IMAGE_ID="fb2f75d1-1157-40b9-af2f-5459685ca089" \
IBMPOWERVS_SERVICE_INSTANCE_ID="e449d86e-c3a0-4c07-959e-8557fdf55482" \
IBMPOWERVS_NETWORK_ID="daf2b616-542b-47ed-8cec-ceaec1e90f4d" \
IBMPOWERVS_VPC_ID="r018-d4a39527-a80d-4de2-ab14-fc6f59232751" \
clusterctl config cluster ibm-powervs-1 --kubernetes-version v1.21.2 \
--target-namespace namespace-1 \
--control-plane-machine-count=1 \
--worker-machine-count=0 \
--from ./cluster-template-powervs-lb.yaml