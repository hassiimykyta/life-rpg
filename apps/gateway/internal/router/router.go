package router

import (
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type CORSOpts struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           int
}

type Options struct {
	CORS CORSOpts
}

func New(d Deps, opts Options) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   opts.CORS.AllowedOrigins,
		AllowedMethods:   opts.CORS.AllowedMethods,
		AllowedHeaders:   opts.CORS.AllowedHeaders,
		ExposedHeaders:   opts.CORS.ExposedHeaders,
		AllowCredentials: opts.CORS.AllowCredentials,
		MaxAge:           opts.CORS.MaxAge,
	}))

	MountAPI(r, d)

	return r
}
