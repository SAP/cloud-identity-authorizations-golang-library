package httpclient

import (
	"reflect"
)

func ConvertInput(input any) map[string]any {
	v := reflect.ValueOf(input)
	kind := v.Kind()
	if kind == reflect.Interface || kind == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		return ConvertInput(v.Elem())
	}
	result := make(map[string]any)
	if kind == reflect.Struct {
		for i := range v.NumField() {
			field := v.Type().Field(i)
			name := field.Tag.Get("ams")
			if name == "" {
				name = field.Name
			}
			result[name] = convertInput(v.Field(i).Interface())
		}
		return result
	}
	if kind == reflect.Map {
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key()
			value := iter.Value()
			if key.Kind() == reflect.String {
				result[key.String()] = convertInput(value.Interface())
			}
		}
		return result
	}

	return nil
}

func convertInput(input any) any {
	v := reflect.ValueOf(input)
	kind := v.Kind()
	if kind == reflect.Interface || kind == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		return ConvertInput(v.Elem())
	}
	if kind == reflect.Struct {
		result := make(map[string]any)
		for i := range v.NumField() {
			field := v.Type().Field(i)
			name := field.Tag.Get("ams")
			if name == "" {
				name = field.Name
			}
			result[name] = convertInput(v.Field(i).Interface())
		}
		return result
	}
	if kind == reflect.Map {
		result := make(map[string]any)
		iter := v.MapRange()
		for iter.Next() {
			key := iter.Key()
			value := iter.Value()
			if key.Kind() == reflect.String {
				result[key.String()] = convertInput(value.Interface())
			}
		}
		return result
	}
	return input
}
