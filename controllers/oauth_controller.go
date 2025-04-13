package controllers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
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
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "error": "Code is required"})
		return
	}
	

	token, err := utils.GoogleOAuthConfig.Exchange(c, code)
	if err != nil {
		fmt.Println("Error exchanging code for token:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to exchange token: " + err.Error()})
		return
	}

	

	client := utils.GoogleOAuthConfig.Client(c, token)
	fmt.Println("Client created with token:", token)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to get user info: " + err.Error()})
		return
	}

	body, err := io.ReadAll(userInfo.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "error": "Failed to read response body: " + err.Error()})
		return
	}
	defer userInfo.Body.Close()
	fmt.Println("User info response body:", string(body))
	// accessTokenExp := time.Now().Add(time.Hour * 24).Unix()
	// accessToken, err := utils.GenerateToken(uint(userInfo.ContentLength), accessTokenExp)
	c.JSON(http.StatusOK, gin.H{"success": true, "data": userInfo})
}
