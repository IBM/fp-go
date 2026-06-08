package testpkg

type Person struct {
	Name string
	Age  int
	// Deprecated: Use Email instead
	OldEmail string
	Email    string
}
