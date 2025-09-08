package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/hassiimykyta/life-rpg/apps/gateway/internal/dto"
	"github.com/hassiimykyta/life-rpg/apps/gateway/pkg/jwt"
	"github.com/hassiimykyta/life-rpg/apps/gateway/pkg/resp"
	authv1 "github.com/hassiimykyta/life-rpg/services/auth/v1"
)

type AuthHandler struct {
	Client authv1.AuthServiceClient
	Jwt    *jwt.Manager
}

func NewAuthHandler(client authv1.AuthServiceClient, jwt *jwt.Manager) *AuthHandler {
	return &AuthHandler{Client: client, Jwt: jwt}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.ERROR(w, r, "bad request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	out, err := h.Client.Register(ctx, &authv1.RegisterRequest{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		resp.ERROR(w, r, "error", http.StatusBadRequest)
		return
	}

	acc, accExp, ref, refExp, err := h.Jwt.IssuePair(out.UserId)
	if err != nil {
		resp.ERROR(w, r, "error creation token")
		return
	}

	resp.OK(w, r, map[string]any{"token": dto.Token{
		AccessToken:      acc,
		ExpiresAt:        accExp,
		RefreshToken:     ref,
		RefreshExpiresAt: refExp,
	}},
		"account created",
		http.StatusCreated,
	)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.ERROR(w, r, "bad request", http.StatusBadRequest)
		return
	}

	if req.Email == "" && req.Username == "" {
		resp.ERROR(w, r, "must provide email or username", http.StatusBadRequest)
		return

	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	loginReq := &authv1.LoginRequest{
		Password: req.Password,
	}
	if req.Email != "" {
		loginReq.Subject = &authv1.LoginRequest_Email{Email: req.Email}
	} else if req.Username != "" {
		loginReq.Subject = &authv1.LoginRequest_Username{Username: req.Username}
	}

	out, err := h.Client.Login(ctx, loginReq)
	if err != nil {
		resp.ERROR(w, r, "auth error", http.StatusBadRequest)
		return
	}

	acc, accExp, ref, refExp, err := h.Jwt.IssuePair(out.UserId)
	if err != nil {
		resp.ERROR(w, r, "error creation token")
		return
	}

	resp.OK(w, r, map[string]any{"token": dto.Token{
		AccessToken:      acc,
		ExpiresAt:        accExp,
		RefreshToken:     ref,
		RefreshExpiresAt: refExp,
	}},
		"ok",
	)
}

func (h *AuthHandler) Availability(w http.ResponseWriter, r *http.Request) {
	var req dto.AvailabilityRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.ERROR(w, r, "bad request", http.StatusBadRequest)
		return
	}

	if req.Email == "" && req.Username == "" {
		resp.ERROR(w, r, "email or username required", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()
	out, err := h.Client.CheckAvailability(ctx, &authv1.CheckAvailabilityRequest{
		Email: req.Email, Username: req.Username,
	})

	if err != nil {
		resp.ERROR(w, r, "bad gateway", http.StatusBadGateway)
		return
	}
	resp.OK(w, r, dto.AvailabilityResponse{
		EmailAvailable:    out.EmailAvailable,
		UsernameAvailable: out.UsernameAvailable,
	}, "ok")
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {

	var req dto.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.ERROR(w, r, "bad request", http.StatusBadRequest)
		return

	}
	if strings.TrimSpace(req.RefreshToken) == "" {
		resp.ERROR(w, r, "refresh token required", http.StatusBadRequest)
		return
	}

	acc, accExp, ref, refExp, err := h.Jwt.Refresh(req.RefreshToken)
	if err != nil {
		resp.ERROR(w, r, "auth error", http.StatusBadRequest)
		return
	}

	resp.OK(w, r, map[string]any{"token": dto.Token{
		AccessToken:      acc,
		ExpiresAt:        accExp,
		RefreshToken:     ref,
		RefreshExpiresAt: refExp,
	}},
		"ok",
	)
}
