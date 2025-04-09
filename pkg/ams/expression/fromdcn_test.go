package expression

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
)

func TestUnmarshalJSON(t *testing.T) {
	t.Run("DCNBool", func(t *testing.T) {
		var ec dcn.Expression
		input := `true`
		expected := Bool(true)
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("DCNNumber", func(t *testing.T) {
		var ec dcn.Expression
		input := `123.45`
		expected := Number(123.45)
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("DCNString", func(t *testing.T) {
		var ec dcn.Expression
		input := `"test"`
		expected := String("test")
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("DCNBoolArray", func(t *testing.T) {
		var ec dcn.Expression
		input := `[true, false, true]`
		expected := BoolArray{true, false, true}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("DCNNumberArray", func(t *testing.T) {
		var ec dcn.Expression
		input := `[1, 2, 3]`
		expected := NumberArray{1, 2, 3}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("DCNStringArray", func(t *testing.T) {
		var ec dcn.Expression
		input := `["a", "b", "c"]`
		expected := StringArray{"a", "b", "c"}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("Variable", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"ref": ["x"]}`
		expected := Variable{Name: "x"}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("And", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["and"], "args": [true, false]}`
		expected := And{
			Args: []Expression{Bool(true), Bool(false)},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("Or", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["or"], "args": [true, false]}`
		expected := Or{
			Args: []Expression{Bool(true), Bool(false)},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("Not", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["not"], "args": [true]}`
		expected := Not{
			Arg: Bool(true),
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("IsNull", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["is_null"], "args": [{"ref": ["x"]}]}`
		expected := IsNull{
			Arg: Variable{Name: "x"},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("IsNotNull", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["is_not_null"], "args": [{"ref": ["x"]}]}`
		expected := IsNotNull{
			Arg: Variable{Name: "x"},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("Like", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["like"], "args": ["test", "pattern"]}`
		expected := NewLike(
			String("test"),
			String("pattern"),
			String(""),
		)
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("NotLike", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["not_like"], "args": ["test", "pattern"]}`
		expected := NewNotLike(
			String("test"),
			String("pattern"),
			String(""),
		)
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("Like with escape", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["like"], "args": ["test", "pattern", "escape"]}`
		expected := NewLike(
			String("test"),
			String("pattern"),
			String("escape"),
		)
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("NotLike with escape", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["not_like"], "args": ["test", "pattern", "escape"]}`
		expected := NewNotLike(
			String("test"),
			String("pattern"),
			String("escape"),
		)
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("Eq", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["eq"], "args": [1, 2]}`
		expected := Eq{
			Args: []Expression{Number(1), Number(2)},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("Ne", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["ne"], "args": [1, 2]}`
		expected := Ne{
			Args: []Expression{Number(1), Number(2)},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("Lt", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["lt"], "args": [1, 2]}`
		expected := Lt{
			Args: []Expression{Number(1), Number(2)},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("Le", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["le"], "args": [1, 2]}`
		expected := Le{
			Args: []Expression{Number(1), Number(2)},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("Gt", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["gt"], "args": [1, 2]}`
		expected := Gt{
			Args: []Expression{Number(1), Number(2)},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("Ge", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["ge"], "args": [1, 2]}`
		expected := Ge{
			Args: []Expression{Number(1), Number(2)},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("Between", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["between"], "args": [1, 2, 3]}`
		expected := Between{
			Args: []Expression{Number(1), Number(2), Number(3)},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("NotBetween", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["not_between"], "args": [1, 2, 3]}`
		expected := NotBetween{
			Args: []Expression{Number(1), Number(2), Number(3)},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("In", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["in"], "args":[{"ref":["x"]}, [1, 2, 3]]}`
		expected := In{
			Args: []Expression{Variable{Name: "x"}, NumberArray{1, 2, 3}},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("NotIn", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["not_in"], "args":[{"ref":["x"]}, [1, 2, 3]]}`
		expected := NotIn{
			Args: []Expression{Variable{Name: "x"}, NumberArray{1, 2, 3}},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("IsRestricted", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["restricted"], "args": [{"ref":["x"]}]}`
		expected := IsRestricted{
			Not:          false,
			VariableName: "x",
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("IsNotRestricted", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["not_restricted"], "args": [{"ref":["x"]}]}`
		expected := IsRestricted{
			Not:          true,
			VariableName: "x",
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, nil)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("Function call", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["custom","function"]}`
		functions := Functions{
			"custom.function": Function{body: Bool(true)},
		}
		expected := Function{
			body: Bool(true),
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		e, err := FromDCN(ec, functions)
		if err != nil {
			t.Fatalf("FromDCN() error = %v", err)
		}
		if !reflect.DeepEqual(e.Expression, expected) {
			t.Errorf("UnmarshalJSON() = %v, expected %v", e.Expression, expected)
		}
	})

	t.Run("Unknown function call", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["unknown","function"],  "args": [{"ref":["x"]}]}`
		functions := Functions{
			"custom.function": Function{body: Bool(true)},
		}

		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		_, err = FromDCN(ec, functions)
		if err == nil {
			t.Fatalf("no error thrown")
		}
	})

	t.Run("Unknown call", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["unknown"],  "args": [{"ref":["x"]}]}`
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		_, err = FromDCN(ec, nil)
		if err == nil {
			t.Fatalf("no error thrown")
		}
	})

	t.Run("Invalid object", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"some": "object"}`
		err := json.Unmarshal([]byte(input), &ec)
		if err == nil {
			t.Fatalf("no error thrown")
		}
	})

	t.Run("Invalid call object", func(t *testing.T) {
		var ec dcn.Expression
		input := `{"call": ["and"], "args": [{"invalid": "object"}]}`
		err := json.Unmarshal([]byte(input), &ec)
		if err == nil {
			t.Fatalf("no error thrown")
		}
	})

}
