package main

import (
	"net/http"

	"github.com/indrenicloud/tricloud-server/broker"

	"github.com/gorilla/mux"
)

func main() {

	mBroker := broker.NewBroker()

	mainRouter := getRouter()
	mainRouter.HandleFunc("/websocket", mBroker.ServeAgentWebsocket)
	go func() {
		r := http.NewServeMux()
		r.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
		r.HandleFunc("/websocket", mBroker.ServeAgentWebsocket)
		http.ListenAndServe(":8081", r)

	}()

	http.ListenAndServe(":8080", mainRouter)

}

// stub for rest api
func getRouter() *mux.Router {

	r := mux.NewRouter()
	return r
}
