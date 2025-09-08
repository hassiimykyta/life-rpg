package clients

import (
	"context"
	"fmt"
	"log"
	"time"

	authv1 "github.com/hassiimykyta/life-rpg/services/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	Auth authv1.AuthServiceClient
}

func NewClients(authAddr string) (*Clients, func() error, error) {
	log.Printf("🔌 [gRPC] dialing %s ...", authAddr)

	conn, err := grpc.NewClient(
		authAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("❌ [gRPC] create client failed: %v", err)
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn.Connect()

	for {
		s := conn.GetState()
		if s == connectivity.Ready {
			break
		}
		log.Printf("⏳ [gRPC] waiting for state: %s", s)
		if !conn.WaitForStateChange(ctx, s) {
			_ = conn.Close()
			return nil, nil, fmt.Errorf("gRPC connect timeout to %s (last state: %s)", authAddr, s)
		}
	}

	log.Printf("✅ [gRPC] connected to %s", authAddr)

	cleanup := func() error {
		log.Printf("🔌 [gRPC] closing %s", authAddr)
		return conn.Close()
	}

	return &Clients{
		Auth: authv1.NewAuthServiceClient(conn),
	}, cleanup, nil
}
