package api

import (
	"errors"
	"net/http"

	"github.com/brunoocrv/go-bid/internal/jsonutils"
	"github.com/brunoocrv/go-bid/internal/services"
	"github.com/brunoocrv/go-bid/internal/usecases/users"
)

func (api *Api) handleSignUp(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[users.CreateUserReq](r)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.CreateUser(r.Context(), data.UserName, data.Email, data.Password, data.Bio)
	if err != nil {
		if errors.Is(err, services.ErrDuplicatedEmailOrUsername) {
			jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, map[string]string{
				"error": "invalid email or username",
			})
			return
		}
	}

	jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]any{
		"user_id": id,
	})
	return
}

func (api *Api) handleSignIn(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[users.SignInUserReq](r)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := api.UserService.SignInUser(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			jsonutils.EncodeJSON(w, r, http.StatusUnauthorized, map[string]string{
				"error": "invalid credentials",
			})
			return
		}

		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	if err := api.Sessions.RenewToken(r.Context()); err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	api.Sessions.Put(r.Context(), "authenticated_user_id", id)

	jsonutils.EncodeJSON(w, r, http.StatusOK, map[string]string{
		"message": "user signed in successfully",
	})
}

func (api *Api) handleSignOut(w http.ResponseWriter, r *http.Request) {

	if err := api.Sessions.RenewToken(r.Context()); err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
		return
	}

	api.Sessions.Remove(r.Context(), "authenticated_user_id")

	jsonutils.EncodeJSON(w, r, http.StatusOK, map[string]string{
		"message": "user signed out successfully",
	})
}
