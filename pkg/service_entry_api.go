package istioapi

import (
	"fmt"
	"github.com/ruiwang47/k8s-istio-client/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"log"
)

func (istioApi *IstioApi) ListServiceEntries(namespace string) *v1alpha3.ServiceEntryList {
	serviceEntryList, err := istioApi.ic.NetworkingV1alpha3().ServiceEntries(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get Gateways in %s namespace: %s", namespace, err)
	}
	return serviceEntryList
}

func (istioApi *IstioApi) GetServiceEntry(namespace string, name string) *v1alpha3.ServiceEntry {
	serviceEntry, err := istioApi.ic.NetworkingV1alpha3().ServiceEntries(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	return serviceEntry
}

func (istioApi *IstioApi) UpdateServiceEntry(namespace string, result *v1alpha3.ServiceEntry) {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, updateErr := istioApi.ic.NetworkingV1alpha3().ServiceEntries(namespace).Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
}

func (istioApi *IstioApi) CreateServiceEntries(namespace string, vs *v1alpha3.ServiceEntry) *v1alpha3.ServiceEntry {
	serviceEntry, err := istioApi.ic.NetworkingV1alpha3().ServiceEntries(namespace).Create(vs)
	if err != nil {
		panic(err)
	}
	return serviceEntry
}

func (istioApi *IstioApi) DeleteServiceEntries(namespace string, name string) {
	err := istioApi.ic.NetworkingV1alpha3().ServiceEntries(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
}
