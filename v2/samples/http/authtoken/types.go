package authtoken

//go:generate go run ../../../ lens --dir . --filename gen_lens.go

import (
	"context"
	"net/http"
	"time"

	H "github.com/IBM/fp-go/v2/context/readerioresult/http"
	O "github.com/IBM/fp-go/v2/option"
)

type (
	Option[A any] = O.Option[A]

	ReaderIOResult[R, A any] = func(context.Context, R) (A, error)

	Client = H.Client
)

type Env struct {
	AuthURL string

	APIKey string

	HTTPClient *http.Client

	Cache *TokenCache

	Now func() time.Time
}

// fp-go:Lens
type authResponse struct {
	AccessToken string `json:"access_token"`

	ExpiresIn int `json:"expires_in"`

	TokenType string `json:"token_type"`

	RefreshToken string `json:"refresh_token,omitempty"`
}

// fp-go:Lens
type CachedToken struct {
	Token string

	ExpiresAt time.Time

	RefreshToken string
}

// fp-go:Lens
type TokenServiceDeps struct {
	AuthURI ReaderIOResult[Env, string]

	APIKey ReaderIOResult[Env, string]

	HTTPClient ReaderIOResult[Env, Client]

	CurrentTime ReaderIOResult[Env, time.Time]

	Cache ReaderIOResult[Env, *TokenCache]
}
