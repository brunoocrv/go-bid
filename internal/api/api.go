package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/brunoocrv/go-bid/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/websocket"
)

type Api struct {
	Router         *chi.Mux
	Sessions       *scs.SessionManager
	WSUpgrader     websocket.Upgrader
	AuctionLobby   services.AuctionLobby
	UserService    services.UserService
	ProductService services.ProductsService
	BidsService    services.BidsService
}
