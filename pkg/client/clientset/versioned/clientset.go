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

package versioned

import (
	authenticationv1alpha1 "github.com/RuiWang14/k8s-istio-client/pkg/client/clientset/versioned/typed/authentication/v1alpha1"
	networkingv1alpha3 "github.com/RuiWang14/k8s-istio-client/pkg/client/clientset/versioned/typed/networking/v1alpha3"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	AuthenticationV1alpha1() authenticationv1alpha1.AuthenticationV1alpha1Interface
	NetworkingV1alpha3() networkingv1alpha3.NetworkingV1alpha3Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	authenticationV1alpha1 *authenticationv1alpha1.AuthenticationV1alpha1Client
	networkingV1alpha3     *networkingv1alpha3.NetworkingV1alpha3Client
}

// AuthenticationV1alpha1 retrieves the AuthenticationV1alpha1Client
func (c *Clientset) AuthenticationV1alpha1() authenticationv1alpha1.AuthenticationV1alpha1Interface {
	return c.authenticationV1alpha1
}

// NetworkingV1alpha3 retrieves the NetworkingV1alpha3Client
func (c *Clientset) NetworkingV1alpha3() networkingv1alpha3.NetworkingV1alpha3Interface {
	return c.networkingV1alpha3
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.authenticationV1alpha1, err = authenticationv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.networkingV1alpha3, err = networkingv1alpha3.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.authenticationV1alpha1 = authenticationv1alpha1.NewForConfigOrDie(c)
	cs.networkingV1alpha3 = networkingv1alpha3.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.authenticationV1alpha1 = authenticationv1alpha1.New(c)
	cs.networkingV1alpha3 = networkingv1alpha3.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
