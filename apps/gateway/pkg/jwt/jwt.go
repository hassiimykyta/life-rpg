package jwt

import (
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
)

const (
	ACCESS  = "access"
	REFRESH = "refresh"
)

type Claims struct {
	UserID    string `json:"sub"`
	TokenType string `json:"typ"`
	jwtlib.RegisteredClaims
}

type Manager struct {
	secret     []byte
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewManager(secret, issuer string, accessTTL, refreshTTL time.Duration) *Manager {
	return &Manager{
		secret:     []byte(secret),
		issuer:     issuer,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

func (m *Manager) IssuePair(userID string) (access string, accessExp int64, refresh string, refreshExp int64, err error) {
	now := time.Now()

	aClaims := Claims{
		UserID:    userID,
		TokenType: ACCESS,
		RegisteredClaims: jwtlib.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			ExpiresAt: jwtlib.NewNumericDate(now.Add(m.accessTTL)),
			IssuedAt:  jwtlib.NewNumericDate(now),
		},
	}
	rClaims := Claims{
		UserID:    userID,
		TokenType: REFRESH,
		RegisteredClaims: jwtlib.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   userID,
			ExpiresAt: jwtlib.NewNumericDate(now.Add(m.refreshTTL)),
			IssuedAt:  jwtlib.NewNumericDate(now),
		},
	}

	access, err = jwtlib.
		NewWithClaims(jwtlib.SigningMethodHS256, aClaims).
		SignedString(m.secret)
	if err != nil {
		return "", 0, "", 0, err
	}

	refresh, err = jwtlib.
		NewWithClaims(jwtlib.SigningMethodHS256, rClaims).
		SignedString(m.secret)
	if err != nil {
		return "", 0, "", 0, err
	}

	return access, aClaims.ExpiresAt.Unix(), refresh, rClaims.ExpiresAt.Unix(), nil
}

func (m *Manager) VerifyAccess(token string) (Claims, error) {
	return m.verify(token, "access")
}

func (m *Manager) VerifyRefresh(token string) (Claims, error) {
	return m.verify(token, "refresh")
}

func (m *Manager) Refresh(refreshToken string) (newAccess string, newAccessExp int64, newRefresh string, newRefreshExp int64, err error) {
	c, err := m.VerifyRefresh(refreshToken)
	if err != nil {
		return "", 0, "", 0, err
	}
	return m.IssuePair(c.UserID)
}

func (m *Manager) verify(token string, wantType string) (Claims, error) {
	var c Claims
	t, err := jwtlib.ParseWithClaims(token, &c, func(t *jwtlib.Token) (interface{}, error) {
		return m.secret, nil
	})
	if err != nil {
		return Claims{}, err
	}
	if !t.Valid {
		return Claims{}, jwtlib.ErrTokenInvalidClaims
	}
	if c.TokenType != wantType {
		return Claims{}, jwtlib.ErrTokenInvalidClaims
	}
	if c.Issuer != m.issuer {
		return Claims{}, jwtlib.ErrTokenInvalidIssuer
	}
	return c, nil
}
