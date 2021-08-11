package v1alpha3

import (
	"github.com/kubernetes-sigs/cluster-api-provider-ibmcloud/api/v1alpha4"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the v1alpha3 IBMVPCCluster receiver to a v1alpha4 IBMVPCCluster.
func (r *IBMVPCCluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.IBMVPCCluster)

	return Convert_v1alpha3_IBMVPCCluster_To_v1alpha4_IBMVPCCluster(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 AWSClusterList receiver to a v1alpha3 AWSClusterList.
func (r *IBMVPCCluster) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.IBMVPCCluster)

	return Convert_v1alpha4_IBMVPCCluster_To_v1alpha3_IBMVPCCluster(src, r, nil)
}

// ConvertTo converts the v1alpha3 IBMVPCClusterList receiver to a v1alpha4 IBMVPCClusterList.
func (r *IBMVPCClusterList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.IBMVPCClusterList)

	return Convert_v1alpha3_IBMVPCClusterList_To_v1alpha4_IBMVPCClusterList(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 IBMVPCClusterList receiver to a v1alpha3 IBMVPCClusterList.
func (r *IBMVPCClusterList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.IBMVPCClusterList)

	return Convert_v1alpha4_IBMVPCClusterList_To_v1alpha3_IBMVPCClusterList(src, r, nil)
}

// ConvertTo converts the v1alpha3 IBMVPCMachine receiver to a v1alpha4 IBMVPCMachine.
func (r *IBMVPCMachine) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.IBMVPCMachine)

	return Convert_v1alpha3_IBMVPCMachine_To_v1alpha4_IBMVPCMachine(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 IBMVPCMachine receiver to a v1alpha3 IBMVPCMachine.
func (r *IBMVPCMachine) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.IBMVPCMachine)

	return Convert_v1alpha4_IBMVPCMachine_To_v1alpha3_IBMVPCMachine(src, r, nil)
}

// ConvertTo converts the v1alpha3 IBMVPCMachineList receiver to a v1alpha4 IBMVPCMachineList.
func (r *IBMVPCMachineList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.IBMVPCMachineList)

	return Convert_v1alpha3_IBMVPCMachineList_To_v1alpha4_IBMVPCMachineList(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 IBMVPCMachineList receiver to a v1alpha3 IBMVPCMachineList.
func (r *IBMVPCMachineList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.IBMVPCMachineList)

	return Convert_v1alpha4_IBMVPCMachineList_To_v1alpha3_IBMVPCMachineList(src, r, nil)
}

// ConvertTo converts the v1alpha3 IBMVPCMachineTemplate receiver to a v1alpha4 IBMVPCMachineTemplate.
func (r *IBMVPCMachineTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.IBMVPCMachineTemplate)

	return Convert_v1alpha3_IBMVPCMachineTemplate_To_v1alpha4_IBMVPCMachineTemplate(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 IBMVPCMachineTemplate receiver to a v1alpha3 IBMVPCMachineTemplate.
func (r *IBMVPCMachineTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.IBMVPCMachineTemplate)

	return Convert_v1alpha4_IBMVPCMachineTemplate_To_v1alpha3_IBMVPCMachineTemplate(src, r, nil)
}

// ConvertTo converts the v1alpha3 IBMVPCMachineTemplateList receiver to a v1alpha4 IBMVPCMachineTemplateList.
func (r *IBMVPCMachineTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha4.IBMVPCMachineTemplateList)

	return Convert_v1alpha3_IBMVPCMachineTemplateList_To_v1alpha4_IBMVPCMachineTemplateList(r, dst, nil)
}

// ConvertFrom converts the v1alpha4 IBMVPCMachineTemplateList receiver to a v1alpha3 IBMVPCMachineTemplateList.
func (r *IBMVPCMachineTemplateList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha4.IBMVPCMachineTemplateList)

	return Convert_v1alpha4_IBMVPCMachineTemplateList_To_v1alpha3_IBMVPCMachineTemplateList(src, r, nil)
}
