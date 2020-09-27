package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("default", "config")
	if err != nil{
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil{
		panic(err.Error())
	}

	stopper := make(chan struct{})
	defer close(stopper)

	factory := informers.NewSharedInformerFactory(clientset, 0)
	nodeInformer := factory.Core().V1().Nodes()
	defer runtime.HandleCrash()

	factory.Start(stopper)

	Lister := nodeInformer.Lister()

	nodeList, err := Lister.List(labels.Everything())
	if err != nil{
		panic(err.Error())
	}

	fmt.Println(nodeList)
	<-stopper
}
