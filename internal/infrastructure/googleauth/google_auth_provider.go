package googleauth

import (
	"context"
	"encoding/base64"
	"encoding/json"
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

// googleIDTokenClaims are the ID token fields used after a code exchange.
// email_verified is a JSON bool in an ID token (unlike the tokeninfo endpoint).
type googleIDTokenClaims struct {
	Issuer        string `json:"iss"`
	Audience      string `json:"aud"`
	Subject       string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	ExpiresAt     int64  `json:"exp"`
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
		// Keep Google's error code (e.g. redirect_uri_mismatch) in the message;
		// its description is often just "Bad Request".
		detail := strings.TrimSpace(strings.Join([]string{payload.Error, payload.ErrorDescription}, ": "))
		detail = strings.Trim(detail, ": ")
		if detail == "" {
			detail = fmt.Sprintf("status %d", res.StatusCode)
		}
		switch payload.Error {
		case "invalid_grant", "invalid_request", "redirect_uri_mismatch":
			return nil, fmt.Errorf("%w: %s", serviceinterface.ErrGoogleAuthInvalidOAuthCode, detail)
		case "invalid_client", "unauthorized_client":
			return nil, fmt.Errorf("%w: %s", serviceinterface.ErrGoogleAuthNotConfigured, detail)
		default:
			return nil, fmt.Errorf("google token exchange failed: %s", detail)
		}
	}

	idToken := strings.TrimSpace(payload.IDToken)
	if idToken == "" {
		return nil, fmt.Errorf("%w: missing id_token", serviceinterface.ErrGoogleAuthInvalidOAuthCode)
	}

	// The ID token came straight from Google's token endpoint over TLS, so per
	// OpenID Connect Core 3.1.3.7 its claims can be checked locally instead of
	// making a second round trip to the tokeninfo endpoint.
	return p.tokenInfoFromIDTokenClaims(idToken, time.Now())
}

func (p *googleAuthProvider) tokenInfoFromIDTokenClaims(
	idToken string,
	now time.Time,
) (*serviceinterface.GoogleTokenInfo, error) {
	invalid := func(reason string) error {
		return fmt.Errorf("%w: %s", serviceinterface.ErrGoogleAuthInvalidOAuthCode, reason)
	}

	parts := strings.Split(idToken, ".")
	if len(parts) != 3 {
		return nil, invalid("malformed id_token")
	}
	rawClaims, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, invalid("malformed id_token payload")
	}

	var claims googleIDTokenClaims
	if err := json.Unmarshal(rawClaims, &claims); err != nil {
		return nil, invalid("malformed id_token claims")
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return nil, invalid("unexpected id_token issuer")
	}
	if strings.TrimSpace(claims.Audience) != p.clientID {
		return nil, invalid("id_token audience mismatch")
	}
	if claims.ExpiresAt == 0 || now.Unix() >= claims.ExpiresAt {
		return nil, invalid("id_token expired")
	}
	if strings.TrimSpace(claims.Subject) == "" {
		return nil, invalid("id_token subject is empty")
	}

	return &serviceinterface.GoogleTokenInfo{
		Subject:       strings.TrimSpace(claims.Subject),
		Email:         strings.TrimSpace(claims.Email),
		EmailVerified: claims.EmailVerified,
		Name:          strings.TrimSpace(claims.Name),
		Picture:       strings.TrimSpace(claims.Picture),
	}, nil
}
