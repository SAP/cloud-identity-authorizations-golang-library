package expression

import (
	"reflect"
	"testing"
)

func TestAnd(t *testing.T) {
	t.Run("Both true", func(t *testing.T) {
		and := And(Bool(true), Bool(true))
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
		and := And(Bool(false), Bool(true))
		result := and.Evaluate(nil)
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})
}

func TestOr(t *testing.T) {
	t.Run("Both true", func(t *testing.T) {
		or := Or(Bool(true), Bool(true))
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
		or := Or(Bool(false), Bool(true))
		result := or.Evaluate(nil)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run(("false or false"), func(t *testing.T) {
		or := Or(Bool(false), Bool(false))
		result := or.Evaluate(nil)
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("UNKNOWN or UNKNOWN", func(t *testing.T) {
		or := Or(Ref("x"), Ref("y"))
		result := or.Evaluate(Input{})

		if !reflect.DeepEqual(result, or) {
			t.Errorf("Expected %v, got %v", or, result)
		}
	})
}

func TestNot(t *testing.T) {
	t.Run("Not true", func(t *testing.T) {
		not := Not(Bool(true))
		result := not.Evaluate(nil)
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		if ToString(not) != "not(true)" {
			t.Errorf("Expected !true, got %v", ToString(not))
		}
	})

	t.Run("Not false", func(t *testing.T) {
		not := Not(Bool(false))
		result := not.Evaluate(nil)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("Not UNKNOWN", func(t *testing.T) {
		not := Not(Ref("x"))
		result := not.Evaluate(Input{})
		if !reflect.DeepEqual(result, not) {
			t.Errorf("Expected %v, got %v", not, result)
		}
	})
}
