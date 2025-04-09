package expression

import (
	"reflect"
	"testing"
)

func TestAnd(t *testing.T) {

	t.Run("Both true", func(t *testing.T) {
		and := And{Args: []Expression{Bool(true), Bool(true)}}
		result := and.Evaluate(nil)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		if ToString(and) != "and(true, true)" {
			t.Errorf("Expected true && true, got %v", ToString(and))
		}
		if IsRestrictable(and) {
			t.Errorf("Expected false, got %v", true)
		}
	})

	t.Run("First false", func(t *testing.T) {
		and := And{Args: []Expression{Bool(false), Bool(true)}}
		result := and.Evaluate(nil)
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("IGNORE and true", func(t *testing.T) {
		and := And{Args: []Expression{IGNORE, Bool(true)}}
		result := and.Evaluate(nil)
		if result != IGNORE {
			t.Errorf("Expected IGNORE, got %v", result)
		}
	})

	t.Run("IGNORE and false", func(t *testing.T) {
		and := And{Args: []Expression{IGNORE, Bool(false)}}
		result := and.Evaluate(nil)
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("IGNORE and UNKNOWN", func(t *testing.T) {

		expected := Eq{Args: []Expression{
			Variable{Name: "x"},
			Bool(true),
		}}
		and := And{Args: []Expression{
			IGNORE,
			expected,
		}}

		result := and.Evaluate(Input{"x": UNKNOWN})

		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("IGNORE and Expression", func(t *testing.T) {
		and := And{Args: []Expression{IGNORE, Variable{Name: "x"}}}
		result := and.Evaluate(Input{"x": UNKNOWN})
		expected := Variable{Name: "x"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("IGNORE and two Expressions", func(t *testing.T) {
		and := And{Args: []Expression{IGNORE, Variable{Name: "x"}, Variable{Name: "y"}}}
		result := and.Evaluate(Input{"x": UNKNOWN, "y": UNKNOWN})
		expected := And{Args: []Expression{Variable{Name: "x"}, Variable{Name: "y"}}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("UNSET and true", func(t *testing.T) {
		and := And{Args: []Expression{UNSET, Bool(true)}}
		result := and.Evaluate(nil)
		if result != UNSET {
			t.Errorf("Expected UNSET, got %v", result)
		}
	})

	t.Run("UNSET and false", func(t *testing.T) {
		and := And{Args: []Expression{UNSET, Bool(false)}}
		result := and.Evaluate(nil)
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("UNSET and UNKNOWN", func(t *testing.T) {
		and := And{Args: []Expression{UNSET, Variable{Name: "x"}}}
		result := and.Evaluate(Input{"x": UNKNOWN})
		expected := UNSET
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

}

func TestOr(t *testing.T) {
	t.Run("Both true", func(t *testing.T) {
		or := Or{Args: []Expression{Bool(true), Bool(true)}}
		result := or.Evaluate(nil)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		if ToString(or) != "or(true, true)" {
			t.Errorf("Expected true || true, got %v", ToString(or))
		}
		if IsRestrictable(or) {
			t.Errorf("Expected false, got %v", true)
		}
	})

	t.Run("First false", func(t *testing.T) {
		or := Or{Args: []Expression{Bool(false), Bool(true)}}
		result := or.Evaluate(nil)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("IGNORE or true", func(t *testing.T) {
		or := Or{Args: []Expression{IGNORE, Bool(true)}}
		result := or.Evaluate(nil)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("IGNORE or false", func(t *testing.T) {
		or := Or{Args: []Expression{IGNORE, Bool(false)}}
		result := or.Evaluate(nil)
		if result != IGNORE {
			t.Errorf("Expected IGNORE, got %v", result)
		}
	})

	t.Run("IGNORE or UNKNOWN", func(t *testing.T) {
		or := Or{Args: []Expression{IGNORE, Variable{Name: "x"}}}
		result := or.Evaluate(Input{"x": UNKNOWN})
		expected := IGNORE
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run("UNSET or true", func(t *testing.T) {
		or := Or{Args: []Expression{UNSET, Bool(true)}}
		result := or.Evaluate(nil)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("UNSET or false", func(t *testing.T) {
		or := Or{Args: []Expression{UNSET, Bool(false)}}
		result := or.Evaluate(nil)
		if result != UNSET {
			t.Errorf("Expected UNSET, got %v", result)
		}
	})

	t.Run("UNSET or UNKNOWN", func(t *testing.T) {
		or := Or{Args: []Expression{UNSET, Variable{Name: "x"}}}
		result := or.Evaluate(Input{"x": UNKNOWN})
		expected := Variable{Name: "x"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
	})

	t.Run(("false or false"), func(t *testing.T) {
		or := Or{Args: []Expression{Bool(false), Bool(false)}}
		result := or.Evaluate(nil)
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("UNKOWN or UNKOWN", func(t *testing.T) {
		or := Or{Args: []Expression{Variable{Name: "x"}, Variable{Name: "y"}}}
		result := or.Evaluate(Input{"x": UNKNOWN, "y": UNKNOWN})

		if !reflect.DeepEqual(result, or) {
			t.Errorf("Expected %v, got %v", or, result)
		}
	})

}

func TestNot(t *testing.T) {
	t.Run("Not true", func(t *testing.T) {
		not := Not{Arg: Bool(true)}
		result := not.Evaluate(nil)
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		if ToString(not) != "not(true)" {
			t.Errorf("Expected !true, got %v", ToString(not))
		}
	})

	t.Run("Not false", func(t *testing.T) {
		not := Not{Arg: Bool(false)}
		result := not.Evaluate(nil)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("Not IGNORE", func(t *testing.T) {
		not := Not{Arg: IGNORE}
		result := not.Evaluate(nil)
		if result != IGNORE {
			t.Errorf("Expected IGNORE, got %v", result)
		}
	})

	t.Run("Not UNSET", func(t *testing.T) {
		not := Not{Arg: UNSET}
		result := not.Evaluate(nil)
		if result != UNSET {
			t.Errorf("Expected UNSET, got %v", result)
		}
	})

	t.Run("Not UNKNOWN", func(t *testing.T) {
		not := Not{Arg: Variable{Name: "x"}}
		result := not.Evaluate(Input{"x": UNKNOWN})
		if !reflect.DeepEqual(result, not) {
			t.Errorf("Expected %v, got %v", not, result)
		}
	})
}
