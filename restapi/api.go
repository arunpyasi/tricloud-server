package restapi

import (
	"context"
	"encoding/json"
	"errors"
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
	r.HandleFunc("/agents", CreateAgent).Methods("POST")
	r.HandleFunc("/agents/{id}", GetAgent).Methods("GET")
	r.HandleFunc("/agents/{id}", UpdateAgent).Methods("PUT")
	r.HandleFunc("/agents/{id}", DeleteAgent).Methods("DELETE")

	r.HandleFunc("/agents/{id}/hostinfo", CreateHostInfo).Methods("POST")
	r.HandleFunc("/agents/{id}/hostinfo", GetHostInfo).Methods("GET")
	r.HandleFunc("/agents/{id}/hostinfo", UpdateHostInfo).Methods("PUT")
	r.HandleFunc("/agents/{id}/cpuinfo", CreateCPUInfo).Methods("POST")
	r.HandleFunc("/agents/{id}/cpuinfo", GetCPUInfo).Methods("GET")
	r.HandleFunc("/agents/{id}/cpuinfo", UpdateCPUInfo).Methods("PUT")

	r.Use(MiddlewareSession, MiddlewareJson)
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
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	var user database.User
	json.Unmarshal(body, &user)
	if id == user.ID {

		database.UpdateUser(id, user)
		updated_users, err := database.GetUser(id)
		resp := GenerateResponse(updated_users, err)
		w.Write(resp)
	} else {
		resp := GenerateResponse(nil, errors.New("Invalid ID"))
		w.Write(resp)
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	database.DeleteUser(id)
	updated_users, err := database.GetAllUsers()
	resp := GenerateResponse(updated_users, err)
	w.Write(resp)

}

func GetAgents(w http.ResponseWriter, r *http.Request) {
	agents, err := database.GetAllAgents()
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	resp := GenerateResponse(agents, err)
	w.Write(resp)
}
func GetAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["id"]
	user, err := database.GetAgent(ID)
	resp := GenerateResponse(user, err)
	w.Write(resp)
}
func CreateAgent(w http.ResponseWriter, r *http.Request) {
	var agent database.Agent
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	json.Unmarshal(body, &agent)
	database.CreateAgent(agent)
	updated_agent, err := database.GetAllAgents()
	resp := GenerateResponse(updated_agent, err)
	w.Write(resp)
}
func UpdateAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var agent database.Agent
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	json.Unmarshal(body, &agent)
	if id == agent.ID {
		database.UpdateAgent(id, agent)
		updated_users, err := database.GetAgent(id)
		resp := GenerateResponse(updated_users, err)
		w.Write(resp)
	} else {
		resp := GenerateResponse(nil, errors.New("Invalid ID"))
		w.Write(resp)
	}
}

func DeleteAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["id"]
	database.DeleteAgent(ID)
	agents, err := database.GetAllAgents()
	resp := GenerateResponse(agents, err)
	w.Write(resp)
}

func GetHostInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	host_info, err := database.GetHostInfo(id)
	resp := GenerateResponse(host_info, err)
	w.Write(resp)
}

func CreateHostInfo(w http.ResponseWriter, r *http.Request) {
	var hostinfo *database.HostInfo
	vars := mux.Vars(r)
	id := vars["id"]
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	json.Unmarshal(body, &hostinfo)
	database.CreateHostInfo(id, hostinfo)
	host_info, err := database.GetHostInfo(id)
	resp := GenerateResponse(host_info, err)
	w.Write(resp)
}

func UpdateHostInfo(w http.ResponseWriter, r *http.Request) {
}

func GetCPUInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	cpu_info, err := database.GetCPUInfo(id)
	resp := GenerateResponse(cpu_info, err)
	w.Write(resp)
}

func CreateCPUInfo(w http.ResponseWriter, r *http.Request) {
	var cpuinfo *database.CPUInfo
	vars := mux.Vars(r)
	id := vars["id"]
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	json.Unmarshal(body, &cpuinfo)
	database.CreateCPUInfo(id, cpuinfo)
	cpu_info, err := database.GetCPUInfo(id)
	resp := GenerateResponse(cpu_info, err)
	w.Write(resp)
}

func UpdateCPUInfo(w http.ResponseWriter, r *http.Request) {

}
