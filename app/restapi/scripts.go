package restapi

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/indrenicloud/tricloud-server/app/database"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

func GetScripts(w http.ResponseWriter, r *http.Request) {

	scripts, err := database.GetAllScripts()
	if err != nil {
		errorResp(w, err)
		return
	}
	generateResp(w, scripts, err)

}

func CreateScript(w http.ResponseWriter, r *http.Request) {

	user, _ := parseUser(r)

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		errorResp(w, err)
		return
	}
	defer r.Body.Close()
	var scriptinfo map[string]string
	err = deJson(body, &scriptinfo)
	if err != nil {
		errorResp(w, err)
		return
	}
	script, err := database.NewScript(user, scriptinfo, false)
	if err != nil {
		errorResp(w, err)
		return
	}
	err = database.CreateScript(script)
	if err == nil {
		mScript.AddScript(script)
	}

	updatedscript, err := database.GetAllScripts()
	generateResp(w, updatedscript, err)

}

func DeleteScript(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	name := vars["name"]
	logg.Warn("Deleting script url")
	err := database.DeleteScript(name)
	if err != nil {
		errorResp(w, err)
	}

	mScript.RemoveScript(name)

	updatedscripts, err := database.GetAllScripts()

	if err != nil {
		errorResp(w, err)
		return
	}

	generateResp(w, updatedscripts, err)

}

func GetScript(w http.ResponseWriter, r *http.Request) {
	// only if user owns agent or superuser
	vars := mux.Vars(r)
	name := vars["name"]
	script, err := database.GetScript(name)

	if err != nil {
		errorResp(w, err)
		return
	}

	generateResp(w, script, err)
}
func RunScript(w http.ResponseWriter, r *http.Request) {
	logg.Debug("RUNNNNNNNNN")
	vars := mux.Vars(r)
	name := vars["name"]
	script, err := database.GetScript(name)

	if err != nil {
		errorResp(w, err)
		return
	}
	mScript.RunScript(script)
	generateResp(w, map[string]string{"ok": "ok"}, nil)
}
