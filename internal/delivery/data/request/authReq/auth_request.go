package authreq

type GoogleLoginRequest struct {
	IDToken string `json:"idToken" binding:"required"`
}
