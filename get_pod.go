package main

import (
	"fmt"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// uses the current context in kubeconfig
	// path-to-kubeconfig -- for example, /root/.kube/config
	config, err := clientcmd.BuildConfigFromFlags("", "config")

	if err != nil{
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)

	if err != nil{
		panic(err.Error())
	}
	// access the API to list pods
	pods, _ := clientset.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{})
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
}
