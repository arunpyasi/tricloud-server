package restapi

import (
	"encoding/json"
	"net/http"

	"github.com/indrenicloud/tricloud-server/app/auth"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

func generateResp(w http.ResponseWriter, data interface{}, err error) {

	if data != nil {
		m := make(map[string]interface{})
		m["status"] = "ok"
		m["data"] = data
		m["err"] = err // this could be used to send deprecated api Warning
		response, err2 := enJson(m)
		if err2 != nil {
			errorResp(w, err2)
			return
		}
		w.Write(response)
		return
	}
	errorResp(w, err)
}

func errorResp(w http.ResponseWriter, err error) {
	logg.Warn(err)
	w.Write([]byte(`{"msg":"` + err.Error() + `","status":"failed"}`))
}

func enJson(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func deJson(raw []byte, out interface{}) error {
	return json.Unmarshal(raw, out)
}

func parseUser(r *http.Request) (string, bool) {
	c := r.Context().Value(ContextUser)
	claims, ok := c.(*auth.MyClaims)
	if !ok {
		return "", false
	}
	return claims.User, claims.Super
}

func isAuthorized(user string, r *http.Request) bool {
	sessuser, super := parseUser(r)
	if super {
		return super
	}
	if sessuser == "" {
		return false
	}
	return user == sessuser
}

func isSuperUser(r *http.Request) bool {
	_, super := parseUser(r)
	return super
}
