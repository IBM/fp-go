package testpkg_test

import "fmt"

// ExamplePerson_second shows another way to use Person
func ExamplePerson_second() {
	p := Person{Name: "Charlie", Age: 35}
	fmt.Println(p.Name)
	// Output: Charlie
}

// ExamplePerson_third demonstrates yet another Person example
func ExamplePerson_third() {
	p := Person{Name: "Diana", Age: 28}
	fmt.Printf("Age: %d", p.Age)
	// Output: Age: 28
}
