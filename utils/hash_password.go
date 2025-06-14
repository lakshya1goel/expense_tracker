package utils

import (
	"crypto/sha256"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	sha := sha256.Sum256([]byte(password))
	bytes, err := bcrypt.GenerateFromPassword(sha[:], bcrypt.DefaultCost)
	return string(bytes), err
}

func VerifyPasswordHash(password, hash string) bool {
	sha := sha256.Sum256([]byte(password))
	err := bcrypt.CompareHashAndPassword([]byte(hash), sha[:])
	return err == nil
}
