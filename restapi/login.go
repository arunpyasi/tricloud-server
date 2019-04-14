package restapi

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/restapi/auth"
	"github.com/indrenicloud/tricloud-server/restapi/database"
)

func RegisterAuthHandlers(r *mux.Router) {

	r.HandleFunc("/signin", SignIn).Methods("POST")
}

func SignIn(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	var logininfo map[string]string
	err = json.Unmarshal(body, logininfo) //todo check fields

	user, err := database.GetUser(logininfo["userid"])

	if err != nil {
		return
	}

	if !auth.IsPasswordCorrect(logininfo["password"], user.Password) {
		return
	}
	var apierr error
	api := auth.NewAPIKey("session", user.ID)

	if api == "" {
		apierr = errors.New("could not generate api")
	}

	resp := GenerateResponse(map[string]string{"api": api}, apierr)
	w.Write(resp)
}
