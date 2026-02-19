package requests

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogoutRequest struct {
	UserID       uint64 `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokensRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
