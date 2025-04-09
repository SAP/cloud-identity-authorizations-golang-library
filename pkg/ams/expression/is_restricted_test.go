package expression

import (
	"reflect"
	"testing"
)

func TestIsRestricted_Evaluate(t *testing.T) {
	t.Run("evaluates always to the same value", func(t *testing.T) {
		e := IsRestricted{Not: true, VariableName: "foo"}
		if got, want := e.Evaluate(nil), Bool(true); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		result := e.Evaluate(Input{
			"foo": UNKNOWN,
		})
		if ToString(result) != "is_not_restricted(foo)" {
			t.Errorf("got %v, want %v", ToString(result), "is_not_restricted(foo)")
		}

		e = IsRestricted{Not: false, VariableName: "foo"}
		if got, want := e.Evaluate(nil), Bool(false); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		result = e.Evaluate(Input{
			"foo": UNKNOWN,
		})
		if ToString(result) != "is_restricted(foo)" {
			t.Errorf("got %v, want %v", ToString(result), "is_restricted(foo)")
		}
	})

	t.Run(" is restrictable", func(t *testing.T) {
		e := IsRestricted{Not: true, VariableName: "foo"}
		if got := IsRestrictable(e); !got {
			t.Errorf("got %v, want %v", got, true)
		}
	})

	t.Run("is restrictable inside AND and OR", func(t *testing.T) {
		var e Expression
		e = And{Args: []Expression{IsRestricted{Not: true, VariableName: "foo"}, Bool(true)}}
		if !IsRestrictable(e) {
			t.Errorf("got %v, want %v", false, true)
		}

		restrEq := Eq{Args: []Expression{Variable{Name: "foo"}, String("bar")}}

		restr := []ExpressionContainer{
			{
				Expression:    restrEq,
				VariableNames: variableSet{"foo": true},
			},
		}
		e = ApplyRestriction(e, restr)
		var expected Expression
		expected = And{Args: []Expression{restrEq, Bool(true)}}
		if !reflect.DeepEqual(e, expected) {
			t.Errorf("got %v, want %v", e, expected)
		}

		e = Or{Args: []Expression{IsRestricted{Not: true, VariableName: "foo"}, Bool(false)}}
		if !IsRestrictable(e) {
			t.Errorf("got %v, want %v", false, true)
		}

		e = ApplyRestriction(e, restr)
		expected = Or{Args: []Expression{restrEq, Bool(false)}}
		if !reflect.DeepEqual(e, expected) {
			t.Errorf("got %v, want %v", e, expected)
		}

		e = Not{Arg: IsRestricted{Not: true, VariableName: "foo"}}
		if !IsRestrictable(e) {
			t.Errorf("got %v, want %v", false, true)
		}
		e = ApplyRestriction(e, restr)
		expected = Not{Arg: restrEq}
		if !reflect.DeepEqual(e, expected) {
			t.Errorf("got %v, want %v", e, expected)
		}
	})

	t.Run("only retricts if variable is in restriction", func(t *testing.T) {
		e := IsRestricted{Not: true, VariableName: "foo"}
		restr := []ExpressionContainer{
			{
				Expression:    Eq{Args: []Expression{Variable{Name: "bar"}, String("baz")}},
				VariableNames: variableSet{"bar": true},
			},
		}
		result := ApplyRestriction(e, restr)
		expected := IsRestricted{Not: true, VariableName: "foo"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, want %v", result, expected)
		}
	})

}
