package effect

import (
	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/result"
)

//go:inline
func Local[C1, C2, A any](acc Reader[C1, C2]) Kleisli[C1, Effect[C2, A], A] {
	return readerreaderioresult.Local[A](acc)
}

//go:inline
func Contramap[C1, C2, A any](acc Reader[C1, C2]) Kleisli[C1, Effect[C2, A], A] {
	return readerreaderioresult.Local[A](acc)
}

//go:inline
func LocalIOK[A, C1, C2 any](f io.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A] {
	return readerreaderioresult.LocalIOK[A](f)
}

//go:inline
func LocalIOResultK[A, C1, C2 any](f ioresult.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A] {
	return readerreaderioresult.LocalIOResultK[A](f)
}

//go:inline
func LocalResultK[A, C1, C2 any](f result.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A] {
	return readerreaderioresult.LocalResultK[A](f)
}

//go:inline
func LocalThunkK[A, C1, C2 any](f thunk.Kleisli[C2, C1]) func(Effect[C1, A]) Effect[C2, A] {
	return readerreaderioresult.LocalReaderIOResultK[A](f)
}

// LocalEffectK transforms the context of an Effect using an Effect-returning function.
// This is the most powerful context transformation function, allowing the transformation
// itself to be effectful (can fail, perform I/O, and access the outer context).
//
// LocalEffectK takes a Kleisli arrow that:
//   - Accepts the outer context C2
//   - Returns an Effect that produces the inner context C1
//   - Can fail with an error during context transformation
//   - Can perform I/O operations during transformation
//
// This is useful when:
//   - Context transformation requires I/O (e.g., loading config from a file)
//   - Context transformation can fail (e.g., validating or parsing context)
//   - Context transformation needs to access the outer context
//
// Type Parameters:
//   - A: The value type produced by the effect
//   - C1: The inner context type (required by the original effect)
//   - C2: The outer context type (provided to the transformed effect)
//
// Parameters:
//   - f: A Kleisli arrow (C2 -> Effect[C2, C1]) that transforms C2 to C1 effectfully
//
// Returns:
//   - A function that transforms Effect[C1, A] to Effect[C2, A]
//
// Example:
//
//	type DatabaseConfig struct {
//		ConnectionString string
//	}
//
//	type AppConfig struct {
//		ConfigPath string
//	}
//
//	// Effect that needs DatabaseConfig
//	dbEffect := effect.Of[DatabaseConfig, string]("query result")
//
//	// Transform AppConfig to DatabaseConfig effectfully
//	// (e.g., load config from file, which can fail)
//	loadConfig := func(app AppConfig) Effect[AppConfig, DatabaseConfig] {
//		return effect.Chain[AppConfig](func(_ AppConfig) Effect[AppConfig, DatabaseConfig] {
//			// Simulate loading config from file (can fail)
//			return effect.Of[AppConfig, DatabaseConfig](DatabaseConfig{
//				ConnectionString: "loaded from " + app.ConfigPath,
//			})
//		})(effect.Of[AppConfig, AppConfig](app))
//	}
//
//	// Apply the transformation
//	transform := effect.LocalEffectK[string, DatabaseConfig, AppConfig](loadConfig)
//	appEffect := transform(dbEffect)
//
//	// Run with AppConfig
//	ioResult := effect.Provide(AppConfig{ConfigPath: "/etc/app.conf"})(appEffect)
//	readerResult := effect.RunSync(ioResult)
//	result, err := readerResult(context.Background())
//
// Comparison with other Local functions:
//   - Local/Contramap: Pure context transformation (C2 -> C1)
//   - LocalIOK: IO-based transformation (C2 -> IO[C1])
//   - LocalIOResultK: IO with error handling (C2 -> IOResult[C1])
//   - LocalReaderIOResultK: Reader-based with IO and errors (C2 -> ReaderIOResult[C1])
//   - LocalEffectK: Full Effect transformation (C2 -> Effect[C2, C1])
//
//go:inline
func LocalEffectK[A, C1, C2 any](f Kleisli[C2, C2, C1]) func(Effect[C1, A]) Effect[C2, A] {
	return readerreaderioresult.LocalReaderReaderIOEitherK[A](f)
}
