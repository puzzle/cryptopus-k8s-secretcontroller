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

package v1alpha1

import (
	v1alpha1 "github.com/puzzle/cryptopus-k8s-secretcontroller/pkg/apis/cryptopussecretcontroller/v1alpha1"
	scheme "github.com/puzzle/cryptopus-k8s-secretcontroller/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// SecretClaimsGetter has a method to return a SecretClaimInterface.
// A group's client should implement this interface.
type SecretClaimsGetter interface {
	SecretClaims(namespace string) SecretClaimInterface
}

// SecretClaimInterface has methods to work with SecretClaim resources.
type SecretClaimInterface interface {
	Create(*v1alpha1.SecretClaim) (*v1alpha1.SecretClaim, error)
	Update(*v1alpha1.SecretClaim) (*v1alpha1.SecretClaim, error)
	UpdateStatus(*v1alpha1.SecretClaim) (*v1alpha1.SecretClaim, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.SecretClaim, error)
	List(opts v1.ListOptions) (*v1alpha1.SecretClaimList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SecretClaim, err error)
	SecretClaimExpansion
}

// secretClaims implements SecretClaimInterface
type secretClaims struct {
	client rest.Interface
	ns     string
}

// newSecretClaims returns a SecretClaims
func newSecretClaims(c *CryptopussecretcontrollerV1alpha1Client, namespace string) *secretClaims {
	return &secretClaims{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the secretClaim, and returns the corresponding secretClaim object, and an error if there is any.
func (c *secretClaims) Get(name string, options v1.GetOptions) (result *v1alpha1.SecretClaim, err error) {
	result = &v1alpha1.SecretClaim{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("secretclaims").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of SecretClaims that match those selectors.
func (c *secretClaims) List(opts v1.ListOptions) (result *v1alpha1.SecretClaimList, err error) {
	result = &v1alpha1.SecretClaimList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("secretclaims").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested secretClaims.
func (c *secretClaims) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("secretclaims").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a secretClaim and creates it.  Returns the server's representation of the secretClaim, and an error, if there is any.
func (c *secretClaims) Create(secretClaim *v1alpha1.SecretClaim) (result *v1alpha1.SecretClaim, err error) {
	result = &v1alpha1.SecretClaim{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("secretclaims").
		Body(secretClaim).
		Do().
		Into(result)
	return
}

// Update takes the representation of a secretClaim and updates it. Returns the server's representation of the secretClaim, and an error, if there is any.
func (c *secretClaims) Update(secretClaim *v1alpha1.SecretClaim) (result *v1alpha1.SecretClaim, err error) {
	result = &v1alpha1.SecretClaim{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("secretclaims").
		Name(secretClaim.Name).
		Body(secretClaim).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *secretClaims) UpdateStatus(secretClaim *v1alpha1.SecretClaim) (result *v1alpha1.SecretClaim, err error) {
	result = &v1alpha1.SecretClaim{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("secretclaims").
		Name(secretClaim.Name).
		SubResource("status").
		Body(secretClaim).
		Do().
		Into(result)
	return
}

// Delete takes name of the secretClaim and deletes it. Returns an error if one occurs.
func (c *secretClaims) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("secretclaims").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *secretClaims) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("secretclaims").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched secretClaim.
func (c *secretClaims) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.SecretClaim, err error) {
	result = &v1alpha1.SecretClaim{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("secretclaims").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
