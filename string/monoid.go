package string

import (
	M "github.com/IBM/fp-go/monoid"
)

// Monoid is the monoid implementing string concatenation
var Monoid = M.MakeMonoid(concat, "")
