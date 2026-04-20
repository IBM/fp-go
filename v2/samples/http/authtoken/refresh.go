package authtoken

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	DefaultSafetyBuffer = 5 * time.Minute
	AuthorizationHeader = "Authorization"
	ContentTypeHeader   = "Content-Type"
	ContentTypeJSON     = "application/json"
	BearerPrefix        = "Bearer "
)

// RefreshToken performs a token refresh using the default safety buffer.
func RefreshToken(env Env) func(context.Context) (CachedToken, error) {
	return RefreshTokenWithBuffer(env, DefaultSafetyBuffer)
}

// RefreshTokenWithBuffer performs a token refresh with a custom safety buffer.
func RefreshTokenWithBuffer(env Env, safetyBuffer time.Duration) func(context.Context) (CachedToken, error) {
	return func(ctx context.Context) (CachedToken, error) {
		authURI := env.AuthURL
		apiKey := env.APIKey
		client := env.HTTPClient
		currentTime := env.Now()

		if client == nil {
			client = http.DefaultClient
		}

		request, err := http.NewRequestWithContext(ctx, http.MethodPost, authURI, nil)
		if err != nil {
			return CachedToken{}, fmt.Errorf("failed to create auth request: %w", err)
		}

		request.Header.Set(AuthorizationHeader, BearerPrefix+apiKey)
		request.Header.Set(ContentTypeHeader, ContentTypeJSON)

		resp, err := client.Do(request)
		if err != nil {
			return CachedToken{}, fmt.Errorf("failed to execute auth request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return CachedToken{}, fmt.Errorf("auth request failed with status %d: %s", resp.StatusCode, string(body))
		}

		var authResp authResponse
		if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
			return CachedToken{}, fmt.Errorf("failed to decode auth response: %w", err)
		}

		lenses := MakeauthResponseLenses()
		accessToken := lenses.AccessToken.Get(authResp)
		expiresIn := lenses.ExpiresIn.Get(authResp)
		refreshToken := lenses.RefreshToken.Get(authResp)

		if accessToken == "" {
			return CachedToken{}, fmt.Errorf("auth response missing access_token")
		}
		if expiresIn <= 0 {
			return CachedToken{}, fmt.Errorf("auth response has invalid expires_in: %d", expiresIn)
		}

		expiresAt := calculateExpirationTime(currentTime, expiresIn, safetyBuffer)

		return CachedToken{
			Token:        accessToken,
			ExpiresAt:    expiresAt,
			RefreshToken: refreshToken,
		}, nil
	}
}

func calculateExpirationTime(currentTime time.Time, expiresInSeconds int, safetyBuffer time.Duration) time.Time {
	expiresIn := time.Duration(expiresInSeconds) * time.Second
	effectiveLifetime := expiresIn - safetyBuffer
	if effectiveLifetime < 0 {
		effectiveLifetime = 0
	}
	return currentTime.Add(effectiveLifetime)
}

// IsTokenExpired checks if a token is expired.
func IsTokenExpired(token CachedToken, currentTime time.Time) bool {
	return !currentTime.Before(token.ExpiresAt)
}

// IsTokenValid checks if a token is valid (not expired).
func IsTokenValid(token CachedToken, currentTime time.Time) bool {
	return currentTime.Before(token.ExpiresAt)
}

// TimeUntilExpiration returns the duration until the token expires.
func TimeUntilExpiration(token CachedToken, currentTime time.Time) time.Duration {
	if IsTokenExpired(token, currentTime) {
		return 0
	}
	return token.ExpiresAt.Sub(currentTime)
}

// FormatTokenForAuth formats a token for use in an Authorization header.
func FormatTokenForAuth(token string) string {
	if strings.HasPrefix(token, BearerPrefix) {
		return token
	}
	return BearerPrefix + token
}
