package generic

import (
	GT "github.com/IBM/fp-go/v2/optics/traversal/generic"
)

func Compose[
	S, HKTS, A, B, HKTA, HKTB any](ab Traversal[A, B, HKTA, HKTB]) func(Traversal[S, A, HKTS, HKTA]) Traversal[S, B, HKTS, HKTB] {
	return GT.Compose[Traversal[A, B, HKTA, HKTB], Traversal[S, A, HKTS, HKTA], Traversal[S, B, HKTS, HKTB]](ab)
}
