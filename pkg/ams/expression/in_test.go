package expression

import (
	"reflect"
	"testing"
)

func TestIn(t *testing.T) {
	t.Run("string variable in constant StringArray", func(t *testing.T) { //nolint:dupl
		in := In(Ref("x"), StringArray{String("a"), String("b")})
		result := in.Evaluate(Input{"x": String("a")})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = in.Evaluate(Input{"x": String("c")})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = in.Evaluate(Input{})
		if !reflect.DeepEqual(result, in) {
			t.Errorf("Expected %+v, got %+v", in, result)
		}
		want := `in({x}, ["a" "b"])`
		if ToString(in) != want {
			t.Errorf("Expected %s, got %v", want, ToString(in))
		}
	})
	t.Run("string variable in variable StringArray", func(t *testing.T) { //nolint:dupl
		in := In(Ref("x"), Ref("y"))
		result := in.Evaluate(Input{"x": String("a"), "y": StringArray{String("a"), String("b")}})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = in.Evaluate(Input{"x": String("c"), "y": StringArray{String("a"), String("b")}})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = in.Evaluate(Input{"y": StringArray{String("a"), String("b")}})
		expected := In(Ref("x"), StringArray{String("a"), String("b")})
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %+v, got %+v", expected, result)
		}

		result = in.Evaluate(Input{"x": String("a")})
		expected = In(String("a"), Ref("y"))
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
		result = result.Evaluate(Input{"y": StringArray{String("a"), String("b")}})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("number variable in constant NumberArray", func(t *testing.T) {
		in := In(Ref("x"), NumberArray{Number(1), Number(2)})
		result := in.Evaluate(Input{"x": Number(1)})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = in.Evaluate(Input{"x": Number(3)})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = in.Evaluate(Input{})
		if !reflect.DeepEqual(result, in) {
			t.Errorf("Expected %v, got %v", in, result)
		}
	})

	t.Run("number variable in variable NumberArray", func(t *testing.T) { //nolint:dupl
		in := In(Ref("x"), Ref("y"))
		result := in.Evaluate(Input{"x": Number(1), "y": NumberArray{Number(1), Number(2)}})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = in.Evaluate(Input{"x": Number(3), "y": NumberArray{Number(1), Number(2)}})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = in.Evaluate(Input{"y": NumberArray{Number(1), Number(2)}})
		expected := In(Ref("x"), NumberArray{Number(1), Number(2)})
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		result = in.Evaluate(Input{"x": Number(1)})
		expected = In(Number(1), Ref("y"))
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
		result = result.Evaluate(Input{"y": NumberArray{Number(1), Number(2)}})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("bool variable in constant BoolArray", func(t *testing.T) {
		in := In(Ref("x"), BoolArray{Bool(true), Bool(false)})
		result := in.Evaluate(Input{"x": Bool(true)})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = in.Evaluate(Input{"x": Bool(false)})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = in.Evaluate(Input{"x": Bool(false)})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = in.Evaluate(Input{"x": Bool(false)})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = in.Evaluate(Input{})
		if !reflect.DeepEqual(result, in) {
			t.Errorf("Expected %v, got %v", in, result)
		}
	})
}

func TestNotIn(t *testing.T) {
	t.Run("string variable not in constant StringArray", func(t *testing.T) { //nolint:dupl
		notIn := NotIn(Ref("x"), StringArray{String("a"), String("b")})
		result := notIn.Evaluate(Input{"x": String("a")})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = notIn.Evaluate(Input{"x": String("c")})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = notIn.Evaluate(Input{})
		if !reflect.DeepEqual(result, notIn) {
			t.Errorf("Expected %v, got %v", notIn, result)
		}

		want := `not_in({x}, ["a" "b"])`
		if ToString(notIn) != want {
			t.Errorf("Expected %s, got %v", want, ToString(notIn))
		}
	})

	t.Run("string variable not in variable StringArray", func(t *testing.T) { //nolint:dupl
		notIn := NotIn(Ref("x"), Ref("y"))
		result := notIn.Evaluate(Input{"x": String("a"), "y": StringArray{String("a"), String("b")}})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = notIn.Evaluate(Input{"x": String("c"), "y": StringArray{String("a"), String("b")}})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = notIn.Evaluate(Input{"y": StringArray{String("a"), String("b")}})
		expected := NotIn(Ref("x"), StringArray{String("a"), String("b")})
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}

		result = notIn.Evaluate(Input{"x": String("a")})
		expected = NotIn(String("a"), Ref("y"))
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
		result = result.Evaluate(Input{"y": StringArray{String("a"), String("b")}})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("In empty array is always false", func(t *testing.T) {
		in := In(Ref("a"), StringArray{})
		result := in.Evaluate(Input{})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("NotIn empty array is always true", func(t *testing.T) {
		notIn := NotIn(Ref("a"), StringArray{})
		result := notIn.Evaluate(Input{})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})
}
