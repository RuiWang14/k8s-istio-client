package istioapi

import (
	versionedclient "github.com/RuiWang14/k8s-istio-client/pkg/client/clientset/versioned"
	"k8s.io/client-go/tools/clientcmd"
	"log"
)

type IstioApi struct {
	ic *versionedclient.Clientset
}

func NewIstioApi(kubeconfig string) *IstioApi {
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

	istioApi := &IstioApi{
		ic: ic,
	}

	return istioApi
}
