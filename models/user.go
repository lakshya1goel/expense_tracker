package models

type User struct {
	ID            uint   `json:"id"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Otp           string `json:"otp"`
	OtpExp        int64  `json:"otp_exp"`
	IsVerified    bool   `json:"is_verified"`
}