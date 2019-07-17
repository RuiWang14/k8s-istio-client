package istioapi

import (
	"flag"
	"github.com/RuiWang14/k8s-istio-client/pkg/apis/networking/v1alpha3"
	"log"
	"testing"

	istiov1alpha3 "istio.io/api/networking/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestVirtualService(t *testing.T) {
	flag.Parse()

	kubeconfig := "/Users/Rui/.kube/config"
	namespace := "default"

	// New VirtualServiceApi
	istioApi := NewIstioApi(kubeconfig)

	vsName := "test-virtualservice"

	// Create VirtualService
	vs := buildVirtualService(vsName)
	result, _ := istioApi.CreateVirtualService(namespace, vs)
	log.Printf("Create VirtualService %+v\n", result.Name)

	// List VirtualServices
	vsList, _ := istioApi.ListVirtualServices(namespace)
	for i := range vsList.Items {
		vs := vsList.Items[i]
		log.Printf("Index: %d VirtualService Hosts: %+v\n", i, vs.Name)
	}

	// Get VirtualService
	updateVs, _ := istioApi.GetVirtualService(namespace, vsName)

	// Update VirtualService
	updateVs.Spec.Hosts[0] = "boot-example-b-svc"
	istioApi.UpdateVirtualService(namespace, updateVs)

	// Delete VirtualService
	istioApi.DeleteVirtualService(namespace, vsName)
	log.Printf("Delete VirtualService %+v\n", vsName)

	// List VirtualServices
	vsList, _ = istioApi.ListVirtualServices(namespace)
	for i := range vsList.Items {
		vs := vsList.Items[i]
		log.Printf("Index: %d VirtualService Hosts: %+v\n", i, vs.Name)
	}

}

func buildVirtualService(name string) *v1alpha3.VirtualService {

	virtualService := &v1alpha3.VirtualService{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: v1alpha3.VirtualServiceSpec{
			VirtualService: istiov1alpha3.VirtualService{
				Hosts: []string{
					"productpage",
				},
				Http: []*istiov1alpha3.HTTPRoute{
					{
						Route: []*istiov1alpha3.HTTPRouteDestination{
							{
								Destination: &istiov1alpha3.Destination{
									Host:   "productpage",
									Subset: "v1",
								},
							},
						},
					},
				},
			},
		},
	}
	return virtualService
}
