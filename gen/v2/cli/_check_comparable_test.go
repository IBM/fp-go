package cli_test

import (
	"fmt"
	"go/types"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestCheckComparableDebug(t *testing.T) {
	cfg := &packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedSyntax,
	}
	pkgs, err := packages.Load(cfg, "net/http")
	if err != nil {
		t.Fatal(err)
	}
	pkg := pkgs[0]
	scope := pkg.Types.Scope()
	obj := scope.Lookup("Server")
	typName := obj.(*types.TypeName)
	named := typName.Type().(*types.Named)
	structType := named.Underlying().(*types.Struct)
	for i := 0; i < structType.NumFields(); i++ {
		f := structType.Field(i)
		fmt.Printf("Field: %-20s Comparable: %v  Type: %s\n", f.Name(), types.Comparable(f.Type()), f.Type().String())
	}
}
