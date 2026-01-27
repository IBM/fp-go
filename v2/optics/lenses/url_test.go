// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package lenses

import (
	"errors"
	"net/url"
	"testing"

	"github.com/IBM/fp-go/v2/option"
	__option "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestUserinfoRefLenses_Username tests the Username lens
func TestUserinfoRefLenses_Username(t *testing.T) {
	lenses := MakeUserinfoRefLenses()

	t.Run("Get username from UserPassword", func(t *testing.T) {
		userinfo := url.UserPassword("john", "secret123")
		username := lenses.Username.Get(userinfo)
		assert.Equal(t, "john", username)
	})

	t.Run("Get username from User", func(t *testing.T) {
		userinfo := url.User("alice")
		username := lenses.Username.Get(userinfo)
		assert.Equal(t, "alice", username)
	})

	t.Run("Set username with password", func(t *testing.T) {
		userinfo := url.UserPassword("john", "secret123")
		updated := lenses.Username.Set("bob")(userinfo)
		assert.Equal(t, "bob", updated.Username())
		// Password should be preserved
		pwd, ok := updated.Password()
		assert.True(t, ok)
		assert.Equal(t, "secret123", pwd)
	})

	t.Run("Set username without password", func(t *testing.T) {
		userinfo := url.User("alice")
		updated := lenses.Username.Set("bob")(userinfo)
		assert.Equal(t, "bob", updated.Username())
		// Should still have no password
		_, ok := updated.Password()
		assert.False(t, ok)
	})
}

// TestUserinfoRefLenses_Password tests the Password lens
func TestUserinfoRefLenses_Password(t *testing.T) {
	lenses := MakeUserinfoRefLenses()

	t.Run("Get password when present", func(t *testing.T) {
		userinfo := url.UserPassword("john", "secret123")
		password := lenses.Password.Get(userinfo)
		assert.Equal(t, "secret123", password)
	})

	t.Run("Get password when absent", func(t *testing.T) {
		userinfo := url.User("alice")
		password := lenses.Password.Get(userinfo)
		assert.Equal(t, "", password)
	})

	t.Run("Set password", func(t *testing.T) {
		userinfo := url.UserPassword("john", "secret123")
		updated := lenses.Password.Set("newpass")(userinfo)
		assert.Equal(t, "john", updated.Username())
		pwd, ok := updated.Password()
		assert.True(t, ok)
		assert.Equal(t, "newpass", pwd)
	})

	t.Run("Set password on user without password", func(t *testing.T) {
		userinfo := url.User("alice")
		updated := lenses.Password.Set("newpass")(userinfo)
		assert.Equal(t, "alice", updated.Username())
		pwd, ok := updated.Password()
		assert.True(t, ok)
		assert.Equal(t, "newpass", pwd)
	})
}

// TestUserinfoRefLenses_UsernameO tests the optional Username lens
func TestUserinfoRefLenses_UsernameO(t *testing.T) {
	lenses := MakeUserinfoRefLenses()

	t.Run("Get non-empty username", func(t *testing.T) {
		userinfo := url.User("john")
		opt := lenses.UsernameO.Get(userinfo)
		assert.True(t, __option.IsSome(opt))
		value := __option.GetOrElse(func() string { return "" })(opt)
		assert.Equal(t, "john", value)
	})

	t.Run("Get empty username", func(t *testing.T) {
		userinfo := url.User("")
		opt := lenses.UsernameO.Get(userinfo)
		assert.True(t, __option.IsNone(opt))
	})

	t.Run("Set Some username", func(t *testing.T) {
		userinfo := url.User("")
		updated := lenses.UsernameO.Set(__option.Some("alice"))(userinfo)
		assert.Equal(t, "alice", updated.Username())
	})

	t.Run("Set None username", func(t *testing.T) {
		userinfo := url.User("john")
		updated := lenses.UsernameO.Set(__option.None[string]())(userinfo)
		assert.Equal(t, "", updated.Username())
	})
}

// TestUserinfoRefLenses_PasswordO tests the optional Password lens
func TestUserinfoRefLenses_PasswordO(t *testing.T) {
	lenses := MakeUserinfoRefLenses()

	t.Run("Get Some password when present", func(t *testing.T) {
		userinfo := url.UserPassword("john", "secret123")
		opt := lenses.PasswordO.Get(userinfo)
		assert.True(t, __option.IsSome(opt))
		value := __option.GetOrElse(func() string { return "" })(opt)
		assert.Equal(t, "secret123", value)
	})

	t.Run("Get None password when absent", func(t *testing.T) {
		userinfo := url.User("alice")
		opt := lenses.PasswordO.Get(userinfo)
		assert.True(t, __option.IsNone(opt))
	})

	t.Run("Set Some password", func(t *testing.T) {
		userinfo := url.User("john")
		updated := lenses.PasswordO.Set(__option.Some("newpass"))(userinfo)
		assert.Equal(t, "john", updated.Username())
		pwd, ok := updated.Password()
		assert.True(t, ok)
		assert.Equal(t, "newpass", pwd)
	})

	t.Run("Set None password removes it", func(t *testing.T) {
		userinfo := url.UserPassword("john", "secret123")
		updated := lenses.PasswordO.Set(__option.None[string]())(userinfo)
		assert.Equal(t, "john", updated.Username())
		_, ok := updated.Password()
		assert.False(t, ok)
	})

	t.Run("Set None password on user without password", func(t *testing.T) {
		userinfo := url.User("alice")
		updated := lenses.PasswordO.Set(__option.None[string]())(userinfo)
		assert.Equal(t, "alice", updated.Username())
		_, ok := updated.Password()
		assert.False(t, ok)
	})
}

// TestUserinfoRefLenses_Composition tests composing lens operations
func TestUserinfoRefLenses_Composition(t *testing.T) {
	lenses := MakeUserinfoRefLenses()

	t.Run("Update both username and password", func(t *testing.T) {
		userinfo := url.UserPassword("john", "secret123")

		// Update username first
		updated1 := lenses.Username.Set("alice")(userinfo)
		// Then update password
		updated2 := lenses.Password.Set("newpass")(updated1)

		assert.Equal(t, "alice", updated2.Username())
		pwd, ok := updated2.Password()
		assert.True(t, ok)
		assert.Equal(t, "newpass", pwd)
	})

	t.Run("Add password to user without password", func(t *testing.T) {
		userinfo := url.User("bob")
		updated := lenses.PasswordO.Set(__option.Some("pass123"))(userinfo)

		assert.Equal(t, "bob", updated.Username())
		pwd, ok := updated.Password()
		assert.True(t, ok)
		assert.Equal(t, "pass123", pwd)
	})

	t.Run("Remove password from user with password", func(t *testing.T) {
		userinfo := url.UserPassword("charlie", "oldpass")
		updated := lenses.PasswordO.Set(__option.None[string]())(userinfo)

		assert.Equal(t, "charlie", updated.Username())
		_, ok := updated.Password()
		assert.False(t, ok)
	})
}

// TestUserinfoRefLenses_EdgeCases tests edge cases
func TestUserinfoRefLenses_EdgeCases(t *testing.T) {
	lenses := MakeUserinfoRefLenses()

	t.Run("Empty username and password", func(t *testing.T) {
		userinfo := url.UserPassword("", "")

		username := lenses.Username.Get(userinfo)
		assert.Equal(t, "", username)

		password := lenses.Password.Get(userinfo)
		assert.Equal(t, "", password)
	})

	t.Run("Special characters in username", func(t *testing.T) {
		userinfo := url.User("user@domain.com")
		username := lenses.Username.Get(userinfo)
		assert.Equal(t, "user@domain.com", username)

		updated := lenses.Username.Set("new@user.com")(userinfo)
		assert.Equal(t, "new@user.com", updated.Username())
	})

	t.Run("Special characters in password", func(t *testing.T) {
		userinfo := url.UserPassword("john", "p@$$w0rd!")
		password := lenses.Password.Get(userinfo)
		assert.Equal(t, "p@$$w0rd!", password)

		updated := lenses.Password.Set("n3w!p@ss")(userinfo)
		pwd, ok := updated.Password()
		assert.True(t, ok)
		assert.Equal(t, "n3w!p@ss", pwd)
	})

	t.Run("Very long username", func(t *testing.T) {
		longUsername := "verylongusernamethatexceedsnormallengthbutshouldbehanded"
		userinfo := url.User(longUsername)
		username := lenses.Username.Get(userinfo)
		assert.Equal(t, longUsername, username)
	})

	t.Run("Very long password", func(t *testing.T) {
		longPassword := "verylongpasswordthatexceedsnormallengthbutshouldbehanded"
		userinfo := url.UserPassword("john", longPassword)
		password := lenses.Password.Get(userinfo)
		assert.Equal(t, longPassword, password)
	})
}

// TestUserinfoRefLenses_Immutability tests that operations return new instances
func TestUserinfoRefLenses_Immutability(t *testing.T) {
	lenses := MakeUserinfoRefLenses()

	t.Run("Setting username returns new instance", func(t *testing.T) {
		original := url.User("john")
		updated := lenses.Username.Set("alice")(original)

		// Original should be unchanged
		assert.Equal(t, "john", original.Username())
		// Updated should have new value
		assert.Equal(t, "alice", updated.Username())
		// Should be different instances
		assert.NotSame(t, original, updated)
	})

	t.Run("Setting password returns new instance", func(t *testing.T) {
		original := url.UserPassword("john", "pass1")
		updated := lenses.Password.Set("pass2")(original)

		// Original should be unchanged
		pwd, _ := original.Password()
		assert.Equal(t, "pass1", pwd)
		// Updated should have new value
		pwd, _ = updated.Password()
		assert.Equal(t, "pass2", pwd)
		// Should be different instances
		assert.NotSame(t, original, updated)
	})

	t.Run("Multiple updates create new instances", func(t *testing.T) {
		original := url.UserPassword("john", "pass1")
		updated1 := lenses.Username.Set("alice")(original)
		updated2 := lenses.Password.Set("pass2")(updated1)

		// Original unchanged
		assert.Equal(t, "john", original.Username())
		pwd, _ := original.Password()
		assert.Equal(t, "pass1", pwd)

		// First update has new username, old password
		assert.Equal(t, "alice", updated1.Username())
		pwd, _ = updated1.Password()
		assert.Equal(t, "pass1", pwd)

		// Second update has new username and password
		assert.Equal(t, "alice", updated2.Username())
		pwd, _ = updated2.Password()
		assert.Equal(t, "pass2", pwd)

		// All different instances
		assert.NotSame(t, original, updated1)
		assert.NotSame(t, updated1, updated2)
		assert.NotSame(t, original, updated2)
	})
}

// TestUserinfoRefLenses_PasswordPresence tests password presence detection
func TestUserinfoRefLenses_PasswordPresence(t *testing.T) {
	lenses := MakeUserinfoRefLenses()

	t.Run("Distinguish between no password and empty password", func(t *testing.T) {
		// User with no password
		userNoPass := url.User("john")
		optNoPass := lenses.PasswordO.Get(userNoPass)
		assert.True(t, __option.IsNone(optNoPass))

		// User with empty password (still has password set)
		userEmptyPass := url.UserPassword("john", "")
		optEmptyPass := lenses.PasswordO.Get(userEmptyPass)
		// Empty password is still Some (password is set, just empty)
		assert.True(t, __option.IsSome(optEmptyPass))
		value := __option.GetOrElse(func() string { return "default" })(optEmptyPass)
		assert.Equal(t, "", value)
	})

	t.Run("Password lens returns empty string for no password", func(t *testing.T) {
		userinfo := url.User("john")
		password := lenses.Password.Get(userinfo)
		assert.Equal(t, "", password)

		// But PasswordO returns None
		opt := lenses.PasswordO.Get(userinfo)
		assert.True(t, __option.IsNone(opt))
	})
}

// TestErrorLenses tests lenses for url.Error
func TestErrorLenses(t *testing.T) {
	lenses := MakeErrorLenses()

	t.Run("Get and Set Op field", func(t *testing.T) {
		urlErr := url.Error{
			Op:  "Get",
			URL: "https://example.com",
			Err: assert.AnError,
		}

		// Test Get
		op := lenses.Op.Get(urlErr)
		assert.Equal(t, "Get", op)

		// Test Set (curried, returns new Error)
		updated := lenses.Op.Set("Post")(urlErr)
		assert.Equal(t, "Post", updated.Op)
		assert.Equal(t, "Get", urlErr.Op) // Original unchanged
	})

	t.Run("Get and Set URL field", func(t *testing.T) {
		urlErr := url.Error{
			Op:  "Get",
			URL: "https://example.com",
			Err: assert.AnError,
		}

		// Test Get
		urlStr := lenses.URL.Get(urlErr)
		assert.Equal(t, "https://example.com", urlStr)

		// Test Set (curried)
		updated := lenses.URL.Set("https://newsite.com")(urlErr)
		assert.Equal(t, "https://newsite.com", updated.URL)
		assert.Equal(t, "https://example.com", urlErr.URL) // Original unchanged
	})

	t.Run("Get and Set Err field", func(t *testing.T) {
		originalErr := assert.AnError
		urlErr := url.Error{
			Op:  "Get",
			URL: "https://example.com",
			Err: originalErr,
		}

		// Test Get
		err := lenses.Err.Get(urlErr)
		assert.Equal(t, originalErr, err)

		// Test Set (curried)
		newErr := errors.New("new error")
		updated := lenses.Err.Set(newErr)(urlErr)
		assert.Equal(t, newErr, updated.Err)
		assert.Equal(t, originalErr, urlErr.Err) // Original unchanged
	})

	t.Run("Optional lenses", func(t *testing.T) {
		urlErr := url.Error{
			Op:  "Get",
			URL: "https://example.com",
			Err: assert.AnError,
		}

		// Test OpO
		opOpt := lenses.OpO.Get(urlErr)
		assert.True(t, __option.IsSome(opOpt))

		// Test with empty Op
		emptyErr := url.Error{Op: "", URL: "test", Err: nil}
		opOpt = lenses.OpO.Get(emptyErr)
		assert.True(t, __option.IsNone(opOpt))

		// Test URLO
		urlOpt := lenses.URLO.Get(urlErr)
		assert.True(t, __option.IsSome(urlOpt))

		// Test ErrO
		errOpt := lenses.ErrO.Get(urlErr)
		assert.True(t, __option.IsSome(errOpt))

		// Test with nil error
		nilErrErr := url.Error{Op: "Get", URL: "test", Err: nil}
		errOpt = lenses.ErrO.Get(nilErrErr)
		assert.True(t, __option.IsNone(errOpt))
	})
}

// TestErrorRefLenses tests reference lenses for url.Error
func TestErrorRefLenses(t *testing.T) {
	lenses := MakeErrorRefLenses()

	t.Run("Get and Set creates new instance", func(t *testing.T) {
		urlErr := &url.Error{
			Op:  "Get",
			URL: "https://example.com",
			Err: assert.AnError,
		}

		// Test Get
		op := lenses.Op.Get(urlErr)
		assert.Equal(t, "Get", op)

		// Test Set (creates copy)
		updated := lenses.Op.Set("Post")(urlErr)
		assert.Equal(t, "Post", updated.Op)
		assert.Equal(t, "Get", urlErr.Op) // Original unchanged
		assert.NotSame(t, urlErr, updated)
	})
}

// TestURLLenses tests lenses for url.URL
func TestURLLenses(t *testing.T) {
	lenses := MakeURLLenses()

	t.Run("Get and Set Scheme", func(t *testing.T) {
		u := url.URL{Scheme: "https", Host: "example.com"}

		scheme := lenses.Scheme.Get(u)
		assert.Equal(t, "https", scheme)

		updated := lenses.Scheme.Set("http")(u)
		assert.Equal(t, "http", updated.Scheme)
		assert.Equal(t, "https", u.Scheme) // Original unchanged
	})

	t.Run("Get and Set Host", func(t *testing.T) {
		u := url.URL{Scheme: "https", Host: "example.com"}

		host := lenses.Host.Get(u)
		assert.Equal(t, "example.com", host)

		updated := lenses.Host.Set("newsite.com")(u)
		assert.Equal(t, "newsite.com", updated.Host)
		assert.Equal(t, "example.com", u.Host)
	})

	t.Run("Get and Set Path", func(t *testing.T) {
		u := url.URL{Scheme: "https", Host: "example.com", Path: "/api/v1"}

		path := lenses.Path.Get(u)
		assert.Equal(t, "/api/v1", path)

		updated := lenses.Path.Set("/api/v2")(u)
		assert.Equal(t, "/api/v2", updated.Path)
		assert.Equal(t, "/api/v1", u.Path)
	})

	t.Run("Get and Set RawQuery", func(t *testing.T) {
		u := url.URL{Scheme: "https", Host: "example.com", RawQuery: "page=1"}

		query := lenses.RawQuery.Get(u)
		assert.Equal(t, "page=1", query)

		updated := lenses.RawQuery.Set("page=2&limit=10")(u)
		assert.Equal(t, "page=2&limit=10", updated.RawQuery)
		assert.Equal(t, "page=1", u.RawQuery)
	})

	t.Run("Get and Set Fragment", func(t *testing.T) {
		u := url.URL{Scheme: "https", Host: "example.com", Fragment: "section1"}

		fragment := lenses.Fragment.Get(u)
		assert.Equal(t, "section1", fragment)

		updated := lenses.Fragment.Set("section2")(u)
		assert.Equal(t, "section2", updated.Fragment)
		assert.Equal(t, "section1", u.Fragment)
	})

	t.Run("Get and Set User", func(t *testing.T) {
		userinfo := url.User("john")
		u := url.URL{Scheme: "https", Host: "example.com", User: userinfo}

		user := lenses.User.Get(u)
		assert.Equal(t, userinfo, user)

		newUser := url.UserPassword("alice", "pass")
		updated := lenses.User.Set(newUser)(u)
		assert.Equal(t, newUser, updated.User)
		assert.Equal(t, userinfo, u.User)
	})

	t.Run("Boolean fields", func(t *testing.T) {
		u := url.URL{Scheme: "https", Host: "example.com", ForceQuery: true}

		forceQuery := lenses.ForceQuery.Get(u)
		assert.True(t, forceQuery)

		updated := lenses.ForceQuery.Set(false)(u)
		assert.False(t, updated.ForceQuery)
		assert.True(t, u.ForceQuery)
	})

	t.Run("Optional lenses", func(t *testing.T) {
		u := url.URL{Scheme: "https", Host: "example.com", Path: "/test"}

		// Non-empty scheme
		schemeOpt := lenses.SchemeO.Get(u)
		assert.True(t, __option.IsSome(schemeOpt))

		// Empty RawQuery
		queryOpt := lenses.RawQueryO.Get(u)
		assert.True(t, __option.IsNone(queryOpt))

		// Set Some
		withQuery := lenses.RawQueryO.Set(__option.Some("q=test"))(u)
		assert.Equal(t, "q=test", withQuery.RawQuery)

		// Set None
		cleared := lenses.RawQueryO.Set(__option.None[string]())(withQuery)
		assert.Equal(t, "", cleared.RawQuery)
	})
}

// TestURLRefLenses tests reference lenses for url.URL
func TestURLRefLenses(t *testing.T) {
	lenses := MakeURLRefLenses()

	t.Run("Creates new instances", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "example.com", Path: "/api"}

		// Test Get
		scheme := lenses.Scheme.Get(u)
		assert.Equal(t, "https", scheme)

		// Test Set (creates copy)
		updated := lenses.Scheme.Set("http")(u)
		assert.Equal(t, "http", updated.Scheme)
		assert.Equal(t, "https", u.Scheme) // Original unchanged
		assert.NotSame(t, u, updated)
	})

	t.Run("Multiple field updates", func(t *testing.T) {
		u := &url.URL{Scheme: "https", Host: "example.com"}

		updated1 := lenses.Path.Set("/api/v1")(u)
		updated2 := lenses.RawQuery.Set("page=1")(updated1)
		updated3 := lenses.Fragment.Set("top")(updated2)

		// Original unchanged
		assert.Equal(t, "", u.Path)
		assert.Equal(t, "", u.RawQuery)
		assert.Equal(t, "", u.Fragment)

		// Final result has all updates
		assert.Equal(t, "/api/v1", updated3.Path)
		assert.Equal(t, "page=1", updated3.RawQuery)
		assert.Equal(t, "top", updated3.Fragment)
	})
}

// TestURLLenses_ComplexScenarios tests complex URL manipulation scenarios
func TestURLLenses_ComplexScenarios(t *testing.T) {
	lenses := MakeURLLenses()

	t.Run("Build URL incrementally", func(t *testing.T) {
		u := url.URL{}

		u = lenses.Scheme.Set("https")(u)
		u = lenses.Host.Set("api.example.com")(u)
		u = lenses.Path.Set("/v1/users")(u)
		u = lenses.RawQuery.Set("limit=10&offset=0")(u)
		u = lenses.Fragment.Set("results")(u)

		assert.Equal(t, "https", u.Scheme)
		assert.Equal(t, "api.example.com", u.Host)
		assert.Equal(t, "/v1/users", u.Path)
		assert.Equal(t, "limit=10&offset=0", u.RawQuery)
		assert.Equal(t, "results", u.Fragment)
	})

	t.Run("Update URL with authentication", func(t *testing.T) {
		u := url.URL{
			Scheme: "https",
			Host:   "example.com",
			Path:   "/api",
		}

		userinfo := url.UserPassword("admin", "secret")
		updated := lenses.User.Set(userinfo)(u)

		assert.NotNil(t, updated.User)
		assert.Equal(t, "admin", updated.User.Username())
		pwd, ok := updated.User.Password()
		assert.True(t, ok)
		assert.Equal(t, "secret", pwd)
	})

	t.Run("Clear optional fields", func(t *testing.T) {
		u := url.URL{
			Scheme:   "https",
			Host:     "example.com",
			Path:     "/api",
			RawQuery: "page=1",
			Fragment: "top",
		}

		// Clear query and fragment
		u = lenses.RawQueryO.Set(option.None[string]())(u)
		u = lenses.FragmentO.Set(option.None[string]())(u)

		assert.Equal(t, "", u.RawQuery)
		assert.Equal(t, "", u.Fragment)
		assert.Equal(t, "https", u.Scheme) // Other fields unchanged
		assert.Equal(t, "example.com", u.Host)
	})
}
