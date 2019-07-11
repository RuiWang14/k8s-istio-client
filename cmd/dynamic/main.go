package main

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

func main() {

	namespace := "default"

	client := buildClient()

	// list the virtual service
	list := listVirtualServices(client, namespace)
	printResource(list)

	// add new virtual service
	//result := createVirtualServices(client, namespace)
	//fmt.Printf("Created deployment %q.\n", result.GetName())

	// update virtual service
	//updateVirtualServices(client, namespace)

	// delete virtual service
	//deleteVirtualServices(client, namespace)
}

func deleteVirtualServices(client dynamic.Interface, namespace string) {
	vsRes := buildVirtualServicesResource()
	deletePolicy := metav1.DeletePropagationForeground
	deleteOptions := &metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}
	if err := client.Resource(vsRes).Namespace(namespace).Delete("demo-productpage", deleteOptions); err != nil {
		panic(err)
	}
}

func updateVirtualServices(client dynamic.Interface, namespace string) {
	fmt.Println("Updating deployment...")
	vsRes := buildVirtualServicesResource()

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {

		result, getErr := client.Resource(vsRes).Namespace(namespace).Get("demo-productpage", metav1.GetOptions{})
		if getErr != nil {
			panic(fmt.Errorf("Failed to get : %v", getErr))
		}

		// switch destination to v2
		// TODO refactor
		http, found, err := unstructured.NestedSlice(result.Object, "spec", "http")
		if err != nil || !found || http == nil {
			panic(fmt.Errorf("virtualservices http not found or error in spec: %v", err))
		}
		routes, found, err := unstructured.NestedSlice(http[0].(map[string]interface{}), "route")
		if err != nil || !found || routes == nil {
			panic(fmt.Errorf("virtualservices routes not found or error in spec: %v", err))
		}

		if err := unstructured.SetNestedField(routes[0].(map[string]interface{}), "v2", "destination", "subset"); err != nil {
			panic(err)
		}
		if err := unstructured.SetNestedField(http[0].(map[string]interface{}), routes, "route"); err != nil {
			panic(err)
		}
		if err := unstructured.SetNestedField(result.Object, http, "spec", "http"); err != nil {
			panic(err)
		}

		_, updateErr := client.Resource(vsRes).Namespace(namespace).Update(result, metav1.UpdateOptions{})
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}
}

func createVirtualServices(client dynamic.Interface, namespace string) *unstructured.Unstructured {
	vsRes := buildVirtualServicesResource()

	vsDeploy := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "networking.istio.io/v1alpha3",
			"kind":       "VirtualService",
			"metadata": map[string]interface{}{
				"name": "demo-productpage",
			},
			"spec": map[string]interface{}{
				"hosts": []string{
					"productpage",
				},
				"http": []map[string]interface{}{
					{
						"route": []map[string]interface{}{
							{
								"destination": map[string]interface{}{
									"host":   "productpage",
									"subset": "v1",
								},
							},
						},
					},
				},
			},
		},
	}

	vs, err := client.Resource(vsRes).Namespace(namespace).Create(vsDeploy, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}

	return vs

}

func listVirtualServices(client dynamic.Interface, namespace string) (list *unstructured.UnstructuredList) {
	virtualServicesRes := buildVirtualServicesResource()

	list, err := client.Resource(virtualServicesRes).Namespace(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	return
}

func buildVirtualServicesResource() (res schema.GroupVersionResource) {
	res = schema.GroupVersionResource{Group: "networking.istio.io", Version: "v1alpha3", Resource: "virtualservices"}
	return
}

func printResource(list *unstructured.UnstructuredList) {
	for _, d := range list.Items {
		fmt.Printf(" * %s \n", d.GetName())
	}
}

func buildClient() (client dynamic.Interface) {
	kubeconfig := "./config"

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}
	client, err = dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	return

}
