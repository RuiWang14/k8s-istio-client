package istioapi

import (
	"fmt"
	"github.com/RuiWang14/k8s-istio-client/pkg/apis/authentication/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/util/retry"
	"log"
)

func (istioApi *IstioApi) ListMeshPolices() *v1alpha1.MeshPolicyList {
	meshPolicyList, err := istioApi.ic.AuthenticationV1alpha1().MeshPolicies().List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get MeshPolices %s", err)
	}
	return meshPolicyList
}

func (istioApi *IstioApi) GetMeshPolicy(name string) *v1alpha1.MeshPolicy {
	meshPolicy, err := istioApi.ic.AuthenticationV1alpha1().MeshPolicies().Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	return meshPolicy
}

func (istioApi *IstioApi) UpdateMeshPolicy(result *v1alpha1.MeshPolicy) {
	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		_, updateErr := istioApi.ic.AuthenticationV1alpha1().MeshPolicies().Update(result)
		return updateErr
	})
	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
}

func (istioApi *IstioApi) CreateMeshPolicy(vs *v1alpha1.MeshPolicy) *v1alpha1.MeshPolicy {
	meshPolicy, err := istioApi.ic.AuthenticationV1alpha1().MeshPolicies().Create(vs)
	if err != nil {
		panic(err)
	}
	return meshPolicy
}

func (istioApi *IstioApi) DeleteMeshPolicy(name string) {
	err := istioApi.ic.AuthenticationV1alpha1().MeshPolicies().Delete(name, &metav1.DeleteOptions{})
	if err != nil {
		panic(err)
	}
}
