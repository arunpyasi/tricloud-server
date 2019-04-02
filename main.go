package main

import (
	"net/http"

	"github.com/indrenicloud/tricloud-server/broker"
	"github.com/indrenicloud/tricloud-server/restapi"
	"github.com/indrenicloud/tricloud-server/restapi/database"

	"github.com/gorilla/mux"
)

var mBroker *broker.Broker

func main() {

	mBroker = broker.NewBroker()

	go listenAgentsConnection()

	r := mux.NewRouter()
	restapi.RegisterAPI(r.PathPrefix("/api").Subrouter())
	restapi.RegisterAuthHandlers(r.PathPrefix("/login").Subrouter())
	defer database.Close()

	r.HandleFunc("/", rootRoute)

	r.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))
	r.HandleFunc("/websocket", mBroker.ServeUserWebsocket)
	http.ListenAndServe(":8080", r)

}

func listenAgentsConnection() {
	mainRouter := mux.NewRouter()
	mainRouter.HandleFunc("/websocket", mBroker.ServeAgentWebsocket)
	http.ListenAndServe(":8081", mainRouter)

}

func rootRoute(h http.ResponseWriter, r *http.Request) {
	http.ServeFile(h, r, "./public/index.html")
}
