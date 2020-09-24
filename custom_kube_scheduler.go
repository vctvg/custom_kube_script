package main

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"time"
)

const schedulerName = "new-scheduler"

type Scheduler struct{
	clientset *kubernetes.Clientset
}

func NewScheduler() Scheduler {
	config, err := clientcmd.BuildConfigFromFlags("", "config")
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return Scheduler{
		clientset: clientset,
	}
}

func (s *Scheduler) WatchNewPod() error {
	watch, err := s.clientset.CoreV1().Pods("").Watch(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for event := range watch.ResultChan() {
		if event.Type != "ADDED" {
			continue
		}
		p := event.Object.(*v1.Pod)
		fmt.Println(p.Namespace, p.Name)

		node, err := s.GetNode()
		if err != nil {
			panic(err.Error())
		}

		err = s.BindPod(p, node)
		if err != nil {
			panic(err.Error())
		}

		message := fmt.Sprintf("Placed pod [%s/%s] on %s\n", p.Namespace, p.Name, node.Name)

		err = s.CreateEvent(p, message)
		if err != nil {
			panic(err.Error())
		}
	}
	return nil
}

func (s *Scheduler) GetNode() (*v1.Node, error) {
	node, err := s.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
    if err != nil{
    	panic(err.Error())
	}
    return &node.Items[rand.Intn(len(node.Items))], nil
}

func (s *Scheduler) BindPod(p *v1.Pod, randomNode *v1.Node) error {
	return s.clientset.CoreV1().Pods(p.Namespace).Bind(context.TODO(), &v1.Binding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      p.Name,
			Namespace: p.Namespace,
		},
		Target: v1.ObjectReference{
			APIVersion: "v1",
			Kind:       "Node",
			Name:       randomNode.Name,
		},
	})
}

func (s *Scheduler) CreateEvent(p *v1.Pod, message string) error {
	timestamp := time.Now().UTC()
	s.clientset.CoreV1().Events(p.Namespace).Create(&v1.Event{
		Count:          1,
		Message:        message,
		Reason:         "Scheduled",
		LastTimestamp:  metav1.NewTime(timestamp),
		FirstTimestamp: metav1.NewTime(timestamp),
		Type:           "Normal",
		Source: v1.EventSource{
			Component: schedulerName,
		},
		InvolvedObject: v1.ObjectReference{
			Kind:      "Pod",
			Name:      p.Name,
			Namespace: p.Namespace,
			UID:       p.UID,
		},
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: p.Name + "-",
		},
	})
	if err != nil{
		panic(err.Err())
	}
	return nil
}

func main() {
	//scheduler := client *kubernetes.Clientset
	scheduler := NewScheduler()
	scheduler.WatchNewPod()
}
