package api

import (
	"errors"
	"net/http"

	"github.com/brunoocrv/go-bid/internal/jsonutils"
	"github.com/brunoocrv/go-bid/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (api *Api) handleSubscribeUserToAuction(w http.ResponseWriter, r *http.Request) {
	rawProductId := chi.URLParam(r, "product_id")

	productId, err := uuid.Parse(rawProductId)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusBadRequest, map[string]string{"error": "invalid product id"})
		return
	}

	if _, err := api.ProductService.GetProductById(r.Context(), productId); err != nil {
		if errors.Is(err, services.ErrProductNotFound) {
			jsonutils.EncodeJSON(w, r, http.StatusNotFound, map[string]string{"error": "product not found"})
			return
		}

		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{"error": "internal server error"})
		return
	}

	userId, ok := api.Sessions.Get(r.Context(), "authenticated_user_id").(uuid.UUID)
	if !ok {
		jsonutils.EncodeJSON(w, r, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	api.AuctionLobby.Lock()
	room, ok := api.AuctionLobby.Rooms[productId]
	api.AuctionLobby.Unlock()

	if !ok {
		jsonutils.EncodeJSON(w, r, http.StatusBadRequest, map[string]string{"message": "room was ended"})
		return
	}

	conn, err := api.WSUpgrader.Upgrade(w, r, nil)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{"error": "could not upgrade connection"})
		return
	}

	client := services.NewClient(room, conn, userId)

	room.Register <- client

	go client.ReadEventLoop()
	go client.WriteEventLoop()

	for {
	}

}
