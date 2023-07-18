package generic

import (
	G "github.com/IBM/fp-go/internal/apply"
)

// MonadApFirst combines two effectful actions, keeping only the result of the first.
func MonadApFirst[GRA ~func(R) GA, GRB ~func(R) GB, GRBA ~func(R) GBA, GA ~func() A, GB ~func() B, GBA ~func() func(B) A, R, A, B any](first GRA, second GRB) GRA {
	return G.MonadApFirst(
		MonadAp[GRB, GRA, GRBA, GB, GA, GBA, R, B, A],
		MonadMap[GRA, GRBA, GA, GBA, R, A, func(B) A],

		first,
		second,
	)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
func ApFirst[GRA ~func(R) GA, GRB ~func(R) GB, GRBA ~func(R) GBA, GA ~func() A, GB ~func() B, GBA ~func() func(B) A, R, A, B any](second GRB) func(GRA) GRA {
	return G.ApFirst(
		MonadAp[GRB, GRA, GRBA, GB, GA, GBA, R, B, A],
		MonadMap[GRA, GRBA, GA, GBA, R, A, func(B) A],

		second,
	)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
func MonadApSecond[GRA ~func(R) GA, GRB ~func(R) GB, GRBB ~func(R) GBB, GA ~func() A, GB ~func() B, GBB ~func() func(B) B, R, A, B any](first GRA, second GRB) GRB {
	return G.MonadApSecond(
		MonadAp[GRB, GRB, GRBB, GB, GB, GBB, R, B, B],
		MonadMap[GRA, GRBB, GA, GBB, R, A, func(B) B],

		first,
		second,
	)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
func ApSecond[GRA ~func(R) GA, GRB ~func(R) GB, GRBB ~func(R) GBB, GA ~func() A, GB ~func() B, GBB ~func() func(B) B, R, A, B any](second GRB) func(GRA) GRB {
	return G.ApSecond(
		MonadAp[GRB, GRB, GRBB, GB, GB, GBB, R, B, B],
		MonadMap[GRA, GRBB, GA, GBB, R, A, func(B) B],

		second,
	)
}
