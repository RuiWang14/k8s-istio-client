package istioapi

import (
	"fmt"
	"github.com/RuiWang14/k8s-istio-client/pkg/apis/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"log"
)

func (istioApi *IstioApi) ListDestinationRules(namespace string) *v1alpha3.DestinationRuleList {
	destinationRuleList, err := istioApi.ic.NetworkingV1alpha3().DestinationRules(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get DestinationRules in %s namespace: %s", namespace, err)
	}
	return destinationRuleList
}

func (istioApi *IstioApi) GetDestinationRule(namespace string, name string) *v1alpha3.DestinationRule {
	destinationRule, err := istioApi.ic.NetworkingV1alpha3().DestinationRules(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	return destinationRule
}

func (istioApi *IstioApi) UpdateDestinationRule(namespace string, result *v1alpha3.DestinationRule) {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, updateErr := istioApi.ic.NetworkingV1alpha3().DestinationRules(namespace).Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
}

func (istioApi *IstioApi) CreateDestinationRule(namespace string, vs *v1alpha3.DestinationRule) *v1alpha3.DestinationRule {
	destinationRule, err := istioApi.ic.NetworkingV1alpha3().DestinationRules(namespace).Create(vs)
	if err != nil {
		panic(err)
	}
	return destinationRule
}

func (istioApi *IstioApi) DeleteDestinationRule(namespace string, name string) {
	err := istioApi.ic.NetworkingV1alpha3().DestinationRules(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
}
