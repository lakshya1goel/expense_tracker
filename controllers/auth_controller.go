package controllers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/database"
	"github.com/lakshya1goel/expense_tracker/dto"
	"github.com/lakshya1goel/expense_tracker/models"
	"github.com/lakshya1goel/expense_tracker/services"
	"github.com/lakshya1goel/expense_tracker/utils"
	"gorm.io/gorm"
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

	var existingUser models.User

	result := database.Db.Where("email = ?", request.Email).First(&existingUser)
	if result.Error == nil  {
		if existingUser.IsVerified {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "User already exists"})
			return
		}

		existingUser.Password = hashedPassword
		if err := database.Db.Save(&existingUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to update user " + err.Error()})
			return
		}

		userResponse := dto.UserResponseDto{
			ID:         existingUser.ID,
			Email:      existingUser.Email,
			IsVerified: existingUser.IsVerified,
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "User updated, please verify your email", "data": userResponse})
		return
	}
	
	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Database error " + result.Error.Error()})
		return
	}

	newUser := models.User{
		Email:    request.Email,
		Password: hashedPassword,
	}

	if err := database.Db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to create user " + err.Error()})
		return
	}

	response := dto.UserResponseDto{
		ID:             newUser.ID,
		Email:          newUser.Email,
		IsVerified:     newUser.IsVerified,
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "message": "User created successfully, please verify your email", "data": response})
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
	if result.Error != nil || !user.IsVerified {
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

	user := models.User{
		Email: request.Email,
	}

	result := database.Db.Where("email = ?", request.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "User not found"})
		return
	}

	otp := utils.GenerateOtp(6)
	if err := services.SendMail(request.Email, "OTP for email verification", otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to send OTP " + err.Error()})
		return
	}

	otpModel := models.Otp{
		Id: user.ID,
		Email: request.Email,
		Otp:    otp,
		OtpExp:    time.Now().Add(time.Minute * 5).Unix(),
	}

	database.Db.Save(&otpModel)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "OTP sent successfully", "data": otp})
}

func VerifyOtp(c *gin.Context) {
	var request dto.VerifyOtpDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": err.Error()})
		return
	}

	otpModel := models.Otp{
		Email: request.Email,
	}

	result := database.Db.Where("email = ?", request.Email).First(&otpModel)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "User not found"})
		return
	}

	if otpModel.Otp != request.Otp || otpModel.OtpExp < time.Now().Unix() {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "error": "Invalid OTP"})
		return
	}

	var user models.User
    result = database.Db.Where("email = ?", request.Email).First(&user)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"success": false, "error": "User not found"})
        return
    }

	if err := database.Db.Model(&user).Update("is_verified", true).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to update user"})
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

	if err := database.Db.Delete(&otpModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to delete OTP " + err.Error()})
		return
	}

	response := dto.UserResponseDto{
		ID:             user.ID,
		Email:          user.Email,
		AccessToken:    accessToken,
		AccessTokenEx:  accessTokenExp,
		RefreshToken:   refreshToken,
		RefreshTokenEx: refreshTokenExp,
		IsVerified:    true,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Email verified successfully", "data": response})
}