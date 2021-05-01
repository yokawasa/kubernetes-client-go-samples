package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	networkingv1alpha3 "istio.io/api/networking/v1alpha3"

	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

// https://istio.io/latest/blog/2019/announcing-istio-client-go/
// https://github.com/istio/client-go/blob/release-1.9/cmd/example/client.go
// https://github.com/Michael754267513/k8s-client-go/blob/k8s-note/istio/VirtualServices/main.go
// api references
// https://pkg.go.dev/istio.io/client-go/pkg/clientset/versioned

var (
	KubeConfig = flag.String("kubeconfig", "", "kubeconfig file")
)

func usage() {
	fmt.Println(usageText)
	os.Exit(0)
}

var usageText = `istio-update-vs [options...]

Options:
-namespace   your namespace
-name        your virtual service name
-destination your destination
`

func main() {

	var (
		namespace   string
		name        string
		destination string
	)
	flag.StringVar(&namespace, "namespace", "default", "namespace")
	flag.StringVar(&name, "name", "", "virtual service name")
	flag.StringVar(&destination, "destination", "", "virtual service HTTP Route destination host name")
	flag.Parse()

	if name == "" || destination == "" {
		usage()
	}

	var (
		httpRouteList            []*networkingv1alpha3.HTTPRoute
		newHttpRouteList         []*networkingv1alpha3.HTTPRoute
		HTTPRouteDestinationList []*networkingv1alpha3.HTTPRouteDestination
		newWeight                int32
	)

	// create kubernetes client
	client, err := newClient(*KubeConfig)
	if err != nil {
		log.Fatalf("Failed to create istio client: %s", err)
	}

	log.Printf("Get VirtualService: %v (namespace: %v)\n", name, namespace)
	// Get VirtualService
	vs, err := client.NetworkingV1alpha3().VirtualServices(namespace).Get(context.TODO(), name, v1.GetOptions{})
	if err != nil {
		log.Fatalf("Failed to get VirtualService %s in %s namespace: %s", name, namespace, err)
	}
	log.Printf("VirtualService Hosts %+v Gateway %+v Http %+v\n", vs.Spec.Hosts, vs.Spec.Gateways, vs.Spec.Http)

	// Update VirtualService
	// How to update VirtualService?
	// - Set weight of your desitnation to 100%
	// - Set the rest of the destination to 0%
	httpRouteList = vs.Spec.GetHttp()
	if httpRouteList != nil {
		for _, httpRoute := range httpRouteList {
			log.Printf("VirtualService httpRoute %+v \n", httpRoute)
			HTTPRouteDestinationList = httpRoute.Route
			var newHTTPRouteDestinationList []*networkingv1alpha3.HTTPRouteDestination
			for _, dest := range HTTPRouteDestinationList {
				newWeight = 0
				if dest.Destination.Host == destination {
					newWeight = 100
				}
				HTTPRouteDestination := &networkingv1alpha3.HTTPRouteDestination{
					Destination: dest.Destination,
					Weight:      newWeight,
				}
				newHTTPRouteDestinationList = append(newHTTPRouteDestinationList, HTTPRouteDestination)
			}

			httpRoute.Route = newHTTPRouteDestinationList
			newHttpRouteList = append(newHttpRouteList, httpRoute)
		}
		vs.Spec.Http = newHttpRouteList

		newvs, err := client.NetworkingV1alpha3().VirtualServices(namespace).Update(context.TODO(), vs, v1.UpdateOptions{})
		if err != nil {
			return
		}
		log.Printf("New VirtualService %+v \n", newvs)
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
