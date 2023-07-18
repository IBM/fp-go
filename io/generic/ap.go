package generic

import (
	G "github.com/IBM/fp-go/internal/apply"
)

const (
	// useParallel is the feature flag to control if we use the parallel or the sequential implementation of ap
	useParallel = true
)

// monadApSeq implements the applicative on a single thread by first executing mab and the ma
func monadApSeq[GA ~func() A, GB ~func() B, GAB ~func() func(A) B, A, B any](mab GAB, ma GA) GB {
	return MakeIO[GB](func() B {
		return mab()(ma())
	})
}

// monadApPar implements the applicative on two threads, the main thread executes mab and the actuall
// apply operation and the second thred computes ma. Communication between the threads happens via a channel
func monadApPar[GA ~func() A, GB ~func() B, GAB ~func() func(A) B, A, B any](mab GAB, ma GA) GB {
	return MakeIO[GB](func() B {
		c := make(chan A)
		go func() {
			c <- ma()
			close(c)
		}()
		return mab()(<-c)
	})
}

// MonadAp implements the `ap` operation. Depending on a feature flag this will be sequential or parallel, the preferred implementation
// is parallel
func MonadAp[GA ~func() A, GB ~func() B, GAB ~func() func(A) B, A, B any](mab GAB, ma GA) GB {
	if useParallel {
		return monadApPar[GA, GB](mab, ma)
	}
	return monadApSeq[GA, GB](mab, ma)
}

// MonadApFirst combines two effectful actions, keeping only the result of the first.
func MonadApFirst[GA ~func() A, GB ~func() B, GBA ~func() func(B) A, A, B any](first GA, second GB) GA {
	return G.MonadApFirst(
		MonadAp[GB, GA, GBA, B, A],
		MonadMap[GA, GBA, A, func(B) A],

		first,
		second,
	)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
func ApFirst[GA ~func() A, GB ~func() B, GBA ~func() func(B) A, A, B any](second GB) func(GA) GA {
	return G.ApFirst(
		MonadAp[GB, GA, GBA, B, A],
		MonadMap[GA, GBA, A, func(B) A],

		second,
	)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
func MonadApSecond[GA ~func() A, GB ~func() B, GBB ~func() func(B) B, A, B any](first GA, second GB) GB {
	return G.MonadApSecond(
		MonadAp[GB, GB, GBB, B, B],
		MonadMap[GA, GBB, A, func(B) B],

		first,
		second,
	)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
func ApSecond[GA ~func() A, GB ~func() B, GBB ~func() func(B) B, A, B any](second GB) func(GA) GB {
	return G.ApSecond(
		MonadAp[GB, GB, GBB, B, B],
		MonadMap[GA, GBB, A, func(B) B],

		second,
	)
}
