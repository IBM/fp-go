package cli

import (
	"fmt"
	"testing"
)

func TestCollectExamplesUsesGoDocNaming(t *testing.T) {
	examples, err := collectExamples("../testpkg")
	if err != nil {
		t.Fatalf("collectExamples() error = %v", err)
	}

	got := map[string]Example{}
	for _, example := range examples {
		got[example.Name] = example
	}

	// Test type example
	if got["ExamplePerson"].Symbol != "Person" {
		t.Fatalf("ExamplePerson symbol = %q, expected %q", got["ExamplePerson"].Symbol, "Person")
	}

	// Test method example
	if got["ExamplePerson_String"].Symbol != "Person.String" {
		t.Fatalf("ExamplePerson_String symbol = %q, expected %q", got["ExamplePerson_String"].Symbol, "Person.String")
	}

	// Test type examples with suffixes
	if got["ExamplePerson_second"].Symbol != "Person" {
		t.Fatalf("ExamplePerson_second symbol = %q, expected %q", got["ExamplePerson_second"].Symbol, "Person")
	}
	if got["ExamplePerson_third"].Symbol != "Person" {
		t.Fatalf("ExamplePerson_third symbol = %q, expected %q", got["ExamplePerson_third"].Symbol, "Person")
	}

	// Test function example with suffix
	if got["ExampleHelloWorld_somesuffix"].Symbol != "HelloWorld" {
		t.Fatalf("ExampleHelloWorld_somesuffix symbol = %q, expected %q", got["ExampleHelloWorld_somesuffix"].Symbol, "HelloWorld")
	}
}

func TestCollectExamplesUsesGoDocNamingArray(t *testing.T) {
	examples, err := collectExamples("../../../v2/array")
	if err != nil {
		t.Fatalf("collectExamples() error = %v", err)
	}

	got := map[string]Example{}
	for _, example := range examples {
		got[example.Name] = example

		fmt.Println(example.Name)
	}
}
