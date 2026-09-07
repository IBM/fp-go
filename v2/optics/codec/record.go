package codec

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/optics/codec/decode"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/record"
)

// AtKey returns a codec that focuses on a single key in a record.
//
// The codec decodes by looking up the given key in the input record. If the
// key is present its value is returned as a successful validation. If the key
// is absent the decode fails with a descriptive message that includes the key
// name. Encoding wraps the value in a new single-entry record keyed on the
// same key.
//
// The key type K must be both comparable (so it can be used as a map key)
// and implement fmt.Stringer (so it can appear in human-readable error messages
// and in the codec name).
//
// Type Parameters:
//   - V: The value type stored at the key
//   - K: The key type; must be comparable and implement fmt.Stringer
//
// Parameters:
//   - key: The map key to focus on
//
// Returns:
//   - A Type[V, record.Record[K,V], record.Record[K,V]] that decodes by
//     looking up key in the input record and encodes by creating a
//     single-entry record for that key.
//
// See Also:
//   - record.Lookup: The underlying key lookup used during decoding
//   - record.Of: The underlying constructor used during encoding
func AtKey[V any, K interface {
	comparable
	fmt.Stringer
}](key K) Type[V, record.Record[K, V], record.Record[K, V]] {
	return MakeType(
		fmt.Sprintf("AtKey[%s]", key),
		Is[V](),
		F.Flow2(
			record.Lookup[V](key),
			decode.FromOption(
				lazy.Of(validation.FailureWithMessage[V](key, fmt.Sprintf("expected key %s is unavailable", key))),
			),
		),
		F.Bind1st(record.Of[K, V], key),
	)
}
