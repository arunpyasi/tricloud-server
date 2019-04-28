package restapi

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/indrenicloud/tricloud-server/app/database/statstore"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/app/auth"
	"github.com/indrenicloud/tricloud-server/app/database"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

func GenerateResponse(data interface{}, err error) []byte {
	var response []byte

	if data != nil || err == nil {
		m := make(map[string]interface{})
		m["status"] = "ok"
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
		logg.Warn(err)
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
	logg.Warn("Deleting user")
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
		logg.Error(err)
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

//{offset}/{noentries}
func GetAgentStatus(w http.ResponseWriter, r *http.Request) {
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

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	var statusparm map[string]int64
	json.Unmarshal(body, &statusparm)

	offset, _ := statusparm["offset"]

	noentries, _ := statusparm["noofentries"]

	GenerateResponse(statstore.GetStats(agent.ID, offset, noentries), nil)
}

func parseUser(r *http.Request) (string, bool) {
	c := r.Context().Value(ContextUser)
	claims, ok := c.(*auth.MyClaims)
	if !ok {
		return "", false
	}
	return claims.User, claims.Super
}
