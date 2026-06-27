package context

import (
	"context"

	"github.com/IBM/fp-go/v2/pair"
)

type (

	// Pair represents a tuple of two values of types A and B.
	// It is used to group two related values together.
	Pair[A, B any] = pair.Pair[A, B]

	// ContextCancel represents a pair of a cancel function and a context.
	// It is used in operations that create new contexts with cancellation capabilities.
	//
	// The first element is the CancelFunc that should be called to release resources.
	// The second element is the new Context that was created.
	ContextCancel = Pair[context.CancelFunc, context.Context]

	Reader[A any] = func(context.Context) A

	Kleisli[A, B any] = func(A) Reader[B]
)
