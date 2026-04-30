package httpclient

import (
	"reflect"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/util"
)

func insertCustomInput(result reqInput, input reflect.Value, path []string) {
	v := input
	kind := v.Kind()
	currentPath := util.StringifyQualifiedName(path)

	if kind == reflect.Invalid {
		return
	}

	// first we resolve pointers and interfaces
	if kind == reflect.Interface || kind == reflect.Pointer {
		if v.IsNil() {
			return
		}
		c, ok := v.Interface().(expression.Constant)
		if ok {
			result[currentPath] = c
			return
		}
		insertCustomInput(result, v.Elem(), path)
		return
	}
	switch kind { //nolint:exhaustive
	case reflect.Struct:
		for i := range v.NumField() {
			fieldValue := v.Field(i)
			field := v.Type().Field(i)
			if !field.IsExported() {
				continue
			}
			name := field.Tag.Get("ams")
			if name == "" {
				name = field.Name
			}
			insertCustomInput(result, fieldValue, append(path, name))
		}
	case reflect.Map:
		if v.IsNil() {
			return
		}
		for _, k := range v.MapKeys() {
			fieldValue := v.MapIndex(k)
			insertCustomInput(result, fieldValue, append(path, k.String()))
		}
	case reflect.Slice, reflect.Array:
		if input.IsNil() {
			return
		}
		result[currentPath] = v.Interface()
	default:
		result[currentPath] = v.Interface()
	}
}
