package restapi

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/indrenicloud/tricloud-server/app/auth"
	"github.com/indrenicloud/tricloud-server/app/database"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

var (
	ErrorLoginIncorrect = errors.New("Incorrect login details")
	ErrorAPI            = errors.New("could not generate api")
)

func SignIn(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	var logininfo map[string]string
	err = deJson(body, &logininfo) //todo check fields
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logg.Info(logininfo)

	user, err := database.GetUser(logininfo["userid"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !auth.IsPasswordCorrect(logininfo["password"], user.Password) {
		http.Error(w, ErrorLoginIncorrect.Error(), http.StatusBadRequest)
		return
	}
	api := auth.NewAPIKey("session", user.ID, user.SuperUser)

	if api == "" {
		http.Error(w, ErrorAPI.Error(), http.StatusBadRequest)
		return
	}

	generateResp(w, map[string]string{"api": api}, nil)

}
