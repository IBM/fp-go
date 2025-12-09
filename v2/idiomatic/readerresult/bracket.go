package readerresult

import "github.com/IBM/fp-go/v2/idiomatic/result"

// Bracket makes sure that a resource is cleaned up in the event of an error. The release action is called regardless of
// whether the body action returns and error or not.
func Bracket[
	R, A, B, ANY any](

	acquire Lazy[ReaderResult[R, A]],
	use Kleisli[R, A, B],
	release func(A, B, error) ReaderResult[R, ANY],
) ReaderResult[R, B] {
	return func(r R) (B, error) {
		// acquire the resource
		a, aerr := acquire()(r)
		if aerr != nil {
			return result.Left[B](aerr)
		}
		b, berr := use(a)(r)
		_, xerr := release(a, b, berr)(r)
		if berr != nil {
			return result.Left[B](berr)
		}
		if xerr != nil {
			return result.Left[B](xerr)
		}
		return result.Of(b)
	}
}

func WithResource[B, R, A, ANY any](
	onCreate Lazy[ReaderResult[R, A]],
	onRelease Kleisli[R, A, ANY],
) Kleisli[R, Kleisli[R, A, B], B] {
	return func(k Kleisli[R, A, B]) ReaderResult[R, B] {
		return func(r R) (B, error) {
			a, aerr := onCreate()(r)
			if aerr != nil {
				return result.Left[B](aerr)
			}
			b, berr := k(a)(r)
			_, xerr := onRelease(a)(r)
			if berr != nil {
				return result.Left[B](berr)
			}
			if xerr != nil {
				return result.Left[B](xerr)
			}
			return result.Of(b)
		}
	}
}
