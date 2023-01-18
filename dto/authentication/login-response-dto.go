package authentication

type LoginResponseDTO struct {
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	AccessToken string `json:"access_token"`
}
