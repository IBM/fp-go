package authtoken

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	H "github.com/IBM/fp-go/v2/context/readerioresult/http"
)

type SimpleResult[A any] = func(context.Context) (A, error)

func MakeAuthURIDep(uri string) SimpleResult[string] {
	return func(ctx context.Context) (string, error) {
		return uri, nil
	}
}

func MakeAPIKeyDep(apiKey string) SimpleResult[string] {
	return func(ctx context.Context) (string, error) {
		return apiKey, nil
	}
}

func MakeHTTPClientDep(client *http.Client) SimpleResult[Client] {
	return func(ctx context.Context) (Client, error) {
		if client == nil {
			client = http.DefaultClient
		}
		return H.MakeClient(client), nil
	}
}

func MakeCurrentTimeDep() SimpleResult[time.Time] {
	return func(ctx context.Context) (time.Time, error) {
		return time.Now(), nil
	}
}

func MakeFixedTimeDep(t time.Time) SimpleResult[time.Time] {
	return func(ctx context.Context) (time.Time, error) {
		return t, nil
	}
}

func MakeDeps(authURI, apiKey string, client *http.Client, currentTime SimpleResult[time.Time]) TokenServiceDeps {
	return TokenServiceDeps{
		AuthURI:     wrapSimpleToReader(MakeAuthURIDep(authURI)),
		APIKey:      wrapSimpleToReader(MakeAPIKeyDep(apiKey)),
		HTTPClient:  wrapSimpleToReaderClient(MakeHTTPClientDep(client)),
		CurrentTime: wrapSimpleToReaderTime(currentTime),
	}
}

func wrapSimpleToReader(f SimpleResult[string]) ReaderIOResult[Env, string] {
	return func(ctx context.Context, env Env) (string, error) {
		return f(ctx)
	}
}

func wrapSimpleToReaderTime(f SimpleResult[time.Time]) ReaderIOResult[Env, time.Time] {
	return func(ctx context.Context, env Env) (time.Time, error) {
		return f(ctx)
	}
}

func wrapSimpleToReaderClient(f SimpleResult[Client]) ReaderIOResult[Env, Client] {
	return func(ctx context.Context, env Env) (Client, error) {
		return f(ctx)
	}
}

func MakeProductionDeps(authURI, apiKey string) TokenServiceDeps {
	return MakeDeps(authURI, apiKey, http.DefaultClient, MakeCurrentTimeDep())
}

func GetAuthURI(deps TokenServiceDeps) ReaderIOResult[Env, string] {
	return deps.AuthURI
}

func GetAPIKey(deps TokenServiceDeps) ReaderIOResult[Env, string] {
	return deps.APIKey
}

func GetHTTPClient(deps TokenServiceDeps) ReaderIOResult[Env, Client] {
	return deps.HTTPClient
}

func GetCurrentTime(deps TokenServiceDeps) ReaderIOResult[Env, time.Time] {
	return deps.CurrentTime
}

func WithAuthURI(deps TokenServiceDeps, uri string) TokenServiceDeps {
	deps.AuthURI = wrapSimpleToReader(MakeAuthURIDep(uri))
	return deps
}

func WithAPIKey(deps TokenServiceDeps, apiKey string) TokenServiceDeps {
	deps.APIKey = wrapSimpleToReader(MakeAPIKeyDep(apiKey))
	return deps
}

func WithHTTPClient(deps TokenServiceDeps, client *http.Client) TokenServiceDeps {
	deps.HTTPClient = func(ctx context.Context, env Env) (Client, error) {
		return H.MakeClient(client), nil
	}
	return deps
}

func WithCurrentTime(deps TokenServiceDeps, currentTime SimpleResult[time.Time]) TokenServiceDeps {
	deps.CurrentTime = wrapSimpleToReaderTime(currentTime)
	return deps
}

func MakeAuthURIFromEnvDep(envVarName string) SimpleResult[string] {
	return func(ctx context.Context) (string, error) {
		value := os.Getenv(envVarName)
		if value == "" {
			return "", fmt.Errorf("environment variable %s is not set or is empty", envVarName)
		}
		return value, nil
	}
}

func MakeAPIKeyFromEnvDep(envVarName string) SimpleResult[string] {
	return func(ctx context.Context) (string, error) {
		value := os.Getenv(envVarName)
		if value == "" {
			return "", fmt.Errorf("environment variable %s is not set or is empty", envVarName)
		}
		return value, nil
	}
}

func MakeAuthURIFromEnvOrDefaultDep(envVarName, defaultValue string) SimpleResult[string] {
	return func(ctx context.Context) (string, error) {
		if value := os.Getenv(envVarName); value != "" {
			return value, nil
		}
		return defaultValue, nil
	}
}

func MakeAPIKeyFromEnvOrDefaultDep(envVarName, defaultValue string) SimpleResult[string] {
	return func(ctx context.Context) (string, error) {
		if value := os.Getenv(envVarName); value != "" {
			return value, nil
		}
		return defaultValue, nil
	}
}

func MakeDepsFromEnv(authURIEnvVar, apiKeyEnvVar string, client *http.Client, currentTime SimpleResult[time.Time]) TokenServiceDeps {
	return TokenServiceDeps{
		AuthURI:     wrapSimpleToReader(MakeAuthURIFromEnvDep(authURIEnvVar)),
		APIKey:      wrapSimpleToReader(MakeAPIKeyFromEnvDep(apiKeyEnvVar)),
		HTTPClient:  wrapSimpleToReaderClient(MakeHTTPClientDep(client)),
		CurrentTime: wrapSimpleToReaderTime(currentTime),
	}
}

func MakeProductionDepsFromEnv(authURIEnvVar, apiKeyEnvVar string) TokenServiceDeps {
	return MakeDepsFromEnv(authURIEnvVar, apiKeyEnvVar, http.DefaultClient, MakeCurrentTimeDep())
}

// MakeEnv creates a complete Env structure.
func MakeEnv(authURL, apiKey string, client *http.Client, cache *TokenCache, now func() time.Time) Env {
	if client == nil {
		client = http.DefaultClient
	}
	return Env{
		AuthURL:    authURL,
		APIKey:     apiKey,
		HTTPClient: client,
		Cache:      cache,
		Now:        now,
	}
}
