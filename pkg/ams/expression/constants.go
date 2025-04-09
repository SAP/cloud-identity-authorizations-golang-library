package expression

type Constant interface {
	Expression
	Equals(Constant) bool
	LessThan(Constant) bool
}

type ArrayConstant interface {
	Contains(Constant) bool
	IsEmpty() bool
	Elements() []Constant
	Constant
}

type Number float64

type String string

type Bool bool

type NumberArray []Number

type StringArray []String

type BoolArray []Bool

const TRUE = Bool(true)
const FALSE = Bool(false)

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
	return UNSET
}

func (n Number) Equals(c Constant) bool {
	return n == c.(Number)
}

func (n Number) LessThan(c Constant) bool {
	return n < c.(Number)
}

func (s String) Equals(c Constant) bool {
	return s == c.(String)
}

func (s String) LessThan(c Constant) bool {
	return s < c.(String)
}

func (b Bool) Equals(c Constant) bool {
	return b == c.(Bool)
}

func (b Bool) LessThan(c Constant) bool {

	return bool(!b && c.(Bool))
}

func (n NumberArray) Contains(c Constant) bool {
	for _, v := range n {
		if v.Equals(c) {
			return true
		}
	}
	return false
}

func (s StringArray) Contains(c Constant) bool {
	for _, v := range s {
		if v.Equals(c) {
			return true
		}
	}
	return false
}

func (b BoolArray) Contains(c Constant) bool {
	for _, v := range b {
		if v.Equals(c) {
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

func (b Bool) Evaluate(input Input) Expression {
	return b
}

func (s StringArray) Equals(c Constant) bool {
	return false
}

func (s StringArray) LessThan(c Constant) bool {
	return false
}

func (b BoolArray) Equals(c Constant) bool {
	return false
}

func (b BoolArray) LessThan(c Constant) bool {
	return false
}

func (n NumberArray) Equals(c Constant) bool {
	return false
}

func (n NumberArray) LessThan(c Constant) bool {
	return false
}
