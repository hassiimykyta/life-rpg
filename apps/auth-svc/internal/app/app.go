package app

import (
	"context"
	"log"
	"net"

	"github.com/hassiimykyta/life-rpg/apps/auth-svc/internal/auth"
	"github.com/hassiimykyta/life-rpg/apps/auth-svc/internal/models"
	"github.com/hassiimykyta/life-rpg/apps/auth-svc/internal/repo"
	"github.com/hassiimykyta/life-rpg/apps/auth-svc/internal/security/password"
	"github.com/hassiimykyta/life-rpg/pkg/config"
	"github.com/hassiimykyta/life-rpg/pkg/db"
	"github.com/hassiimykyta/life-rpg/pkg/kafka"
	"github.com/hassiimykyta/life-rpg/pkg/ulid"
	authv1 "github.com/hassiimykyta/life-rpg/services/auth/v1"
	"google.golang.org/grpc"
	"gorm.io/gorm/logger"
)

type App struct {
	cfg   *config.Config
	db    *db.Conn
	grpc  *grpc.Server
	lis   net.Listener
	kafka *kafka.ProducerFactory
}

func New() (*App, error) {
	cfg, err := config.Load(config.WithDB(), config.WithJWT())
	if err != nil {
		return nil, err
	}

	producer := kafka.NewProducerFactory(kafka.ProducerFactoryConfig{
		Brokers: []string{"kafka:9092"},
	})
	defer producer.Close()

	conn, err := db.Open(db.Options{
		DSN:           cfg.DB.DSN,
		MaxOpen:       cfg.DB.MaxOpen,
		MaxIdle:       cfg.DB.MaxIdle,
		MaxIdleTime:   cfg.DB.MaxIdleTime,
		LogLevel:      logger.Info,
		SingularTable: true,
	})
	if err != nil {
		return nil, err
	}

	if err := conn.Gorm.AutoMigrate(&models.Identity{}); err != nil {
		return nil, err
	}

	repository := repo.NewIdentityRepo(conn.Gorm)
	idgen := ulid.NewULIDGenerator()
	hasher := password.Bcrypt{}

	svc := auth.New(repository, hasher, idgen, producer)

	s := grpc.NewServer()
	authv1.RegisterAuthServiceServer(s, svc)

	lis, err := net.Listen("tcp", ":"+cfg.App.Port)
	if err != nil {
		return nil, err
	}

	return &App{
		cfg:   cfg,
		db:    conn,
		grpc:  s,
		lis:   lis,
		kafka: producer,
	}, nil

}

func (a *App) Start() error {
	log.Printf("auth-svc listening on %s (env=%s)", a.lis.Addr(), a.cfg.App.Env)
	return a.grpc.Serve(a.lis)
}

func (a *App) Stop(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		a.grpc.GracefulStop()
		close(stopped)
	}()
	select {
	case <-stopped:
	case <-ctx.Done():
		a.grpc.Stop()
	}

	if a.db != nil && a.db.SQL != nil {
		_ = a.db.SQL.Close()
	}
	return nil
}
