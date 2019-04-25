package restapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/indrenicloud/tricloud-server/app/auth"
	"github.com/indrenicloud/tricloud-server/app/logg"
)

// when using refelection to find type  using custom type avoids collision in contex.value
type contextkey int

const ContextUser contextkey = iota

var (
	ErrorNotAuthorized = errors.New("Not authorized")
)

//
func MiddlewareJson(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		next.ServeHTTP(w, r)
	})
}

//COrs middleware
func MiddlewareCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
			return
		} else {
			h.ServeHTTP(w, r)
		}
	})
}

// MiddlewareSession checks the session for request and tags username to request context
func MiddlewareSession(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawapi := r.Header.Get("Api-key")
		if rawapi == "" {
			logg.Warn("Token Not set")
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}
		token := auth.ParseAPIKey(rawapi)
		if !token.Valid {
			logg.Warn("Token invalid")
			http.Error(w, "Not Authorized", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(*auth.MyClaims)
		logg.Info(claims)

		if !ok {
			logg.Warn("Cannot cast clam type")
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}

		ctx := context.WithValue(r.Context(), ContextUser, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logg.Info(r.URL.Path)
		h.ServeHTTP(w, r)
	})
}
