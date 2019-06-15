package restapi

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/indrenicloud/tricloud-server/app/database/statstore"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/app/database"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	// only if superuser
	if !isSuperUser(r) {
		errorResp(w, ErrorNotAuthorized)
		return
	}

	users, err := database.GetAllUsers()
	if err != nil {
		errorResp(w, err)
		return
	}
	generateResp(w, users, err)

}
func GetUser(w http.ResponseWriter, r *http.Request) {
	// only if superuser or itself
	vars := mux.Vars(r)
	ID := vars["id"]

	user, err := database.GetUser(ID)
	if !isAuthorized(user.ID, r) {
		errorResp(w, ErrorNotAuthorized)
		return
	}

	generateResp(w, user, err)
}
func CreateUser(w http.ResponseWriter, r *http.Request) {
	// only if superuser
	if !isSuperUser(r) {
		errorResp(w, ErrorNotAuthorized)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		errorResp(w, err)
		return
	}
	defer r.Body.Close()
	var userinfo map[string]string
	err = deJson(body, &userinfo)
	if err != nil {
		errorResp(w, err)
		return
	}
	usr, err := database.NewUser(userinfo, false)
	if err != nil {
		errorResp(w, err)
		return
	}
	database.CreateUser(usr)
	updatedusers, err := database.GetAllUsers()
	generateResp(w, updatedusers, err)

}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// only if superuser or that user but cannot change superuser flag
	vars := mux.Vars(r)
	id := vars["id"]
	if !isAuthorized(id, r) {
		errorResp(w, ErrorNotAuthorized)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		errorResp(w, err)
		return
	}
	defer r.Body.Close()

	var userinfo map[string]string
	err = deJson(body, &userinfo)
	if err != nil {
		errorResp(w, err)
		return
	}
	userinfo["id"] = id

	err = database.UpdateUser(userinfo)
	if err != nil {
		errorResp(w, err)
		return
	}
	updateduser, err := database.GetUser(id)
	if err != nil {
		errorResp(w, err)
		return
	}
	updateduser.Password = ""
	generateResp(w, updateduser, nil)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// only if superuser

	if !isSuperUser(r) {
		errorResp(w, ErrorNotAuthorized)
		return
	}

	vars := mux.Vars(r)
	id := vars["id"]
	logg.Warn("Deleting user")
	err := database.DeleteUser(id)
	if err != nil {
		errorResp(w, err)
	}

	updatedusers, err := database.GetAllUsers()
	if err != nil {
		errorResp(w, err)
		return
	}

	generateResp(w, updatedusers, err)

}

func GetAgents(w http.ResponseWriter, r *http.Request) {
	// only if owns agent
	user, _ := parseUser(r)

	agents, err := database.GetAllUserAgents(user)
	if err != nil {
		errorResp(w, err)
		return
	}
	activeagents := cbroker.GetActiveAgents(user)
	for _, agent := range agents {
		connid, ok := activeagents[agent.ID]
		if ok {
			agent.Active = true
			agent.ActiveConnid = connid
		}
	}
	generateResp(w, agents, err)
}

func GetAgent(w http.ResponseWriter, r *http.Request) {
	// only if user owns agent or superuser
	vars := mux.Vars(r)
	ID := vars["id"]
	agent, err := database.GetAgent(ID)

	if err != nil {
		errorResp(w, err)
		return
	}

	if !isAuthorized(agent.Owner, r) {
		errorResp(w, ErrorNotAuthorized)
		return
	}

	activeagents := cbroker.GetActiveAgents(agent.Owner)
	connid, ok := activeagents[agent.ID]
	if ok {
		agent.Active = true
		agent.ActiveConnid = connid
	}

	generateResp(w, agent, err)
}

func DeleteAgent(w http.ResponseWriter, r *http.Request) {
	// only if user owns agent or superuser
	vars := mux.Vars(r)
	ID := vars["id"]
	agent, err := database.GetAgent(ID)

	if err != nil {
		errorResp(w, err)
		return
	}

	if !isAuthorized(agent.Owner, r) {
		errorResp(w, ErrorNotAuthorized)
		return
	}

	database.DeleteAgent(ID, agent.Owner)
	cbroker.RemoveAgent(ID, agent.Owner)

	agents, err := database.GetAllAgents()
	generateResp(w, agents, err)
}

func GetApiKeys(w http.ResponseWriter, r *http.Request) {
	// if user
	user, _ := parseUser(r)
	keys, err := database.GetapiKeys(user)

	generateResp(w, keys, err)
}

func AddApiKeys(w http.ResponseWriter, r *http.Request) {
	// if user
	user, _ := parseUser(r)
	err := database.AddapiKey(user, "agent")
	if err != nil {
		errorResp(w, err)
		return
	}
	generateResp(w, "ok", err)
}

func RemoveApiKeys(w http.ResponseWriter, r *http.Request) {

	// if user
	vars := mux.Vars(r)
	key := vars["key"]

	user, _ := parseUser(r)

	err := database.RemoveapiKey(user, key)
	if err != nil {
		errorResp(w, err)
		return
	}
	generateResp(w, "ok", err)
}

//{offset}/{noentries}
func GetAgentStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["id"]

	agent, err := database.GetAgent(ID)
	if err != nil {
		errorResp(w, err)
		return
	}

	if !isAuthorized(agent.Owner, r) {
		errorResp(w, ErrorNotAuthorized)
		return
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		errorResp(w, err)
		return
	}
	defer r.Body.Close()

	var statusparm map[string]int64
	err = deJson(body, &statusparm)
	if err != nil {
		errorResp(w, err)
		return
	}

	offset, _ := statusparm["offset"]

	noentries, _ := statusparm["noofentries"]

	generateResp(w, statstore.GetStats(agent.ID, offset, noentries), nil)
}

func GetAgentAlerts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ID := vars["id"]

	agent, err := database.GetAgent(ID)
	if err != nil {
		errorResp(w, err)
		return
	}

	if !isAuthorized(agent.Owner, r) {
		errorResp(w, ErrorNotAuthorized)
		return
	}
	bytbyt, err := statstore.GetAlert([]byte(ID))

	generateResp(w, bytbyt, err)
}
