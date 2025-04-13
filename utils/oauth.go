package utils

import (
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "os"
)

var GoogleOAuthConfig *oauth2.Config

func InitGoogleOAuth() {
    GoogleOAuthConfig = &oauth2.Config{
        ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
        ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
        RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
        Scopes: []string{
            "https://www.googleapis.com/auth/userinfo.email",
            "https://www.googleapis.com/auth/userinfo.profile",
        },
        Endpoint:     google.Endpoint,
    }
}	