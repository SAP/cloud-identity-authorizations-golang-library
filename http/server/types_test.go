package server

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func TestAuthorizationRequestInputUnmarshalInitializesNilMap(t *testing.T) {
	t.Run("unmarshal straight forward", func(t *testing.T) {
		payload := []byte(`{"action":"read","resource":"r1","input":{"name":"alice","active":true,"age":42}}`)

		var req AuthorizationRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			t.Fatalf("unexpected unmarshal error: %v", err)
		}
		if req.Input == nil {
			t.Fatal("expected input map to be initialized")
		}
		if got := req.Input["name"]; got != expression.String("alice") {
			t.Fatalf("expected name to be %v, got %v", expression.String("alice"), got)
		}
		if got := req.Input["active"]; got != expression.Bool(true) {
			t.Fatalf("expected active to be %v, got %v", expression.Bool(true), got)
		}
		if got := req.Input["age"]; got != expression.Number(42) {
			t.Fatalf("expected age to be %v, got %v", expression.Number(42), got)
		}
	})
	t.Run("array typed input", func(t *testing.T) {
		payload := []byte(`{"input":{"names":["alice","bob"],"ages":[42,43],"actives":[true,false]}}`)

		var req AuthorizationRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			t.Fatalf("unexpected unmarshal error: %v", err)
		}
		got := req.Input
		want := Input{
			"names":   expression.StringArray{"alice", "bob"},
			"ages":    expression.NumberArray{42, 43},
			"actives": expression.BoolArray{true, false},
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("expected input to be %v, got %v", want, got)
		}
	})

	t.Run("error on non-object input", func(t *testing.T) {
		payload := []byte(`{"input":"not an object"}`)

		var req AuthorizationRequest
		err := json.Unmarshal(payload, &req)
		if err == nil {
			t.Fatal("expected unmarshal error, got nil")
		}
	})

	t.Run("error on object values", func(t *testing.T) {
		payload := []byte(`{"input":{"name":{"first":"alice","last":"smith"}}}`)

		var req AuthorizationRequest
		err := json.Unmarshal(payload, &req)
		if err == nil {
			t.Fatal("expected unmarshal error, got nil")
		}
	})

	t.Run("error on array with non-uniform types", func(t *testing.T) {
		payload := []byte(`{"input":{"mixed":["alice", 42, true]}}`)

		var req AuthorizationRequest
		err := json.Unmarshal(payload, &req)
		if err == nil {
			t.Fatal("expected unmarshal error, got nil")
		}
		payload = []byte(`{"input":{"mixed":[ 42, true, "alice"]}}`)
		err = json.Unmarshal(payload, &req)
		if err == nil {
			t.Fatal("expected unmarshal error, got nil")
		}
		payload = []byte(`{"input":{"mixed":[ true, "alice", 42]}}`)
		err = json.Unmarshal(payload, &req)
		if err == nil {
			t.Fatal("expected unmarshal error, got nil")
		}
	})

	t.Run("error on array with non-primitive types", func(t *testing.T) {
		payload := []byte(`{"input":{"mixed":[{"name":"alice"},{"name":"bob"}]}}`)

		var req AuthorizationRequest
		err := json.Unmarshal(payload, &req)
		if err == nil {
			t.Fatal("expected unmarshal error, got nil")
		}
	})

	t.Run("unmarshal from marshalled expression.Input", func(t *testing.T) {
		input := expression.Input{
			"string_field": expression.String("value"),
			"number_field": expression.Number(42),
			"bool_field":   expression.Bool(true),
			"string_array": expression.StringArray{"a", "b"},
			"number_array": expression.NumberArray{1, 2},
			"bool_array":   expression.BoolArray{true, false},
			"empty_array":  expression.EmptyArray{},
		}
		data, err := json.Marshal(input)
		if err != nil {
			t.Fatalf("unexpected marshal error: %v", err)
		}
		payload := []byte(`{"input":` + string(data) + `}`)
		var req AuthorizationRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			t.Fatalf("unexpected unmarshal error: %v", err)
		}
		if !reflect.DeepEqual(expression.Input(req.Input), input) {
			t.Fatalf("expected input to be %v, got %v", input, req.Input)
		}
	})
}
