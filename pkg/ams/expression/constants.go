package expression

import (
	"fmt"
	"reflect"
)

type Type byte

const (
	TypeString Type = iota
	TypeNumber
	TypeBool
	TypeStringArray
	TypeNumberArray
	TypeBoolArray
)

type Constant interface {
	Expression
	equals(c Constant) bool
	LessThan(c Constant) bool
}

type Number float64

type String string

type Bool bool

const (
	TRUE  = Bool(true)
	FALSE = Bool(false)
)

func ConstantFrom(v any) (Constant, error) {
	switch v := v.(type) {
	case string:
		return String(v), nil
	case float64:
		return Number(v), nil
	case int:
		return Number(v), nil
	case int64:
		return Number(v), nil
	case uint:
		return Number(v), nil
	case uint64:
		return Number(v), nil
	case int8:
		return Number(v), nil
	case int16:
		return Number(v), nil
	case int32:
		return Number(v), nil
	case uint8:
		return Number(v), nil
	case uint16:
		return Number(v), nil
	case uint32:
		return Number(v), nil
	case bool:
		return Bool(v), nil
	}
	reflectV := reflect.ValueOf(v)
	switch reflectV.Kind() { //nolint:exhaustive
	case reflect.Interface, reflect.Pointer:
		if reflectV.IsNil() {
			return nil, fmt.Errorf("unsupported constant nil")
		}
		return ConstantFrom(reflectV.Elem().Interface())
	case reflect.Slice, reflect.Array:
		if reflectV.Len() == 0 {
			return EmptyArray{}, nil
		}
		firstElement, err := ConstantFrom(reflectV.Index(0).Interface())
		if err != nil {
			return nil, err
		}
		var ok bool
		switch firstElement.(type) {
		case String:
			result := make([]String, reflectV.Len())
			for i := range reflectV.Len() {
				element, err := ConstantFrom(reflectV.Index(i).Interface())
				if err != nil {
					return nil, err
				}
				result[i], ok = element.(String)
				if !ok {
					return nil, fmt.Errorf("mixed types in array: %T and %T", firstElement, element)
				}
			}
			return StringArray(result), nil
		case Number:
			result := make([]Number, reflectV.Len())
			for i := range reflectV.Len() {
				element, err := ConstantFrom(reflectV.Index(i).Interface())
				if err != nil {
					return nil, err
				}
				result[i], ok = element.(Number)
				if !ok {
					return nil, fmt.Errorf("mixed types in array: %T and %T", firstElement, element)
				}
			}
			return NumberArray(result), nil
		case Bool:
			result := make([]Bool, reflectV.Len())
			for i := range reflectV.Len() {
				element, err := ConstantFrom(reflectV.Index(i).Interface())
				if err != nil {
					return nil, err
				}
				result[i], ok = element.(Bool)
				if !ok {
					return nil, fmt.Errorf("mixed types in array: %T and %T", firstElement, element)
				}
			}
			return BoolArray(result), nil
		default:
			return nil, fmt.Errorf("unsupported array element type: %T", firstElement)
		}
	}
	return nil, fmt.Errorf("unsupported constant type: %T", v)
}

func (n Number) equals(c Constant) bool {
	return n == c.(Number) //nolint:forcetypeassert
}

func (n Number) LessThan(c Constant) bool {
	n2 := c.(Number) //nolint:forcetypeassert
	return n < n2    //nolint:forcetypeassert
}

func (n Number) String() string {
	return fmt.Sprintf("%v", float64(n))
}

func (s String) equals(c Constant) bool {
	return s == c.(String) //nolint:forcetypeassert
}

func (s String) LessThan(c Constant) bool {
	return s < c.(String) //nolint:forcetypeassert
}

func (b Bool) equals(c Constant) bool {
	return b == c.(Bool) //nolint:forcetypeassert
}

func (b Bool) LessThan(c Constant) bool {
	return bool(!b && c.(Bool)) //nolint:forcetypeassert
}

func (b Bool) String() string {
	return fmt.Sprintf("%v", bool(b))
}

func (n Number) Evaluate(input Input) Expression {
	return n
}

func (s String) Evaluate(input Input) Expression {
	return s
}

func (s String) String() string {
	return `"` + string(s) + `"`
}

func (b Bool) Evaluate(input Input) Expression {
	return b
}
