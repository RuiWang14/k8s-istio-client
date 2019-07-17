package istioapi

import (
	"fmt"
	"github.com/RuiWang14/k8s-istio-client/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"log"
)

func (istioApi *IstioApi) ListGateways(namespace string) *v1alpha3.GatewayList {
	gatewayList, err := istioApi.ic.NetworkingV1alpha3().Gateways(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get Gateways in %s namespace: %s", namespace, err)
	}
	return gatewayList
}

func (istioApi *IstioApi) GetGateway(namespace string, name string) *v1alpha3.Gateway {
	gateway, err := istioApi.ic.NetworkingV1alpha3().Gateways(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	return gateway
}

func (istioApi *IstioApi) UpdateGateway(namespace string, result *v1alpha3.Gateway) {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, updateErr := istioApi.ic.NetworkingV1alpha3().Gateways(namespace).Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
}

func (istioApi *IstioApi) CreateGateway(namespace string, vs *v1alpha3.Gateway) *v1alpha3.Gateway {
	gateway, err := istioApi.ic.NetworkingV1alpha3().Gateways(namespace).Create(vs)
	if err != nil {
		panic(err)
	}
	return gateway
}

func (istioApi *IstioApi) DeleteGateway(namespace string, name string) {
	err := istioApi.ic.NetworkingV1alpha3().Gateways(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
}
