package authtoken

import (
	"context"

	IO "github.com/IBM/fp-go/v2/io"
	O "github.com/IBM/fp-go/v2/option"
)

// GetToken returns a valid token, refreshing only if needed.
func GetToken(env Env) func(context.Context) (CachedToken, error) {
	return func(ctx context.Context) (CachedToken, error) {
		currentTime := env.Now()
		cache := env.Cache

		cachedOpt := cache.GetIfValid(currentTime)()
		if token, ok := O.Unwrap(cachedOpt); ok {
			return token, nil
		}

		newToken, err := RefreshToken(env)(ctx)
		if err != nil {
			return CachedToken{}, err
		}

		cache.Set(newToken)()
		return newToken, nil
	}
}

// GetTokenIO returns an IO-lifted version of GetToken.
func GetTokenIO(env Env) func(context.Context) IO.IO[CachedToken] {
	return func(ctx context.Context) IO.IO[CachedToken] {
		return func() CachedToken {
			token, _ := GetToken(env)(ctx)
			return token
		}
	}
}

// GetTokenString returns just the access token string.
func GetTokenString(env Env) func(context.Context) (string, error) {
	return func(ctx context.Context) (string, error) {
		token, err := GetToken(env)(ctx)
		if err != nil {
			return "", err
		}
		return token.Token, nil
	}
}

// GetAuthorizationHeader returns a "Bearer <token>" header value.
func GetAuthorizationHeader(env Env) func(context.Context) (string, error) {
	return func(ctx context.Context) (string, error) {
		token, err := GetTokenString(env)(ctx)
		if err != nil {
			return "", err
		}
		return FormatTokenForAuth(token), nil
	}
}
