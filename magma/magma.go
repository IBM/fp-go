package magma

type Magma[A any] interface {
	Concat(x A, y A) A
}

type magma[A any] struct {
	c func(A, A) A
}

func (self magma[A]) Concat(x A, y A) A {
	return self.c(x, y)
}

func MakeMagma[A any](c func(A, A) A) Magma[A] {
	return magma[A]{c: c}
}

func Reverse[A any](m Magma[A]) Magma[A] {
	return MakeMagma(func(x A, y A) A {
		return m.Concat(y, y)
	})
}

func filterFirst[A any](p func(A) bool, c func(A, A) A, x A, y A) A {
	if p(x) {
		return c(x, y)
	}
	return y
}

func filterSecond[A any](p func(A) bool, c func(A, A) A, x A, y A) A {
	if p(y) {
		return c(x, y)
	}
	return x
}

func FilterFirst[A any](p func(A) bool) func(Magma[A]) Magma[A] {
	return func(m Magma[A]) Magma[A] {
		c := m.Concat
		return MakeMagma(func(x A, y A) A {
			return filterFirst(p, c, x, y)
		})
	}
}

func FilterSecond[A any](p func(A) bool) func(Magma[A]) Magma[A] {
	return func(m Magma[A]) Magma[A] {
		c := m.Concat
		return MakeMagma(func(x A, y A) A {
			return filterSecond(p, c, x, y)
		})
	}
}

func first[A any](x A, y A) A {
	return x
}

func second[A any](x A, y A) A {
	return y
}

func First[A any]() Magma[A] {
	return MakeMagma(first[A])
}

func Second[A any]() Magma[A] {
	return MakeMagma(second[A])
}

func endo[A any](f func(A) A, c func(A, A) A, x A, y A) A {
	return c(f(x), f(y))
}

func Endo[A any](f func(A) A) func(Magma[A]) Magma[A] {
	return func(m Magma[A]) Magma[A] {
		c := m.Concat
		return MakeMagma(func(x A, y A) A {
			return endo(f, c, x, y)
		})
	}
}
