package controllers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/database"
	"github.com/lakshya1goel/expense_tracker/dto"
	"github.com/lakshya1goel/expense_tracker/models"
	"github.com/lakshya1goel/expense_tracker/services"
	"github.com/lakshya1goel/expense_tracker/utils"
)

func Register(c *gin.Context) {
	var request dto.RegisterDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to hash password " + err.Error()})
		return
	}

	user := models.User{
		Email:    request.Email,
		Password: hashedPassword,
	}

	result := database.Db.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to create user " + result.Error.Error()})
		return
	}

	accessTokenExp := time.Now().Add(time.Hour * 24).Unix()
	accessToken, err := utils.GenerateToken(user.ID, accessTokenExp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to generate access token " + err.Error()})
		return
	}

	refreshTokenExp := time.Now().Add(time.Hour * 24 * 30).Unix()
	refreshToken, err := utils.GenerateToken(user.ID, refreshTokenExp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to generate refresh token " + err.Error()})
		return
	}

	response := dto.UserResponseDto{
		ID:             user.ID,
		Email:          user.Email,
		AccessToken:    accessToken,
		AccessTokenEx:  accessTokenExp,
		RefreshToken:   refreshToken,
		RefreshTokenEx: refreshTokenExp,
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "User created successfully", "data": response})
}

func Login(c *gin.Context) {
	var request dto.LoginDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	user := models.User{
		Email: request.Email,
	}

	result := database.Db.Where("email = ?", request.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "User not found"})
		return
	}

	if !utils.VerifyPasswordHash(user.Password, request.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Invalid password"})
		return
	}

	accessTokenExp := time.Now().Add(time.Hour * 24).Unix()
	accessToken, err := utils.GenerateToken(user.ID, accessTokenExp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to generate access token " + err.Error()})
		return
	}

	refreshTokenExp := time.Now().Add(time.Hour * 24 * 30).Unix()
	refreshToken, err := utils.GenerateToken(user.ID, refreshTokenExp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to generate refresh token " + err.Error()})
		return
	}

	response := dto.UserResponseDto{
		ID:             user.ID,
		Email:          user.Email,
		AccessToken:    accessToken,
		AccessTokenEx:  accessTokenExp,
		RefreshToken:   refreshToken,
		RefreshTokenEx: refreshTokenExp,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User logged in successfully", "data": response})
}

func SendOtp(c *gin.Context) {
	var request dto.SendOtpDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	otp := utils.GenerateOtp(6)
	services.SendMail(request.Email, "OTP for email verification", otp)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "OTP sent successfully", "data": otp})
}

func VerifyOtp(c *gin.Context) {
	var request dto.VerifyOtpDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	user := models.User{
		Email: request.Email,
	}

	result := database.Db.Where("email = ?", request.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "User not found"})
		return
	}

	if user.Otp != request.Otp {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Invalid OTP"})
		return
	}

	user.IsVerified = true
	database.Db.Save(&user)

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Email verified successfully"})
}
