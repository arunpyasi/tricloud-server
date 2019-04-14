package restapi

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/indrenicloud/tricloud-server/restapi/auth"
	"github.com/indrenicloud/tricloud-server/restapi/database"
)

func RegisterAgent(h http.ResponseWriter, r *http.Request) {

	token := auth.ParseAPIKey(r.Header.Get("Api-key"))
	claims, ok := token.Claims.(auth.MyClaims)

	if !ok || !token.Valid {
		http.Error(h, "not authorized", http.StatusUnauthorized)
		return
	}

	agentid, err := database.CreateAgent(claims.User)
	if err != nil {
		http.Error(h, "not authorized", http.StatusUnauthorized)
	}

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		http.Error(h, "some error ¯\\_(ツ)_/¯", http.StatusUnauthorized)
	}
	defer r.Body.Close()

	var userinfo map[string]string
	err = json.Unmarshal(body, &userinfo)
	if err != nil {
		http.Error(h, "some error ¯\\_(ツ)_/¯", http.StatusUnauthorized)
	}

	database.UpdateSystemInfo(claims.User, userinfo)

	if err != nil {
		http.Error(h, "some error ¯\\_(ツ)_/¯", http.StatusUnauthorized)
	}
	h.Write(GenerateResponse(map[string]string{"id": agentid}, nil))
}
