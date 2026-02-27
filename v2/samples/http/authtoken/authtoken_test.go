package authtoken_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/IBM/fp-go/v2/samples/http/authtoken"
)

func fakeAuthServerWithToken(t *testing.T, token string, expiresIn int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"access_token": "%s", "expires_in": %d, "token_type": "Bearer"}`, token, expiresIn)
	}))
}

func fakeAuthServerError(t *testing.T, statusCode int, body string) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(body))
	}))
}

func TestGetToken_RefreshesWhenCacheEmpty(t *testing.T) {
	server := fakeAuthServerWithToken(t, "fresh-token", 3600)
	defer server.Close()

	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	cache := authtoken.MakeTokenCache()()
	env := authtoken.MakeEnv(server.URL, "test-api-key", server.Client(), cache, func() time.Time { return fixedTime })

	ctx := context.Background()
	token, err := authtoken.GetToken(env)(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token.Token != "fresh-token" {
		t.Errorf("expected token 'fresh-token', got '%s'", token.Token)
	}
}

func TestGetToken_UsesCachedToken(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"access_token": "token-%d", "expires_in": 3600, "token_type": "Bearer"}`, callCount)
	}))
	defer server.Close()

	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	cache := authtoken.MakeTokenCache()()
	env := authtoken.MakeEnv(server.URL, "test-api-key", server.Client(), cache, func() time.Time { return fixedTime })

	ctx := context.Background()

	// First call - should hit server
	token1, _ := authtoken.GetToken(env)(ctx)
	// Second call - should use cache
	token2, _ := authtoken.GetToken(env)(ctx)

	if callCount != 1 {
		t.Errorf("expected 1 server call, got %d", callCount)
	}
	if token1.Token != token2.Token {
		t.Errorf("expected same token, got '%s' and '%s'", token1.Token, token2.Token)
	}
}

func TestGetToken_RefreshesExpiredToken(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"access_token": "token-%d", "expires_in": 3600, "token_type": "Bearer"}`, callCount)
	}))
	defer server.Close()

	currentTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	cache := authtoken.MakeTokenCache()()
	env := authtoken.MakeEnv(server.URL, "test-api-key", server.Client(), cache, func() time.Time { return currentTime })

	ctx := context.Background()

	// First call
	token1, _ := authtoken.GetToken(env)(ctx)

	// Advance time past expiration
	currentTime = currentTime.Add(2 * time.Hour)
	env = authtoken.MakeEnv(server.URL, "test-api-key", server.Client(), cache, func() time.Time { return currentTime })

	// Second call - should refresh
	token2, _ := authtoken.GetToken(env)(ctx)

	if callCount != 2 {
		t.Errorf("expected 2 server calls, got %d", callCount)
	}
	if token1.Token == token2.Token {
		t.Errorf("expected different tokens after expiration")
	}
}

func TestGetToken_HandlesServerError(t *testing.T) {
	server := fakeAuthServerError(t, http.StatusUnauthorized, "invalid credentials")
	defer server.Close()

	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	cache := authtoken.MakeTokenCache()()
	env := authtoken.MakeEnv(server.URL, "bad-api-key", server.Client(), cache, func() time.Time { return fixedTime })

	ctx := context.Background()
	_, err := authtoken.GetToken(env)(ctx)

	if err == nil {
		t.Fatal("expected error for unauthorized request")
	}
}

func TestGetTokenString_ReturnsTokenValue(t *testing.T) {
	server := fakeAuthServerWithToken(t, "my-access-token", 3600)
	defer server.Close()

	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	cache := authtoken.MakeTokenCache()()
	env := authtoken.MakeEnv(server.URL, "test-api-key", server.Client(), cache, func() time.Time { return fixedTime })

	ctx := context.Background()
	tokenStr, err := authtoken.GetTokenString(env)(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tokenStr != "my-access-token" {
		t.Errorf("expected 'my-access-token', got '%s'", tokenStr)
	}
}

func TestGetAuthorizationHeader_ReturnsBearerFormat(t *testing.T) {
	server := fakeAuthServerWithToken(t, "my-token", 3600)
	defer server.Close()

	fixedTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	cache := authtoken.MakeTokenCache()()
	env := authtoken.MakeEnv(server.URL, "test-api-key", server.Client(), cache, func() time.Time { return fixedTime })

	ctx := context.Background()
	header, err := authtoken.GetAuthorizationHeader(env)(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "Bearer my-token"
	if header != expected {
		t.Errorf("expected '%s', got '%s'", expected, header)
	}
}
