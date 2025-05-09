package expression

import (
	"fmt"
	"reflect"
	"testing"
)

func TestUnknownIgnore(t *testing.T) { //nolint:dupl,maintidx
	empty := map[string]bool{}
	xl := map[string]bool{"x": true}
	yl := map[string]bool{"y": true}
	t.Run("eq(unset,x) => unset => false", func(t *testing.T) {
		e := Eq(Ref("y"), Ref("x"))
		result := UnknownIgnore(e, xl, empty)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})
	t.Run("eq(ignore,x) => ignore => true", func(t *testing.T) {
		e := Eq(Ref("y"), Ref("x"))
		result := UnknownIgnore(e, xl, yl)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})
	t.Run("eq(unset,ignore) => unset => false", func(t *testing.T) {
		e := Eq(Ref("x"), Ref("y"))
		result := UnknownIgnore(e, empty, yl)
		if result != FALSE {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("is_null(unset) => true", func(t *testing.T) {
		e := IsNull(Ref("x"))
		result := UnknownIgnore(e, empty, empty)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})
	t.Run("is_not_null(unset) => false", func(t *testing.T) {
		e := IsNotNull(Ref("x"))
		result := UnknownIgnore(e, empty, empty)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})
	t.Run("is_null(ignore) => ignore => true", func(t *testing.T) {
		e := IsNull(Ref("x"))
		result := UnknownIgnore(e, empty, xl)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})
	t.Run("is_not_null(ignore) => ignore => true", func(t *testing.T) {
		e := IsNotNull(Ref("x"))
		result := UnknownIgnore(e, empty, xl)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("false and ignore => false", func(t *testing.T) {
		e := And(Bool(false), Ref("x"))
		result := UnknownIgnore(e, empty, xl)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})
	t.Run("not(false and ignore) => not(false) or not(ignore) => true or ignore => true", func(t *testing.T) {
		e := Not(And(Ref("x"), Bool(false)))
		result := UnknownIgnore(e, empty, xl)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("false and unset =>  false", func(t *testing.T) {
		e := And(Bool(false), Ref("x"))
		result := UnknownIgnore(e, empty, empty)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})
	t.Run("not(false and unset) => true", func(t *testing.T) {
		e := Not(And(Bool(false), Ref("x")))
		result := UnknownIgnore(e, empty, empty)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})
	t.Run("true and unset => unset => false", func(t *testing.T) {
		e := And(TRUE, Ref("x"))
		result := UnknownIgnore(e, empty, empty)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})
	t.Run("not(true and unset) => unset => false", func(t *testing.T) {
		e := Not(And(TRUE, Ref("x")))
		result := UnknownIgnore(e, empty, empty)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})
	t.Run("ignore and unset => unset => false", func(t *testing.T) {
		e := And(Ref("x"), Ref("y"))
		result := UnknownIgnore(e, empty, xl)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})
	t.Run("not(ignore and unset) => not(ignore) or not(unset) => ignore or unset => ignore => true", func(t *testing.T) {
		e := Not(And(Ref("x"), Ref("y")))
		result := UnknownIgnore(e, empty, xl)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})
	t.Run("ignore and true => ignore => true", func(t *testing.T) {
		e := And(Ref("x"), TRUE)
		result := UnknownIgnore(e, empty, xl)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})
	t.Run("not(ignore and true) => not(ignore) or not(true) => ignore or false => ignore => true", func(t *testing.T) {
		e := Not(And(Ref("x"), TRUE))
		result := UnknownIgnore(e, empty, xl)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("ignore and y => y", func(t *testing.T) {
		e := And(Ref("x"), Ref("y"))
		result := UnknownIgnore(e, yl, xl)
		if result != Ref("y") {
			t.Errorf("Expected y, got %v", result)
		}
	})
	t.Run("not(ignore and y) => not(ignore) or not(y) => ignore => true", func(t *testing.T) {
		e := Not(And(Ref("x"), Ref("y")))
		result := UnknownIgnore(e, yl, xl)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("unset and y => unset => false", func(t *testing.T) {
		e := And(Ref("x"), Ref("y"))
		result := UnknownIgnore(e, yl, empty)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})
	t.Run("not(unset and y) => not(unset) or not(y) => unset or not(y) =>not(y)", func(t *testing.T) {
		e := Not(And(Ref("x"), Ref("y")))
		result := UnknownIgnore(e, yl, empty)
		want := Not(Ref("y"))
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Expected not(y), got %v", result)
		}
	})
	t.Run("false or unset => unset => false", func(t *testing.T) {
		e := Or(Bool(false), Ref("x"))
		result := UnknownIgnore(e, empty, empty)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})
	t.Run("not(false or unset) => not(false) and not(unset) => true and unset => unset => false", func(t *testing.T) {
		e := Not(Or(Bool(false), Ref("x")))
		result := UnknownIgnore(e, empty, empty)
		if result != FALSE {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("true or unset => true", func(t *testing.T) {
		e := Or(TRUE, Ref("x"))
		result := UnknownIgnore(e, empty, empty)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("not(true or unset) => not(true) and not(unset) => false and unset => unset => false", func(t *testing.T) {
		e := Not(Or(TRUE, Ref("x")))
		result := UnknownIgnore(e, empty, empty)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
		e = Not(e)
		result = UnknownIgnore(e, empty, empty)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("ignore or false=> ignore => true", func(t *testing.T) {
		e := Or(Ref("x"), Bool(false))
		result := UnknownIgnore(e, empty, xl)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})
	t.Run("not(ignore or false) => not(ignore) and not(false) => ignore and true => ignore => true", func(t *testing.T) {
		e := Not(Or(Ref("x"), Bool(false)))
		result := UnknownIgnore(e, empty, xl)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
		e = Not(e)
		result = UnknownIgnore(e, empty, xl)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("ignore or unset => ignore => true", func(t *testing.T) {
		e := Or(Ref("x"), Ref("y"))
		result := UnknownIgnore(e, empty, xl)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})
	t.Run("not(ignore or unset) => not(ignore) and not(unset) => ignore and unset => unset => false", func(t *testing.T) {
		e := Not(Or(Ref("x"), Ref("y")))
		result := UnknownIgnore(e, empty, xl)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
		e = Not(e)
		result = UnknownIgnore(e, empty, xl)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("ignore or y => ignore => true", func(t *testing.T) {
		e := Or(Ref("x"), Ref("y"))
		result := UnknownIgnore(e, yl, xl)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})
	t.Run("not(ignore or y) => not(ignore) and not(y) => ignore and not(y) => not(y)", func(t *testing.T) {
		e := Not(Or(Ref("x"), Ref("y")))
		result := UnknownIgnore(e, yl, xl)
		want := Not(Ref("y"))
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Expected not(y), got %v", result)
		}
		e = Not(e)
		result = UnknownIgnore(e, yl, xl)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("unset or y => y", func(t *testing.T) {
		e := Or(Ref("x"), Ref("y"))
		result := UnknownIgnore(e, yl, empty)
		want := Ref("y")
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Expected y, got %v", result)
		}
	})

	t.Run("not(unset or y) => not(unset) and not(y) => unset and not(y) => unset => false", func(t *testing.T) {
		e := Not(Or(Ref("x"), Ref("y")))
		result := UnknownIgnore(e, yl, empty)
		if !reflect.DeepEqual(result, FALSE) {
			t.Errorf("Expected false, got %v", result)
		}
		e = Not(e)
		result = UnknownIgnore(e, yl, empty)
		want := Ref("y")
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Expected y, got %v", result)
		}
	})
}

func TestNullifyExcept(t *testing.T) {
	// basically exactly the same as the test for UnknownIgnore but skipping all the ones that include ignore

	empty := map[string]bool{}
	// xl := map[string]bool{"x": true}
	yl := map[string]bool{"y": true}

	t.Run("eq(unset,y) => unset => false", func(t *testing.T) {
		e := Eq(Ref("y"), Ref("x"))
		result := NullifyExcept(e, yl)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("is_null(unset) => true", func(t *testing.T) {
		e := IsNull(Ref("x"))
		result := NullifyExcept(e, empty)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})
	t.Run("is_not_null(unset) => false", func(t *testing.T) {
		e := IsNotNull(Ref("x"))
		result := NullifyExcept(e, empty)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("false and unset =>  false", func(t *testing.T) {
		e := And(Bool(false), Ref("x"))
		result := NullifyExcept(e, empty)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("not(false and unset) => true", func(t *testing.T) {
		e := Not(And(Bool(false), Ref("x")))
		result := NullifyExcept(e, empty)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("true and unset => unset => false", func(t *testing.T) {
		e := And(TRUE, Ref("x"))
		result := NullifyExcept(e, empty)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("not(true and unset) => unset => false", func(t *testing.T) {
		e := Not(And(TRUE, Ref("x")))
		result := NullifyExcept(e, empty)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("unset and y => unset => false", func(t *testing.T) {
		e := And(Ref("x"), Ref("y"))
		result := NullifyExcept(e, yl)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("not(unset and y) => not(unset) or not(y) =>unset or not(y) =>not(y)", func(t *testing.T) {
		e := Not(And(Ref("x"), Ref("y")))
		result := NullifyExcept(e, yl)
		want := Not(Ref("y"))
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Expected not(y), got %v", result)
		}
	})

	t.Run("false or unset => unset => false", func(t *testing.T) {
		e := Or(Bool(false), Ref("x"))
		result := NullifyExcept(e, empty)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("not(false or unset) => not(false) and not(unset) => true and unset => unset => false", func(t *testing.T) {
		e := Not(Or(Bool(false), Ref("x")))
		result := NullifyExcept(e, empty)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("true or unset => true", func(t *testing.T) {
		e := Or(TRUE, Ref("x"))
		result := NullifyExcept(e, empty)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("not(true or unset) => not(true) and not(unset) => false and unset => unset => false", func(t *testing.T) {
		e := Not(Or(TRUE, Ref("x")))
		result := NullifyExcept(e, empty)
		if result != FALSE {
			t.Errorf("Expected false, got %v", result)
		}
		e = Not(e)
		result = NullifyExcept(e, empty)
		if result != TRUE {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("unset or y => y", func(t *testing.T) {
		e := Or(Ref("x"), Ref("y"))
		result := NullifyExcept(e, yl)
		want := Ref("y")
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Expected y, got %v", result)
		}
	})

	t.Run("not(unset or y) => not(unset) and not(y) => unset and not(y) => unset => false", func(t *testing.T) {
		e := Not(Or(Ref("x"), Ref("y")))
		result := NullifyExcept(e, yl)
		if !reflect.DeepEqual(result, FALSE) {
			t.Errorf("Expected false, got %v", result)
		}
		e = Not(e)
		result = NullifyExcept(e, yl)
		want := Ref("y")
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Expected y, got %v", result)
		}
	})
}

func TestTrivialCases(t *testing.T) {
	// basically exactly the same as the test for UnknownIgnore but skipping all the ones that include ignore

	ignore := map[string]bool{}
	unknown := map[string]bool{"x": true}

	type test struct {
		expr Expression
		want Expression
	}

	tests := []test{
		{And(TRUE, TRUE), TRUE},
		{Not(And(TRUE, TRUE)), FALSE},
		{And(TRUE, FALSE), FALSE},
		{Or(FALSE, FALSE), FALSE},
		{Not(Or(FALSE, FALSE)), TRUE},
		{Or(FALSE, TRUE), TRUE},
		{IsNotNull(Ref("x")), IsNotNull(Ref("x"))},
		{IsNull(Ref("x")), IsNull(Ref("x"))},
		{Eq(Ref("x"), Number(1)), Eq(Ref("x"), Number(1))},
		{Function("test", nil, []Expression{}), Function("test", nil, []Expression{})},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("UNKNOWN [x] IGNORE [] on %v", tc.expr), func(t *testing.T) {
			result := UnknownIgnore(tc.expr, unknown, ignore)
			if !reflect.DeepEqual(result, tc.want) {
				t.Errorf("Expected %v, got %v", tc.want, result)
			}
		})
		t.Run(fmt.Sprintf("NullifyExcept [x] on %v", tc.expr), func(t *testing.T) {
			result := NullifyExcept(tc.expr, unknown)
			if !reflect.DeepEqual(result, tc.want) {
				t.Errorf("Expected %v, got %v", tc.want, result)
			}
		})
	}

	t.Run("evaluate on dummy expression", func(t *testing.T) {
		if null.Evaluate(Input{}) != nil {
			t.Errorf("Expected nil, got %v", null.Evaluate(nil))
		}
	})
}
