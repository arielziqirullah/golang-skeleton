package authentication

type RefreshTokenRequestDTO struct {
	TokenRefresh string `json:"token_refresh" form:"token_refresh" binding:"required"`
}
