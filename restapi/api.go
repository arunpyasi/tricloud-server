package restapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/restapi/database"
)

// when using refelection to find type  using custom type avoids collision in contex.value
//type key int
//const UserType key = iota

func MiddlewareJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}

// MiddlewareSession checks the session for request and tags username to request context
func MiddlewareSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		usr, err := database.GetUserFromSession(r)
		if err != nil {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", usr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RegisterAPI(r *mux.Router) {
	fmt.Println("Welcome to TriCloud REST_API")

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
}

func GenerateResponse(data []byte, err error) []byte {
	var response []byte

	fmt.Print(string(data))
	fmt.Print(err)
	if data != nil || err == nil {
		m := make(map[string]interface{})
		json.Unmarshal(data, &m)
		m["status"] = "ok"
		response, _ = json.Marshal(m)
	} else {
		response = []byte(`{"msg":"` + err.Error() + `","status":"failed"}`)
	}
	return response
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := database.GetAllUsers()
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	resp := GenerateResponse(users, err)
	w.Write(resp)

}
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["id"]
	user, err := database.GetUser(ID)
	resp := GenerateResponse(user, err)
	w.Write(resp)
}
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user database.User
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	json.Unmarshal(body, &user)
	database.CreateUser(user)
	updated_users, err := database.GetAllUsers()
	resp := GenerateResponse(updated_users, err)
	w.Write(resp)

}
func UpdateUser(w http.ResponseWriter, r *http.Request) {}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["id"]
	database.DeleteUser(ID)
	updated_users, err := database.GetAllUsers()
	resp := GenerateResponse(updated_users, err)
	w.Write(resp)

}

func GetAgents(w http.ResponseWriter, r *http.Request)   {}
func GetAgent(w http.ResponseWriter, r *http.Request)    {}
func CreateAgent(w http.ResponseWriter, r *http.Request) {}
func UpdateAgent(w http.ResponseWriter, r *http.Request) {}
func DeleteAgent(w http.ResponseWriter, r *http.Request) {}
