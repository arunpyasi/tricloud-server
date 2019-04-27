package restapi

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/app/auth"
	"github.com/indrenicloud/tricloud-server/app/database"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

func RegisterAgent(h http.ResponseWriter, r *http.Request) {

	rawapi := r.Header.Get("Api-key")
	if rawapi == "" {
		logg.Warn("Token Not set @ agent conn")
		http.Error(h, "Not Authorized", http.StatusUnauthorized)
		return
	}
	token := auth.ParseAPIKey(rawapi)
	if token == nil {
		http.Error(h, "Invalid api key", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*auth.MyClaims)

	if !ok || !token.Valid {
		http.Error(h, "not authorized", http.StatusUnauthorized)
		return
	}

	agentid, err := database.CreateAgent(claims.User)
	if err != nil {
		http.Error(h, "not authorized", http.StatusUnauthorized)
	}

	h.Write(GenerateResponse(agentid, nil))
}

func UpdateSystemInfo(h http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		h.Write([]byte("couldnot read body error"))
		return
	}
	defer r.Body.Close()

	var agentinfo map[string]string
	json.Unmarshal(body, &agentinfo)
	logg.Info(agentinfo)
	/*
		delete(agentinfo, "bootTime")
		delete(agentinfo, "procs")
		delete(agentinfo, "uptime")
	*/

	err = database.UpdateSystemInfo(key, agentinfo)

	if err != nil {
		logg.Warn(err)
		h.Write([]byte("Db error"))
		return
	}

}
