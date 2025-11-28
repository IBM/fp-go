package pair

import "github.com/IBM/fp-go/v2/function"

func MonadSequence[L, A, HKTA, HKTPA any](
	mmap func(HKTA, Kleisli[L, A, A]) HKTPA,
	fas Pair[L, HKTA],
) HKTPA {
	return mmap(Tail(fas), FromHead[A](Head(fas)))
}

func MonadTraverse[L, A, HKTA, HKTPA any](
	mmap func(HKTA, Kleisli[L, A, A]) HKTPA,
	f func(A) HKTA,
	fas Pair[L, A],
) HKTPA {
	return mmap(f(Tail(fas)), FromHead[A](Head(fas)))
}

func Sequence[L, A, HKTA, HKTPA any](
	mmap func(Kleisli[L, A, A]) func(HKTA) HKTPA,
) func(Pair[L, HKTA]) HKTPA {
	fh := function.Flow2(
		Head[L, HKTA],
		FromHead[A, L],
	)
	return func(fas Pair[L, HKTA]) HKTPA {
		return mmap(fh(fas))(Tail(fas))
	}
}

func Traverse[L, A, HKTA, HKTPA any](
	mmap func(Kleisli[L, A, A]) func(HKTA) HKTPA,
) func(func(A) HKTA) func(Pair[L, A]) HKTPA {
	fh := function.Flow2(
		Head[L, A],
		FromHead[A, L],
	)
	return func(f func(A) HKTA) func(Pair[L, A]) HKTPA {
		ft := function.Flow2(
			Tail[L, A],
			f,
		)
		return func(fas Pair[L, A]) HKTPA {
			return mmap(fh(fas))(ft(fas))
		}
	}
}
