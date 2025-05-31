package dto

type RegisterDto struct {
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

type LoginDto struct {
	Email    string `json:"email"`
	Mobile   string `json:"mobile"`
	Password string `json:"password"`
}

type UserResponseDto struct {
	ID               uint   `json:"id"`
	Email            string `json:"email"`
	Mobile           string `json:"mobile"`
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	AccessTokenEx    int64  `json:"access_token_exp"`
	RefreshTokenEx   int64  `json:"refresh_token_exp"`
	IsEmailVerified  bool   `json:"is_email_verified"`
	IsMobileVerified bool   `json:"is_mobile_verified"`
}

type SendOtpDto struct {
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
}

type VerifyMailDto struct {
	Email    string `json:"email"`
	EmailOtp string `json:"otp"`
}

type VerifyMobileDto struct {
	Mobile    string `json:"mobile"`
	MobileOtp string `json:"otp"`
}

type RefreshTokenDto struct {
	RefreshToken string `json:"refresh_token"`
}

type AccessTokenResponseDto struct {
	ID            uint   `json:"id"`
	AccessToken   string `json:"access_token"`
	AccessTokenEx int64  `json:"access_token_exp"`
}
