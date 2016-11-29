package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/FourSigma/kubemon/handlers"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/tools/clientcmd"

	ghandlers "github.com/gorilla/handlers"
)

var (
	kubeconfig = flag.String("kubeconfig", "./config", "absolute path to the kubeconfig file")
)

func main() {
	flag.Parse()
	// uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("KubeConfig loaded....")

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Client init....")

	log.Println("Starting server -- http://localhost:" + "9180")
	http.ListenAndServe(":9180", ghandlers.CORS(ghandlers.IgnoreOptions())(handlers.GetRoutes(clientset)))

}
