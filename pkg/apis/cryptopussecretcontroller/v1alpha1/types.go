/*
Copyright 2017 The Kubernetes Authors.

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

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecretClaim is a specification for a SecretClaim resource
type SecretClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SecretClaimSpec   `json:"spec"`
	Status SecretClaimStatus `json:"status"`
}

// SecretClaimSpec is the spec for a Foo resource
type SecretClaimSpec struct {
	Id int32 `json:"id"`
	SecretName string `json:"secretName"`
	RefreshTime int32 `json:"refreshTime"`
	CryptopusSecret string `json:"cryptopusSecret"`
}

// SecretClaimStatus is the status for a Foo resource
type SecretClaimStatus struct {
	LastUpdate int32 `json:"lastUpdate"`
	Phase string `json:"phase", default:"init"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecretClaimList is a list of Foo resources
type SecretClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []SecretClaim `json:"items"`
}
