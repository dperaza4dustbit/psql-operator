/*
Copyright 2022.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// PSQLInstanceSpec defines the desired state of PSQLInstance
type PSQLInstanceSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// DatabaseName is the database name of the PSQLInstance.
	DatabaseName string `json:"databasename,omitempty"`
	// UserName is the user name of the PSQLInstance.
	UserName string `json:"username,omitempty"`
}

// Binding defines the service binding status pointing to binding secret
type Binding struct {

	//Name is the name of the binding secret
	Name string `json:"name,omitempty"`
}

// PSQLInstanceStatus defines the observed state of PSQLInstance
type PSQLInstanceStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Binding is the object pointing to binding secret
	Binding Binding `json:"binding,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PSQLInstance is the Schema for the psqlinstances API
type PSQLInstance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PSQLInstanceSpec   `json:"spec,omitempty"`
	Status PSQLInstanceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PSQLInstanceList contains a list of PSQLInstance
type PSQLInstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PSQLInstance `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PSQLInstance{}, &PSQLInstanceList{})
}
