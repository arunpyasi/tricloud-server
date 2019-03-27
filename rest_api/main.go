package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func MiddlewareJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}
func final(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing finalHandler")
	w.Write([]byte("OK"))
}

func main() {
	fmt.Println("Welcome to TriCloud REST_API")
	r := mux.NewRouter() // Here, r is router
	agent = append(agent, Agent{ID: "1", Name: "Lenovo-B50", OS: "Ubuntu/Linux"})

	r.HandleFunc("/users", GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")

	r.HandleFunc("/agents", GetAgents).Methods("GET")
	r.HandleFunc("/agent/{id}", GetAgent).Methods("GET")
	r.HandleFunc("/agent/{id}", CreateAgent).Methods("POST")
	r.HandleFunc("/agent/{id}", UpdateAgent).Methods("PUT")
	r.HandleFunc("/agent/{id}", DeleteAgent).Methods("DELETE")
	r.Use(MiddlewareJson)
	log.Fatal(http.ListenAndServe(":8000", r)) //listening and serving router
}

type Agent struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	OS   string `json:"os,omitempty"`
}

var agent []Agent

func GetUsers(w http.ResponseWriter, r *http.Request)   {}
func GetUser(w http.ResponseWriter, r *http.Request)    {}
func CreateUser(w http.ResponseWriter, r *http.Request) {}
func UpdateUser(w http.ResponseWriter, r *http.Request) {}
func DeleteUser(w http.ResponseWriter, r *http.Request) {}

func GetAgents(w http.ResponseWriter, r *http.Request) {
	j, _ := json.Marshal(agent)
	w.Write(j)
}

func GetAgent(w http.ResponseWriter, r *http.Request)    {}
func CreateAgent(w http.ResponseWriter, r *http.Request) {}
func UpdateAgent(w http.ResponseWriter, r *http.Request) {}
func DeleteAgent(w http.ResponseWriter, r *http.Request) {}
