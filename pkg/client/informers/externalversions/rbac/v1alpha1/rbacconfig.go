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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	rbacv1alpha1 "github.com/ruiwang47/k8s-istio-client/pkg/apis/rbac/v1alpha1"
	versioned "github.com/ruiwang47/k8s-istio-client/pkg/client/clientset/versioned"
	internalinterfaces "github.com/ruiwang47/k8s-istio-client/pkg/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/ruiwang47/k8s-istio-client/pkg/client/listers/rbac/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// RbacConfigInformer provides access to a shared informer and lister for
// RbacConfigs.
type RbacConfigInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.RbacConfigLister
}

type rbacConfigInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewRbacConfigInformer constructs a new informer for RbacConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewRbacConfigInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredRbacConfigInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredRbacConfigInformer constructs a new informer for RbacConfig type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredRbacConfigInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.RbacV1alpha1().RbacConfigs(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.RbacV1alpha1().RbacConfigs(namespace).Watch(options)
			},
		},
		&rbacv1alpha1.RbacConfig{},
		resyncPeriod,
		indexers,
	)
}

func (f *rbacConfigInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredRbacConfigInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *rbacConfigInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&rbacv1alpha1.RbacConfig{}, f.defaultInformer)
}

func (f *rbacConfigInformer) Lister() v1alpha1.RbacConfigLister {
	return v1alpha1.NewRbacConfigLister(f.Informer().GetIndexer())
}
