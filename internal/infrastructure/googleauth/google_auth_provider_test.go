package googleauth

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"
)

type testRoundTripFunc func(req *http.Request) (*http.Response, error)

func (f testRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestGetTokenInfoByIDToken_Success(t *testing.T) {
	client := &http.Client{
		Timeout: time.Second,
		Transport: testRoundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Host != "oauth2.googleapis.com" {
				return nil, errors.New("unexpected host")
			}
			payload := `{"aud":"client-id","sub":"google-sub","email":"user@example.com","email_verified":"true","name":"User Name","picture":"https://example.com/avatar.png"}`
			return httpResp(http.StatusOK, payload), nil
		}),
	}

	provider := NewGoogleAuthProvider("client-id", "client-secret", client)
	tokenInfo, err := provider.GetTokenInfoByIDToken(context.Background(), "id-token")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tokenInfo == nil {
		t.Fatal("expected token info, got nil")
	}
	if tokenInfo.Subject != "google-sub" {
		t.Fatalf("expected subject google-sub, got %s", tokenInfo.Subject)
	}
	if !tokenInfo.EmailVerified {
		t.Fatal("expected email verified true")
	}
}

func TestGetTokenInfoByOAuthCode_InvalidGrant(t *testing.T) {
	client := &http.Client{
		Timeout: time.Second,
		Transport: testRoundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Host != "oauth2.googleapis.com" || !strings.HasPrefix(req.URL.Path, "/token") {
				return nil, errors.New("unexpected request")
			}
			payload := `{"error":"invalid_grant","error_description":"Bad Request"}`
			return httpResp(http.StatusBadRequest, payload), nil
		}),
	}

	provider := NewGoogleAuthProvider("client-id", "client-secret", client)
	tokenInfo, err := provider.GetTokenInfoByOAuthCode(context.Background(), "bad-code", "postmessage")
	if !errors.Is(err, serviceinterface.ErrGoogleAuthInvalidOAuthCode) {
		t.Fatalf("expected ErrGoogleAuthInvalidOAuthCode, got %v", err)
	}
	if tokenInfo != nil {
		t.Fatal("expected nil token info for invalid oauth code")
	}
}

func fakeIDToken(t *testing.T, claims map[string]any) string {
	t.Helper()
	raw, err := json.Marshal(claims)
	if err != nil {
		t.Fatalf("marshal claims: %v", err)
	}
	return "header." + base64.RawURLEncoding.EncodeToString(raw) + ".signature"
}

func tokenEndpointClient(t *testing.T, status int, body string) *http.Client {
	t.Helper()
	return &http.Client{
		Timeout: time.Second,
		Transport: testRoundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Host != "oauth2.googleapis.com" || req.URL.Path != "/token" {
				// The tokeninfo endpoint must no longer be called after a code exchange.
				return nil, errors.New("unexpected request to " + req.URL.String())
			}
			return httpResp(status, body), nil
		}),
	}
}

func TestGetTokenInfoByOAuthCode_UsesIDTokenClaimsWithoutTokenInfoCall(t *testing.T) {
	idToken := fakeIDToken(t, map[string]any{
		"iss":            "https://accounts.google.com",
		"aud":            "client-id",
		"sub":            "google-sub",
		"email":          "user@example.com",
		"email_verified": true,
		"name":           "User Name",
		"picture":        "https://example.com/avatar.png",
		"exp":            time.Now().Add(time.Hour).Unix(),
	})
	client := tokenEndpointClient(t, http.StatusOK, `{"id_token":"`+idToken+`"}`)

	provider := NewGoogleAuthProvider("client-id", "client-secret", client)
	info, err := provider.GetTokenInfoByOAuthCode(context.Background(), "code", "https://app.example.com/auth/callback")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if info.Subject != "google-sub" || info.Email != "user@example.com" || !info.EmailVerified {
		t.Fatalf("unexpected token info: %+v", info)
	}
}

func TestGetTokenInfoByOAuthCode_RejectsInvalidIDTokenClaims(t *testing.T) {
	base := map[string]any{
		"iss": "accounts.google.com",
		"aud": "client-id",
		"sub": "google-sub",
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	cases := map[string]map[string]any{
		"wrong audience": {"aud": "other-client"},
		"wrong issuer":   {"iss": "https://evil.example.com"},
		"expired":        {"exp": time.Now().Add(-time.Minute).Unix()},
		"missing sub":    {"sub": ""},
	}
	for name, override := range cases {
		t.Run(name, func(t *testing.T) {
			claims := map[string]any{}
			for k, v := range base {
				claims[k] = v
			}
			for k, v := range override {
				claims[k] = v
			}
			client := tokenEndpointClient(t, http.StatusOK, `{"id_token":"`+fakeIDToken(t, claims)+`"}`)
			provider := NewGoogleAuthProvider("client-id", "client-secret", client)
			_, err := provider.GetTokenInfoByOAuthCode(context.Background(), "code", "postmessage")
			if !errors.Is(err, serviceinterface.ErrGoogleAuthInvalidOAuthCode) {
				t.Fatalf("expected ErrGoogleAuthInvalidOAuthCode, got %v", err)
			}
		})
	}
}

func TestGetTokenInfoByOAuthCode_KeepsGoogleErrorCode(t *testing.T) {
	client := tokenEndpointClient(t, http.StatusBadRequest, `{"error":"redirect_uri_mismatch","error_description":"Bad Request"}`)
	provider := NewGoogleAuthProvider("client-id", "client-secret", client)
	_, err := provider.GetTokenInfoByOAuthCode(context.Background(), "code", "https://app.example.com/auth/callback")
	if !errors.Is(err, serviceinterface.ErrGoogleAuthInvalidOAuthCode) {
		t.Fatalf("expected ErrGoogleAuthInvalidOAuthCode, got %v", err)
	}
	if !strings.Contains(err.Error(), "redirect_uri_mismatch") {
		t.Fatalf("expected error to mention redirect_uri_mismatch, got %v", err)
	}
}

func TestGetTokenInfoByOAuthCode_InvalidClientIsConfigError(t *testing.T) {
	client := tokenEndpointClient(t, http.StatusUnauthorized, `{"error":"invalid_client","error_description":"The OAuth client was not found."}`)
	provider := NewGoogleAuthProvider("client-id", "client-secret", client)
	_, err := provider.GetTokenInfoByOAuthCode(context.Background(), "code", "postmessage")
	if !errors.Is(err, serviceinterface.ErrGoogleAuthNotConfigured) {
		t.Fatalf("expected ErrGoogleAuthNotConfigured, got %v", err)
	}
}

func httpResp(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}
