package istioapi

import (
	"fmt"
	"github.com/RuiWang14/k8s-istio-client/pkg/apis/authentication/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"log"
)

func (istioApi *IstioApi) ListPolices(namespace string) *v1alpha1.PolicyList {
	policyList, err := istioApi.ic.AuthenticationV1alpha1().Policies(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get Polices %s", err)
	}
	return policyList
}

func (istioApi *IstioApi) GetPolicy(namespace string, name string) *v1alpha1.Policy {
	policy, err := istioApi.ic.AuthenticationV1alpha1().Policies(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	return policy
}

func (istioApi *IstioApi) UpdatePolicy(namespace string, policy *v1alpha1.Policy) {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, updateErr := istioApi.ic.AuthenticationV1alpha1().Policies(namespace).Update(policy)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
}

func (istioApi *IstioApi) CreatePolicy(namespace string, vs *v1alpha1.Policy) *v1alpha1.Policy {
	policy, err := istioApi.ic.AuthenticationV1alpha1().Policies(namespace).Create(vs)
	if err != nil {
		panic(err)
	}
	return policy
}

func (istioApi *IstioApi) DeletePolicy(namespace string, name string) {
	err := istioApi.ic.AuthenticationV1alpha1().Policies(namespace).Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
}
