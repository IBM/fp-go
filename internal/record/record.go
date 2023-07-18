package record

func Reduce[M ~map[K]V, K comparable, V, R any](r M, f func(R, V) R, initial R) R {
	current := initial
	for _, v := range r {
		current = f(current, v)
	}
	return current
}

func ReduceWithIndex[M ~map[K]V, K comparable, V, R any](r M, f func(K, R, V) R, initial R) R {
	current := initial
	for k, v := range r {
		current = f(k, current, v)
	}
	return current
}

func ReduceRef[M ~map[K]V, K comparable, V, R any](r M, f func(R, *V) R, initial R) R {
	current := initial
	for _, v := range r {
		current = f(current, &v) // #nosec G601
	}
	return current
}

func ReduceRefWithIndex[M ~map[K]V, K comparable, V, R any](r M, f func(K, R, *V) R, initial R) R {
	current := initial
	for k, v := range r {
		current = f(k, current, &v) // #nosec G601
	}
	return current
}
