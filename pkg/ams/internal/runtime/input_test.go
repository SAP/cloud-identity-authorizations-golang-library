package runtime

import (
	_ "embed"
	"encoding/json"
	"reflect"
	"slices"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

//go:embed testdata/input_test.json
var inputTestJSON []byte

func TestInputValidation(t *testing.T) {
	var r Runtime
	err := json.Unmarshal(inputTestJSON, &r)
	if err != nil {
		t.Fatalf("Failed to unmarshal input test JSON: %v", err)
	}

	t.Run("sanitize input unknown fields", func(t *testing.T) {
		input := expression.Input{
			"unknown_field": expression.String("value"),
		}
		unknownFields, wrongTypeFields := r.SanitizeInput(input, "any_resource")
		want := []string{"unknown_field"}
		if !reflect.DeepEqual(unknownFields, want) {
			t.Errorf("SanitizeInput() unknownFields = %v, want %v", unknownFields, want)
		}
		if len(wrongTypeFields) != 0 {
			t.Errorf("SanitizeInput() wrongTypeFields = %v, want empty", wrongTypeFields)
		}
		wantInput := expression.Input{}
		if !reflect.DeepEqual(input, wantInput) {
			t.Errorf("SanitizeInput() input = %v, want %v", input, wantInput)
		}
	})

	t.Run("sanitize input wrong type fields", func(t *testing.T) {
		input := expression.Input{
			"$app.string": expression.Number(123),
		}
		unknownFields, wrongTypeFields := r.SanitizeInput(input, "any_resource")
		want := []string{"$app.string"}
		if !reflect.DeepEqual(wrongTypeFields, want) {
			t.Errorf("SanitizeInput() wrongTypeFields = %v, want %v", wrongTypeFields, want)
		}
		if len(unknownFields) != 0 {
			t.Errorf("SanitizeInput() unknownFields = %v, want empty", unknownFields)
		}
		wantInput := expression.Input{}
		if !reflect.DeepEqual(input, wantInput) {
			t.Errorf("SanitizeInput() input = %v, want %v", input, wantInput)
		}
	})

	t.Run("sanitize wrong resource input", func(t *testing.T) {
		input := expression.Input{
			"$app.string":      expression.String("value"),
			"$app.number":      expression.Number(123),
			"$resource.string": expression.String("value"),
			"$resource.name":   expression.String("name"),
		}
		unknownFields, wrongTypeFields := r.SanitizeInput(input, "res1")
		wantUnknown := []string{"$resource.name"}
		wantWrongType := []string{}
		if !reflect.DeepEqual(unknownFields, wantUnknown) {
			t.Errorf("SanitizeInput() unknownFields = %v, want %v", unknownFields, wantUnknown)
		}
		if !reflect.DeepEqual(wrongTypeFields, wantWrongType) {
			t.Errorf("SanitizeInput() wrongTypeFields = %v, want %v", wrongTypeFields, wantWrongType)
		}
		unknownFields, wrongTypeFields = r.SanitizeInput(input, "res2")
		wantUnknown = []string{"$resource.string"}
		wantWrongType = []string{}
		if !reflect.DeepEqual(unknownFields, wantUnknown) {
			t.Errorf("SanitizeInput() unknownFields = %v, want %v", unknownFields, wantUnknown)
		}
		if !reflect.DeepEqual(wrongTypeFields, wantWrongType) {
			t.Errorf("SanitizeInput() wrongTypeFields = %v, want %v", wrongTypeFields, wantWrongType)
		}
	})

	t.Run("sanitize input with empty arrays", func(t *testing.T) {
		input := expression.Input{
			"$app.string_array":      expression.EmptyArray{},
			"$app.number_array":      expression.EmptyArray{},
			"$app.boolean_array":     expression.EmptyArray{},
			"$resource.string_array": expression.EmptyArray{},
		}
		unknownFields, wrongTypeFields := r.SanitizeInput(input, "res1")
		wantUnknown := []string{}
		wantWrongType := []string{}
		if !reflect.DeepEqual(unknownFields, wantUnknown) {
			t.Errorf("SanitizeInput() unknownFields = %v, want %v", unknownFields, wantUnknown)
		}
		if !reflect.DeepEqual(wrongTypeFields, wantWrongType) {
			t.Errorf("SanitizeInput() wrongTypeFields = %v, want %v", wrongTypeFields, wantWrongType)
		}
	})

	t.Run("sanitize input all invalid types", func(t *testing.T) {
		input := expression.Input{
			"$app.string":        expression.Number(123),
			"$app.number":        expression.String("value"),
			"$app.boolean":       expression.String("value"),
			"$app.string_array":  expression.NumberArray{1, 2, 3},
			"$app.number_array":  expression.TRUE,
			"$app.boolean_array": expression.StringArray{"true", "false"},
		}
		unknownFields, wrongTypeFields := r.SanitizeInput(input, "res1")
		wantUnknown := []string{}
		wantWrongType := []string{
			"$app.string",
			"$app.number",
			"$app.boolean",
			"$app.string_array",
			"$app.number_array",
			"$app.boolean_array",
		}
		if !compareStringSlices(unknownFields, wantUnknown) {
			t.Errorf("SanitizeInput() unknownFields = %v, want %v", unknownFields, wantUnknown)
		}
		if !compareStringSlices(wrongTypeFields, wantWrongType) {
			t.Errorf("SanitizeInput() wrongTypeFields = %v, want %v", wrongTypeFields, wantWrongType)
		}
		wantInput := expression.Input{}
		if !reflect.DeepEqual(input, wantInput) {
			t.Errorf("SanitizeInput() input = %v, want %v", input, wantInput)
		}
	})
}

func compareStringSlices(a, b []string) bool {
	slices.Sort(a)
	slices.Sort(b)
	return reflect.DeepEqual(a, b)
}

func TestCreateInput(t *testing.T) {
	type TestStruct struct {
		StringField string   `ams:"string_field"`
		NumberField float64  `ams:"number_field"`
		BoolField   bool     `ams:"bool_field"`
		StringArray []string `ams:"string_array"`
		NotTagged   bool
		unexported  string `ams:"unexported"` // This field should be ignored
	}

	t.Run("create input from struct", func(t *testing.T) {
		testData := TestStruct{
			StringField: "test",
			NumberField: 42,
			BoolField:   true,
			StringArray: []string{"a", "b", "c"},
			NotTagged:   true, // This field should be ignored
			unexported:  "should be ignored",
		}

		input := CreateInput(testData)

		wantInput := expression.Input{
			"$app.string_field":      expression.String("test"),
			"$app.number_field":      expression.Number(42),
			"$app.bool_field":        expression.Bool(true),
			"$app.string_array":      expression.StringArray{"a", "b", "c"},
			"$resource.string_field": expression.String("test"),
			"$resource.number_field": expression.Number(42),
			"$resource.bool_field":   expression.Bool(true),
			"$resource.string_array": expression.StringArray{"a", "b", "c"},
		}
		if !reflect.DeepEqual(input, wantInput) {
			t.Errorf("CreateInput() = %+v, want %+v", input, wantInput)
		}
	})

	t.Run("create input from map", func(t *testing.T) {
		testData := map[string]interface{}{
			"string_field": any("test"),
			"bool_field":   true,
			"string_array": []any{"a", 1, true},     // Mixed types, should be ignored
			"number_array": []any{1, 2.0, int64(3)}, // Mixed numeric types, should be interpreted as numbers
		}
		input := CreateInput(testData)
		wantInput := expression.Input{
			"$app.string_field":      expression.String("test"),
			"$app.bool_field":        expression.Bool(true),
			"$app.number_array":      expression.NumberArray{1, 2.0, 3},
			"$resource.string_field": expression.String("test"),
			"$resource.bool_field":   expression.Bool(true),
			"$resource.number_array": expression.NumberArray{1, 2.0, 3},
		}
		if !reflect.DeepEqual(input, wantInput) {
			t.Errorf("CreateInput() = %+v, want %+v", input, wantInput)
		}
	})
}
