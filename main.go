package main

import (
	"log"
	"net/http"

	"github.com/indrenicloud/tricloud-server/broker"
	"github.com/indrenicloud/tricloud-server/restapi"
	"github.com/indrenicloud/tricloud-server/restapi/auth"
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
	defer database.DB.Close()

	r.HandleFunc("/", rootRoute)
	r.HandleFunc("/websocket", mBroker.ServeUserWebsocket)

	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))

	r.Use(Logger)

	log.Println(http.ListenAndServe(":8080", r))

}

func listenAgentsConnection() {
	agentsRouter := mux.NewRouter()
	agentsRouter.HandleFunc("/websocket/{key}", mBroker.ServeAgentWebsocket)
	agentsRouter.HandleFunc("/registeragent", restapi.RegisterAgent).Methods("PUT")
	agentsRouter.HandleFunc("/updatesysinfo/{key}", restapi.RegisterAgent).Methods("POST")
	log.Println(http.ListenAndServe(":8081", agentsRouter))
}

func rootRoute(h http.ResponseWriter, r *http.Request) {
	token := auth.ParseAPIKey(r.Header.Get("Api-key"))
	_, ok := token.Claims.(auth.MyClaims)
	if !ok || !token.Valid {
		http.ServeFile(h, r, "./public/login.html")
		return
	}
	http.ServeFile(h, r, "./public/dashboard.html")
}

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("INFOLOG", r.URL.Path)
		h.ServeHTTP(w, r)
	})
}
