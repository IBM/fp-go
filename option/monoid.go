package option

import (
	F "github.com/ibm/fp-go/function"
	M "github.com/ibm/fp-go/monoid"
	S "github.com/ibm/fp-go/semigroup"
)

func Semigroup[A any]() func(S.Semigroup[A]) S.Semigroup[Option[A]] {
	return func(s S.Semigroup[A]) S.Semigroup[Option[A]] {
		concat := s.Concat
		return S.MakeSemigroup(
			func(x, y Option[A]) Option[A] {
				return MonadFold(x, F.Constant(y), func(left A) Option[A] {
					return MonadFold(y, F.Constant(x), func(right A) Option[A] {
						return Some(concat(left, right))
					})
				})
			},
		)
	}
}

// Monoid returning the left-most non-`None` value. If both operands are `Some`s then the inner values are
// concatenated using the provided `Semigroup`
//
// | x       | y       | concat(x, y)       |
// | ------- | ------- | ------------------ |
// | none    | none    | none               |
// | some(a) | none    | some(a)            |
// | none    | some(b) | some(b)            |
// | some(a) | some(b) | some(concat(a, b)) |
func Monoid[A any]() func(S.Semigroup[A]) M.Monoid[Option[A]] {
	sg := Semigroup[A]()
	return func(s S.Semigroup[A]) M.Monoid[Option[A]] {
		return M.MakeMonoid(sg(s).Concat, None[A]())
	}
}
