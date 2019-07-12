package main

import (
	"flag"
	versionedclient "github.com/RuiWang14/k8s-istio-client/pkg/client/clientset/versioned"
	informers "github.com/RuiWang14/k8s-istio-client/pkg/client/informers/externalversions"

	"github.com/RuiWang14/k8s-istio-client/pkg/signals"
	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	"log"
	"time"
)

func main() {
	flag.Parse()

	kubeconfig := "/Users/Rui/.kube/config"
	namespace := "default"
	if len(kubeconfig) == 0 || len(namespace) == 0 {
		log.Fatalf("Environment variables KUBECONFIG and NAMESPACE need to be set")
	}
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("Failed to create k8s rest client: %s", err)
	}

	ic, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		log.Fatalf("Failed to create istio client: %s", err)
	}
	// Test VirtualServices
	vsList, err := ic.NetworkingV1alpha3().VirtualServices(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Failed to get VirtualService in %s namespace: %s", namespace, err)
	}
	for i := range vsList.Items {
		vs := vsList.Items[i]
		log.Printf("Index: %d VirtualService Hosts: %+v\n", i, vs.Spec.GetHosts())
	}

	// set up signals so we handle the first shutdown signal gracefully
	stopCh := signals.SetupSignalHandler()

	kubeClient, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		glog.Fatalf("Error building kubernetes clientset: %s", err.Error())
	}

	exampleClient, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		klog.Fatalf("Error building example clientset: %s", err.Error())
	}

	exampleInformerFactory := informers.NewSharedInformerFactory(exampleClient, time.Second*30)

	controller := NewController(kubeClient, exampleClient,
		exampleInformerFactory.Networking().V1alpha3().VirtualServices())

	go exampleInformerFactory.Start(stopCh)

	if err = controller.Run(1, stopCh); err != nil {
		glog.Fatalf("Error running controller: %s", err.Error())
	}

}
