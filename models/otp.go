package models

type Otp struct {
	Id        uint   `json:"id" gorm:"primaryKey"`
	Email     string `json:"email" gorm:"unique"`
	Otp       string `json:"otp"`
	OtpExp    int64  `json:"otp_exp"`
}