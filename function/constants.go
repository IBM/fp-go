package function

var ConstTrue = Constant(true)
var ConstFalse = Constant(false)

// ConstNil returns nil
func ConstNil[A any]() *A {
	return (*A)(nil)
}
