package readerio

import (
	"context"

	"github.com/IBM/fp-go/v2/reader"
	RIO "github.com/IBM/fp-go/v2/readerio"
)

// SequenceReader transforms a ReaderIO containing a Reader into a Reader containing a ReaderIO.
// This "flips" the nested structure, allowing you to provide the Reader's environment first,
// then get a ReaderIO that can be executed with a context.
//
// Type transformation:
//
//	From: ReaderIO[Reader[R, A]]
//	      = func(context.Context) func() func(R) A
//
//	To:   Reader[R, ReaderIO[A]]
//	      = func(R) func(context.Context) func() A
//
// This is useful for point-free style programming where you want to partially apply
// the Reader's environment before dealing with the context.
//
// Type Parameters:
//   - R: The environment type that the Reader depends on
//   - A: The value type
//
// Parameters:
//   - ma: A ReaderIO containing a Reader
//
// Returns:
//   - A Reader that produces a ReaderIO when given an environment
//
// Example:
//
//	type Config struct {
//	    Timeout int
//	}
//
//	// A computation that produces a Reader
//	getMultiplier := func(ctx context.Context) IO[func(Config) int] {
//	    return func() func(Config) int {
//	        return func(cfg Config) int {
//	            return cfg.Timeout * 2
//	        }
//	    }
//	}
//
//	// Sequence it to apply Config first
//	sequenced := SequenceReader[Config, int](getMultiplier)
//	cfg := Config{Timeout: 30}
//	result := sequenced(cfg)(t.Context())() // Returns 60
//
//go:inline
func SequenceReader[R, A any](ma ReaderIO[Reader[R, A]]) Reader[R, ReaderIO[A]] {
	return RIO.SequenceReader(ma)
}

// TraverseReader applies a Reader-based transformation to a ReaderIO, introducing a new environment dependency.
//
// This function takes a Reader-based Kleisli arrow and returns a function that can transform
// a ReaderIO. The result allows you to provide the Reader's environment (R) first, which then
// produces a ReaderIO that depends on the context.
//
// Type transformation:
//
//	From: ReaderIO[A]
//	      = func(context.Context) func() A
//
//	With: reader.Kleisli[R, A, B]
//	      = func(A) func(R) B
//
//	To:   func(ReaderIO[A]) func(R) ReaderIO[B]
//	      = func(ReaderIO[A]) func(R) func(context.Context) func() B
//
// This enables transforming values within a ReaderIO using environment-dependent logic.
//
// Type Parameters:
//   - R: The environment type that the Reader depends on
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - f: A Reader-based Kleisli arrow that transforms A to B using environment R
//
// Returns:
//   - A function that takes a ReaderIO[A] and returns a function from R to ReaderIO[B]
//
// Example:
//
//	type Config struct {
//	    Multiplier int
//	}
//
//	// A Reader-based transformation
//	multiply := func(x int) func(Config) int {
//	    return func(cfg Config) int {
//	        return x * cfg.Multiplier
//	    }
//	}
//
//	// Apply TraverseReader
//	traversed := TraverseReader[Config, int, int](multiply)
//	computation := Of(10)
//	result := traversed(computation)
//
//	// Provide Config to get final result
//	cfg := Config{Multiplier: 5}
//	finalResult := result(cfg)(t.Context())() // Returns 50
//
//go:inline
func TraverseReader[R, A, B any](
	f reader.Kleisli[R, A, B],
) func(ReaderIO[A]) Kleisli[R, B] {
	return RIO.TraverseReader[context.Context](f)
}
