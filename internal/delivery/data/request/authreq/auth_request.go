package authreq

type GoogleLoginRequest struct {
	IDToken string `json:"idToken" binding:"required"`
}

type GoogleOAuthCodeLoginRequest struct {
	Code        string `json:"code" binding:"required"`
	RedirectURI string `json:"redirectUri" binding:"required"`
}
