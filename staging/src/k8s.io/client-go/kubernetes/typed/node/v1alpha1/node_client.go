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

package v1alpha1

import (
	"time"

	v1alpha1 "k8s.io/api/node/v1alpha1"
	rand "k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes/scheme"
	rest "k8s.io/client-go/rest"
)

type NodeV1alpha1Interface interface {
	RESTClient() rest.Interface
	RESTClients() []rest.Interface
	RuntimeClassesGetter
}

// NodeV1alpha1Client is used to interact with features provided by the node.k8s.io group.
type NodeV1alpha1Client struct {
	restClients []rest.Interface
}

func (c *NodeV1alpha1Client) RuntimeClasses() RuntimeClassInterface {
	return newRuntimeClasses(c)
}

// NewForConfig creates a new NodeV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*NodeV1alpha1Client, error) {
	configs := rest.CopyConfigs(c)
	if err := setConfigDefaults(configs); err != nil {
		return nil, err
	}

	clients := make([]rest.Interface, len(configs.GetAllConfigs()))
	for i, config := range configs.GetAllConfigs() {
		client, err := rest.RESTClientFor(config)
		if err != nil {
			return nil, err
		}
		clients[i] = client
	}

	return &NodeV1alpha1Client{clients}, nil
}

// NewForConfigOrDie creates a new NodeV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *NodeV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new NodeV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *NodeV1alpha1Client {
	clients := []rest.Interface{c}
	return &NodeV1alpha1Client{clients}
}

func setConfigDefaults(configs *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion

	for _, config := range configs.GetAllConfigs() {
		config.GroupVersion = &gv
		config.APIPath = "/apis"
		config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

		if config.UserAgent == "" {
			config.UserAgent = rest.DefaultKubernetesUserAgent()
		}
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *NodeV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}

	max := len(c.restClients)
	if max == 0 {
		return nil
	}
	if max == 1 {
		return c.restClients[0]
	}

	rand.Seed(time.Now().UnixNano())
	ran := rand.IntnRange(0, max-1)
	return c.restClients[ran]
}

// RESTClients returns all RESTClient that are used to communicate
// with all API servers by this client implementation.
func (c *NodeV1alpha1Client) RESTClients() []rest.Interface {
	if c == nil {
		return nil
	}

	return c.restClients
}