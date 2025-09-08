package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/hassiimykyta/life-rpg/apps/gateway/internal/middleware"
)

func MountAPI(r *chi.Mux, d Deps) {
	r.Route("/api", func(api chi.Router) {
		r.Use(middleware.JSONMiddleware)

		api.Route("/v1", func(v1 chi.Router) {
			v1.Route("/auth", func(auth chi.Router) {
				auth.Post("/register", d.Handlers.AuthHandler.Register)
				auth.Post("/login", d.Handlers.AuthHandler.Login)
				auth.Post("/refresh", d.Handlers.AuthHandler.Refresh)
				auth.Post("/availability", d.Handlers.AuthHandler.Availability)
			})

		})
	})
}
