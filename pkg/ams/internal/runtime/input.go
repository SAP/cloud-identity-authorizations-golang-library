package runtime

import (
	"reflect"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/util"
)

func (r *Runtime) SanitizeInput(input expression.Input, resource string) ([]string, []string) {
	unknownFields, wrongTypeFields := r.Resources[WildCard].Attributes.validateInput(input)

	deleteUnknownFields := func() {
		for _, field := range unknownFields {
			delete(input, field)
		}
	}
	defer deleteUnknownFields()

	resourceDef, ok := r.Resources[resource]
	if !ok {
		return unknownFields, wrongTypeFields
	}
	unknownFieldsResource, wrongTypeFieldsResource := resourceDef.Attributes.validateInput(input)

	// intersect unknownFields with unknownFieldsResource
	intersection := make(map[string]struct{})
	for _, field := range unknownFields {
		intersection[field] = struct{}{}
	}
	unknownFields = make([]string, 0)
	for _, field := range unknownFieldsResource {
		if _, ok := intersection[field]; ok {
			unknownFields = append(unknownFields, field)
		} else {
			delete(intersection, field)
		}
	}
	// wrongTypeFields is the union of wrongTypeFieldsGlobal and wrongTypeFieldsResource
	wrongTypeFields = append(wrongTypeFields, wrongTypeFieldsResource...)

	return unknownFields, wrongTypeFields
}

func (a Attributes) validateInput(input expression.Input) ([]string, []string) {
	unknownFields := make([]string, 0)
	wrongTypeFields := make([]string, 0)

	for k, v := range input {
		attr, ok := a[k]
		if !ok {
			unknownFields = append(unknownFields, k)
			continue
		}
		switch attr.Type {
		case attrTypeString:
			if _, ok := v.(expression.String); !ok {
				wrongTypeFields = append(wrongTypeFields, k)
				delete(input, k)
			}
		case attrTypeNumber:
			if _, ok := v.(expression.Number); !ok {
				wrongTypeFields = append(wrongTypeFields, k)
				delete(input, k)
			}
		case attrTypeBool:
			if _, ok := v.(expression.Bool); !ok {
				wrongTypeFields = append(wrongTypeFields, k)
				delete(input, k)
			}
		case attrTypeArray:
			if _, ok := v.(expression.EmptyArray); ok {
				continue
			}
			switch attr.ElementType {
			case attrTypeString:
				if _, ok := v.(expression.StringArray); !ok {
					wrongTypeFields = append(wrongTypeFields, k)
					delete(input, k)
				}
			case attrTypeNumber:
				if _, ok := v.(expression.NumberArray); !ok {
					wrongTypeFields = append(wrongTypeFields, k)
					delete(input, k)
				}
			case attrTypeBool:
				if _, ok := v.(expression.BoolArray); !ok {
					wrongTypeFields = append(wrongTypeFields, k)
					delete(input, k)
				}
			case attrTypeArray:
				wrongTypeFields = append(wrongTypeFields, k)
				delete(input, k)
			}
		}
	}
	return unknownFields, wrongTypeFields
}

func CreateInput(input any) expression.Input {
	result := make(expression.Input)
	reflectValue := reflect.ValueOf(input)
	InsertInput(result, reflectValue, []string{"$app"})
	InsertInput(result, reflectValue, []string{"$resource"})
	return result
}

func InsertInput(input expression.Input, val reflect.Value, path []string) {
	switch val.Kind() { //nolint:exhaustive
	case reflect.Struct:
		for i := range val.NumField() {
			field := val.Type().Field(i)
			fieldValue := val.Field(i)
			if field.PkgPath != "" { // unexported field
				continue
			}
			name := field.Tag.Get("ams")
			if name == "" {
				continue
			}
			InsertInput(input, fieldValue, append(path, name))
		}
	case reflect.Map:
		iter := val.MapRange()
		for iter.Next() {
			key := iter.Key()
			value := iter.Value()
			if key.Kind() != reflect.String {
				continue
			}
			InsertInput(input, value, append(path, key.String()))
		}
	case reflect.Interface, reflect.Pointer:
		if val.IsNil() {
			return
		}
		InsertInput(input, val.Elem(), path)
	default:
		constant, err := expression.ConstantFrom(val.Interface())
		if err != nil {
			return
		}
		currentPath := util.StringifyQualifiedName(path)
		input[currentPath] = constant
	}
}
