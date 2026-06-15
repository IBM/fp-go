package testpkg_test

import "fmt"

// ExamplePerson demonstrates creating a Person
func ExamplePerson() {
	p := Person{Name: "Alice", Age: 30}
	fmt.Printf("%s is %d years old", p.Name, p.Age)
	// Output: Alice is 30 years old
}

// ExamplePerson_String demonstrates the String method
func ExamplePerson_String() {
	p := Person{Name: "Bob", Age: 25}
	fmt.Println(p.String())
	// Output: Bob (25)
}

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s (%d)", p.Name, p.Age)
}
