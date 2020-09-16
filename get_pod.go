package main

import (
	"fmt"
	"context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "config")
	if err != nil{
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil{
		panic(err.Error())
	}

	pods, _ := clientset.CoreV1().Pods("kube-system").List(context.TODO(), metav1.ListOptions{})
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
}
