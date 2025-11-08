package api

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/brunoocrv/go-bid/internal/jsonutils"
	"github.com/brunoocrv/go-bid/internal/services"
	"github.com/brunoocrv/go-bid/internal/usecases/products"
	"github.com/google/uuid"
)

func (api *Api) handleCreateProduct(w http.ResponseWriter, r *http.Request) {
	data, problems, err := jsonutils.DecodeValidJSON[products.CreateProductReq](r)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	userId, ok := api.Sessions.Get(r.Context(), "authenticated_user_id").(uuid.UUID)
	if !ok {
		jsonutils.EncodeJSON(w, r, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}

	id, err := api.ProductService.CreateProduct(r.Context(), userId, data.Name, data.Description, data.BasePrice, data.AuctionEnd)
	if err != nil {
		jsonutils.EncodeJSON(w, r, http.StatusInternalServerError, map[string]string{"error": "failed to create product auction"})
		return
	}

	ctx, _ := context.WithDeadline(context.Background(), data.AuctionEnd)
	// defer cancel()

	auctionRoom := services.NewAuctionRoom(ctx, id, api.BidsService)

	go auctionRoom.Run()

	api.AuctionLobby.Lock()
	api.AuctionLobby.Rooms[id] = auctionRoom
	api.AuctionLobby.Unlock()

	jsonutils.EncodeJSON(w, r, http.StatusCreated, map[string]any{"id": id, "message": "auction room was created"})
	slog.Info("auction has begun for product", "name", data.Name, "id", id)
}
