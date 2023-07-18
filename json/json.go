package json

import (
	"encoding/json"

	E "github.com/ibm/fp-go/either"
)

// Unmarshal parses a JSON data structure from bytes
func Unmarshal[A any](data []byte) E.Either[error, A] {
	return E.TryCatchError(func() (A, error) {
		var result A
		err := json.Unmarshal(data, &result)
		return result, err
	})
}

// Marshal converts a data structure to json
func Marshal[A any](a A) E.Either[error, []byte] {
	return E.TryCatchError(func() ([]byte, error) {
		return json.Marshal(a)
	})
}
