package generic

import (
	ET "github.com/ibm/fp-go/either"
	G "github.com/ibm/fp-go/internal/apply"
)

// MonadApFirst combines two effectful actions, keeping only the result of the first.
func MonadApFirst[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], GBA ~func() ET.Either[E, func(B) A], E, A, B any](first GA, second GB) GA {
	return G.MonadApFirst(
		MonadAp[GB, GA, GBA, E, B, A],
		MonadMap[GA, GBA, E, A, func(B) A],

		first,
		second,
	)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
func ApFirst[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], GBA ~func() ET.Either[E, func(B) A], E, A, B any](second GB) func(GA) GA {
	return G.ApFirst(
		MonadAp[GB, GA, GBA, E, B, A],
		MonadMap[GA, GBA, E, A, func(B) A],

		second,
	)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
func MonadApSecond[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], GBB ~func() ET.Either[E, func(B) B], E, A, B any](first GA, second GB) GB {
	return G.MonadApSecond(
		MonadAp[GB, GB, GBB, E, B, B],
		MonadMap[GA, GBB, E, A, func(B) B],

		first,
		second,
	)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
func ApSecond[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], GBB ~func() ET.Either[E, func(B) B], E, A, B any](second GB) func(GA) GB {
	return G.ApSecond(
		MonadAp[GB, GB, GBB, E, B, B],
		MonadMap[GA, GBB, E, A, func(B) B],

		second,
	)
}
