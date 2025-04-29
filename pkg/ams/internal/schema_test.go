package internal

import (
	_ "embed"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

//go:embed testfiles/schema/example.dcn
var exampleSchema []byte

//go:embed testfiles/schema/simple.dcn
var simpleSchema []byte

//go:embed testfiles/schema/variables_with_quotes.dcn
var variablesWithQuotes []byte

type NumberWrapper struct {
	NumberValue uint64 `ams:"number_value"`
}

type DeeperNested struct {
	NestedNumberArrayValue []uint        `ams:"nested_number_array_value"`
	DotInName              NumberWrapper `ams:"dot.in.name"`
}

type ExampleInput struct {
	StringValue      string       `ams:"string_value"`
	NumberValue      int          `ams:"number_value"`
	BoolValue        bool         `ams:"bool_value"`
	StringArrayValue []string     `ams:"string_array_value"`
	NumberArrayValue []float32    `ams:"number_array_value"`
	BoolArrayValue   []bool       `ams:"bool_array_value"`
	DeeperNested     DeeperNested `ams:"deeper_nested"`
	unexpected       string       `ams:"unexpected"`
}

type SimpleInput struct {
	NumberValue      *float64
	BoolArrayValue   []bool
	NumberArrayValue []int
}

type SimpleEnv struct {
	EnvN expression.Constant
}

func TestExampleSchema(t *testing.T) { //nolint:maintidx
	var schema Schema

	var ss []dcn.Schema
	err := json.Unmarshal(exampleSchema, &ss)
	if err != nil {
		t.Fatalf("Error parsing schema: %v", err)
	}
	schema = SchemaFromDCN(ss)

	t.Run("contains $dcl.action and $dcl.resource", func(t *testing.T) {
		input := expression.Input{
			"$dcl.action":   expression.String("read"),
			"$dcl.resource": expression.String("example"),
		}

		schema.PurgeInvalidInput(input)

		if len(input) != 2 {
			t.Errorf("Expected 2 field, got %d", len(input))
		}
		if action, ok := input["$dcl.action"]; !ok {
			t.Errorf("Expected $dcl.action to be present")
			if action != expression.String("read") {
				t.Errorf("Expected $dcl.action to be 'read', got %v", action)
			}
		}
		if resource, ok := input["$dcl.resource"]; !ok {
			t.Errorf("Expected $dcl.resource to be present")
			if resource != expression.String("example") {
				t.Errorf("Expected $dcl.resource to be 'example', got %v", resource)
			}
		}
	})

	t.Run("removes invalid fields", func(t *testing.T) {
		input := expression.Input{
			"$dcl.action":   expression.String("read"),
			"$dcl.resource": expression.String("example"),
			"invalid":       expression.String("invalid"),
		}

		schema.PurgeInvalidInput(input)

		if len(input) != 2 {
			t.Errorf("Expected 2 field, got %d", len(input))
		}
		if _, ok := input["invalid"]; ok {
			t.Errorf("Expected 'invalid' to be removed")
		}
	})

	t.Run("removes structure typed input", func(t *testing.T) {
		input := expression.Input{
			"$dcl.action":        expression.String("read"),
			"$app.deeper_nested": expression.String("example"),
		}
		schema.PurgeInvalidInput(input)
		want := expression.Input{
			"$dcl.action": expression.String("read"),
		}
		if !reflect.DeepEqual(input, want) {
			t.Errorf("Expected %v, got %v", want, input)
		}
	})

	t.Run("removes wrongly typed fields", func(t *testing.T) {
		input := expression.Input{
			"$app.string_value":       expression.Number(42),
			"$app.number_value":       expression.String("42"),
			"$app.bool_value":         expression.String("true"),
			"$app.string_array_value": expression.String("42"),
			"$app.number_array_value": expression.Number(42),
			"$app.bool_array_value":   expression.Bool(true),
		}

		schema.PurgeInvalidInput(input)

		if len(input) != 0 {
			t.Errorf("Expected 0 field, got %+v", input)
		}
	})

	t.Run("keeps correctly typed fields", func(t *testing.T) {
		input := expression.Input{
			"$app.string_value":       expression.String("42"),
			"$app.number_value":       expression.Number(42),
			"$app.bool_value":         expression.Bool(true),
			"$app.string_array_value": expression.StringArray{"42"},
			"$app.number_array_value": expression.NumberArray{42},
			"$app.bool_array_value":   expression.BoolArray{true},
		}

		schema.PurgeInvalidInput(input)

		if len(input) != 6 {
			t.Errorf("Expected 6 field, got %+v", input)
		}
	})

	t.Run("Generates input from Custom Structure", func(t *testing.T) {
		i := ExampleInput{
			StringValue:      "string_value",
			NumberValue:      42,
			BoolValue:        true,
			StringArrayValue: []string{"string_array_value"},
			NumberArrayValue: []float32{42},
			BoolArrayValue:   nil, // should not be in the output
			DeeperNested: DeeperNested{
				NestedNumberArrayValue: []uint{42},
				DotInName: NumberWrapper{
					NumberValue: 42,
				},
			},
			unexpected: "unexpected", // should not be in the output
		}

		got := schema.CustomInput("read", "data", i, nil)
		want := expression.Input{
			"$dcl.action":                                     expression.String("read"),
			"$dcl.resource":                                   expression.String("data"),
			"$app.string_value":                               expression.String("string_value"),
			"$app.number_value":                               expression.Number(42),
			"$app.bool_value":                                 expression.Bool(true),
			"$app.string_array_value":                         expression.StringArray{"string_array_value"},
			"$app.number_array_value":                         expression.NumberArray{42},
			"$app.deeper_nested.nested_number_array_value":    expression.NumberArray{42},
			"$app.deeper_nested.\"dot.in.name\".number_value": expression.Number(42),
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %v, got %v", want, got)
		}

		got = schema.CustomInput("read", "data", &i, nil)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})

	t.Run("Generates input from map ", func(t *testing.T) {
		i := map[string]any{
			"string_value": "string_value",
			"number_value": nil,
			"deeper_nested": map[string]any{
				"nested_number_array_value": []uint{42},
				"dot.in.name":               map[string]any(nil),
			},
		}
		got := schema.CustomInput("read", "data", i, nil)
		want := expression.Input{
			"$dcl.action":       expression.String("read"),
			"$dcl.resource":     expression.String("data"),
			"$app.string_value": expression.String("string_value"),
			"$app.deeper_nested.nested_number_array_value": expression.NumberArray{42},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})

	t.Run("Ignores wrongly typed fields", func(t *testing.T) {
		customMap := map[string]any{
			"unexpected":   "unexpected",
			"string_value": 42,
			"bool_value":   SimpleEnv{},
			"number_value": []interface{}{
				"x",
			},
			"string_array_value": []interface{}{
				"x",
			},
			"number_array_value": []interface{}{
				42,
				41.0,
				uint(9),
			},
			"bool_array_value": []interface{}{
				true,
			},
		}

		got := schema.CustomInput("read", "data", customMap, nil)
		want := expression.Input{
			"$dcl.action":             expression.String("read"),
			"$dcl.resource":           expression.String("data"),
			"$app.string_array_value": expression.StringArray{"x"},
			"$app.number_array_value": expression.NumberArray{42, 41.0, 9},
			"$app.bool_array_value":   expression.BoolArray{true},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})

	t.Run("Ignores nil slices", func(t *testing.T) {
		customMap := map[string][]interface{}{
			"string_array_value": nil,
			"number_array_value": nil,
			"bool_array_value":   nil,
		}
		got := schema.CustomInput("read", "data", customMap, nil)
		want := expression.Input{
			"$dcl.action":   expression.String("read"),
			"$dcl.resource": expression.String("data"),
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})

	t.Run("Ignores clices containing nil", func(t *testing.T) {
		customMap := map[string][]interface{}{
			"string_array_value": {nil},
			"number_array_value": {nil},
			"bool_array_value":   {nil},
		}
		got := schema.CustomInput("read", "data", customMap, nil)
		want := expression.Input{
			"$dcl.action":   expression.String("read"),
			"$dcl.resource": expression.String("data"),
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})

	t.Run("Ignores slices containing wrong typed values", func(t *testing.T) {
		customMap := map[string][]interface{}{
			"string_array_value": {"x", 42},
			"number_array_value": {42, "x"},
			"bool_array_value":   {true, "x"},
		}
		got := schema.CustomInput("read", "data", customMap, nil)
		want := expression.Input{
			"$dcl.action":   expression.String("read"),
			"$dcl.resource": expression.String("data"),
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}

func TestSimpleSchema(t *testing.T) {
	var schema Schema

	var ss []dcn.Schema
	err := json.Unmarshal(simpleSchema, &ss)
	if err != nil {
		t.Fatalf("Error parsing schema: %v", err)
	}
	schema = SchemaFromDCN(ss)

	t.Run("Generates input from Custom Structure", func(t *testing.T) {
		fortytwo := 42.0
		app := SimpleInput{
			NumberValue:      &fortytwo,
			BoolArrayValue:   []bool{true},
			NumberArrayValue: []int{42},
		}
		env := SimpleEnv{
			EnvN: expression.Number(3),
		}
		got := schema.CustomInput("read", "data", app, env)
		want := expression.Input{
			"$dcl.action":           expression.String("read"),
			"$dcl.resource":         expression.String("data"),
			"$app.NumberValue":      expression.Number(42),
			"$app.BoolArrayValue":   expression.BoolArray{true},
			"$app.NumberArrayValue": expression.NumberArray{42},
			"$env.EnvN":             expression.Number(3),
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %v, got %v", want, got)
		}

		env.EnvN = nil

		got = schema.CustomInput("read", "data", app, env)
		want = expression.Input{
			"$dcl.action":           expression.String("read"),
			"$dcl.resource":         expression.String("data"),
			"$app.NumberValue":      expression.Number(42),
			"$app.BoolArrayValue":   expression.BoolArray{true},
			"$app.NumberArrayValue": expression.NumberArray{42},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %v, got %v", want, got)
		}

		env.EnvN = nil
		app.NumberValue = nil

		got = schema.CustomInput("read", "data", app, env)
		want = expression.Input{
			"$dcl.action":           expression.String("read"),
			"$dcl.resource":         expression.String("data"),
			"$app.BoolArrayValue":   expression.BoolArray{true},
			"$app.NumberArrayValue": expression.NumberArray{42},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}

func TestVariablesWithQuotes(t *testing.T) {
	var schema Schema

	var ss []dcn.Schema
	err := json.Unmarshal(variablesWithQuotes, &ss)
	if err != nil {
		t.Fatalf("Error parsing schema: %v", err)
	}
	schema = SchemaFromDCN(ss)

	t.Run("Generates input from Custom Structure", func(t *testing.T) {
		app := map[string]any{
			"\"quoted2\"": map[string]any{
				"findme": "findme",
			},
		}
		got := schema.CustomInput("read", "data", app, nil)
		want := expression.Input{
			"$dcl.action":                     expression.String("read"),
			"$dcl.resource":                   expression.String("data"),
			"$app.\"\\\"quoted2\\\"\".findme": expression.String("findme"),
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})

	// t.Run("schema.Set ignores undefined fields", func(t *testing.T) {
	// 	input := expression.Input{
	// 		"$dcl.action":                     expression.String("read"),
	// 		"$dcl.resource":                   expression.String("data"),
	// 		"$app.\"\\\"quoted2\\\"\".findme": expression.String("findme"),
	// 	}

	// 	schema.Set(input, "$app", expression.UNKNOWN)
	// 	schema.Set(input, "$app.not_defined", expression.UNKNOWN)

	// 	want := expression.Input{
	// 		"$dcl.action":                     expression.String("read"),
	// 		"$dcl.resource":                   expression.String("data"),
	// 		"$app.\"\\\"quoted2\\\"\".findme": expression.UNKNOWN,
	// 	}

	// 	if !reflect.DeepEqual(input, want) {
	// 		t.Errorf("Expected %v, got %v", want, input)
	// 	}
	// })

	t.Run("type mapping edge cases", func(t *testing.T) {
		got := mapType("Structure")
		want := STRUCTURE
		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
		got = mapType("asdifguh")
		want = UNDEFINED
		if got != want {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}
