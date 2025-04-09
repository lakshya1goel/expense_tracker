package utils

import (
	"math/rand"
	"time"
)

func GenerateOtp(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	otp := make([]byte, length)
	for i := range otp {
		otp[i] = byte(rand.Intn(10) + 48)
	}
	return string(otp)
}
