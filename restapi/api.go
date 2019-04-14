package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/indrenicloud/tricloud-server/restapi/auth"

	"net/http"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/restapi/database"
)

// when using refelection to find type  using custom type avoids collision in contex.value
type contextkey int

const ContextUser contextkey = iota

var (
	ErrorNotAuthorized = errors.New("Not authorized")
)

func MiddlewareJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}

// MiddlewareSession checks the session for request and tags username to request context
func MiddlewareSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token := auth.ParseAPIKey(r.Header.Get("Api-key"))
		claims, ok := token.Claims.(auth.MyClaims)

		if !ok || !token.Valid {
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ContextUser, claims)
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
	r.HandleFunc("/agents/{id}", GetAgent).Methods("GET")
	r.HandleFunc("/agents/{id}", DeleteAgent).Methods("DELETE")
	r.Use(MiddlewareSession, MiddlewareJson)
}

func GenerateResponse(data interface{} /*datatype string,*/, err error) []byte {
	var response []byte

	if data != nil || err == nil {
		m := make(map[string]interface{})
		m["status"] = "ok"
		/*m["datatype"] = datatype*/
		m["data"] = data
		response, _ = json.Marshal(m)
	} else {
		response = []byte(`{"msg":"` + err.Error() + `","status":"failed"}`)
	}
	return response
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	// only if superuser
	if _, super := parseUser(r); !super {
		w.Write(GenerateResponse(nil, ErrorNotAuthorized))
		return
	}

	users, err := database.GetAllUsers()
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	resp := GenerateResponse(users, err)
	w.Write(resp)

}
func GetUser(w http.ResponseWriter, r *http.Request) {
	// only if superuser or itself
	vars := mux.Vars(r)
	ID := vars["id"]
	if apiuser, super := parseUser(r); !super {
		if ID != apiuser {
			w.Write(GenerateResponse(nil, ErrorNotAuthorized))
			return
		}
	}
	user, err := database.GetUser(ID)
	if err != nil {
		w.Write(GenerateResponse(nil, err))
		return
	}
	resp := GenerateResponse(user, err)
	w.Write(resp)
}
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// only if superuser
	if _, super := parseUser(r); !super {
		w.Write(GenerateResponse(nil, ErrorNotAuthorized))
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()
	var userinfo map[string]interface{}
	err = json.Unmarshal(body, &userinfo)
	if err != nil {
		w.Write(GenerateResponse(nil, err))
	}
	usr, err := database.NewUser(userinfo, false)
	if err != nil {
		w.Write(GenerateResponse(nil, err))
	}
	database.CreateUser(usr)
	updatedusers, err := database.GetAllUsers()
	w.Write(GenerateResponse(updatedusers, err))

}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// only if superuser or that user but cannot change superuser flag
	vars := mux.Vars(r)
	id := vars["id"]
	if apiuser, super := parseUser(r); !super {
		if id != apiuser {
			w.Write(GenerateResponse(nil, ErrorNotAuthorized))
			return
		}
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	var userinfo map[string]interface{}
	json.Unmarshal(body, &userinfo)
	userinfo["id"] = id

	database.UpdateUser(userinfo)
	updated_users, err := database.GetUser(id)
	resp := GenerateResponse(updated_users, err)
	w.Write(resp)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// only if superuser
	if _, super := parseUser(r); !super {
		w.Write(GenerateResponse(nil, ErrorNotAuthorized))
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	database.DeleteUser(id)
	updated_users, err := database.GetAllUsers()
	resp := GenerateResponse(updated_users, err)
	w.Write(resp)

}

func GetAgents(w http.ResponseWriter, r *http.Request) {
	// only if owns agent
	user, _ := parseUser(r)

	agents, err := database.GetAllUserAgents(user)
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	resp := GenerateResponse(agents, err)
	w.Write(resp)
}
func GetAgent(w http.ResponseWriter, r *http.Request) {
	// only if user owns agent or superuser
	vars := mux.Vars(r)
	ID := vars["id"]
	agent, err := database.GetAgent(ID)

	if err != nil {
		//not found return
		w.Write(GenerateResponse(nil, err))
		return
	}

	user, super := parseUser(r)
	if !super {
		if user != agent.Owner {
			w.Write(GenerateResponse(nil, ErrorNotAuthorized))
			return
		}
	}

	resp := GenerateResponse(user, err)
	w.Write(resp)
}

func DeleteAgent(w http.ResponseWriter, r *http.Request) {
	// only if user owns agent or superuser
	vars := mux.Vars(r)
	ID := vars["id"]

	agent, err := database.GetAgent(ID)

	if err != nil {
		//not found return
		w.Write(GenerateResponse(nil, err))
		return
	}

	user, super := parseUser(r)
	if !super {
		if user != agent.Owner {
			w.Write(GenerateResponse(nil, ErrorNotAuthorized))
			return
		}
	}

	database.DeleteAgent(ID)
	agents, err := database.GetAllAgents()
	resp := GenerateResponse(agents, err)
	w.Write(resp)
}

func parseUser(r *http.Request) (string, bool) {
	c := r.Context().Value(ContextUser)
	claims, ok := c.(auth.MyClaims)
	if !ok {
		return "", false
	}
	user, err := database.GetUser(claims.User)
	if err != nil {
		return "", false
	}
	return claims.User, user.SuperUser
}
