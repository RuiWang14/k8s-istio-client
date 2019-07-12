package main

import (
	"flag"
	versionedclient "github.com/RuiWang14/k8s-istio-client/pkg/client/clientset/versioned"
	informers "github.com/RuiWang14/k8s-istio-client/pkg/client/informers/externalversions"

	"github.com/RuiWang14/k8s-istio-client/pkg/signals"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"time"
)

func main() {
	flag.Parse()

	kubeconfig := "/Users/Rui/.kube/config"
	namespace := "default"

	// rest config
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Failed to create k8s rest client: %s", err)
	}

	// istio client
	ic, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		log.Fatalf("Failed to create istio client: %s", err)
	}

	// kube client for event board cast
	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		log.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	// List VirtualServices
	vsList, err := ic.NetworkingV1alpha3().VirtualServices(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get VirtualService in %s namespace: %s", namespace, err)
	}
	for i := range vsList.Items {
		vs := vsList.Items[i]
		log.Printf("Index: %d VirtualService Hosts: %+v\n", i, vs.Spec.GetHosts())
	}

	// Test Customer istio controller
	stopCh := signals.SetupSignalHandler()

	informerFactory := informers.NewSharedInformerFactory(ic, time.Second*30)

	controller := NewController(kubeClient, ic,
		informerFactory.Networking().V1alpha3().VirtualServices())

	go informerFactory.Start(stopCh)

	if err = controller.Run(1, stopCh); err != nil {
		log.Fatalf("Error running controller: %s", err.Error())
	}

}
