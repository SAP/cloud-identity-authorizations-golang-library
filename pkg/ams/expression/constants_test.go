package expression

import (
	"reflect"
	"testing"
)

func TestDCNNumber_Equals(t *testing.T) {
	tests := []struct {
		name     string
		n        Number
		c        Constant
		expected bool
	}{
		{"Equal numbers", Number(5), Number(5), true},
		{"Different numbers", Number(5), Number(6), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.Equals(tt.c); got != tt.expected {
				t.Errorf("DCNNumber.Equals() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestDCNNumber_LessThan(t *testing.T) {
	tests := []struct {
		name     string
		n        Number
		c        Constant
		expected bool
	}{
		{"Less than", Number(5), Number(6), true},
		{"Not less than", Number(6), Number(5), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.LessThan(tt.c); got != tt.expected {
				t.Errorf("DCNNumber.LessThan() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestDCNString_Equals(t *testing.T) {
	tests := []struct {
		name     string
		s        String
		c        Constant
		expected bool
	}{
		{"Equal strings", String("test"), String("test"), true},
		{"Different strings", String("test"), String("different"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Equals(tt.c); got != tt.expected {
				t.Errorf("DCNString.Equals() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestDCNString_LessThan(t *testing.T) {
	tests := []struct {
		name     string
		s        String
		c        Constant
		expected bool
	}{
		{"Less than", String("a"), String("b"), true},
		{"Not less than", String("b"), String("a"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.LessThan(tt.c); got != tt.expected {
				t.Errorf("DCNString.LessThan() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestDCNBool_Equals(t *testing.T) {
	tests := []struct {
		name     string
		b        Bool
		c        Constant
		expected bool
	}{
		{"Equal bools", Bool(true), Bool(true), true},
		{"Different bools", Bool(true), Bool(false), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.Equals(tt.c); got != tt.expected {
				t.Errorf("DCNBool.Equals() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestDCNBool_LessThan(t *testing.T) {
	tests := []struct {
		name     string
		b        Bool
		c        Constant
		expected bool
	}{
		{"Less than", Bool(false), Bool(true), true},
		{"Not less than", Bool(true), Bool(false), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.LessThan(tt.c); got != tt.expected {
				t.Errorf("DCNBool.LessThan() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestDCNArrayConstant_Equals(t *testing.T) {
	tests := []struct {
		name     string
		arr      Constant
		c        Constant
		expected bool
	}{
		{"Equal arrays", NumberArray{1, 2, 3}, NumberArray{1, 2, 3}, false},
		{"Equal string arrays", StringArray{"a", "b", "c"}, StringArray{"a", "b", "c"}, false},
		{"Equal bool arrays", BoolArray{true, false}, BoolArray{true, false}, false},
		{"Different arrays", NumberArray{1, 2, 3}, NumberArray{1, 2, 4}, false},
		{"Different string arrays", StringArray{"a", "b", "c"}, StringArray{"a", "b", "d"}, false},
		{"Different bool arrays", BoolArray{true, false}, BoolArray{true, true}, false},
		{"Different BoolPlus", UNSET, IGNORE, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arr.Equals(tt.c); got != tt.expected {
				t.Errorf("DCNArrayConstant.Equals() = %v, expected %v", got, tt.expected)
			}
			if got := tt.arr.LessThan(tt.arr); got != tt.expected {
				t.Errorf("DCNArrayConstant.LessThan() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestDCNNumberArray_Contains(t *testing.T) {
	tests := []struct {
		name     string
		arr      NumberArray
		c        Constant
		expected bool
	}{
		{"Contains number", NumberArray{1, 2, 3}, Number(2), true},
		{"Does not contain number", NumberArray{1, 2, 3}, Number(4), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arr.Contains(tt.c); got != tt.expected {
				t.Errorf("DCNNumberArray.Contains() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestDCNStringArray_Contains(t *testing.T) {
	tests := []struct {
		name     string
		arr      StringArray
		c        Constant
		expected bool
	}{
		{"Contains string", StringArray{"a", "b", "c"}, String("b"), true},
		{"Does not contain string", StringArray{"a", "b", "c"}, String("d"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arr.Contains(tt.c); got != tt.expected {
				t.Errorf("DCNStringArray.Contains() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestDCNBoolArray_Contains(t *testing.T) {
	tests := []struct {
		name     string
		arr      BoolArray
		c        Constant
		expected bool
	}{
		{"Contains bool", BoolArray{true, false}, Bool(true), true},
		{"Does not contain bool", BoolArray{false}, Bool(true), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arr.Contains(tt.c); got != tt.expected {
				t.Errorf("DCNBoolArray.Contains() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func TestDCNArrayConstant_Elements(t *testing.T) {
	tests := []struct {
		name     string
		arr      ArrayConstant
		expected []Constant
	}{
		{"Number array", NumberArray{1, 2, 3}, []Constant{Number(1), Number(2), Number(3)}},
		{"String array", StringArray{"a", "b", "c"}, []Constant{String("a"), String("b"), String("c")}},
		{"Bool array", BoolArray{true, false}, []Constant{Bool(true), Bool(false)}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.arr.Elements(); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("DCNArrayConstant.Elements() = %v, expected %v", got, tt.expected)
			}
		})
	}
}
