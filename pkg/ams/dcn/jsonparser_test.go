package dcn

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestUnmarshalJSON(t *testing.T) { //nolint:maintidx
	t.Run("boolean constant", func(t *testing.T) {
		var ec Expression
		input := `true`
		expected := Expression{
			Constant: true,
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		if !reflect.DeepEqual(ec, expected) {
			t.Errorf("UnmarshalJSON() = %+v, expected %+v", ec, expected)
		}
		marshaled, err := json.Marshal(ec)
		if err != nil {
			t.Fatalf("MarshalJSON() error = %v", err)
		}
		if string(marshaled) != input {
			t.Errorf("MarshalJSON() = %s, expected %s", string(marshaled), input)
		}
	})

	t.Run("number constant", func(t *testing.T) {
		var ec Expression
		input := `123.45`
		expected := Expression{
			Constant: 123.45,
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		if !reflect.DeepEqual(ec, expected) {
			t.Errorf("UnmarshalJSON() = %+v, expected %+v", ec, expected)
		}
		marshaled, err := json.Marshal(ec)
		if err != nil {
			t.Fatalf("MarshalJSON() error = %v", err)
		}
		if string(marshaled) != input {
			t.Errorf("MarshalJSON() = %s, expected %s", string(marshaled), input)
		}
	})

	t.Run("string constant", func(t *testing.T) {
		var ec Expression
		input := `"hello"`
		expected := Expression{
			Constant: "hello",
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		if !reflect.DeepEqual(ec, expected) {
			t.Errorf("UnmarshalJSON() = %+v, expected %+v", ec, expected)
		}
		marshaled, err := json.Marshal(ec)
		if err != nil {
			t.Fatalf("MarshalJSON() error = %v", err)
		}
		if string(marshaled) != input {
			t.Errorf("MarshalJSON() = %s, expected %s", string(marshaled), input)
		}
	})

	t.Run("boolean array constant", func(t *testing.T) {
		var ec Expression
		input := `[true,false]`
		expected := Expression{
			Constant: []bool{true, false},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		if !reflect.DeepEqual(ec, expected) {
			t.Errorf("UnmarshalJSON() = %+v, expected %+v", ec, expected)
		}

		marshaled, err := json.Marshal(ec)
		if err != nil {
			t.Fatalf("MarshalJSON() error = %v", err)
		}

		if string(marshaled) != input {
			t.Errorf("MarshalJSON() = %s, expected %s", string(marshaled), input)
		}
	})

	t.Run("number array constant", func(t *testing.T) {
		var ec Expression
		input := `[123.45,678.9]`
		expected := Expression{
			Constant: []float64{123.45, 678.9},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		if !reflect.DeepEqual(ec, expected) {
			t.Errorf("UnmarshalJSON() = %+v, expected %+v", ec, expected)
		}
		marshaled, err := json.Marshal(ec)
		if err != nil {
			t.Fatalf("MarshalJSON() error = %v", err)
		}
		if string(marshaled) != input {
			t.Errorf("MarshalJSON() = %s, expected %s", string(marshaled), input)
		}
	})

	t.Run("string array constant", func(t *testing.T) {
		var ec Expression
		input := `["hello","world"]`
		expected := Expression{
			Constant: []string{"hello", "world"},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		if !reflect.DeepEqual(ec, expected) {
			t.Errorf("UnmarshalJSON() = %+v, expected %+v", ec, expected)
		}
		marshaled, err := json.Marshal(ec)
		if err != nil {
			t.Fatalf("MarshalJSON() error = %v", err)
		}
		if string(marshaled) != input {
			t.Errorf("MarshalJSON() = %s, expected %s", string(marshaled), input)
		}
	})

	t.Run("call object", func(t *testing.T) {
		var ec Expression
		input := `{"call":["and"],"args":[{"call":["eq"],"args":[{"ref":["$dcl.action"]},"read"]}]}`

		expected := Expression{
			Call: []string{"and"},
			Args: []Expression{
				{
					Call: []string{"eq"},
					Args: []Expression{
						{
							Ref: []string{"$dcl.action"},
						},
						{
							Constant: "read",
						},
					},
				},
			},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		if !reflect.DeepEqual(ec, expected) {
			t.Errorf("UnmarshalJSON() = %+v, expected %+v", ec, expected)
		}
		marshaled, err := json.Marshal(ec)
		if err != nil {
			t.Fatalf("MarshalJSON() error = %v", err)
		}
		if string(marshaled) != input {
			t.Errorf("MarshalJSON() = %s, expected %s", string(marshaled), input)
		}
	})

	t.Run("variable object", func(t *testing.T) {
		var ec Expression
		input := `{"ref":["$dcl.action"]}`
		expected := Expression{
			Ref: []string{"$dcl.action"},
		}
		err := json.Unmarshal([]byte(input), &ec)
		if err != nil {
			t.Fatalf("UnmarshalJSON() error = %v", err)
		}
		if !reflect.DeepEqual(ec, expected) {
			t.Errorf("UnmarshalJSON() = %+v, expected %+v", ec, expected)
		}
		marshaled, err := json.Marshal(ec)
		if err != nil {
			t.Fatalf("MarshalJSON() error = %v", err)
		}
		if string(marshaled) != input {
			t.Errorf("MarshalJSON() = %s, expected %s", string(marshaled), input)
		}
	})

	t.Run("invalid input", func(t *testing.T) {
		var ec Expression
		input := `{"invalid": "input"}`
		err := json.Unmarshal([]byte(input), &ec)
		if err == nil {
			t.Errorf("UnmarshalJSON() error = %v, expected error", err)
		}
	})

	t.Run("marshal empty expression", func(t *testing.T) {
		var ec Expression

		_, err := json.Marshal(ec)
		if err == nil {
			t.Errorf("MarshalJSON() error = %v, expected error", err)
		}
	})
}
