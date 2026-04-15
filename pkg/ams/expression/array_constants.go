package expression

import "fmt"

type ArrayConstant interface {
	Contains(c Constant) bool
	Elements() []Constant
	IsEmpty() bool
	Constant
}

type NumberArray []Number

type StringArray []String

type BoolArray []Bool

type EmptyArray struct{}

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

func (b BoolArray) String() string {
	return fmt.Sprintf("%v", b.AsBool())
}

func (n NumberArray) Evaluate(input Input) Expression {
	return n
}

func (n NumberArray) String() string {
	return fmt.Sprintf("%v", n.AsFloat())
}

func (s StringArray) Evaluate(input Input) Expression {
	return s
}

func (s StringArray) String() string {
	return fmt.Sprintf("%v", s.Elements())
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

func (n EmptyArray) Contains(c Constant) bool {
	return false
}

func (n EmptyArray) IsEmpty() bool {
	return true
}

func (n EmptyArray) Elements() []Constant {
	return []Constant{}
}

func (n EmptyArray) Evaluate(input Input) Expression {
	return n
}

func (n EmptyArray) String() string {
	return "[]"
}

func (n EmptyArray) equals(c Constant) bool {
	return false
}

func (n EmptyArray) LessThan(c Constant) bool {
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
	return EmptyArray{}
}
