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
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if request.Email == "" || request.Mobile == "" || request.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "All fields are required"})
		return
	}

	hashedPassword, err := utils.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to hash password " + err.Error()})
		return
	}

	var existingUser models.User
	result := database.Db.Where("email = ? OR mobile = ?", request.Email, request.Mobile).First(&existingUser)
	if result.Error == nil {
		if existingUser.IsEmailVerified && existingUser.IsMobileVerified {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Email or mobile already exists"})
			return
		}

		existingUser.Password = hashedPassword
		existingUser.Email = request.Email
		existingUser.Mobile = request.Mobile
		if err := database.Db.Save(&existingUser).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update user " + err.Error()})
			return
		}

		emailOtp := utils.GenerateOtp(6)
		mobileOtp := utils.GenerateOtp(6)
		if err := services.SendMail(request.Email, "OTP for email verification", emailOtp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to send OTP " + err.Error()})
			return
		}

		if err := services.SendSms("+91"+request.Mobile, mobileOtp); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to send OTP " + err.Error()})
			return
		}

		otpModel := models.Otp{
			Email:     request.Email,
			Mobile:    request.Mobile,
			EmailOtp:  emailOtp,
			MobileOtp: mobileOtp,
			OtpExp:    time.Now().Add(time.Minute * 5).Unix(),
		}

		database.Db.Save(&otpModel)
		c.JSON(http.StatusOK, gin.H{"success": true, "message": "OTP sent successfully"})
		return
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Database error " + result.Error.Error()})
		return
	}

	newUser := models.User{
		Email:    request.Email,
		Mobile:   request.Mobile,
		Password: hashedPassword,
	}

	if err := database.Db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create user " + err.Error()})
		return
	}

	emailOtp := utils.GenerateOtp(6)
	mobileOtp := utils.GenerateOtp(6)
	if err := services.SendMail(request.Email, "OTP for email verification", emailOtp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to send OTP " + err.Error()})
		return
	}

	if err := services.SendSms("+91"+request.Mobile, mobileOtp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to send OTP " + err.Error()})
		return
	}

	otpModel := models.Otp{
		Email:     request.Email,
		Mobile:    request.Mobile,
		EmailOtp:  emailOtp,
		MobileOtp: mobileOtp,
		OtpExp:    time.Now().Add(time.Minute * 5).Unix(),
	}

	database.Db.Save(&otpModel)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "OTP sent successfully"})
}

func Login(c *gin.Context) {
	var request dto.LoginDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if request.Email == "" && request.Mobile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Email or mobile is required"})
		return
	}

	var user models.User
	var result *gorm.DB
	if request.Email != "" {
		result = database.Db.Where("email = ?", request.Email).First(&user)
	} else if request.Mobile != "" {
		result = database.Db.Where("mobile = ?", request.Mobile).First(&user)
	}
	if result.Error != nil {
		if result.Error != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Database error " + result.Error.Error()})
			return
		} else if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
			return
		}
	} else if !user.IsEmailVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Email not verified"})
		return
	} else if !user.IsMobileVerified {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Mobile not verified"})
		return
	}

	if !utils.VerifyPasswordHash(request.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid password"})
		return
	}

	accessTokenExp := time.Now().Add(time.Hour * 24).Unix()
	accessToken, err := utils.GenerateToken(user.ID, accessTokenExp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate access token " + err.Error()})
		return
	}

	refreshTokenExp := time.Now().Add(time.Hour * 24 * 30).Unix()
	refreshToken, err := utils.GenerateToken(user.ID, refreshTokenExp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate refresh token " + err.Error()})
		return
	}

	response := dto.UserResponseDto{
		ID:               user.ID,
		Email:            user.Email,
		Mobile:           user.Mobile,
		IsEmailVerified:  user.IsEmailVerified,
		IsMobileVerified: user.IsMobileVerified,
		AccessToken:      accessToken,
		AccessTokenEx:    accessTokenExp,
		RefreshToken:     refreshToken,
		RefreshTokenEx:   refreshTokenExp,
	}

	status, err := utils.CreatePrivateGroup(user.ID)
	if err != nil {
		c.JSON(status, gin.H{"success": false, "message": "Failed to create private group " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User logged in successfully", "data": response})
}

func SendOtp(c *gin.Context) {
	var request dto.SendOtpDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if request.Email == "" && request.Mobile == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Email or mobile is required"})
		return
	}

	var user models.User
	if request.Email != "" {
		result := database.Db.Where("email = ?", request.Email).First(&user)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Email not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Database error: " + result.Error.Error()})
			return
		}
	}

	if request.Mobile != "" {
		result := database.Db.Where("mobile = ? and id = ?", request.Mobile, user.ID).First(&user)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Mobile number not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Database error: " + result.Error.Error()})
			return
		}
	}

	emailOtp := utils.GenerateOtp(6)
	mobileOtp := utils.GenerateOtp(6)
	if err := services.SendMail(request.Email, "OTP for email verification", emailOtp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to send OTP " + err.Error()})
		return
	}

	if err := services.SendSms("+91"+request.Mobile, mobileOtp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to send OTP " + err.Error()})
		return
	}

	otpModel := models.Otp{
		Email:     request.Email,
		Mobile:    request.Mobile,
		EmailOtp:  emailOtp,
		MobileOtp: mobileOtp,
		OtpExp:    time.Now().Add(time.Minute * 5).Unix(),
	}

	database.Db.Save(&otpModel)
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "OTP sent successfully"})
}

func VerifyMail(c *gin.Context) {
	var request dto.VerifyMailDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	var otpModel models.Otp
	result := database.Db.Where("email = ?", request.Email).First(&otpModel)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}

	if otpModel.EmailOtp != request.EmailOtp || otpModel.OtpExp < time.Now().Unix() {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid OTP"})
		return
	}

	var user models.User
	result = database.Db.Where("email = ?", request.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}

	if err := database.Db.Model(&user).Update("is_email_verified", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update user"})
		return
	}
	if err := database.Db.Model(&otpModel).Where("id = ?", otpModel.Id).Update("email_otp", "").Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update OTP"})
		return
	}

	response := dto.UserResponseDto{
		ID:              user.ID,
		Email:           user.Email,
		IsEmailVerified: true,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Email verified successfully", "data": response})
}

func VerifyMobile(c *gin.Context) {
	var request dto.VerifyMobileDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	var otpModel models.Otp
	result := database.Db.Where("mobile = ?", request.Mobile).First(&otpModel)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}

	if otpModel.MobileOtp != request.MobileOtp || otpModel.OtpExp < time.Now().Unix() {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Invalid OTP"})
		return
	}

	var user models.User
	result = database.Db.Where("mobile = ?", request.Mobile).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}

	if err := database.Db.Model(&user).Update("is_mobile_verified", true).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update user"})
		return
	}

	accessTokenExp := time.Now().Add(time.Hour * 24).Unix()
	accessToken, err := utils.GenerateToken(user.ID, accessTokenExp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate access token " + err.Error()})
		return
	}

	refreshTokenExp := time.Now().Add(time.Hour * 24 * 30).Unix()
	refreshToken, err := utils.GenerateToken(user.ID, refreshTokenExp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate refresh token " + err.Error()})
		return
	}

	if err := database.Db.Delete(&otpModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to delete OTP " + err.Error()})
		return
	}

	response := dto.UserResponseDto{
		ID:               user.ID,
		Mobile:           user.Mobile,
		IsMobileVerified: true,
		AccessToken:      accessToken,
		AccessTokenEx:    accessTokenExp,
		RefreshToken:     refreshToken,
		RefreshTokenEx:   refreshTokenExp,
	}

	status, err := utils.CreatePrivateGroup(user.ID)
	if err != nil {
		c.JSON(status, gin.H{"success": false, "message": "Failed to create private group " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Mobile verified successfully", "data": response})
}

func GetAccessTokenFromRefreshToken(c *gin.Context) {
	var request dto.RefreshTokenDto

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	if request.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Refresh token is required"})
		return
	}

	claims, err := utils.VerifyToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid authentication token"})
		return
	}

	userId := claims["user_id"]

	var user models.User
	result := database.Db.First(&user, userId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}

	accessTokenExp := time.Now().Add(time.Hour * 24).Unix()
	accessToken, err := utils.GenerateToken(user.ID, accessTokenExp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate access token " + err.Error()})
		return
	}

	response := dto.AccessTokenResponseDto{
		ID:            user.ID,
		AccessToken:   accessToken,
		AccessTokenEx: accessTokenExp,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Access token generated successfully", "data": response})
}

func GetUserDetails(c *gin.Context) {
	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Unauthorized"})
		return
	}

	var user models.User
	result := database.Db.First(&user, userId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"success": false, "message": "User not found"})
		return
	}

	response := dto.UserResponseDto{
		ID:               user.ID,
		Email:            user.Email,
		Mobile:           user.Mobile,
		IsEmailVerified:  user.IsEmailVerified,
		IsMobileVerified: user.IsMobileVerified,
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User details fetched successfully", "data": response})
}
