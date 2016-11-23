package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/FourSigma/kubemon/handlers"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/tools/clientcmd"
)

var (
	kubeconfig = flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
)

func main() {
	flag.Parse()
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// for {
	// 	pods, err := clientset.Core().Pods(api.NamespaceDefault).List(api.ListOptions{})
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
	// 	time.Sleep(time.Millisecond)
	// }
	fmt.Println("Starting server -- http://localhost:" + "9180")
	http.ListenAndServe(":9180", handlers.GetRoutes(clientset))

}
