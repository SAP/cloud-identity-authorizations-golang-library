package expression

import "fmt"

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
	case bool:
		return Bool(v)
	case []string:
		result := make(StringArray, len(v))
		for i, s := range v {
			result[i] = String(s)
		}
		return result
	case []float64:
		result := make(NumberArray, len(v))
		for i, n := range v {
			result[i] = Number(n)
		}
		return result
	case []bool:
		result := make(BoolArray, len(v))
		for i, b := range v {
			result[i] = Bool(b)
		}
		return result
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
