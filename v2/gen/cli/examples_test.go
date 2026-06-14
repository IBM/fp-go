package cli

import "testing"

func TestIsValidExampleName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		// Valid cases
		{"package example", "Example", true},
		{"package example with suffix", "Example_suffix", true},
		{"package example with underscore suffix", "Example_immutability_struct", true},
		{"package example with multiple underscores", "Example_foo_bar_baz", true},
		{"type example", "ExampleFoo", true},
		{"type example with suffix", "ExampleFoo_suffix", true},
		{"type method example", "ExampleFoo_Bar", true},
		{"type method with suffix", "ExampleFoo_Bar_suffix", true},

		// Invalid cases
		{"no Example prefix", "TestFoo", false},
		{"lowercase after Example", "Examplefoo", false},
		{"empty suffix", "Example_", false},
		{"uppercase in package suffix", "Example_Suffix", false},
		{"uppercase in middle of suffix", "Example_suffIx", false},
		{"too many parts", "ExampleFoo_Bar_Baz_Qux", false},
		{"type with lowercase suffix", "ExampleFoo_bar", true}, // This is valid: example for Foo with suffix "bar"
		{"uppercase in method suffix", "ExampleFoo_Bar_Suffix", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isValidExampleName(tt.input)
			if result != tt.expected {
				t.Errorf("isValidExampleName(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

// Made with Bob
