package initializer

import (
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)
type Kube struct {
   client *kubernetes.Clientset
   config *rest.Config
   namespace string
}

var k = &Kube{}

func CreatClient() (*kubernetes.Clientset, *rest.Config) {
	home, _ := os.UserHomeDir()
	kubeConfigPath := filepath.Join(home, ".kube/config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)

	if err != nil {
		panic(err.Error())
	}

	client := kubernetes.NewForConfigOrDie(config)
	k.client = client
	k.config = config
	k.namespace = "default"
	fmt.Println("client Created successfully")
	return client, config
}