package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var KubeConfig = flag.String("kubeconfig", "", "kubeconfig file")

func main() {
	client, err := newClient(*KubeConfig)
	if err != nil {
		fmt.Printf("[ERROR] Failed to create client: %s\n", err)
		os.Exit(1)
	}

	namespace := "default"
	appname := "httpbin"

	svc, err := getService(appname, namespace, client)
	if err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
		os.Exit(1)
	}

	pods, err := getPodsForService(svc, namespace, client)
	if err != nil {
		fmt.Printf("[ERROR]: %v\n", err)
		os.Exit(1)
	}

	for _, pod := range pods.Items {
		fmt.Printf("Pod Name %s \n", pod.Name)
		fmt.Printf("Pod Namespace %s \n", pod.Namespace)
		//Dump labels
		labels := pod.Labels
		for k, v := range labels {
			fmt.Printf("Pod %s Namespace %s Label %s:%s\n", pod.Name, pod.Namespace, k, v)
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

// ref: https://stackoverflow.com/questions/41545123/how-to-get-pods-under-the-service-with-client-go-the-client-library-of-kubernete
func getService(deployment string, namespace string, client kubernetes.Interface) (*v1.Service, error) {
	listOptions := metav1.ListOptions{}
	services, err := client.CoreV1().Services(namespace).List(context.TODO(), listOptions)
	if err != nil {
		fmt.Printf("[ERROR] Failed to list services: %s\n", err)
		return nil, err
	}
	for _, service := range services.Items {
		if strings.Contains(service.Name, deployment) {
			fmt.Printf("Service name: %v\n", service.Name)
			return &service, nil
		}
	}
	return nil, errors.New("Cannot find service for deployment")
}

func getPodsForService(svc *v1.Service, namespace string, client kubernetes.Interface) (*v1.PodList, error) {
	set := labels.Set(svc.Spec.Selector)
	listOptions := metav1.ListOptions{LabelSelector: set.AsSelector().String()}
	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), listOptions)
	return pods, err
}
