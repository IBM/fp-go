package prism

import (
	"net/url"
	"testing"

	"github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestURLPrisms tests prisms for url.URL
func TestURLPrisms(t *testing.T) {
	prisms := MakeURLPrisms()

	t.Run("Scheme prism", func(t *testing.T) {
		u := url.URL{Scheme: "https", Host: "example.com"}
		opt := prisms.Scheme.GetOption(u)
		assert.True(t, option.IsSome(opt))

		emptyU := url.URL{Scheme: "", Host: "example.com"}
		opt = prisms.Scheme.GetOption(emptyU)
		assert.True(t, option.IsNone(opt))

		// ReverseGet
		constructed := prisms.Scheme.ReverseGet("ftp")
		assert.Equal(t, "ftp", constructed.Scheme)
	})

	t.Run("Host prism", func(t *testing.T) {
		u := url.URL{Scheme: "https", Host: "example.com"}
		opt := prisms.Host.GetOption(u)
		assert.True(t, option.IsSome(opt))
		value := option.GetOrElse(func() string { return "" })(opt)
		assert.Equal(t, "example.com", value)
	})

	t.Run("Path prism", func(t *testing.T) {
		u := url.URL{Scheme: "https", Host: "example.com", Path: "/api"}
		opt := prisms.Path.GetOption(u)
		assert.True(t, option.IsSome(opt))

		emptyPath := url.URL{Scheme: "https", Host: "example.com", Path: ""}
		opt = prisms.Path.GetOption(emptyPath)
		assert.True(t, option.IsNone(opt))
	})

	t.Run("User prism", func(t *testing.T) {
		userinfo := url.User("john")
		u := url.URL{Scheme: "https", Host: "example.com", User: userinfo}
		opt := prisms.User.GetOption(u)
		assert.True(t, option.IsSome(opt))

		noUser := url.URL{Scheme: "https", Host: "example.com", User: nil}
		opt = prisms.User.GetOption(noUser)
		assert.True(t, option.IsNone(opt))
	})

	t.Run("Boolean prisms", func(t *testing.T) {
		u := url.URL{Scheme: "https", Host: "example.com", ForceQuery: true}
		opt := prisms.ForceQuery.GetOption(u)
		assert.True(t, option.IsSome(opt))

		noForce := url.URL{Scheme: "https", Host: "example.com", ForceQuery: false}
		opt = prisms.ForceQuery.GetOption(noForce)
		assert.True(t, option.IsNone(opt))
	})
}

// TestErrorPrisms tests prisms for url.Error
func TestErrorPrisms(t *testing.T) {
	prisms := MakeErrorPrisms()

	t.Run("Op prism", func(t *testing.T) {
		urlErr := url.Error{Op: "Get", URL: "test", Err: nil}
		opt := prisms.Op.GetOption(urlErr)
		assert.True(t, option.IsSome(opt))

		emptyErr := url.Error{Op: "", URL: "test", Err: nil}
		opt = prisms.Op.GetOption(emptyErr)
		assert.True(t, option.IsNone(opt))

		// ReverseGet
		constructed := prisms.Op.ReverseGet("Post")
		assert.Equal(t, "Post", constructed.Op)
	})

	t.Run("URL prism", func(t *testing.T) {
		urlErr := url.Error{Op: "Get", URL: "https://example.com", Err: nil}
		opt := prisms.URL.GetOption(urlErr)
		assert.True(t, option.IsSome(opt))

		emptyErr := url.Error{Op: "Get", URL: "", Err: nil}
		opt = prisms.URL.GetOption(emptyErr)
		assert.True(t, option.IsNone(opt))
	})

	t.Run("Err prism", func(t *testing.T) {
		urlErr := url.Error{Op: "Get", URL: "test", Err: assert.AnError}
		opt := prisms.Err.GetOption(urlErr)
		assert.True(t, option.IsSome(opt))

		nilErr := url.Error{Op: "Get", URL: "test", Err: nil}
		opt = prisms.Err.GetOption(nilErr)
		assert.True(t, option.IsNone(opt))
	})
}
