package app

import (
	"context"
	"fmt"
	"log"

	"github.com/go-chi/chi/v5"
	"github.com/hassiimykyta/life-rpg/apps/gateway/internal/clients"
	"github.com/hassiimykyta/life-rpg/apps/gateway/internal/handlers"
	"github.com/hassiimykyta/life-rpg/apps/gateway/internal/router"
	"github.com/hassiimykyta/life-rpg/apps/gateway/pkg/jwt"
	"github.com/hassiimykyta/life-rpg/pkg/config"
	"github.com/hassiimykyta/life-rpg/pkg/helpers"
	"github.com/hassiimykyta/life-rpg/pkg/httpserver"
)

type App struct {
	cfg     *config.Config
	server  *httpserver.Server
	cleanup func() error
}

func initRouter(cfg *config.Config, cli *clients.Clients, jwtMgr *jwt.Manager) *chi.Mux {
	return router.New(
		router.Deps{
			Handlers: router.Handlers{
				AuthHandler: handlers.NewAuthHandler(cli.Auth, jwtMgr),
			},
		},
		router.Options{
			CORS: router.CORSOpts{
				AllowedOrigins:   cfg.CORS.AllowedOrigins,
				AllowedMethods:   cfg.CORS.AllowedMethods,
				AllowedHeaders:   cfg.CORS.AllowedHeaders,
				ExposedHeaders:   cfg.CORS.ExposedHeaders,
				AllowCredentials: cfg.CORS.AllowCredentials,
				MaxAge:           cfg.CORS.MaxAge,
			},
		},
	)
}

func New() (*App, error) {
	// Gateway нужен CORS и JWT (для мидлварей/хендлеров)
	cfg, err := config.Load(config.WithCORS(), config.WithJWT())
	if err != nil {
		return nil, err
	}

	authAddr := helpers.GetEnv("AUTH_SVC_ADDR", "auth-svc:8081")

	cli, cleanup, err := clients.NewClients(authAddr)
	if err != nil {
		return nil, err
	}

	jwtMgr := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL)
	r := initRouter(cfg, cli, jwtMgr)

	addr := fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port)
	srv, err := httpserver.New(httpserver.Options{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.App.ReadTimeout,
		WriteTimeout: cfg.App.WriteTimeout,
		IdleTimeout:  cfg.App.IdleTimeout,
	})
	if err != nil {
		_ = cleanup()
		return nil, err
	}

	return &App{
		cfg:     cfg,
		server:  srv,
		cleanup: cleanup,
	}, nil
}

func (a *App) Start() {
	log.Printf("▶ gateway listening on http://%s:%s (env=%s)", a.cfg.App.Host, a.cfg.App.Port, a.cfg.App.Env)
	a.server.Start()
}

func (a *App) Stop(ctx context.Context) error {
	if a.cleanup != nil {
		_ = a.cleanup()
	}
	return a.server.Shutdown(ctx)
}
