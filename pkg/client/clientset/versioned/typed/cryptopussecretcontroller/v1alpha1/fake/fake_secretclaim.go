/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/puzzle/cryptopus-k8s-secretcontroller/pkg/apis/cryptopussecretcontroller/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSecretClaims implements SecretClaimInterface
type FakeSecretClaims struct {
	Fake *FakeCryptopussecretcontrollerV1alpha1
	ns   string
}

var secretclaimsResource = schema.GroupVersionResource{Group: "cryptopussecretcontroller.puzzle.ch", Version: "v1alpha1", Resource: "secretclaims"}

var secretclaimsKind = schema.GroupVersionKind{Group: "cryptopussecretcontroller.puzzle.ch", Version: "v1alpha1", Kind: "SecretClaim"}

// Get takes name of the secretClaim, and returns the corresponding secretClaim object, and an error if there is any.
func (c *FakeSecretClaims) Get(name string, options v1.GetOptions) (result *v1alpha1.SecretClaim, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(secretclaimsResource, c.ns, name), &v1alpha1.SecretClaim{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SecretClaim), err
}

// List takes label and field selectors, and returns the list of SecretClaims that match those selectors.
func (c *FakeSecretClaims) List(opts v1.ListOptions) (result *v1alpha1.SecretClaimList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(secretclaimsResource, secretclaimsKind, c.ns, opts), &v1alpha1.SecretClaimList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.SecretClaimList{ListMeta: obj.(*v1alpha1.SecretClaimList).ListMeta}
	for _, item := range obj.(*v1alpha1.SecretClaimList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested secretClaims.
func (c *FakeSecretClaims) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(secretclaimsResource, c.ns, opts))

}

// Create takes the representation of a secretClaim and creates it.  Returns the server's representation of the secretClaim, and an error, if there is any.
func (c *FakeSecretClaims) Create(secretClaim *v1alpha1.SecretClaim) (result *v1alpha1.SecretClaim, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(secretclaimsResource, c.ns, secretClaim), &v1alpha1.SecretClaim{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SecretClaim), err
}

// Update takes the representation of a secretClaim and updates it. Returns the server's representation of the secretClaim, and an error, if there is any.
func (c *FakeSecretClaims) Update(secretClaim *v1alpha1.SecretClaim) (result *v1alpha1.SecretClaim, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(secretclaimsResource, c.ns, secretClaim), &v1alpha1.SecretClaim{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SecretClaim), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSecretClaims) UpdateStatus(secretClaim *v1alpha1.SecretClaim) (*v1alpha1.SecretClaim, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(secretclaimsResource, "status", c.ns, secretClaim), &v1alpha1.SecretClaim{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SecretClaim), err
}

// Delete takes name of the secretClaim and deletes it. Returns an error if one occurs.
func (c *FakeSecretClaims) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(secretclaimsResource, c.ns, name), &v1alpha1.SecretClaim{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSecretClaims) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(secretclaimsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.SecretClaimList{})
	return err
}

// Patch applies the patch and returns the patched secretClaim.
func (c *FakeSecretClaims) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SecretClaim, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(secretclaimsResource, c.ns, name, pt, data, subresources...), &v1alpha1.SecretClaim{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.SecretClaim), err
}
