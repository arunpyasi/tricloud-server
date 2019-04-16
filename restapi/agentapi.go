package restapi

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/restapi/auth"
	"github.com/indrenicloud/tricloud-server/restapi/database"
)

func RegisterAgent(h http.ResponseWriter, r *http.Request) {

	token := auth.ParseAPIKey(r.Header.Get("Api-key"))
	claims, ok := token.Claims.(*auth.MyClaims)

	if !ok || !token.Valid {
		http.Error(h, "not authorized", http.StatusUnauthorized)
		return
	}

	agentid, err := database.CreateAgent(claims.User)
	if err != nil {
		http.Error(h, "not authorized", http.StatusUnauthorized)
	}

	h.Write(GenerateResponse(map[string]string{"id": agentid}, nil))
}

func UpdateSystemInfo(h http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		http.Error(h, "couldnot read body¯", http.StatusUnauthorized)
	}
	defer r.Body.Close()

	var agentinfo map[string]string
	err = json.Unmarshal(body, &agentinfo)
	if err != nil {
		http.Error(h, "could not unmarsel sys info", http.StatusUnauthorized)
	}

	database.UpdateSystemInfo(key, agentinfo)

	if err != nil {
		http.Error(h, "db error", http.StatusUnauthorized)
	}

}
