package models

type Otp struct {
	Id        uint   `json:"id" gorm:"primaryKey"`
	Email     string `json:"email" gorm:"unique"`
	EmailOtp  string `json:"otp"`
	Mobile    string `json:"mobile_no"`
	MobileOtp string `json:"mobile_otp"`
	OtpExp    int64  `json:"otp_exp"`
}