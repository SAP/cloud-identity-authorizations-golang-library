package expression

type Constant interface {
	Expression
	equals(c Constant) bool
	LessThan(c Constant) bool
}

type ArrayConstant interface {
	Contains(c Constant) bool
	Elements() []Constant
	IsEmpty() bool
	Constant
}

type Number float64

type String string

type Bool bool

type NumberArray []Number

type StringArray []String

type BoolArray []Bool

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

func (n NumberArray) Contains(c Constant) bool {
	for _, v := range n {
		if v.equals(c) {
			return true
		}
	}
	return false
}

func (s StringArray) Contains(c Constant) bool {
	for _, v := range s {
		if v.equals(c) {
			return true
		}
	}
	return false
}

func (b BoolArray) Contains(c Constant) bool {
	for _, v := range b {
		if v.equals(c) {
			return true
		}
	}
	return false
}

func (n NumberArray) IsEmpty() bool {
	return len(n) == 0
}

func (s StringArray) IsEmpty() bool {
	return len(s) == 0
}

func (b BoolArray) IsEmpty() bool {
	return len(b) == 0
}
func (n NumberArray) Elements() []Constant {
	result := make([]Constant, len(n))
	for i, v := range n {
		result[i] = v
	}
	return result
}

func (s StringArray) Elements() []Constant {
	result := make([]Constant, len(s))
	for i, v := range s {
		result[i] = v
	}
	return result
}

func (b BoolArray) Elements() []Constant {
	result := make([]Constant, len(b))
	for i, v := range b {
		result[i] = v
	}
	return result
}

func (s StringArray) AsString() []string {
	result := make([]string, len(s))
	for i, s := range s {
		result[i] = string(s)
	}
	return result
}

func (b BoolArray) AsBool() []bool {
	result := make([]bool, len(b))
	for i, b := range b {
		result[i] = bool(b)
	}
	return result
}

func (n NumberArray) AsFloat() []float64 {
	result := make([]float64, len(n))
	for i, n := range n {
		result[i] = float64(n)
	}
	return result
}

func (b BoolArray) Evaluate(input Input) Expression {
	return b
}

func (n NumberArray) Evaluate(input Input) Expression {
	return n
}

func (s StringArray) Evaluate(input Input) Expression {
	return s
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

func (s StringArray) equals(c Constant) bool {
	return false
}

func (s StringArray) LessThan(c Constant) bool {
	return false
}

func (b BoolArray) equals(c Constant) bool {
	return false
}

func (b BoolArray) LessThan(c Constant) bool {
	return false
}

func (n NumberArray) equals(c Constant) bool {
	return false
}

func (n NumberArray) LessThan(c Constant) bool {
	return false
}

func ArrayFrom[T string | float64 | bool](v []T) ArrayConstant {
	switch vals := any(v).(type) {
	case []string:
		result := make(StringArray, len(vals))
		for i, s := range vals {
			result[i] = String(s)
		}
		return result
	case []float64:
		result := make(NumberArray, len(vals))
		for i, n := range vals {
			result[i] = Number(n)
		}
		return result
	case []bool:
		result := make(BoolArray, len(vals))
		for i, b := range vals {
			result[i] = Bool(b)
		}
		return result
	}
	return nil
}
