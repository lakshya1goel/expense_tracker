package utils

import (
	"math/rand"
	"time"
)

func GenerateOtp(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = byte(r.Intn(10) + 48)
	}
	return string(otp)
}
