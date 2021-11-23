module sigs.k8s.io/cluster-api-provider-ibmcloud

go 1.16

require (
	github.com/IBM-Cloud/bluemix-go v0.0.0-20200921095234-26d1d0148c62
	github.com/IBM-Cloud/power-go-client v1.0.78
	github.com/IBM/go-sdk-core/v5 v5.8.1
	github.com/IBM/vpc-go-sdk v0.13.0
	github.com/go-logr/logr v0.4.0
	github.com/golang-jwt/jwt v3.2.2+incompatible
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.17.0
	github.com/pkg/errors v0.9.1
	github.com/ppc64le-cloud/powervs-utils v0.0.0-20210106101518-5d3f965b0344
	k8s.io/api v0.22.2
	k8s.io/apimachinery v0.22.2
	k8s.io/client-go v0.22.2
	k8s.io/klog/v2 v2.9.0
	k8s.io/kube-openapi v0.0.0-20211110012726-3cc51fd1e909 // indirect
	k8s.io/utils v0.0.0-20210819203725-bdf08cb9a70a
	sigs.k8s.io/cluster-api v0.0.0-20210526191338-0e1629b75111
	sigs.k8s.io/controller-runtime v0.10.3
)

replace sigs.k8s.io/cluster-api => sigs.k8s.io/cluster-api v0.4.4
