// +build !ignore_autogenerated

/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha4

import (
	"k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IBMPowerVSCluster) DeepCopyInto(out *IBMPowerVSCluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IBMPowerVSCluster.
func (in *IBMPowerVSCluster) DeepCopy() *IBMPowerVSCluster {
	if in == nil {
		return nil
	}
	out := new(IBMPowerVSCluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IBMPowerVSCluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IBMPowerVSClusterList) DeepCopyInto(out *IBMPowerVSClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]IBMPowerVSCluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IBMPowerVSClusterList.
func (in *IBMPowerVSClusterList) DeepCopy() *IBMPowerVSClusterList {
	if in == nil {
		return nil
	}
	out := new(IBMPowerVSClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IBMPowerVSClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IBMPowerVSClusterSpec) DeepCopyInto(out *IBMPowerVSClusterSpec) {
	*out = *in
	out.ControlPlaneEndpoint = in.ControlPlaneEndpoint
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IBMPowerVSClusterSpec.
func (in *IBMPowerVSClusterSpec) DeepCopy() *IBMPowerVSClusterSpec {
	if in == nil {
		return nil
	}
	out := new(IBMPowerVSClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IBMPowerVSClusterStatus) DeepCopyInto(out *IBMPowerVSClusterStatus) {
	*out = *in
	in.APIEndpoint.DeepCopyInto(&out.APIEndpoint)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IBMPowerVSClusterStatus.
func (in *IBMPowerVSClusterStatus) DeepCopy() *IBMPowerVSClusterStatus {
	if in == nil {
		return nil
	}
	out := new(IBMPowerVSClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IBMPowerVSMachine) DeepCopyInto(out *IBMPowerVSMachine) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IBMPowerVSMachine.
func (in *IBMPowerVSMachine) DeepCopy() *IBMPowerVSMachine {
	if in == nil {
		return nil
	}
	out := new(IBMPowerVSMachine)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IBMPowerVSMachine) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IBMPowerVSMachineList) DeepCopyInto(out *IBMPowerVSMachineList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]IBMPowerVSMachine, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IBMPowerVSMachineList.
func (in *IBMPowerVSMachineList) DeepCopy() *IBMPowerVSMachineList {
	if in == nil {
		return nil
	}
	out := new(IBMPowerVSMachineList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IBMPowerVSMachineList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IBMPowerVSMachineSpec) DeepCopyInto(out *IBMPowerVSMachineSpec) {
	*out = *in
	if in.ProviderID != nil {
		in, out := &in.ProviderID, &out.ProviderID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IBMPowerVSMachineSpec.
func (in *IBMPowerVSMachineSpec) DeepCopy() *IBMPowerVSMachineSpec {
	if in == nil {
		return nil
	}
	out := new(IBMPowerVSMachineSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IBMPowerVSMachineStatus) DeepCopyInto(out *IBMPowerVSMachineStatus) {
	*out = *in
	if in.Addresses != nil {
		in, out := &in.Addresses, &out.Addresses
		*out = make([]v1.NodeAddress, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IBMPowerVSMachineStatus.
func (in *IBMPowerVSMachineStatus) DeepCopy() *IBMPowerVSMachineStatus {
	if in == nil {
		return nil
	}
	out := new(IBMPowerVSMachineStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IBMPowerVSMachineTemplate) DeepCopyInto(out *IBMPowerVSMachineTemplate) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IBMPowerVSMachineTemplate.
func (in *IBMPowerVSMachineTemplate) DeepCopy() *IBMPowerVSMachineTemplate {
	if in == nil {
		return nil
	}
	out := new(IBMPowerVSMachineTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IBMPowerVSMachineTemplate) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IBMPowerVSMachineTemplateList) DeepCopyInto(out *IBMPowerVSMachineTemplateList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]IBMPowerVSMachineTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IBMPowerVSMachineTemplateList.
func (in *IBMPowerVSMachineTemplateList) DeepCopy() *IBMPowerVSMachineTemplateList {
	if in == nil {
		return nil
	}
	out := new(IBMPowerVSMachineTemplateList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *IBMPowerVSMachineTemplateList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IBMPowerVSMachineTemplateResource) DeepCopyInto(out *IBMPowerVSMachineTemplateResource) {
	*out = *in
	in.Spec.DeepCopyInto(&out.Spec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IBMPowerVSMachineTemplateResource.
func (in *IBMPowerVSMachineTemplateResource) DeepCopy() *IBMPowerVSMachineTemplateResource {
	if in == nil {
		return nil
	}
	out := new(IBMPowerVSMachineTemplateResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IBMPowerVSMachineTemplateSpec) DeepCopyInto(out *IBMPowerVSMachineTemplateSpec) {
	*out = *in
	in.Template.DeepCopyInto(&out.Template)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IBMPowerVSMachineTemplateSpec.
func (in *IBMPowerVSMachineTemplateSpec) DeepCopy() *IBMPowerVSMachineTemplateSpec {
	if in == nil {
		return nil
	}
	out := new(IBMPowerVSMachineTemplateSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *IBMPowerVSMachineTemplateStatus) DeepCopyInto(out *IBMPowerVSMachineTemplateStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new IBMPowerVSMachineTemplateStatus.
func (in *IBMPowerVSMachineTemplateStatus) DeepCopy() *IBMPowerVSMachineTemplateStatus {
	if in == nil {
		return nil
	}
	out := new(IBMPowerVSMachineTemplateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PowerVSAPIEndpoint) DeepCopyInto(out *PowerVSAPIEndpoint) {
	*out = *in
	if in.Address != nil {
		in, out := &in.Address, &out.Address
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PowerVSAPIEndpoint.
func (in *PowerVSAPIEndpoint) DeepCopy() *PowerVSAPIEndpoint {
	if in == nil {
		return nil
	}
	out := new(PowerVSAPIEndpoint)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PowerVSVIP) DeepCopyInto(out *PowerVSVIP) {
	*out = *in
	if in.Address != nil {
		in, out := &in.Address, &out.Address
		*out = new(string)
		**out = **in
	}
	if in.ExternalAddress != nil {
		in, out := &in.ExternalAddress, &out.ExternalAddress
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PowerVSVIP.
func (in *PowerVSVIP) DeepCopy() *PowerVSVIP {
	if in == nil {
		return nil
	}
	out := new(PowerVSVIP)
	in.DeepCopyInto(out)
	return out
}
