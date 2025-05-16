package dto

type GoogleUserInfo struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Verified bool   `json:"email_verified"`
}

type UserInfo struct {
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type PhoneData struct {
	PhoneNumbers []struct {
		Value string `json:"value"`
	} `json:"phoneNumbers"`
}
