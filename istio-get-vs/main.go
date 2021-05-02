package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var KubeConfig = flag.String("kubeconfig", "", "kubeconfig file")

func usage() {
	fmt.Println(usageText)
	os.Exit(1)
}

var usageText = `istio-get-vs [options...]

Options:
-n string  Namespace By default "default"
-s string  VirtualService name
-h         help message
`

func main() {

	var (
		namespace string
		vsname    string
	)
	flag.StringVar(&namespace, "n", "default", "namespace")
	flag.StringVar(&vsname, "s", "", "VirtualService name")
	flag.Usage = usage
	flag.Parse()

	client, err := newIstioClient(*KubeConfig)
	if err != nil {
		fmt.Printf("[ERROR] Failed to create istio client: %s\n", err)
		os.Exit(1)
	}

	if vsname != "" {
		log.Printf("Get VirtualService: %v (namespace: %v)\n", vsname, namespace)
		// Get VirtualService
		vs, err := client.NetworkingV1alpha3().VirtualServices(namespace).Get(
			context.TODO(),
			vsname,
			v1.GetOptions{})
		if err != nil {
			fmt.Printf("[ERROR] Failed to get VirtualService %s in %s namespace: %s\n", vsname, namespace, err)
			os.Exit(1)
		}
		fmt.Printf("VirtualService Hosts %+v Gateway %+v Http %+v\n",
			vs.Spec.Hosts,
			vs.Spec.Gateways,
			vs.Spec.Http)

	} else {
		fmt.Printf("List VirtualServiceis (namespace: %v)\n", namespace)
		// Print all VirtualServices in the namespace
		vsList, err := client.NetworkingV1alpha3().VirtualServices(namespace).List(
			context.TODO(),
			v1.ListOptions{})
		if err != nil {
			fmt.Printf("[ERROR] Failed to get VirtualService in %s namespace: %s\n", namespace, err)
			os.Exit(1)
		}
		for i := range vsList.Items {
			vs := vsList.Items[i]
			fmt.Printf("VirtualService Hosts %+v Gateway %+v Http %+v\n",
				vs.Spec.Hosts,
				vs.Spec.Gateways,
				vs.Spec.Http)
		}
	}
}

func newIstioClient(kubeConfigPath string) (versionedclient.Interface, error) {
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
	return versionedclient.NewForConfig(kubeConfig)
}
