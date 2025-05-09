package expression

import (
	"reflect"
	"testing"
)

func TestIsNull(t *testing.T) { //nolint:dupl
	t.Run("variable is null", func(t *testing.T) {
		isNull := IsNull(Ref("x"))

		result := isNull.Evaluate(Input{"x": String("a")})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = isNull.Evaluate(Input{})
		if !reflect.DeepEqual(result, isNull) {
			t.Errorf("Expected %v, got %v", isNull, result)
		}
		want := "is_null({x})"
		if ToString(isNull) != want {
			t.Errorf("Expected %s, got %v", want, ToString(isNull))
		}
	})

	t.Run("constant is null", func(t *testing.T) {
		isNull := IsNull(String("a"))
		result := isNull.Evaluate(nil)
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})
}

func TestIsNotNull(t *testing.T) { //nolint:dupl
	t.Run("variable is not null", func(t *testing.T) {
		isNotNull := IsNotNull(Ref("x"))

		result := isNotNull.Evaluate(Input{"x": String("a")})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = isNotNull.Evaluate(Input{})
		if !reflect.DeepEqual(result, isNotNull) {
			t.Errorf("Expected %v, got %v", isNotNull, result)
		}

		want := "is_not_null({x})"
		if ToString(isNotNull) != want {
			t.Errorf("Expected %s, got %v", want, ToString(isNotNull))
		}
	})

	t.Run("constant is not null", func(t *testing.T) {
		isNotNull := IsNotNull(String("a"))
		result := isNotNull.Evaluate(nil)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})
}
