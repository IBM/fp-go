package cli

import "testing"

func TestCollectExamplesUsesGoDocNaming(t *testing.T) {
	examples, err := collectExamples("../testpkg")
	if err != nil {
		t.Fatalf("collectExamples() error = %v", err)
	}

	got := map[string]Example{}
	for _, example := range examples {
		got[example.Name] = example
	}

	if got["ExamplePerson"].Symbol != "" {
		t.Fatalf("ExamplePerson symbol = %q, expected empty package-level symbol from [`doc.NewFromFiles()`](cli/examples.go:22)", got["ExamplePerson"].Symbol)
	}
	if got["ExamplePerson_String"].Symbol != "" {
		t.Fatalf("ExamplePerson_String symbol = %q, expected empty package-level symbol from [`doc.NewFromFiles()`](cli/examples.go:22)", got["ExamplePerson_String"].Symbol)
	}
	if got["ExamplePerson_second"].Symbol != "" {
		t.Fatalf("ExamplePerson_second symbol = %q, expected empty package-level symbol from [`doc.NewFromFiles()`](cli/examples.go:22)", got["ExamplePerson_second"].Symbol)
	}
	if got["ExamplePerson_third"].Symbol != "" {
		t.Fatalf("ExamplePerson_third symbol = %q, expected empty package-level symbol from [`doc.NewFromFiles()`](cli/examples.go:22)", got["ExamplePerson_third"].Symbol)
	}
}
