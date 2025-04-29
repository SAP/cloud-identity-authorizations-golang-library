package expression

import (
	"reflect"
	"testing"
)

func TestIsRestricted_Evaluate(t *testing.T) {
	t.Run("evaluates always to the same value", func(t *testing.T) {
		e := NotRestricted(Ref("foo"))
		if got, want := e.Evaluate(nil), Bool(true); got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		e = Restricted(Ref("foo"))
		if got, want := e.Evaluate(nil), Bool(false); got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run(" is restrictable", func(t *testing.T) {
		e := NotRestricted(Ref("foo"))
		if got := IsRestrictable(e); !got {
			t.Errorf("got %v, want %v", got, true)
		}
	})

	t.Run("is restrictable inside AND and OR", func(t *testing.T) {
		var e Expression
		e = And(NotRestricted(Ref("foo")), Bool(true))
		if !IsRestrictable(e) {
			t.Errorf("got %v, want %v", false, true)
		}

		restrEq := Eq(Ref("foo"), String("bar"))

		restr := []ExpressionContainer{
			{
				Expression: restrEq,
				References: referenceSet{"foo": true},
			},
		}
		e = ApplyRestriction(e, restr)
		var expected Expression
		expected = And(restrEq, Bool(true))
		if !reflect.DeepEqual(e, expected) {
			t.Errorf("got %v, want %v", e, expected)
		}

		e = Or(NotRestricted(Ref("foo")), Bool(false))
		if !IsRestrictable(e) {
			t.Errorf("got %v, want %v", false, true)
		}

		e = ApplyRestriction(e, restr)
		expected = Or(restrEq, Bool(false))
		if !reflect.DeepEqual(e, expected) {
			t.Errorf("got %v, want %v", e, expected)
		}

		e = Not(NotRestricted(Ref("foo")))
		if !IsRestrictable(e) {
			t.Errorf("got %v, want %v", false, true)
		}
		e = ApplyRestriction(e, restr)
		expected = Not(restrEq)
		if !reflect.DeepEqual(e, expected) {
			t.Errorf("got %v, want %v", e, expected)
		}
	})

	t.Run("only retricts if variable is in restriction", func(t *testing.T) {
		e := NotRestricted(Ref("foo"))
		restr := []ExpressionContainer{
			{
				Expression: Eq(Ref("bar"), String("baz")),
				References: referenceSet{"bar": true},
			},
		}
		result := ApplyRestriction(e, restr)
		expected := NotRestricted(Ref("foo"))
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("got %v, want %v", result, expected)
		}
	})
}
