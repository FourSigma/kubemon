package handlers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/FourSigma/kubemon/models"
	"github.com/gorilla/mux"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/pkg/api"
	v1 "k8s.io/client-go/1.4/pkg/api/v1"
)

func GetRoutes(client *kubernetes.Clientset) *mux.Router {

	wh := &WebHandler{client: client}

	r := mux.NewRouter()
	r.Headers("Access-Control-Allow-Origin", "*")
	r.Headers("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	r.HandleFunc("/pods", wh.GetPodListHandler).Methods("GET")             //List of Pods
	r.HandleFunc("/pods", wh.DeleteAllPodsHandler).Methods("OPTIONS")      //Delete Pods
	r.HandleFunc("/pods/{id}", wh.DeletePodByIdHandler).Methods("OPTIONS") //Delete Pod by Id
	r.HandleFunc("/pods/{id}/logs", wh.PodLogHandler)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("dist/"))))
	return r
}

type WebHandler struct {
	client *kubernetes.Clientset
}

func (w *WebHandler) PodLogHandler(rw http.ResponseWriter, r *http.Request) {
	flusher, ok := rw.(http.Flusher)
	if !ok {
		http.Error(rw, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	fmt.Println("\n\nClient subscribed\n\n")
	podId := mux.Vars(r)["id"]
	fmt.Println("PodId", podId)

	rw.Header().Set("Content-Type", "text/event-stream")
	rw.Header().Set("Cache-Control", "no-cache")
	rw.Header().Set("Connection", "keep-alive")

	rc, err := handleLog(w.client, podId)
	if err != nil {
		fmt.Fprintln(rw, err)
		return
	}

	defer rc.Close() //Close the PodLog Reader
	buf := bufio.NewReader(rc)

	for {
		b, _, err := buf.ReadLine()
		if err != nil {
			log.Println("Buffer read line error -- ", err)
			break
		}
		fmt.Fprintln(rw, string(b), "\n")
		flusher.Flush()
	}

}

func (w *WebHandler) DeletePodByIdHandler(rw http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := w.client.Core().Pods(api.NamespaceDefault).Delete(id, &api.DeleteOptions{})
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Fprintln(rw, "Deleted: ", id)
	return
}

func (w *WebHandler) DeleteAllPodsHandler(rw http.ResponseWriter, r *http.Request) {
	err := w.client.Core().Pods(api.NamespaceDefault).DeleteCollection(&api.DeleteOptions{}, api.ListOptions{})
	if err != nil {
		fmt.Println("An ERROR has occured....")
		fmt.Fprintln(rw, err)
		return
	}
	fmt.Fprintln(rw, "DELETING ALL PODS....")
	return
}

func (w *WebHandler) GetPodListHandler(rw http.ResponseWriter, r *http.Request) {
	pods, err := w.client.Core().Pods(api.NamespaceDefault).List(api.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	pl := make([]models.Pod, len(pods.Items))

	for i, _ := range pods.Items {
		pl[i].Pod = &pods.Items[i]
	}
	p, err := json.Marshal(pl)

	fmt.Fprintln(rw, string(p))

}

//Copied and modified from kubectl implementation
func handleLog(c *kubernetes.Clientset, podID string) (rc io.ReadCloser, err error) {

	req := c.Core().Pods(api.NamespaceDefault).GetLogs(podID, &v1.PodLogOptions{Follow: true, Timestamps: true})
	rc, err = req.Stream()
	if err != nil {
		fmt.Println("Log streaming error...")
		return nil, err
	}
	//defer readCloser.Close()
	//_, err = io.Copy(out, readCloser)
	return rc, nil
}
