package expression

import (
	"reflect"
	"testing"
)

func TestIn(t *testing.T) {
	t.Run("string variable in constant StringArray", func(t *testing.T) {
		in := In{Args: []Expression{Variable{Name: "x"}, StringArray{String("a"), String("b")}}}
		result := in.Evaluate(Input{"x": String("a")})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = in.Evaluate(Input{"x": String("c")})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = in.Evaluate(Input{"x": UNKNOWN})
		if !reflect.DeepEqual(result, in) {
			t.Errorf("Expected %v, got %v", in, result)
		}
		result = in.Evaluate(Input{"x": IGNORE})
		if result != IGNORE {
			t.Errorf("Expected IGNORE, got %v", result)
		}
		result = in.Evaluate(Input{"x": UNSET})
		if result != UNSET {
			t.Errorf("Expected UNSET, got %v", result)
		}
		if ToString(in) != "in(x, [a b])" {
			t.Errorf("Expected in(x, [a b]), got %v", ToString(in))
		}

	})
	t.Run("string variable in variable StringArray", func(t *testing.T) {
		in := In{Args: []Expression{Variable{Name: "x"}, Variable{Name: "y"}}}
		result := in.Evaluate(Input{"x": String("a"), "y": StringArray{String("a"), String("b")}})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = in.Evaluate(Input{"x": String("c"), "y": StringArray{String("a"), String("b")}})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = in.Evaluate(Input{"x": UNKNOWN, "y": StringArray{String("a"), String("b")}})
		expected := In{Args: []Expression{Variable{Name: "x"}, StringArray{String("a"), String("b")}}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
		result = in.Evaluate(Input{"x": IGNORE, "y": StringArray{String("a"), String("b")}})
		if result != IGNORE {
			t.Errorf("Expected IGNORE, got %v", result)
		}
		result = in.Evaluate(Input{"x": UNSET, "y": StringArray{String("a"), String("b")}})
		if result != UNSET {
			t.Errorf("Expected UNSET, got %v", result)
		}

		result = in.Evaluate(Input{"x": String("a"), "y": UNKNOWN})
		expected = In{Args: []Expression{String("a"), Variable{Name: "y"}}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
		result = result.Evaluate(Input{"y": StringArray{String("a"), String("b")}})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("number variable in constant NumberArray", func(t *testing.T) {
		in := In{Args: []Expression{Variable{Name: "x"}, NumberArray{Number(1), Number(2)}}}
		result := in.Evaluate(Input{"x": Number(1)})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = in.Evaluate(Input{"x": Number(3)})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = in.Evaluate(Input{"x": UNKNOWN})
		if !reflect.DeepEqual(result, in) {
			t.Errorf("Expected %v, got %v", in, result)
		}
		result = in.Evaluate(Input{"x": IGNORE})
		if result != IGNORE {
			t.Errorf("Expected IGNORE, got %v", result)
		}
		result = in.Evaluate(Input{"x": UNSET})
		if result != UNSET {
			t.Errorf("Expected UNSET, got %v", result)
		}
	})

	t.Run("number variable in variable NumberArray", func(t *testing.T) {
		in := In{Args: []Expression{Variable{Name: "x"}, Variable{Name: "y"}}}
		result := in.Evaluate(Input{"x": Number(1), "y": NumberArray{Number(1), Number(2)}})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = in.Evaluate(Input{"x": Number(3), "y": NumberArray{Number(1), Number(2)}})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = in.Evaluate(Input{"x": UNKNOWN, "y": NumberArray{Number(1), Number(2)}})
		expected := In{Args: []Expression{Variable{Name: "x"}, NumberArray{Number(1), Number(2)}}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
		result = in.Evaluate(Input{"x": IGNORE, "y": NumberArray{Number(1), Number(2)}})
		if result != IGNORE {
			t.Errorf("Expected IGNORE, got %v", result)
		}
		result = in.Evaluate(Input{"x": UNSET, "y": NumberArray{Number(1), Number(2)}})
		if result != UNSET {
			t.Errorf("Expected UNSET, got %v", result)
		}

		result = in.Evaluate(Input{"x": Number(1), "y": UNKNOWN})
		expected = In{Args: []Expression{Number(1), Variable{Name: "y"}}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
		result = result.Evaluate(Input{"y": NumberArray{Number(1), Number(2)}})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("bool variable in constant BoolArray", func(t *testing.T) {
		in := In{Args: []Expression{Variable{Name: "x"}, BoolArray{Bool(true), Bool(false)}}}
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
		result = in.Evaluate(Input{"x": UNKNOWN})
		if !reflect.DeepEqual(result, in) {
			t.Errorf("Expected %v, got %v", in, result)
		}
		result = in.Evaluate(Input{"x": IGNORE})
		if result != IGNORE {
			t.Errorf("Expected IGNORE, got %v", result)
		}
		result = in.Evaluate(Input{"x": UNSET})
		if result != UNSET {
			t.Errorf("Expected UNSET, got %v", result)
		}
	})
}

func TestNotIn(t *testing.T) {
	t.Run("string variable not in constant StringArray", func(t *testing.T) {
		notIn := NotIn{Args: []Expression{Variable{Name: "x"}, StringArray{String("a"), String("b")}}}
		result := notIn.Evaluate(Input{"x": String("a")})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = notIn.Evaluate(Input{"x": String("c")})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = notIn.Evaluate(Input{"x": UNKNOWN})
		if !reflect.DeepEqual(result, notIn) {
			t.Errorf("Expected %v, got %v", notIn, result)
		}
		result = notIn.Evaluate(Input{"x": IGNORE})
		if result != IGNORE {
			t.Errorf("Expected IGNORE, got %v", result)
		}
		result = notIn.Evaluate(Input{"x": UNSET})
		if result != UNSET {
			t.Errorf("Expected UNSET, got %v", result)
		}
		if ToString(notIn) != "not_in(x, [a b])" {
			t.Errorf("Expected not_in(x, [a b]), got %v", ToString(notIn))
		}
	})

	t.Run("string variable not in variable StringArray", func(t *testing.T) {
		notIn := NotIn{Args: []Expression{Variable{Name: "x"}, Variable{Name: "y"}}}
		result := notIn.Evaluate(Input{"x": String("a"), "y": StringArray{String("a"), String("b")}})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = notIn.Evaluate(Input{"x": String("c"), "y": StringArray{String("a"), String("b")}})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = notIn.Evaluate(Input{"x": UNKNOWN, "y": StringArray{String("a"), String("b")}})
		expected := NotIn{Args: []Expression{Variable{Name: "x"}, StringArray{String("a"), String("b")}}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
		result = notIn.Evaluate(Input{"x": IGNORE, "y": StringArray{String("a"), String("b")}})
		if result != IGNORE {
			t.Errorf("Expected IGNORE, got %v", result)
		}
		result = notIn.Evaluate(Input{"x": UNSET, "y": StringArray{String("a"), String("b")}})
		if result != UNSET {
			t.Errorf("Expected UNSET, got %v", result)
		}

		result = notIn.Evaluate(Input{"x": String("a"), "y": UNKNOWN})
		expected = NotIn{Args: []Expression{String("a"), Variable{Name: "y"}}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("Expected %v, got %v", expected, result)
		}
		result = result.Evaluate(Input{"y": StringArray{String("a"), String("b")}})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("In empty array is always false", func(t *testing.T) {
		in := In{Args: []Expression{Variable{Name: "a"}, StringArray{}}}
		result := in.Evaluate(Input{
			"a": UNKNOWN,
		})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("NotIn empty array is always true", func(t *testing.T) {
		notIn := NotIn{Args: []Expression{Variable{Name: "a"}, StringArray{}}}
		result := notIn.Evaluate(Input{
			"a": UNKNOWN,
		})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})
}
