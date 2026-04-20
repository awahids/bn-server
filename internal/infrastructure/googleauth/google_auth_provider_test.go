package googleauth

import (
	"bytes"
	"context"
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

func httpResp(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}
