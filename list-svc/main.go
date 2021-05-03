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

var KubeConfig = flag.String("kubeconfig", "", "kubeconfig file")

func main() {
	flag.Parse()

	client, err := newClient(*KubeConfig)
	if err != nil {
		fmt.Printf("[ERROR] Failed to create client: %s\n", err)
		os.Exit(1)
	}

	// Get all services in all namespace
	services, err := client.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("[ERROR] Failed to list service: %s\n", err)
		os.Exit(1)
	}

	for _, service := range services.Items {
		fmt.Printf("Service name= %s namespace=%s type=%s \n", service.Name, service.Namespace, service.Spec.Type)
		if service.Spec.Type == "ClusterIP" {
			for _, port := range service.Spec.Ports {
				if port.Name != "" {
					fmt.Printf("ClientIP Port: name=%s port=%d\n", port.Name, port.Port)
				}
			}
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
