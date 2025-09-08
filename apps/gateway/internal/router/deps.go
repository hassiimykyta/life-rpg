package router

import (
	"github.com/hassiimykyta/life-rpg/apps/gateway/internal/handlers"
	"github.com/hassiimykyta/life-rpg/pkg/redisx"
)

type Handlers struct {
	AuthHandler *handlers.AuthHandler
}

type Deps struct {
	Handlers Handlers
	RDB      *redisx.Cache
}
