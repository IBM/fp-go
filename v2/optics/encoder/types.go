package encoder

import "github.com/IBM/fp-go/v2/reader"

type (
	Encoder[O, A any] = reader.Reader[A, O]
)
