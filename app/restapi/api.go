package restapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/app/broker"
	"github.com/indrenicloud/tricloud-server/app/logg"
	"github.com/indrenicloud/tricloud-server/app/monitor"
)

var cbroker *broker.Broker
var sMonitor *monitor.Monitor

// Main Router
func GetMainRouter(b *broker.Broker) *mux.Router {
	cbroker = b
	sMonitor = monitor.NewMonitor(b.GetEventManager())
	go sMonitor.Run()
	r := mux.NewRouter()
	registerUserAPI(r.PathPrefix("/api").Subrouter())
	registerAuthHandlers(r.PathPrefix("/login").Subrouter())
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	r.HandleFunc("/websocket/{apikey}", b.ServeUserWebsocket)
	r.HandleFunc("/", rootRoute)
	r.Use(Logger)
	return r
}

func registerAuthHandlers(r *mux.Router) {
	r.HandleFunc("/signin", SignIn).Methods("POST")
	r.Use(MiddlewareJson)
}

func registerUserAPI(r *mux.Router) {
	logg.Info("Welcome to TriCloud REST_API")

	r.HandleFunc("/users", GetUsers).Methods("GET")
	r.HandleFunc("/users", CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")
	r.HandleFunc("/users/{id}/alerts", GetAlerts).Methods("GET")
	r.HandleFunc("/user/api", GetApiKeys).Methods("GET")
	r.HandleFunc("/user/api", AddApiKeys).Methods("POST")
	r.HandleFunc("/user/api/{key}", RemoveApiKeys).Methods("DELETE")

	r.HandleFunc("/agents", GetAgents).Methods("GET")
	r.HandleFunc("/agents/{id}", GetAgent).Methods("GET")
	r.HandleFunc("/agents/{id}", DeleteAgent).Methods("DELETE")
	r.HandleFunc("/agents/{id}/status", GetAgentStatus).Methods("POST")

	r.HandleFunc("/websites", GetWebsites).Methods("GET")
	r.HandleFunc("/websites", CreateWebsite).Methods("POST")
	r.HandleFunc("/websites/{name}", GetWebsite).Methods("GET")
	r.HandleFunc("/websites/{name}", DeleteWebsite).Methods("DELETE")
	r.HandleFunc("/scripts", GetScripts).Methods("GET")
	r.HandleFunc("/scripts", CreateScript).Methods("POST")
	r.HandleFunc("/scripts/{name}", GetScript).Methods("GET")
	r.HandleFunc("/scripts/{name}", DeleteScript).Methods("DELETE")

	r.Use(MiddlewareSession, MiddlewareJson)
}

func GetAgentRouter(b *broker.Broker) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/registeragent", RegisterAgent).Methods("POST")
	r.HandleFunc("/updatesysinfo/{key}", UpdateSystemInfo).Methods("PUT")
	r.HandleFunc("/websocket/{key}", b.ServeAgentWebsocket)

	return r
}
