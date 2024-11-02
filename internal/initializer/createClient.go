package initializer

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)
type Kube struct {
   Client *kubernetes.Clientset
   Config *rest.Config
   Namespace string
   Pod  v1.PodInterface
}

var K = &Kube{}

func CreatClient() (*kubernetes.Clientset, *rest.Config) {
	home, _ := os.UserHomeDir()
	kubeConfigPath := filepath.Join(home, ".kube/config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)

	if err != nil {
		panic(err.Error())
	}

	client := kubernetes.NewForConfigOrDie(config)
	K.Client = client
	K.Config = config
	K.Namespace = "default"
	K.Pod = client.CoreV1().Pods(K.Namespace)
	fmt.Println("client Created successfully")
	return client, config
}