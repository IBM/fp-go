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
	"net/url"
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestHostname_ValueBased tests the Hostname lens for url.URL (value-based)
func TestHostname_ValueBased(t *testing.T) {
	lenses := MakeURLLenses()

	t.Run("get hostname from URL with port", func(t *testing.T) {
		u := url.URL{Host: "example.com:8080"}
		hostname := lenses.Hostname.Get(u)
		assert.Equal(t, "example.com", hostname)
	})

	t.Run("get hostname from URL without port", func(t *testing.T) {
		u := url.URL{Host: "example.com"}
		hostname := lenses.Hostname.Get(u)
		assert.Equal(t, "example.com", hostname)
	})

	t.Run("get hostname from URL with IPv6 and port", func(t *testing.T) {
		u := url.URL{Host: "[::1]:8080"}
		hostname := lenses.Hostname.Get(u)
		assert.Equal(t, "::1", hostname)
	})

	t.Run("get hostname from URL with IPv6 without port", func(t *testing.T) {
		u := url.URL{Host: "[::1]"}
		hostname := lenses.Hostname.Get(u)
		assert.Equal(t, "[::1]", hostname)
	})

	t.Run("get hostname from empty Host", func(t *testing.T) {
		u := url.URL{Host: ""}
		hostname := lenses.Hostname.Get(u)
		assert.Equal(t, "", hostname)
	})

	t.Run("set hostname preserves port", func(t *testing.T) {
		u := url.URL{Host: "example.com:8080"}
		updated := lenses.Hostname.Set("newhost.com")(u)
		assert.Equal(t, "newhost.com:8080", updated.Host)
		assert.Equal(t, "example.com:8080", u.Host) // original unchanged
	})

	t.Run("set hostname when no port exists", func(t *testing.T) {
		u := url.URL{Host: "example.com"}
		updated := lenses.Hostname.Set("newhost.com")(u)
		assert.Equal(t, "newhost.com", updated.Host)
	})

	t.Run("set hostname with IPv6 preserves port", func(t *testing.T) {
		u := url.URL{Host: "[::1]:8080"}
		updated := lenses.Hostname.Set("::2")(u)
		assert.Equal(t, "[::2]:8080", updated.Host)
	})

	t.Run("set hostname on empty Host", func(t *testing.T) {
		u := url.URL{Host: ""}
		updated := lenses.Hostname.Set("example.com")(u)
		assert.Equal(t, "example.com", updated.Host)
	})

	t.Run("set empty hostname", func(t *testing.T) {
		u := url.URL{Host: "example.com:8080"}
		updated := lenses.Hostname.Set("")(u)
		assert.Equal(t, ":8080", updated.Host)
	})
}

// TestHostname_ReferenceBased tests the Hostname lens for *url.URL (reference-based)
func TestHostname_ReferenceBased(t *testing.T) {
	lenses := MakeURLRefLenses()

	t.Run("get hostname from URL with port", func(t *testing.T) {
		u := &url.URL{Host: "example.com:8080"}
		hostname := lenses.Hostname.Get(u)
		assert.Equal(t, "example.com", hostname)
	})

	t.Run("get hostname from URL without port", func(t *testing.T) {
		u := &url.URL{Host: "example.com"}
		hostname := lenses.Hostname.Get(u)
		assert.Equal(t, "example.com", hostname)
	})

	t.Run("get hostname from URL with IPv6 and port", func(t *testing.T) {
		u := &url.URL{Host: "[::1]:8080"}
		hostname := lenses.Hostname.Get(u)
		assert.Equal(t, "::1", hostname)
	})

	t.Run("get hostname from URL with IPv6 without port", func(t *testing.T) {
		u := &url.URL{Host: "[::1]"}
		hostname := lenses.Hostname.Get(u)
		assert.Equal(t, "[::1]", hostname)
	})

	t.Run("set hostname preserves port", func(t *testing.T) {
		u := &url.URL{Host: "example.com:8080"}
		original := u.Host
		updated := lenses.Hostname.Set("newhost.com")(u)
		assert.Equal(t, "newhost.com:8080", updated.Host)
		assert.NotSame(t, u, updated)     // reference-based creates a copy for safety
		assert.Equal(t, original, u.Host) // original unchanged
	})

	t.Run("set hostname when no port exists", func(t *testing.T) {
		u := &url.URL{Host: "example.com"}
		updated := lenses.Hostname.Set("newhost.com")(u)
		assert.Equal(t, "newhost.com", updated.Host)
		assert.NotSame(t, u, updated)          // creates a copy
		assert.Equal(t, "example.com", u.Host) // original unchanged
	})

	t.Run("set hostname with IPv6 preserves port", func(t *testing.T) {
		u := &url.URL{Host: "[::1]:8080"}
		updated := lenses.Hostname.Set("::2")(u)
		assert.Equal(t, "[::2]:8080", updated.Host)
	})

	t.Run("set hostname on empty Host", func(t *testing.T) {
		u := &url.URL{Host: ""}
		updated := lenses.Hostname.Set("example.com")(u)
		assert.Equal(t, "example.com", updated.Host)
	})
}

// TestHostnameO_ValueBased tests the HostnameO optional lens for url.URL
func TestHostnameO_ValueBased(t *testing.T) {
	lenses := MakeURLLenses()

	t.Run("get Some hostname from URL with port", func(t *testing.T) {
		u := url.URL{Host: "example.com:8080"}
		hostname := lenses.HostnameO.Get(u)
		assert.Equal(t, O.Some("example.com"), hostname)
	})

	t.Run("get Some hostname from URL without port", func(t *testing.T) {
		u := url.URL{Host: "example.com"}
		hostname := lenses.HostnameO.Get(u)
		assert.Equal(t, O.Some("example.com"), hostname)
	})

	t.Run("get None from empty Host", func(t *testing.T) {
		u := url.URL{Host: ""}
		hostname := lenses.HostnameO.Get(u)
		assert.Equal(t, O.None[string](), hostname)
	})

	t.Run("set Some hostname preserves port", func(t *testing.T) {
		u := url.URL{Host: "example.com:8080"}
		updated := lenses.HostnameO.Set(O.Some("newhost.com"))(u)
		assert.Equal(t, "newhost.com:8080", updated.Host)
	})

	t.Run("set None clears hostname but preserves port", func(t *testing.T) {
		u := url.URL{Host: "example.com:8080"}
		updated := lenses.HostnameO.Set(O.None[string]())(u)
		assert.Equal(t, ":8080", updated.Host)
	})

	t.Run("set Some hostname on empty Host", func(t *testing.T) {
		u := url.URL{Host: ""}
		updated := lenses.HostnameO.Set(O.Some("example.com"))(u)
		assert.Equal(t, "example.com", updated.Host)
	})

	t.Run("set None on empty Host", func(t *testing.T) {
		u := url.URL{Host: ""}
		updated := lenses.HostnameO.Set(O.None[string]())(u)
		assert.Equal(t, "", updated.Host)
	})
}

// TestHostnameO_ReferenceBased tests the HostnameO optional lens for *url.URL
func TestHostnameO_ReferenceBased(t *testing.T) {
	lenses := MakeURLRefLenses()

	t.Run("get Some hostname from URL with port", func(t *testing.T) {
		u := &url.URL{Host: "example.com:8080"}
		hostname := lenses.HostnameO.Get(u)
		assert.Equal(t, O.Some("example.com"), hostname)
	})

	t.Run("get Some hostname from URL without port", func(t *testing.T) {
		u := &url.URL{Host: "example.com"}
		hostname := lenses.HostnameO.Get(u)
		assert.Equal(t, O.Some("example.com"), hostname)
	})

	t.Run("get None from empty Host", func(t *testing.T) {
		u := &url.URL{Host: ""}
		hostname := lenses.HostnameO.Get(u)
		assert.Equal(t, O.None[string](), hostname)
	})

	t.Run("set Some hostname preserves port", func(t *testing.T) {
		u := &url.URL{Host: "example.com:8080"}
		updated := lenses.HostnameO.Set(O.Some("newhost.com"))(u)
		assert.Equal(t, "newhost.com:8080", updated.Host)
		assert.NotSame(t, u, updated)               // creates a copy
		assert.Equal(t, "example.com:8080", u.Host) // original unchanged
	})

	t.Run("set None clears hostname but preserves port", func(t *testing.T) {
		u := &url.URL{Host: "example.com:8080"}
		updated := lenses.HostnameO.Set(O.None[string]())(u)
		assert.Equal(t, ":8080", updated.Host)
	})

	t.Run("set Some hostname on empty Host", func(t *testing.T) {
		u := &url.URL{Host: ""}
		updated := lenses.HostnameO.Set(O.Some("example.com"))(u)
		assert.Equal(t, "example.com", updated.Host)
	})
}

// TestPort_ValueBased tests the Port lens for url.URL (value-based)
func TestPort_ValueBased(t *testing.T) {
	lenses := MakeURLLenses()

	t.Run("get port from URL with port", func(t *testing.T) {
		u := url.URL{Host: "example.com:8080"}
		port := lenses.Port.Get(u)
		assert.Equal(t, "8080", port)
	})

	t.Run("get empty port from URL without port", func(t *testing.T) {
		u := url.URL{Host: "example.com"}
		port := lenses.Port.Get(u)
		assert.Equal(t, "", port)
	})

	t.Run("get port from IPv6 with port", func(t *testing.T) {
		u := url.URL{Host: "[::1]:8080"}
		port := lenses.Port.Get(u)
		assert.Equal(t, "8080", port)
	})

	t.Run("get empty port from IPv6 without port", func(t *testing.T) {
		u := url.URL{Host: "[::1]"}
		port := lenses.Port.Get(u)
		assert.Equal(t, "", port)
	})

	t.Run("get empty port from empty Host", func(t *testing.T) {
		u := url.URL{Host: ""}
		port := lenses.Port.Get(u)
		assert.Equal(t, "", port)
	})

	t.Run("set port preserves hostname", func(t *testing.T) {
		u := url.URL{Host: "example.com:8080"}
		updated := lenses.Port.Set("9090")(u)
		assert.Equal(t, "example.com:9090", updated.Host)
		assert.Equal(t, "example.com:8080", u.Host) // original unchanged
	})

	t.Run("set port when no port exists", func(t *testing.T) {
		u := url.URL{Host: "example.com"}
		updated := lenses.Port.Set("8080")(u)
		assert.Equal(t, "example.com:8080", updated.Host)
	})

	t.Run("set port with IPv6 preserves hostname", func(t *testing.T) {
		u := url.URL{Host: "[::1]:8080"}
		updated := lenses.Port.Set("9090")(u)
		assert.Equal(t, "[::1]:9090", updated.Host)
	})

	t.Run("set port on empty Host", func(t *testing.T) {
		u := url.URL{Host: ""}
		updated := lenses.Port.Set("8080")(u)
		assert.Equal(t, ":8080", updated.Host)
	})

	t.Run("set empty port removes port", func(t *testing.T) {
		u := url.URL{Host: "example.com:8080"}
		updated := lenses.Port.Set("")(u)
		assert.Equal(t, "example.com:", updated.Host)
	})

	t.Run("set port with hostname only", func(t *testing.T) {
		u := url.URL{Host: "example.com"}
		updated := lenses.Port.Set("443")(u)
		assert.Equal(t, "example.com:443", updated.Host)
	})
}

// TestPort_ReferenceBased tests the Port lens for *url.URL (reference-based)
func TestPort_ReferenceBased(t *testing.T) {
	lenses := MakeURLRefLenses()

	t.Run("get port from URL with port", func(t *testing.T) {
		u := &url.URL{Host: "example.com:8080"}
		port := lenses.Port.Get(u)
		assert.Equal(t, "8080", port)
	})

	t.Run("get empty port from URL without port", func(t *testing.T) {
		u := &url.URL{Host: "example.com"}
		port := lenses.Port.Get(u)
		assert.Equal(t, "", port)
	})

	t.Run("get port from IPv6 with port", func(t *testing.T) {
		u := &url.URL{Host: "[::1]:8080"}
		port := lenses.Port.Get(u)
		assert.Equal(t, "8080", port)
	})

	t.Run("get empty port from IPv6 without port", func(t *testing.T) {
		u := &url.URL{Host: "[::1]"}
		port := lenses.Port.Get(u)
		assert.Equal(t, "", port)
	})

	t.Run("set port preserves hostname", func(t *testing.T) {
		u := &url.URL{Host: "example.com:8080"}
		original := u.Host
		updated := lenses.Port.Set("9090")(u)
		assert.Equal(t, "example.com:9090", updated.Host)
		assert.NotSame(t, u, updated)     // reference-based creates a copy for safety
		assert.Equal(t, original, u.Host) // original unchanged
	})

	t.Run("set port when no port exists", func(t *testing.T) {
		u := &url.URL{Host: "example.com"}
		updated := lenses.Port.Set("8080")(u)
		assert.Equal(t, "example.com:8080", updated.Host)
		assert.NotSame(t, u, updated)          // creates a copy
		assert.Equal(t, "example.com", u.Host) // original unchanged
	})

	t.Run("set port with IPv6 preserves hostname", func(t *testing.T) {
		u := &url.URL{Host: "[::1]:8080"}
		updated := lenses.Port.Set("9090")(u)
		assert.Equal(t, "[::1]:9090", updated.Host)
	})

	t.Run("set port on empty Host", func(t *testing.T) {
		u := &url.URL{Host: ""}
		updated := lenses.Port.Set("8080")(u)
		assert.Equal(t, ":8080", updated.Host)
	})
}

// TestPortO_ValueBased tests the PortO optional lens for url.URL
func TestPortO_ValueBased(t *testing.T) {
	lenses := MakeURLLenses()

	t.Run("get Some port from URL with port", func(t *testing.T) {
		u := url.URL{Host: "example.com:8080"}
		port := lenses.PortO.Get(u)
		assert.Equal(t, O.Some("8080"), port)
	})

	t.Run("get None from URL without port", func(t *testing.T) {
		u := url.URL{Host: "example.com"}
		port := lenses.PortO.Get(u)
		assert.Equal(t, O.None[string](), port)
	})

	t.Run("get None from empty Host", func(t *testing.T) {
		u := url.URL{Host: ""}
		port := lenses.PortO.Get(u)
		assert.Equal(t, O.None[string](), port)
	})

	t.Run("get Some port from IPv6 with port", func(t *testing.T) {
		u := url.URL{Host: "[::1]:8080"}
		port := lenses.PortO.Get(u)
		assert.Equal(t, O.Some("8080"), port)
	})

	t.Run("set Some port preserves hostname", func(t *testing.T) {
		u := url.URL{Host: "example.com:8080"}
		updated := lenses.PortO.Set(O.Some("9090"))(u)
		assert.Equal(t, "example.com:9090", updated.Host)
	})

	t.Run("set None removes port but preserves hostname", func(t *testing.T) {
		u := url.URL{Host: "example.com:8080"}
		updated := lenses.PortO.Set(O.None[string]())(u)
		assert.Equal(t, "example.com:", updated.Host)
	})

	t.Run("set Some port on URL without port", func(t *testing.T) {
		u := url.URL{Host: "example.com"}
		updated := lenses.PortO.Set(O.Some("8080"))(u)
		assert.Equal(t, "example.com:8080", updated.Host)
	})

	t.Run("set None on URL without port", func(t *testing.T) {
		u := url.URL{Host: "example.com"}
		updated := lenses.PortO.Set(O.None[string]())(u)
		assert.Equal(t, "example.com:", updated.Host)
	})

	t.Run("set Some port on empty Host", func(t *testing.T) {
		u := url.URL{Host: ""}
		updated := lenses.PortO.Set(O.Some("8080"))(u)
		assert.Equal(t, ":8080", updated.Host)
	})
}

// TestPortO_ReferenceBased tests the PortO optional lens for *url.URL
func TestPortO_ReferenceBased(t *testing.T) {
	lenses := MakeURLRefLenses()

	t.Run("get Some port from URL with port", func(t *testing.T) {
		u := &url.URL{Host: "example.com:8080"}
		port := lenses.PortO.Get(u)
		assert.Equal(t, O.Some("8080"), port)
	})

	t.Run("get None from URL without port", func(t *testing.T) {
		u := &url.URL{Host: "example.com"}
		port := lenses.PortO.Get(u)
		assert.Equal(t, O.None[string](), port)
	})

	t.Run("get None from empty Host", func(t *testing.T) {
		u := &url.URL{Host: ""}
		port := lenses.PortO.Get(u)
		assert.Equal(t, O.None[string](), port)
	})

	t.Run("get Some port from IPv6 with port", func(t *testing.T) {
		u := &url.URL{Host: "[::1]:8080"}
		port := lenses.PortO.Get(u)
		assert.Equal(t, O.Some("8080"), port)
	})

	t.Run("set Some port preserves hostname", func(t *testing.T) {
		u := &url.URL{Host: "example.com:8080"}
		updated := lenses.PortO.Set(O.Some("9090"))(u)
		assert.Equal(t, "example.com:9090", updated.Host)
		assert.NotSame(t, u, updated)               // creates a copy
		assert.Equal(t, "example.com:8080", u.Host) // original unchanged
	})

	t.Run("set None removes port but preserves hostname", func(t *testing.T) {
		u := &url.URL{Host: "example.com:8080"}
		updated := lenses.PortO.Set(O.None[string]())(u)
		assert.Equal(t, "example.com:", updated.Host)
	})

	t.Run("set Some port on URL without port", func(t *testing.T) {
		u := &url.URL{Host: "example.com"}
		updated := lenses.PortO.Set(O.Some("8080"))(u)
		assert.Equal(t, "example.com:8080", updated.Host)
	})

	t.Run("set Some port on empty Host", func(t *testing.T) {
		u := &url.URL{Host: ""}
		updated := lenses.PortO.Set(O.Some("8080"))(u)
		assert.Equal(t, ":8080", updated.Host)
	})
}

// TestHostnamePort_Integration tests interaction between Hostname and Port lenses
func TestHostnamePort_Integration(t *testing.T) {
	lenses := MakeURLLenses()

	t.Run("set hostname then port", func(t *testing.T) {
		u := url.URL{Host: "example.com:8080"}
		updated := lenses.Port.Set("9090")(lenses.Hostname.Set("newhost.com")(u))
		assert.Equal(t, "newhost.com:9090", updated.Host)
	})

	t.Run("set port then hostname", func(t *testing.T) {
		u := url.URL{Host: "example.com:8080"}
		updated := lenses.Hostname.Set("newhost.com")(lenses.Port.Set("9090")(u))
		assert.Equal(t, "newhost.com:9090", updated.Host)
	})

	t.Run("clear hostname preserves port", func(t *testing.T) {
		u := url.URL{Host: "example.com:8080"}
		updated := lenses.Hostname.Set("")(u)
		assert.Equal(t, ":8080", updated.Host)
		assert.Equal(t, "8080", lenses.Port.Get(updated))
	})

	t.Run("clear port preserves hostname", func(t *testing.T) {
		u := url.URL{Host: "example.com:8080"}
		updated := lenses.Port.Set("")(u)
		assert.Equal(t, "example.com:", updated.Host)
		assert.Equal(t, "example.com", lenses.Hostname.Get(updated))
	})
}

// TestHostnamePort_EdgeCases tests edge cases for Hostname and Port lenses
func TestHostnamePort_EdgeCases(t *testing.T) {
	lenses := MakeURLLenses()

	t.Run("hostname with special characters", func(t *testing.T) {
		u := url.URL{Host: "sub-domain.example.com:8080"}
		hostname := lenses.Hostname.Get(u)
		assert.Equal(t, "sub-domain.example.com", hostname)
	})

	t.Run("port with leading zeros", func(t *testing.T) {
		u := url.URL{Host: "example.com:0080"}
		port := lenses.Port.Get(u)
		assert.Equal(t, "0080", port)
	})

	t.Run("IPv4 address as hostname", func(t *testing.T) {
		u := url.URL{Host: "192.168.1.1:8080"}
		hostname := lenses.Hostname.Get(u)
		assert.Equal(t, "192.168.1.1", hostname)
	})

	t.Run("IPv6 loopback with port", func(t *testing.T) {
		u := url.URL{Host: "[::1]:8080"}
		hostname := lenses.Hostname.Get(u)
		port := lenses.Port.Get(u)
		assert.Equal(t, "::1", hostname)
		assert.Equal(t, "8080", port)
	})

	t.Run("IPv6 full address with port", func(t *testing.T) {
		u := url.URL{Host: "[2001:db8::1]:8080"}
		hostname := lenses.Hostname.Get(u)
		port := lenses.Port.Get(u)
		assert.Equal(t, "2001:db8::1", hostname)
		assert.Equal(t, "8080", port)
	})

	t.Run("set IPv6 hostname preserves port", func(t *testing.T) {
		u := url.URL{Host: "[::1]:8080"}
		updated := lenses.Hostname.Set("2001:db8::1")(u)
		assert.Equal(t, "[2001:db8::1]:8080", updated.Host)
	})

	t.Run("malformed host handled gracefully", func(t *testing.T) {
		u := url.URL{Host: "example.com:"}
		hostname := lenses.Hostname.Get(u)
		port := lenses.Port.Get(u)
		assert.Equal(t, "example.com", hostname)
		assert.Equal(t, "", port)
	})
}
