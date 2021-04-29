package main

import (
	"context"
	"flag"
	"log"
	"os"

	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// https://istio.io/latest/blog/2019/announcing-istio-client-go/
// https://github.com/istio/client-go/blob/release-1.9/cmd/example/client.go
// api references
// https://pkg.go.dev/istio.io/client-go/pkg/clientset/versioned

var (
	KubeConfig = flag.String("kubeconfig", "", "kubeconfig file")
)

func main() {

	var (
		namespace string
		name      string
	)
	flag.StringVar(&namespace, "namespace", "default", "namespace")
	flag.StringVar(&name, "name", "", "virtual service name")
	flag.Parse()

	// create kubernetes client
	client, err := newClient(*KubeConfig)
	if err != nil {
		log.Fatalf("Failed to create istio client: %s", err)
	}

	if name != "" {
		log.Printf("Get VirtualService: %v (namespace: %v)\n", name, namespace)
		// Get VirtualService
		vs, err := client.NetworkingV1alpha3().VirtualServices(namespace).Get(context.TODO(), name, v1.GetOptions{})
		if err != nil {
			log.Fatalf("Failed to get VirtualService %s in %s namespace: %s", name, namespace, err)
		}
		log.Printf("VirtualService Hosts %+v Gateway %+v Http %+v\n", vs.Spec.Hosts, vs.Spec.Gateways, vs.Spec.Http)

	} else {
		log.Printf("List VirtualServiceis (namespace: %v)\n", namespace)
		// Print all VirtualServices
		vsList, err := client.NetworkingV1alpha3().VirtualServices(namespace).List(context.TODO(), v1.ListOptions{})
		if err != nil {
			log.Fatalf("Failed to get VirtualService in %s namespace: %s", namespace, err)
		}
		for i := range vsList.Items {
			vs := vsList.Items[i]
			log.Printf("VirtualService Hosts %+v Gateway %+v Http %+v\n", vs.Spec.Hosts, vs.Spec.Gateways, vs.Spec.Http)
		}
	}
}

func newClient(kubeConfigPath string) (versionedclient.Interface, error) {
	if kubeConfigPath == "" {
		kubeConfigPath = os.Getenv("KUBECONFIG")
	}
	if kubeConfigPath == "" {
		// use default path(.kube/config) if kubeconfig path is not set
		kubeConfigPath = clientcmd.RecommendedHomeFile
	}
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}
	//return kubernetes.NewForConfig(kubeConfig)
	return versionedclient.NewForConfig(kubeConfig)
}
