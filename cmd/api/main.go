package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/pgxstore"
	"github.com/alexedwards/scs/v2"
	"github.com/brunoocrv/go-bid/internal/api"
	"github.com/brunoocrv/go-bid/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func init() {
	gob.Register(uuid.UUID{})
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		panic(err)
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("GOBID_DB_USER"),
		os.Getenv("GOBID_DB_PASSWORD"),
		os.Getenv("GOBID_DB_HOST"),
		os.Getenv("GOBID_DB_PORT"),
		os.Getenv("GOBID_DB_NAME"),
	))
	if err != nil {
		fmt.Println("error connecting to database:", err)
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		fmt.Println("error pinging database:", err)
		panic(err)
	}

	session := scs.New()
	session.Store = pgxstore.New(pool)
	session.Lifetime = 24 * time.Hour
	session.Cookie.HttpOnly = true
	session.Cookie.SameSite = http.SameSiteLaxMode

	api := api.Api{
		Router:   chi.NewMux(),
		Sessions: session,
		WSUpgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		AuctionLobby: services.AuctionLobby{
			Rooms: make(map[uuid.UUID]*services.AuctionRoom),
		},
		UserService:    services.NewUserService(pool),
		ProductService: services.NewProductsService(pool),
		BidsService:    services.NewBidsService(pool),
	}

	api.BindRoutes()

	fmt.Println("god-bid service started at port :8080")
	if err := http.ListenAndServe("localhost:8080", api.Router); err != nil {
		fmt.Println("error starting server:", err)
		panic(err)
	}

}
