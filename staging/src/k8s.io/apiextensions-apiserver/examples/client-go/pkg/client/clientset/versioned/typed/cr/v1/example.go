/*
Copyright The Kubernetes Authors.
Copyright 2020 Authors of Arktos - file modified.

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

package v1

import (
	"time"

	v1 "k8s.io/apiextensions-apiserver/examples/client-go/pkg/apis/cr/v1"
	scheme "k8s.io/apiextensions-apiserver/examples/client-go/pkg/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ExamplesGetter has a method to return a ExampleInterface.
// A group's client should implement this interface.
type ExamplesGetter interface {
	Examples(namespace string) ExampleInterface
	ExamplesWithMultiTenancy(namespace string, tenant string) ExampleInterface
}

// ExampleInterface has methods to work with Example resources.
type ExampleInterface interface {
	Create(*v1.Example) (*v1.Example, error)
	Update(*v1.Example) (*v1.Example, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.Example, error)
	List(opts metav1.ListOptions) (*v1.ExampleList, error)
	Watch(opts metav1.ListOptions) watch.AggregatedWatchInterface
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Example, err error)
	ExampleExpansion
}

// examples implements ExampleInterface
type examples struct {
	client  rest.Interface
	clients []rest.Interface
	ns      string
	te      string
}

// newExamples returns a Examples
func newExamples(c *CrV1Client, namespace string) *examples {
	return newExamplesWithMultiTenancy(c, namespace, "default")
}

func newExamplesWithMultiTenancy(c *CrV1Client, namespace string, tenant string) *examples {
	return &examples{
		client:  c.RESTClient(),
		clients: c.RESTClients(),
		ns:      namespace,
		te:      tenant,
	}
}

// Get takes name of the example, and returns the corresponding example object, and an error if there is any.
func (c *examples) Get(name string, options metav1.GetOptions) (result *v1.Example, err error) {
	result = &v1.Example{}
	err = c.client.Get().
		Tenant(c.te).
		Namespace(c.ns).
		Resource("examples").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)

	return
}

// List takes label and field selectors, and returns the list of Examples that match those selectors.
func (c *examples) List(opts metav1.ListOptions) (result *v1.ExampleList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.ExampleList{}
	err = c.client.Get().
		Tenant(c.te).Namespace(c.ns).
		Resource("examples").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)

	return
}

// Watch returns a watch.Interface that watches the requested examples.
func (c *examples) Watch(opts metav1.ListOptions) watch.AggregatedWatchInterface {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	aggWatch := watch.NewAggregatedWatcher()
	for _, client := range c.clients {
		watcher, err := client.Get().
			Tenant(c.te).
			Namespace(c.ns).
			Resource("examples").
			VersionedParams(&opts, scheme.ParameterCodec).
			Timeout(timeout).
			Watch()
		aggWatch.AddWatchInterface(watcher, err)
	}
	return aggWatch
}

// Create takes the representation of a example and creates it.  Returns the server's representation of the example, and an error, if there is any.
func (c *examples) Create(example *v1.Example) (result *v1.Example, err error) {
	result = &v1.Example{}
	err = c.client.Post().
		Tenant(c.te).
		Namespace(c.ns).
		Resource("examples").
		Body(example).
		Do().
		Into(result)

	return
}

// Update takes the representation of a example and updates it. Returns the server's representation of the example, and an error, if there is any.
func (c *examples) Update(example *v1.Example) (result *v1.Example, err error) {
	result = &v1.Example{}
	err = c.client.Put().
		Tenant(c.te).
		Namespace(c.ns).
		Resource("examples").
		Name(example.Name).
		Body(example).
		Do().
		Into(result)

	return
}

// Delete takes name of the example and deletes it. Returns an error if one occurs.
func (c *examples) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Tenant(c.te).
		Namespace(c.ns).
		Resource("examples").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *examples) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Tenant(c.te).
		Namespace(c.ns).
		Resource("examples").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched example.
func (c *examples) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Example, err error) {
	result = &v1.Example{}
	err = c.client.Patch(pt).
		Tenant(c.te).
		Namespace(c.ns).
		Resource("examples").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)

	return
}