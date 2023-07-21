package generic

import (
	"sync"
	"time"

	F "github.com/IBM/fp-go/function"
	C "github.com/IBM/fp-go/internal/chain"
)

// type IO[A any] = func() A

func MakeIO[GA ~func() A, A any](f func() A) GA {
	return f
}

func Of[GA ~func() A, A any](a A) GA {
	return MakeIO[GA](F.Constant(a))
}

func FromIO[GA ~func() A, A any](a GA) GA {
	return a
}

// FromImpure converts a side effect without a return value into a side effect that returns any
func FromImpure[GA ~func() any, IMP ~func()](f IMP) GA {
	return MakeIO[GA](func() any {
		f()
		return nil
	})
}

func MonadOf[GA ~func() A, A any](a A) GA {
	return MakeIO[GA](F.Constant(a))
}

func MonadMap[GA ~func() A, GB ~func() B, A, B any](fa GA, f func(A) B) GB {
	return MakeIO[GB](func() B {
		return F.Pipe1(fa(), f)
	})
}

func Map[GA ~func() A, GB ~func() B, A, B any](f func(A) B) func(GA) GB {
	return F.Bind2nd(MonadMap[GA, GB, A, B], f)
}

func MonadMapTo[GA ~func() A, GB ~func() B, A, B any](fa GA, b B) GB {
	return MonadMap[GA, GB](fa, F.Constant1[A](b))
}

func MapTo[GA ~func() A, GB ~func() B, A, B any](b B) func(GA) GB {
	return F.Bind2nd(MonadMapTo[GA, GB, A, B], b)
}

// MonadChain composes computations in sequence, using the return value of one computation to determine the next computation.
func MonadChain[GA ~func() A, GB ~func() B, A, B any](fa GA, f func(A) GB) GB {
	return MakeIO[GB](func() B {
		return F.Pipe1(fa(), f)()
	})
}

// Chain composes computations in sequence, using the return value of one computation to determine the next computation.
func Chain[GA ~func() A, GB ~func() B, A, B any](f func(A) GB) func(GA) GB {
	return F.Bind2nd(MonadChain[GA, GB, A, B], f)
}

// MonadChainTo composes computations in sequence, ignoring the return value of the first computation
func MonadChainTo[GA ~func() A, GB ~func() B, A, B any](fa GA, fb GB) GB {
	return MonadChain(fa, F.Constant1[A](fb))
}

// ChainTo composes computations in sequence, ignoring the return value of the first computation
func ChainTo[GA ~func() A, GB ~func() B, A, B any](fb GB) func(GA) GB {
	return F.Bind2nd(MonadChainTo[GA, GB, A, B], fb)
}

// MonadChainFirst composes computations in sequence, using the return value of one computation to determine the next computation and
// keeping only the result of the first.
func MonadChainFirst[GA ~func() A, GB ~func() B, A, B any](fa GA, f func(A) GB) GA {
	return C.MonadChainFirst(MonadChain[GA, GA, A, A], MonadMap[GB, GA, B, A], fa, f)
}

// ChainFirst composes computations in sequence, using the return value of one computation to determine the next computation and
// keeping only the result of the first.
func ChainFirst[GA ~func() A, GB ~func() B, A, B any](f func(A) GB) func(GA) GA {
	return C.ChainFirst(MonadChain[GA, GA, A, A], MonadMap[GB, GA, B, A], f)
}

func ApSeq[GB ~func() B, GAB ~func() func(A) B, GA ~func() A, B, A any](ma GA) func(GAB) GB {
	return F.Bind2nd(MonadApSeq[GA, GB, GAB, A, B], ma)
}

func ApPar[GB ~func() B, GAB ~func() func(A) B, GA ~func() A, B, A any](ma GA) func(GAB) GB {
	return F.Bind2nd(MonadApPar[GA, GB, GAB, A, B], ma)
}

func Ap[GB ~func() B, GAB ~func() func(A) B, GA ~func() A, B, A any](ma GA) func(GAB) GB {
	return F.Bind2nd(MonadAp[GA, GB, GAB, A, B], ma)
}

func Flatten[GA ~func() A, GAA ~func() GA, A any](mma GAA) GA {
	return mma()
}

// Memoize computes the value of the provided IO monad lazily but exactly once
func Memoize[GA ~func() A, A any](ma GA) GA {
	// synchronization primitives
	var once sync.Once
	var result A
	// callback
	gen := func() {
		result = ma()
	}
	// returns our memoized wrapper
	return func() A {
		once.Do(gen)
		return result
	}
}

// Delay creates an operation that passes in the value after some delay
func Delay[GA ~func() A, A any](delay time.Duration) func(GA) GA {
	return func(ga GA) GA {
		return MakeIO[GA](func() A {
			time.Sleep(delay)
			return ga()
		})
	}
}

// Now returns the current timestamp
func Now[GA ~func() time.Time]() GA {
	return MakeIO[GA](time.Now)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[GA ~func() A, A any](gen func() GA) GA {
	return MakeIO[GA](func() A {
		return gen()()
	})
}
