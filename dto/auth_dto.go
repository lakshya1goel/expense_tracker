package dto

type RegisterDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponseDto struct {
	ID            uint   `json:"id"`
	Email         string `json:"email"`
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	AccessTokenEx int64  `json:"access_token_exp"`
	RefreshTokenEx int64  `json:"refresh_token_exp"`
}

type SendOtpDto struct {
	Email string `json:"email"`
}

type VerifyOtpDto struct {
	Email string `json:"email"`
	Otp   string `json:"otp"`
}
