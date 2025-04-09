package internal

import (
	"reflect"
	"strings"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/util"
)

type Schema struct {
	tenantSchemas map[string]string
	inputTypes    map[string]InputType
	normalized    map[string][]string
}

type InputType byte

const (
	STRING InputType = iota
	BOOLEAN
	NUMBER
	STRING_ARRAY
	NUMBER_ARRAY
	BOOLEAN_ARRAY
	STRUCTURE
	UNDEFINED
)

func SchemaFromDCN(sc []dcn.Schema) Schema {
	result := Schema{
		tenantSchemas: make(map[string]string),
		inputTypes: map[string]InputType{
			"$dcl":          STRUCTURE,
			"$dcl.action":   STRING,
			"$dcl.resource": STRING,
		},
		normalized: make(map[string][]string),
	}

	for _, s := range sc {
		if s.Tenant != "" {
			tenantPackage := util.StringifyReference(s.QualifiedName[:len(s.QualifiedName)-1])
			result.tenantSchemas[tenantPackage] = s.Tenant
		}
		if s.Definition.Nested != nil {
			result.buildSchemaAttributes(s.Definition, []string{})
		}

	}
	return result
}

func mapType(dcnType string) InputType {
	switch dcnType {
	case "String":
		return STRING
	case "Boolean":
		return BOOLEAN
	case "Number":
		return NUMBER
	case "String[]":
		return STRING_ARRAY
	case "Boolean[]":
		return BOOLEAN_ARRAY
	case "Number[]":
		return NUMBER_ARRAY
	case "Structure":
		return STRUCTURE
	}
	return UNDEFINED

}

func (s *Schema) buildSchemaAttributes(a dcn.SchemaAttribute, path []string) {
	for k, v := range a.Nested {
		newPath := append(path, k)
		if v.Nested != nil {
			s.inputTypes[util.StringifyReference(newPath)] = STRUCTURE
			s.buildSchemaAttributes(v, newPath)
		} else {
			s.inputTypes[util.StringifyReference(newPath)] = mapType(v.Type)
		}
	}
}

func (s Schema) Set(input expression.Input, val string, value expression.Wildcard) {
	t, ok := s.inputTypes[val]
	if !ok {
		return
	}
	if t == STRUCTURE {
		for k, it := range s.inputTypes {
			if strings.HasPrefix(k, val) && it != STRUCTURE {
				input[k] = value
			}
		}
	} else {
		input[val] = value
	}
}

func (s Schema) GetTenantForQualifiedName(qn dcn.QualifiedName) string {
	p := util.StringifyReference(qn[:len(qn)-1])
	tenant, ok := s.tenantSchemas[p]
	if !ok {
		return ""
	}
	return tenant
}

func (s Schema) CustomInput(action, resource string, input any, env any) expression.Input {
	result := expression.Input{
		"$dcl.action":   expression.String(action),
		"$dcl.resource": expression.String(resource),
	}

	s.convertCustomInput(result, reflect.ValueOf(input), []string{"$app"})

	s.convertCustomInput(result, reflect.ValueOf(env), []string{"$env"})

	return result
}

func (s Schema) convertCustomInput(result expression.Input, input reflect.Value, path []string) {

	v := input
	kind := v.Kind()
	currentPath := util.StringifyReference(path)
	shouldBeType, ok := s.inputTypes[currentPath]
	if !ok {
		return
	}
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
		s.convertCustomInput(result, v.Elem(), path)
		return
	}

	switch shouldBeType {
	case STRUCTURE:
		switch kind {
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				fieldValue := v.Field(i)
				field := v.Type().Field(i)
				if !field.IsExported() {
					continue
				}
				name := field.Tag.Get("ams")
				if name == "" {
					name = field.Name
				}
				s.convertCustomInput(result, fieldValue, append(path, name))
			}
		case reflect.Map:
			if v.IsNil() {
				return
			}
			for _, k := range v.MapKeys() {
				fieldValue := v.MapIndex(k)
				s.convertCustomInput(result, fieldValue, append(path, k.String()))
			}
		}
	case STRING:
		if kind == reflect.String {
			result[currentPath] = expression.String(v.String())
		}
	case BOOLEAN:
		if kind == reflect.Bool {
			result[currentPath] = expression.Bool(v.Bool())
		}
	case NUMBER:
		switch kind {
		case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result[currentPath] = expression.Number(v.Uint())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result[currentPath] = expression.Number(v.Int())
		case reflect.Float32, reflect.Float64:
			result[currentPath] = expression.Number(v.Float())
		}
	case STRING_ARRAY:
		if kind == reflect.Slice || kind == reflect.Array {
			if v.IsNil() {
				return
			}
			r := expression.StringArray{}
			for i := 0; i < v.Len(); i++ {
				vi := v.Index(i)
				for vi.Kind() == reflect.Interface || vi.Kind() == reflect.Pointer {
					if vi.IsNil() {
						return
					}
					vi = vi.Elem()
				}
				if vi.Kind() != reflect.String {
					return
				}
				r = append(r, expression.String(vi.String()))
			}
			result[currentPath] = r
		}
	case NUMBER_ARRAY:
		if kind == reflect.Slice || kind == reflect.Array {
			if v.IsNil() {
				return
			}
			r := expression.NumberArray{}
			for i := 0; i < v.Len(); i++ {
				vi := v.Index(i)
				for vi.Kind() == reflect.Interface || vi.Kind() == reflect.Pointer {
					if vi.IsNil() {
						return
					}
					vi = vi.Elem()
				}
				switch vi.Kind() {
				case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
					r = append(r, expression.Number(vi.Uint()))
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					r = append(r, expression.Number(vi.Int()))
				case reflect.Float32, reflect.Float64:
					r = append(r, expression.Number(vi.Float()))
				default:
					return
				}
			}
			result[currentPath] = r
		}
	case BOOLEAN_ARRAY:
		if kind == reflect.Slice || kind == reflect.Array {
			if v.IsNil() {
				return
			}
			r := expression.BoolArray{}
			for i := 0; i < v.Len(); i++ {
				vi := v.Index(i)

				for vi.Kind() == reflect.Interface || vi.Kind() == reflect.Pointer {
					if vi.IsNil() {
						return
					}
					vi = vi.Elem()
				}
				if vi.Kind() != reflect.Bool {
					return
				}
				r = append(r, expression.Bool(vi.Bool()))
			}
			result[currentPath] = r
		}
	}

}
func (s Schema) PurgeInvalidInput(input expression.Input) {
	for k, v := range input {
		t, ok := s.inputTypes[k]
		if !ok {
			delete(input, k)
			continue
		}
		_, ok = v.(expression.Wildcard)
		if ok {
			continue
		}
		switch t {
		case STRING:
			_, ok := v.(expression.String)
			if !ok {
				delete(input, k)
				continue
			}
		case BOOLEAN:
			_, ok := v.(expression.Bool)
			if !ok {
				delete(input, k)
				continue
			}
		case NUMBER:
			_, ok := v.(expression.Number)
			if !ok {
				delete(input, k)
				continue
			}
		case STRING_ARRAY:
			_, ok := v.(expression.StringArray)
			if !ok {
				delete(input, k)
				continue
			}
		case NUMBER_ARRAY:
			_, ok := v.(expression.NumberArray)
			if !ok {
				delete(input, k)
				continue
			}
		case BOOLEAN_ARRAY:
			_, ok := v.(expression.BoolArray)
			if !ok {
				delete(input, k)
				continue
			}
		}

	}
}
