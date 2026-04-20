package googleauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"
)

type googleAuthProvider struct {
	clientID     string
	clientSecret string
	httpClient   *http.Client
}

type googleTokenInfoResponse struct {
	IssuedTo      string `json:"issued_to"`
	Audience      string `json:"aud"`
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified,string"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
}

type googleTokenExchangeResponse struct {
	IDToken          string `json:"id_token"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func NewGoogleAuthProvider(clientID, clientSecret string, httpClient *http.Client) serviceinterface.GoogleAuthProvider {
	client := httpClient
	if client == nil {
		client = &http.Client{Timeout: 12 * time.Second}
	}

	return &googleAuthProvider{
		clientID:     strings.TrimSpace(clientID),
		clientSecret: strings.TrimSpace(clientSecret),
		httpClient:   client,
	}
}

func (p *googleAuthProvider) GetTokenInfoByIDToken(
	ctx context.Context,
	idToken string,
) (*serviceinterface.GoogleTokenInfo, error) {
	idToken = strings.TrimSpace(idToken)
	if idToken == "" {
		return nil, serviceinterface.ErrGoogleAuthInvalidIDToken
	}
	if p.clientID == "" {
		return nil, serviceinterface.ErrGoogleAuthNotConfigured
	}

	endpoint := "https://oauth2.googleapis.com/tokeninfo?id_token=" + url.QueryEscape(idToken)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	res, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: tokeninfo status %d", serviceinterface.ErrGoogleAuthInvalidIDToken, res.StatusCode)
	}

	var payload googleTokenInfoResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	audience := strings.TrimSpace(payload.Audience)
	if audience == "" {
		audience = strings.TrimSpace(payload.IssuedTo)
	}
	if audience != p.clientID {
		return nil, fmt.Errorf("%w: audience mismatch", serviceinterface.ErrGoogleAuthInvalidIDToken)
	}
	if strings.TrimSpace(payload.Subject) == "" {
		return nil, fmt.Errorf("%w: subject is empty", serviceinterface.ErrGoogleAuthInvalidIDToken)
	}

	return &serviceinterface.GoogleTokenInfo{
		Subject:       strings.TrimSpace(payload.Subject),
		Email:         strings.TrimSpace(payload.Email),
		EmailVerified: payload.EmailVerified,
		Name:          strings.TrimSpace(payload.Name),
		Picture:       strings.TrimSpace(payload.Picture),
	}, nil
}

func (p *googleAuthProvider) GetTokenInfoByOAuthCode(
	ctx context.Context,
	code, redirectURI string,
) (*serviceinterface.GoogleTokenInfo, error) {
	if p.clientID == "" || p.clientSecret == "" {
		return nil, serviceinterface.ErrGoogleAuthNotConfigured
	}

	code = strings.TrimSpace(code)
	redirectURI = strings.TrimSpace(redirectURI)
	if code == "" || redirectURI == "" {
		return nil, serviceinterface.ErrGoogleAuthInvalidOAuthCode
	}
	if redirectURI != "postmessage" {
		parsed, err := url.ParseRequestURI(redirectURI)
		if err != nil || (parsed.Scheme != "https" && parsed.Scheme != "http") {
			return nil, serviceinterface.ErrGoogleAuthInvalidOAuthCode
		}
	}

	form := url.Values{}
	form.Set("code", code)
	form.Set("client_id", p.clientID)
	form.Set("client_secret", p.clientSecret)
	form.Set("grant_type", "authorization_code")
	form.Set("redirect_uri", redirectURI)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://oauth2.googleapis.com/token",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var payload googleTokenExchangeResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		if payload.Error == "invalid_grant" || payload.Error == "invalid_request" {
			if payload.ErrorDescription != "" {
				return nil, fmt.Errorf("%w: %s", serviceinterface.ErrGoogleAuthInvalidOAuthCode, payload.ErrorDescription)
			}
			return nil, serviceinterface.ErrGoogleAuthInvalidOAuthCode
		}
		if payload.ErrorDescription != "" {
			return nil, fmt.Errorf("google token exchange failed: %s", payload.ErrorDescription)
		}
		if payload.Error != "" {
			return nil, fmt.Errorf("google token exchange failed: %s", payload.Error)
		}
		return nil, fmt.Errorf("google token exchange failed with status %d", res.StatusCode)
	}

	idToken := strings.TrimSpace(payload.IDToken)
	if idToken == "" {
		return nil, fmt.Errorf("%w: missing id_token", serviceinterface.ErrGoogleAuthInvalidOAuthCode)
	}

	tokenInfo, err := p.GetTokenInfoByIDToken(ctx, idToken)
	if err != nil {
		if errors.Is(err, serviceinterface.ErrGoogleAuthInvalidIDToken) {
			return nil, fmt.Errorf("%w: invalid id token returned by provider", serviceinterface.ErrGoogleAuthInvalidOAuthCode)
		}
		return nil, err
	}

	return tokenInfo, nil
}
