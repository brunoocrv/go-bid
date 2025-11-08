package api

import (
	"fmt"
	"net/http"

	"github.com/brunoocrv/go-bid/internal/jsonutils"
	"github.com/gorilla/csrf"
)

func (api *Api) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		exists := api.Sessions.Exists(r.Context(), "authenticated_user_id")
		fmt.Printf("DEBUG: Auth middleware - session exists: %v\n", exists)
		if !api.Sessions.Exists(r.Context(), "authenticated_user_id") {
			jsonutils.EncodeJSON(w, r, http.StatusUnauthorized, map[string]string{
				"message": "must to be signed in",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (api *Api) HandleGetCSREFToken(w http.ResponseWriter, r *http.Request) {
	token := csrf.Token(r)

	jsonutils.EncodeJSON(w, r, http.StatusOK, map[string]string{
		"csrf_token": token,
	})
}
