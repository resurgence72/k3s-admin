package service

import (
	"fmt"
	"k3s-admin/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var K8s k8s

type k8s struct {
	ClientSet *kubernetes.Clientset
}

// 初始化k7s client
func (k *k8s) Init() error {
	conf, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		fmt.Println("k8s config build failed ", err)
		return err
	}
	clientset, err := kubernetes.NewForConfig(conf)
	if err != nil {
		fmt.Println("k8s new config build failed ", err)
		return err
	}
	k.ClientSet = clientset
	K8s = *k
	return nil
}
