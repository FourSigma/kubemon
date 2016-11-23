package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/FourSigma/kubemon/models"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"

	"k8s.io/client-go/1.4/kubernetes"
	"k8s.io/client-go/1.4/pkg/api"
	v1 "k8s.io/client-go/1.4/pkg/api/v1"
)

func GetRoutes(client *kubernetes.Clientset) *mux.Router {

	wh := &WebHandler{client: client}

	r := mux.NewRouter()
	r.HandleFunc("/websocket", wh.WebSocketHandler)
	r.HandleFunc("/pods", wh.GetPodListHandler).Methods("GET")           //List of Pods
	r.HandleFunc("/pods/delete", wh.DeleteAllPodsHandler).Methods("GET") //Pod Details
	r.HandleFunc("/pods/{id}", wh.DeletePodByIdHandler).Methods("GET")   //Pod Details
	r.HandleFunc("/pods/:id/logs", PodLogHandler)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("static/"))))
	return r
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func MyWebSocketHandler(rw http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("\n\nClient subscribed\n\n", conn)

	var i int = 0
	for {
		time.Sleep(time.Millisecond)
		err := conn.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(i)))
		if err != nil {
			log.Println(err)
			conn.Close()
			break
		}
		fmt.Println(i)
		i++
	}

}

type WebHandler struct {
	client *kubernetes.Clientset
}

func (w *WebHandler) WebSocketHandler(rw http.ResponseWriter, r *http.Request) {
	// conn, err := upgrader.Upgrade(rw, r, nil)
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }
	fmt.Println("\n\nClient subscribed\n\n")

	for {

		err := handleLog(w.client, "com-core.svc.http.core-58963565-neaoj", rw)
		if err != nil {
			fmt.Fprintln(rw, "Socket ERROR")
			return
		}
	}
	// select {
	// 	err := conn.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(i)))
	// 	if err != nil {
	// 		log.Println(err)
	// 		conn.Close()
	// 		break
	// 	}
	// 	fmt.Println(i)
	// 	i++
	// }

}

func (w *WebHandler) Error(rw http.ResponseWriter, r *http.Request) {
	return
}
func (w *WebHandler) DeletePodByIdHandler(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := w.client.Core().Pods(api.NamespaceDefault).Delete(id, &api.DeleteOptions{})
	if err != nil {
		fmt.Fprintln(rw, err)
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
	pl := make([]*models.Pod, len(pods.Items))
	for i, _ := range pods.Items {
		pl[i] = &models.Pod{}
		pl[i].Pod = &pods.Items[i]
	}
	p, err := json.Marshal(pl)

	//fmt.Fprintln(rw, len(pods.Items), string(p))
	fmt.Fprintln(rw, string(p))

}

//Copied from kubectl implementation
func handleLog(c *kubernetes.Clientset, podID string, out io.Writer) error {

	// TODO: transform this into a PodLogOptions call
	req := c.Core().Pods(api.NamespaceDefault).GetLogs(podID, &v1.PodLogOptions{Follow: true, Timestamps: true})

	readCloser, err := req.Stream()
	if err != nil {
		fmt.Println("WE have an error")
		return err
	}

	defer readCloser.Close()
	_, err = io.Copy(out, readCloser)
	return err
}
func PodIdHandler(rw http.ResponseWriter, r *http.Request) {

}

func PodLogHandler(rw http.ResponseWriter, r *http.Request) {

}
