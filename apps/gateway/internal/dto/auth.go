package dto

type RegisterRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}

type AvailabilityRequest struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
}

type AvailabilityResponse struct {
	EmailAvailable    bool `json:"email_available"`
	UsernameAvailable bool `json:"username_available"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}
