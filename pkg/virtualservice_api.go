package istioapi

import (
	"fmt"
	"github.com/ruiwang47/k8s-istio-client/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"log"
)

func (istioApi *IstioApi) ListVirtualServices(namespace string) (*v1alpha3.VirtualServiceList, error) {
	vsList, err := istioApi.ic.NetworkingV1alpha3().VirtualServices(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get VirtualService in %s namespace: %s", namespace, err)
		return nil, ErrorListIstioResource
	}
	return vsList, nil
}

func (istioApi *IstioApi) GetVirtualService(namespace string, name string) (*v1alpha3.VirtualService, error) {
	vs, err := istioApi.ic.NetworkingV1alpha3().VirtualServices(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, ErrorGetIstioResource
	}
	return vs, nil
}

func (istioApi *IstioApi) UpdateVirtualService(namespace string, result *v1alpha3.VirtualService) error {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, updateErr := istioApi.ic.NetworkingV1alpha3().VirtualServices(namespace).Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
		return ErrorUpdateIstioResource
	}

	return nil
}

func (istioApi *IstioApi) CreateVirtualService(namespace string, vs *v1alpha3.VirtualService) (*v1alpha3.VirtualService, error) {
	result, err := istioApi.ic.NetworkingV1alpha3().VirtualServices(namespace).Create(vs)
	if err != nil {
		return nil, ErrorCreateIstioResource
	}
	return result, nil
}

func (istioApi *IstioApi) DeleteVirtualService(namespace string, name string) error {
	err := istioApi.ic.NetworkingV1alpha3().VirtualServices(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		return ErrorDeleteIstioResource
	}
	return nil
}
