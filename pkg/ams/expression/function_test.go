package expression

import (
	"reflect"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
)

func TestFunction(t *testing.T) {
	t.Run("simple function", func(t *testing.T) {
		funcDCN := []dcn.Function{
			{
				QualifiedName: []string{"func1"},
				Result: dcn.Expression{
					Ref: []string{"a"},
				},
			},
		}
		functions, err := FunctionsFromDCN(funcDCN)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(functions) != 1 {
			t.Fatalf("unexpected number of functions: %d", len(functions))
		}
		f, ok := functions["func1"]
		if !ok {
			t.Fatalf("function not found")
		}
		if !reflect.DeepEqual(f.body, Reference{Name: "a"}) {
			t.Fatalf("unexpected body: %v", f.body)
		}

	})

	t.Run("function calling another function", func(t *testing.T) {
		funcDCN := []dcn.Function{
			{
				QualifiedName: []string{"pkg", "func1"},
				Result: dcn.Expression{
					Ref: []string{"a"},
				},
			},
			{
				QualifiedName: []string{"pkg", "func2"},
				Result: dcn.Expression{
					Call: []string{"is_not_null"},
					Args: []dcn.Expression{
						{
							Call: []string{"pkg", "func1"},
						},
					},
				},
			},
		}

		functions, err := FunctionsFromDCN(funcDCN)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(functions) != 2 {
			t.Fatalf("unexpected number of functions: %d", len(functions))
		}
		f1, ok := functions["pkg.func1"]
		if !ok {
			t.Fatalf("function not found")
		}
		if !reflect.DeepEqual(f1.body, Reference{Name: "a"}) {
			t.Fatalf("unexpected body: %v", f1.body)
		}
		f2, ok := functions["pkg.func2"]
		if !ok {
			t.Fatalf("function not found")
		}
		expected := IsNotNull{
			Arg: Function{
				body: Reference{Name: "a"},
			},
		}
		if !reflect.DeepEqual(f2.body, expected) {
			t.Fatalf("unexpected body: %v", f2.body)
		}

		result := f2.Evaluate(Input{
			"a": Bool(true),
		})
		if !reflect.DeepEqual(result, Bool(true)) {
			t.Fatalf("unexpected result: %v", result)
		}
	})

	t.Run("cycle in function calls", func(t *testing.T) {
		funcDCN := []dcn.Function{
			{
				QualifiedName: []string{"pkg", "func1"},
				Result: dcn.Expression{
					Call: []string{"pkg", "func2"},
				},
			},
			{
				QualifiedName: []string{"pkg", "func2"},
				Result: dcn.Expression{
					Call: []string{"pkg", "func1"},
				},
			},
		}
		_, err := FunctionsFromDCN(funcDCN)
		if err == nil {
			t.Fatalf("expected error")
		}
	})

	t.Run("error in function body", func(t *testing.T) {
		funcDCN := []dcn.Function{
			{
				QualifiedName: []string{"func1"},
				Result: dcn.Expression{
					Call: []string{"unknown"},
				},
			},
		}
		_, err := FunctionsFromDCN(funcDCN)
		if err == nil {
			t.Fatalf("expected error")
		}
	})
}
