package istioapi

import (
	"fmt"
	"github.com/ruiwang47/k8s-istio-client/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"log"
)

func (istioApi *IstioApi) ListEnvoyFilters(namespace string) *v1alpha3.EnvoyFilterList {
	envoyFilterList, err := istioApi.ic.NetworkingV1alpha3().EnvoyFilters(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get EnvoyFilters in %s namespace: %s", namespace, err)
	}
	return envoyFilterList
}

func (istioApi *IstioApi) GetEnvoyFilter(namespace string, name string) *v1alpha3.EnvoyFilter {
	envoyFilter, err := istioApi.ic.NetworkingV1alpha3().EnvoyFilters(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	return envoyFilter
}

func (istioApi *IstioApi) UpdateEnvoyFilter(namespace string, result *v1alpha3.EnvoyFilter) {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, updateErr := istioApi.ic.NetworkingV1alpha3().EnvoyFilters(namespace).Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
}

func (istioApi *IstioApi) CreateEnvoyFilter(namespace string, vs *v1alpha3.EnvoyFilter) *v1alpha3.EnvoyFilter {
	envoyFilter, err := istioApi.ic.NetworkingV1alpha3().EnvoyFilters(namespace).Create(vs)
	if err != nil {
		panic(err)
	}
	return envoyFilter
}

func (istioApi *IstioApi) DeleteEnvoyFilter(namespace string, name string) {
	err := istioApi.ic.NetworkingV1alpha3().EnvoyFilters(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
}
