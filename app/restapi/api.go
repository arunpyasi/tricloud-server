package restapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/app/broker"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

// Main Router
func GetMainRouter(b *broker.Broker) *mux.Router {

	r := mux.NewRouter()
	registerUserAPI(r.PathPrefix("/api").Subrouter())
	registerAuthHandlers(r.PathPrefix("/login").Subrouter())
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
	r.HandleFunc("/websocket", b.ServeUserWebsocket)
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
	r.HandleFunc("/user/api", GetApiKeys).Methods("GET")
	r.HandleFunc("/user/api", AddApiKeys).Methods("PUT")
	r.HandleFunc("/user/api/{key}", RemoveApiKeys).Methods("DELETE")

	r.HandleFunc("/agents", GetAgents).Methods("GET")
	r.HandleFunc("/agents/{id}", GetAgent).Methods("GET")
	r.HandleFunc("/agents/{id}", DeleteAgent).Methods("DELETE")
	r.HandleFunc("/agents/status/{id}", GetAgentStatus).Methods("POST")
	r.Use(MiddlewareSession, MiddlewareJson)
}

func GetAgentRouter(b *broker.Broker) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/registeragent", RegisterAgent).Methods("POST")
	r.HandleFunc("/updatesysinfo/{key}", UpdateSystemInfo).Methods("PUT")
	r.HandleFunc("/websocket/{key}", b.ServeAgentWebsocket)

	return r
}
