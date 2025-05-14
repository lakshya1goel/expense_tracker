package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lakshya1goel/expense_tracker/database"
	"github.com/lakshya1goel/expense_tracker/dto"
	"github.com/lakshya1goel/expense_tracker/models"
	"github.com/lakshya1goel/expense_tracker/utils"
	"golang.org/x/oauth2"
)

func GoogleSignIn(c *gin.Context) {
	url := utils.GoogleOAuthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
	fmt.Println("Redirecting to Google OAuth URL:", url)
}

func GoogleCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Code is required"})
		return
	}

	token, err := utils.GoogleOAuthConfig.Exchange(c, code)
	if err != nil {
		fmt.Println("Error exchanging code for token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to exchange token: " + err.Error()})
		return
	}

	client := utils.GoogleOAuthConfig.Client(c, token)
	fmt.Println("Client created with token:", token)
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to get user info: " + err.Error()})
		return
	}
	defer userInfoResp.Body.Close()

	var userInfo struct {
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		Picture       string `json:"picture"`
	}
	body, _ := io.ReadAll(userInfoResp.Body)
	json.Unmarshal(body, &userInfo)
	fmt.Println("User info response body:", string(body))

	phoneResp, err := client.Get("https://people.googleapis.com/v1/people/me?personFields=phoneNumbers")
	var phone string
	if err == nil {
		defer phoneResp.Body.Close()
		var phoneData struct {
			PhoneNumbers []struct {
				Value string `json:"value"`
			} `json:"phoneNumbers"`
		}
		phoneBody, _ := io.ReadAll(phoneResp.Body)
		json.Unmarshal(phoneBody, &phoneData)
		if len(phoneData.PhoneNumbers) > 0 {
			phone = phoneData.PhoneNumbers[0].Value
		}
	}

	var user models.User
	result := database.Db.Where("email = ?", userInfo.Email).First(&user)
	if result.Error == nil {
		if user.IsEmailVerified && user.IsMobileVerified {
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

			c.JSON(http.StatusOK, gin.H{"success": true, "message": "User logged in successfully", "data": response})
			return
		} else {
			if !user.IsEmailVerified {
				user.IsEmailVerified = userInfo.VerifiedEmail
			}
			if !user.IsMobileVerified {
				user.IsMobileVerified = phone != ""
			}
			if err := database.Db.Save(&user).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to update user: " + err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"success": true, "message": "User updated successfully", "data": user})
			return
		}
	}

	password := (userInfo.Name + userInfo.Email + time.Now().String())
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to hash password " + err.Error()})
		return
	}

	newUser := models.User{
		Email:            userInfo.Email,
		Mobile:           phone,
		Password:         hashedPassword,
		IsEmailVerified:  userInfo.VerifiedEmail,
		IsMobileVerified: phone != "",
	}

	if !newUser.IsEmailVerified || !newUser.IsMobileVerified {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Email or mobile number not verified"})
		return
	}

	if err := database.Db.Create(&newUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to create user: " + err.Error()})
		return
	}

	accessTokenExp := time.Now().Add(time.Hour * 24).Unix()
	accessToken, err := utils.GenerateToken(newUser.ID, accessTokenExp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate access token: " + err.Error()})
		return
	}
	refreshTokenExp := time.Now().Add(time.Hour * 24 * 30).Unix()
	refreshToken, err := utils.GenerateToken(newUser.ID, refreshTokenExp)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Failed to generate refresh token: " + err.Error()})
		return
	}
	response := dto.UserResponseDto{
		ID:               newUser.ID,
		Email:            newUser.Email,
		Mobile:           newUser.Mobile,
		IsEmailVerified:  newUser.IsEmailVerified,
		IsMobileVerified: newUser.IsMobileVerified,
		AccessToken:      accessToken,
		AccessTokenEx:    accessTokenExp,
		RefreshToken:     refreshToken,
		RefreshTokenEx:   refreshTokenExp,
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "User logged in successfully", "data": response})
}
