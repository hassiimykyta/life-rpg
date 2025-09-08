package auth

import (
	"context"
	"strings"

	"github.com/hassiimykyta/life-rpg/apps/auth-svc/internal/models"
	"github.com/hassiimykyta/life-rpg/apps/auth-svc/internal/repo"
	"github.com/hassiimykyta/life-rpg/apps/auth-svc/internal/security/password"
	"github.com/hassiimykyta/life-rpg/pkg/ulid"
	authv1 "github.com/hassiimykyta/life-rpg/services/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	authv1.UnimplementedAuthServiceServer
	repo *repo.IdentityRepo
	hash password.Hasher
	ids  *ulid.ULIDGenerator
}

func New(r *repo.IdentityRepo, h password.Hasher, g *ulid.ULIDGenerator) *Service {
	return &Service{repo: r, hash: h, ids: g}
}

func normIdentifier(ide string) string {
	return strings.TrimSpace(strings.ToLower(ide))
}

func (s *Service) Register(ctx context.Context, in *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	email := normIdentifier(in.GetEmail())
	username := normIdentifier(in.GetUsername())
	password := in.GetPassword()

	if email == "" || username == "" || len(password) < 6 {
		return nil, status.Error(codes.InvalidArgument, "invalid input")
	}

	id, err := s.ids.New()
	if err != nil {
		return nil, status.Error(codes.Internal, "id generation failed")
	}

	h, err := s.hash.Hash(password)
	if err != nil {
		return nil, status.Error(codes.Internal, "hash generation failed")
	}

	err = s.repo.Create(ctx, models.Identity{
		UserId:       id,
		Email:        email,
		Username:     username,
		PasswordHash: h,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "create identity failed")
	}

	return &authv1.RegisterResponse{
		UserId: id,
	}, nil

}

func (s *Service) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	password := in.GetPassword()

	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "password required")
	}

	var (
		ide models.Identity
		err error
	)

	switch sub := in.Subject.(type) {
	case *authv1.LoginRequest_Email:
		email := normIdentifier(sub.Email)
		if email == "" {
			return nil, status.Error(codes.InvalidArgument, "email required")
		}
		ide, err = s.repo.FindByEmail(ctx, email)
	case *authv1.LoginRequest_Username:
		username := normIdentifier(sub.Username)
		if username == "" {
			return nil, status.Error(codes.InvalidArgument, "username required")
		}
		ide, err = s.repo.FindByUsername(ctx, username)

	case *authv1.LoginRequest_UserId:
		userId := sub.UserId
		if userId == "" {
			return nil, status.Error(codes.InvalidArgument, "user id required")
		}
		ide, err = s.repo.FindByUserID(ctx, userId)
	default:
		return nil, status.Error(codes.InvalidArgument, "oneof subject required")
	}

	if err != nil {
		return nil, status.Error(codes.Internal, "lookup failed")
	}

	if !s.hash.Compare(ide.PasswordHash, password) {
		return nil, status.Error(codes.PermissionDenied, "invalid credentials")
	}

	return &authv1.LoginResponse{
		UserId: ide.UserId,
	}, nil
}

func (s *Service) Resolve(ctx context.Context, in *authv1.ResolveRequest) (*authv1.ResolveResponse, error) {
	var (
		ide models.Identity
		err error
	)

	switch sub := in.Key.(type) {
	case *authv1.ResolveRequest_Email:
		email := normIdentifier(sub.Email)
		if email == "" {
			return nil, status.Error(codes.InvalidArgument, "email required")
		}
		ide, err = s.repo.FindByEmail(ctx, email)
	case *authv1.ResolveRequest_Username:
		username := normIdentifier(sub.Username)
		if username == "" {
			return nil, status.Error(codes.InvalidArgument, "username required")
		}
		ide, err = s.repo.FindByUsername(ctx, username)

	case *authv1.ResolveRequest_UserId:
		userId := sub.UserId
		if userId == "" {
			return nil, status.Error(codes.InvalidArgument, "user id required")
		}
		ide, err = s.repo.FindByUserID(ctx, userId)
	default:
		return nil, status.Error(codes.InvalidArgument, "oneof subject required")
	}

	if err != nil {
		return nil, status.Error(codes.Internal, "lookup failed")
	}

	return &authv1.ResolveResponse{
		UserId:   ide.UserId,
		Email:    ide.Email,
		Username: ide.Username,
	}, nil
}
func (s *Service) CheckAvailability(ctx context.Context, in *authv1.CheckAvailabilityRequest) (*authv1.CheckAvailabilityResponse, error) {
	email := normIdentifier(in.GetEmail())
	username := normIdentifier(in.GetUsername())

	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "email required")
	}
	if username == "" {
		return nil, status.Error(codes.InvalidArgument, "username required")
	}

	EmailAvailable := true
	if _, err := s.repo.FindByEmail(ctx, email); err == nil {
		EmailAvailable = false
	}

	UsernameAvailable := true
	if _, err := s.repo.FindByUsername(ctx, username); err == nil {
		UsernameAvailable = false
	}

	return &authv1.CheckAvailabilityResponse{
		EmailAvailable:    EmailAvailable,
		UsernameAvailable: UsernameAvailable,
	}, nil
}
