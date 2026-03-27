package authtoken

import (
	"time"

	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
	IORef "github.com/IBM/fp-go/v2/ioref"
	O "github.com/IBM/fp-go/v2/option"
)

// Cache is a generic thread-safe cache that stores an optional value using IORef.
type Cache[T any] struct {
	ref IORef.IORef[Option[T]]
}

// MakeCache creates a new empty cache for type T.
func MakeCache[T any]() IO.IO[*Cache[T]] {
	return F.Pipe1(
		IORef.MakeIORef(O.None[T]()),
		IO.Map(func(ref IORef.IORef[Option[T]]) *Cache[T] {
			return &Cache[T]{ref: ref}
		}),
	)
}

// Get retrieves the currently cached value, if any.
func (c *Cache[T]) Get() IO.IO[Option[T]] {
	return IORef.Read(c.ref)
}

// Set stores a new value in the cache, replacing any existing value.
func (c *Cache[T]) Set(value T) IO.IO[T] {
	return F.Pipe2(
		c.ref,
		IORef.Write(O.Some(value)),
		IO.Map(func(opt Option[T]) T {
			return value
		}),
	)
}

// Clear removes the cached value, setting the cache back to None.
func (c *Cache[T]) Clear() IO.IO[Option[T]] {
	return IORef.Write(O.None[T]())(c.ref)
}

// GetIf retrieves the cached value only if it satisfies the predicate.
func (c *Cache[T]) GetIf(predicate func(T) bool) IO.IO[Option[T]] {
	return F.Pipe1(
		c.Get(),
		IO.Map(O.Filter(predicate)),
	)
}

// Update applies a transformation function to the cached value.
func (c *Cache[T]) Update(f func(T) T) IO.IO[Option[T]] {
	return IORef.Modify(O.Map(f))(c.ref)
}

// TokenCache is a specialized cache for CachedToken with token-specific methods.
type TokenCache struct {
	cache *Cache[CachedToken]
}

// MakeTokenCache creates a new empty token cache.
func MakeTokenCache() IO.IO[*TokenCache] {
	return F.Pipe1(
		MakeCache[CachedToken](),
		IO.Map(func(cache *Cache[CachedToken]) *TokenCache {
			return &TokenCache{cache: cache}
		}),
	)
}

// Get retrieves the currently cached token, if any.
func (c *TokenCache) Get() IO.IO[Option[CachedToken]] {
	return c.cache.Get()
}

// Set stores a new token in the cache, replacing any existing token.
func (c *TokenCache) Set(token CachedToken) IO.IO[CachedToken] {
	return c.cache.Set(token)
}

// Clear removes the cached token, setting the cache back to None.
func (c *TokenCache) Clear() IO.IO[Option[CachedToken]] {
	return c.cache.Clear()
}

// IsExpired checks if the cached token is expired based on the current time.
func (c *TokenCache) IsExpired(currentTime time.Time) IO.IO[bool] {
	return F.Pipe1(
		c.Get(),
		IO.Map(O.Fold(
			F.Constant(true),
			func(token CachedToken) bool {
				return !currentTime.Before(token.ExpiresAt)
			},
		)),
	)
}

// GetIfValid retrieves the cached token only if it's not expired.
func (c *TokenCache) GetIfValid(currentTime time.Time) IO.IO[Option[CachedToken]] {
	return c.cache.GetIf(func(token CachedToken) bool {
		return currentTime.Before(token.ExpiresAt)
	})
}

// Update applies a transformation function to the cached token.
func (c *TokenCache) Update(f func(CachedToken) CachedToken) IO.IO[Option[CachedToken]] {
	return c.cache.Update(f)
}
