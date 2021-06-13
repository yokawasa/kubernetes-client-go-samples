package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeConfig = flag.String("kubeconfig", "", "kubeconfig file")
	flag.Parse()

	// create kubernetes client
	client, err := newClient(*kubeConfig)
	if err != nil {
		fmt.Printf("[ERROR] Failed to create client: %s\n", err)
		os.Exit(1)
	}

	list, err := client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("[ERROR] listing nodes: %v\n", err)
		os.Exit(1)
	}

	for _, node := range list.Items {
		fmt.Printf("Node name %s \n", node.Name)
		//Dump labels
		labels := node.Labels
		for k, v := range labels {
			fmt.Printf("Node %s Label %s:%s\n", node.Name, k, v)
		}
	}
}

func newClient(kubeConfigPath string) (kubernetes.Interface, error) {
	if kubeConfigPath == "" {
		kubeConfigPath = os.Getenv("KUBECONFIG")
	}
	if kubeConfigPath == "" {
		kubeConfigPath = clientcmd.RecommendedHomeFile // use default path(.kube/config)
	}
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(kubeConfig)
}
