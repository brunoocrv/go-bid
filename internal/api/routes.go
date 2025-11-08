package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (api *Api) BindRoutes() {
	api.Router.Use(middleware.RequestID, middleware.Logger, middleware.Recoverer, api.Sessions.LoadAndSave)

	// csrfMiddleware := csrf.Protect(
	// 	[]byte(os.Getenv("GOBID_CSRF_TOKEN")),
	// 	csrf.Secure(false), // dev only
	// )

	// api.Router.Use(csrfMiddleware)

	api.Router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Get("/csrf-token", api.HandleGetCSREFToken)

			r.Route("/users", func(r chi.Router) {
				r.Post("/sign-up", api.handleSignUp)
				r.Post("/sign-in", api.handleSignIn)
				r.With(api.AuthMiddleware).Post("/sign-out", api.handleSignOut)
			})

			r.With(api.AuthMiddleware).Route("/products", func(r chi.Router) {
				r.Post("/create", api.handleCreateProduct)
				r.Get("/ws/subscribe/{product_id}", api.handleSubscribeUserToAuction)
			})
		})
	})
}
