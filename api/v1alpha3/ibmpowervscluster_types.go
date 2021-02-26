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

package v1alpha3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// IBMPowerVSClusterSpec defines the desired state of IBMPowerVSCluster
type IBMPowerVSClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// CloudInstanceID is the id of the power cloud instance where the vsi instance will get deployed
	CloudInstanceID string `json:"cloudInstanceID"`

	// Network is network ID used for the VSI
	Network string `json:"network"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`
}

// IBMPowerVSClusterStatus defines the observed state of IBMPowerVSCluster
type IBMPowerVSClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Ready       bool               `json:"ready"`
	APIEndpoint PowerVSAPIEndpoint `json:"apiEndpoint,omitempty"`
}

type PowerVSAPIEndpoint struct {
	Address *string `json:"address"`
	//InternalAddress *string `json:"internalAddress"`
	// PortID is the ID for the network port gets created
	//PortID *string `json:"portID"`
}

// +kubebuilder:subresource:status
// +kubebuilder:object:root=true

// IBMPowerVSCluster is the Schema for the ibmpowervsclusters API
type IBMPowerVSCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IBMPowerVSClusterSpec   `json:"spec,omitempty"`
	Status IBMPowerVSClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// IBMPowerVSClusterList contains a list of IBMPowerVSCluster
type IBMPowerVSClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IBMPowerVSCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&IBMPowerVSCluster{}, &IBMPowerVSClusterList{})
}