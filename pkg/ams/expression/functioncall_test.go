package expression

import (
	"reflect"
	"testing"
)

func TestFunctionCallWithoutArgs(t *testing.T) {
	t.Run("FunctionCall without args", func(t *testing.T) {
		fr := NewFunctionRegistry()

		fc := Function("test", fr, []Expression{})

		fr.RegisterExpressionFunction("test", Eq(Ref("a"), Ref("b")))

		result := fc.Evaluate(Input{"b": FALSE})
		want := Eq(Ref("a"), FALSE)
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Expected %v, got %v", want, result)
		}
	})

	t.Run("call of unknown function", func(t *testing.T) {
		fr := NewFunctionRegistry()

		fc := Function("unknown", fr, []Expression{})

		result := fc.Evaluate(Input{"b": FALSE})
		want := fc
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Expected %v, got %v", want, result)
		}
	})
}
