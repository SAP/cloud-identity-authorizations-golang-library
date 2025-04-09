package expression

import (
	"reflect"
	"testing"
)

func TestIsNull(t *testing.T) {
	t.Run("variable is null", func(t *testing.T) {
		isNull := IsNull{Arg: Variable{Name: "x"}}
		result := isNull.Evaluate(Input{"x": UNSET})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = isNull.Evaluate(Input{"x": String("a")})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = isNull.Evaluate(Input{"x": UNKNOWN})
		if !reflect.DeepEqual(result, isNull) {
			t.Errorf("Expected %v, got %v", isNull, result)
		}
		result = isNull.Evaluate(Input{"x": IGNORE})
		if result != IGNORE {
			t.Errorf("Expected IGNORE, got %v", result)
		}
		if ToString(isNull) != "is_null(x)" {
			t.Errorf("Expected is_null(x), got %v", ToString(isNull))
		}
	})

	t.Run("constant is null", func(t *testing.T) {
		isNull := IsNull{Arg: String("a")}
		result := isNull.Evaluate(nil)
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})
}

func TestIsNotNull(t *testing.T) {
	t.Run("variable is not null", func(t *testing.T) {
		isNotNull := IsNotNull{Arg: Variable{Name: "x"}}
		result := isNotNull.Evaluate(Input{"x": UNSET})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = isNotNull.Evaluate(Input{"x": String("a")})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = isNotNull.Evaluate(Input{"x": UNKNOWN})
		if !reflect.DeepEqual(result, isNotNull) {
			t.Errorf("Expected %v, got %v", isNotNull, result)
		}
		result = isNotNull.Evaluate(Input{"x": IGNORE})
		if result != IGNORE {
			t.Errorf("Expected IGNORE, got %v", result)
		}
		if ToString(isNotNull) != "is_not_null(x)" {
			t.Errorf("Expected is_not_null(x), got %v", ToString(isNotNull))
		}
	})

	t.Run("constant is not null", func(t *testing.T) {
		isNotNull := IsNotNull{Arg: String("a")}
		result := isNotNull.Evaluate(nil)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})
}
