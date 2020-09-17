package main

import (
	"fmt"
	//v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

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

	deployment, err := clientset.AppsV1beta1().Deployments("kube-system").Get(context.TODO(),"coredns", metav1.GetOptions{})
	if err != nil{
		panic(err.Error())
	}

	fmt.Println(deployment)
}
