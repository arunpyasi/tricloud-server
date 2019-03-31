package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/rest_api/database"
)

func MiddlewareJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}

func main() {
	fmt.Println("Welcome to TriCloud REST_API")
	r := mux.NewRouter() // Here, r is router

	r.HandleFunc("/users", GetUsers).Methods("GET")
	r.HandleFunc("/users", CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", GetUser).Methods("GET")
	r.HandleFunc("/users/{id}", UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", DeleteUser).Methods("DELETE")

	r.HandleFunc("/agents", GetAgents).Methods("GET")
	r.HandleFunc("/agent", CreateAgent).Methods("POST")
	r.HandleFunc("/agent/{id}", GetAgent).Methods("GET")
	r.HandleFunc("/agent/{id}", UpdateAgent).Methods("PUT")
	r.HandleFunc("/agent/{id}", DeleteAgent).Methods("DELETE")
	r.Use(MiddlewareJson)
	defer database.Conn.Close()
	log.Fatal(http.ListenAndServe(":8000", r)) //listening and serving
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetAllUsers()
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	w.Write(users)

}
func GetUser(w http.ResponseWriter, r *http.Request) {}
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user database.User
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	json.Unmarshal(body, &user)
	database.CreateUser(user)
	updated_users, _ := database.GetAllUsers()
	w.Write(updated_users)

}
func UpdateUser(w http.ResponseWriter, r *http.Request) {}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["id"]
	database.DeleteUser(ID)
	updated_users, _ := database.GetAllUsers()
	w.Write(updated_users)

}

func GetAgents(w http.ResponseWriter, r *http.Request) {

}

func GetAgent(w http.ResponseWriter, r *http.Request) {}
func CreateAgent(w http.ResponseWriter, r *http.Request) {
}
func UpdateAgent(w http.ResponseWriter, r *http.Request) {}
func DeleteAgent(w http.ResponseWriter, r *http.Request) {}
