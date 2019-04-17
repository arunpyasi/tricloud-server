package restapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/restapi/auth"
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

		if !token.Valid {
			log.Println("token invalid")
			http.Error(w, "not authorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*auth.MyClaims)
		log.Println(claims)

		if !ok {
			http.Error(w, "INternal err", http.StatusInternalServerError)
		}

		if !token.Valid {
			log.Println("claims not valid")
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
	r.HandleFunc("/user/api", GetApiKeys).Methods("GET")
	r.HandleFunc("/user/api", AddApiKeys).Methods("PUT")
	r.HandleFunc("/user/api/{key}", RemoveApiKeys).Methods("DELETE")

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
		w.Write(GenerateResponse(nil, err))
		fmt.Printf("error: %s", err)
		return
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
		w.Write(GenerateResponse(nil, err))
	}
	defer r.Body.Close()
	var userinfo map[string]string
	err = json.Unmarshal(body, &userinfo)
	if err != nil {
		w.Write(GenerateResponse(nil, err))
		return
	}
	usr, err := database.NewUser(userinfo, false)
	if err != nil {
		w.Write(GenerateResponse(nil, err))
		return
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

	var userinfo map[string]string
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
	log.Print("should delete")
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

	resp := GenerateResponse(agent, err)
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

func GetApiKeys(w http.ResponseWriter, r *http.Request) {
	// if user
	user, _ := parseUser(r)
	keys, err := database.GetapiKeys(user)
	GenerateResponse(keys, err)
}

func AddApiKeys(w http.ResponseWriter, r *http.Request) {
	// if user
	user, _ := parseUser(r)
	err := database.AddapiKey(user, "agent")
	if err == nil {
		GenerateResponse("ok", nil)
		return
	}
	GenerateResponse(nil, err)
}

func RemoveApiKeys(w http.ResponseWriter, r *http.Request) {

	// if user
	vars := mux.Vars(r)
	key := vars["key"]

	user, _ := parseUser(r)

	err := database.RemoveapiKey(user, key)
	if err == nil {
		GenerateResponse("ok", nil)
		return
	}
	GenerateResponse(nil, err)

}

func parseUser(r *http.Request) (string, bool) {
	c := r.Context().Value(ContextUser)
	claims, ok := c.(*auth.MyClaims)
	if !ok {
		return "", false
	}
	return claims.User, claims.Super
}
