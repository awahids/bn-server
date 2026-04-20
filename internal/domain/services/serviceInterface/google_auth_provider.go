package serviceinterface

import (
	"context"
	"errors"
)

var (
	ErrGoogleAuthInvalidIDToken   = errors.New("invalid google id token")
	ErrGoogleAuthInvalidOAuthCode = errors.New("invalid google oauth code")
	ErrGoogleAuthNotConfigured    = errors.New("google auth provider is not configured")
)

type GoogleTokenInfo struct {
	Subject       string
	Email         string
	EmailVerified bool
	Name          string
	Picture       string
}

type GoogleAuthProvider interface {
	GetTokenInfoByIDToken(ctx context.Context, idToken string) (*GoogleTokenInfo, error)
	GetTokenInfoByOAuthCode(ctx context.Context, code, redirectURI string) (*GoogleTokenInfo, error)
}
