package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

var KubeConfig = flag.String("kubeconfig", "", "kubeconfig file")

func usage() {
	fmt.Println(usageText)
	os.Exit(0)
}

var usageText = `istio-update-vs [options...]

Options:
-n string  Namespace By default "default"
-s string  VirtualService name
-desthost string  HTTP Route destination host name in a VirtualService
-destsubset string  HTTP Route destination subset name in a VirtualService
-h         help message
`

func main() {

	var (
		namespace  string
		vsname     string
		desthost   string
		destsubset string
	)
	flag.StringVar(&vsname, "s", "", "virtual service name")
	flag.StringVar(&namespace, "n", "default", "namespace")
	flag.StringVar(&desthost, "desthost", "", "virtual service HTTP Route destination host name")
	flag.StringVar(&destsubset, "destsubset", "", "virtual service HTTP Route destination subset name")
	flag.Usage = usage
	flag.Parse()

	if vsname == "" || desthost == "" || destsubset == "" {
		usage()
	}

	var (
		httpRouteList            []*networkingv1alpha3.HTTPRoute
		newHttpRouteList         []*networkingv1alpha3.HTTPRoute
		HTTPRouteDestinationList []*networkingv1alpha3.HTTPRouteDestination
	)

	client, err := newIstioClient(*KubeConfig)
	if err != nil {
		fmt.Printf("[ERROR] Failed to create istio client: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Get VirtualService: %v (namespace: %v)\n", vsname, namespace)
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

	httpRouteList = vs.Spec.GetHttp()
	if httpRouteList == nil || len(httpRouteList) < 1 {
		fmt.Printf("No resource found in VirtualService %+v\n", vsname)
		return
	}
	// Update VirtualService
	// How to update VirtualService?
	// - Set weight of your desitnation to 100%
	// - Set the rest of the destination to 0%
	for _, httpRoute := range httpRouteList {
		fmt.Printf("VirtualService httpRoute %+v \n", httpRoute)
		HTTPRouteDestinationList = httpRoute.Route
		var newHTTPRouteDestinationList []*networkingv1alpha3.HTTPRouteDestination
		for _, dest := range HTTPRouteDestinationList {
			newWeight := 0
			if dest.Destination.Host == desthost && dest.Destination.Subset == destsubset {
				newWeight = 100
			}
			HTTPRouteDestination := &networkingv1alpha3.HTTPRouteDestination{
				Destination: dest.Destination,
				Weight:      int32(newWeight),
			}
			newHTTPRouteDestinationList = append(newHTTPRouteDestinationList, HTTPRouteDestination)
		}
		httpRoute.Route = newHTTPRouteDestinationList
		newHttpRouteList = append(newHttpRouteList, httpRoute)
	}
	vs.Spec.Http = newHttpRouteList

	newvs, err := client.NetworkingV1alpha3().VirtualServices(namespace).Update(
		context.TODO(),
		vs,
		v1.UpdateOptions{})
	if err != nil {
		fmt.Printf("[ERROR] Failed to update VirtualService %s in %s namespace: %s", vsname, namespace, err)
		os.Exit(1)
	}
	fmt.Printf("New VirtualService %+v \n", newvs)
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
