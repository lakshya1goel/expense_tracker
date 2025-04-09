package models

type User struct {
	ID            uint   `json:"id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	AccessTokenEx int64  `json:"access_token_exp"`
	RefreshTokenEx int64  `json:"refresh_token_exp"`
}