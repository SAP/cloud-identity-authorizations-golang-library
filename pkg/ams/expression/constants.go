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

func ConstantFrom(v any) Constant {
	switch v := v.(type) {
	case string:
		return String(v)
	case float64:
		return Number(v)
	case int:
		return Number(v)
	case int64:
		return Number(v)
	case uint:
		return Number(v)
	case uint64:
		return Number(v)
	case int8:
		return Number(v)
	case int16:
		return Number(v)
	case int32:
		return Number(v)
	case uint8:
		return Number(v)
	case uint16:
		return Number(v)
	case uint32:
		return Number(v)
	case bool:
		return Bool(v)
	case []string:
		return ArrayFrom(v)
	case []float64:
		return ArrayFrom(v)
	case []bool:
		return ArrayFrom(v)
	}
	reflectV := reflect.ValueOf(v)
	switch reflectV.Kind() { //nolint:exhaustive
	case reflect.Interface, reflect.Pointer:
		if reflectV.IsNil() {
			return nil
		}
		return ConstantFrom(reflectV.Elem().Interface())
	}
	return nil
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
