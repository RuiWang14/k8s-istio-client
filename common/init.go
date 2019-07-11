package common

import (
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// 初始化k8s客户端
func InitClient() (clientset *kubernetes.Clientset, err error) {
	var (
		kubeconfig []byte
		restConf   *rest.Config
	)

	// 读kubeconfig文件
	if kubeconfig, err = ioutil.ReadFile("./config"); err != nil {
		goto END
	}
	// 生成rest client配置
	if restConf, err = clientcmd.RESTConfigFromKubeConfig(kubeconfig); err != nil {
		goto END
	}
	// 生成clientset配置
	if clientset, err = kubernetes.NewForConfig(restConf); err != nil {
		goto END
	}
END:
	return
}
