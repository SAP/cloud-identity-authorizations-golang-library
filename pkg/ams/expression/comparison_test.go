package expression

import (
	"reflect"
	"testing"
)

const DCNFalse = Bool(false)
const DCNTrue = Bool(true)

func TestEq(t *testing.T) {
	t.Run("TestEq with constants", func(t *testing.T) {
		e := Eq{Args: []Expression{Number(1), Number(2)}}
		if got, want := ToString(e), "eq(1, 2)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("TestEq with variable and constant", func(t *testing.T) {
		e := Eq{Args: []Expression{Variable{Name: "a"}, Number(1)}}
		if got, want := ToString(e), "eq(a, 1)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("TestEq with variables", func(t *testing.T) {
		e := Eq{Args: []Expression{Variable{Name: "a"}, Variable{Name: "b"}}}
		if got, want := ToString(e), "eq(a, b)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := ToString(e.Evaluate(Input{"a": Number(1), "b": UNKNOWN})), "eq(1, b)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1)}), UNSET; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1), "b": Number(1)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1), "b": Number(2)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2), "b": Number(1)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		if got, want := e.Evaluate(Input{"a": String("a"), "b": IGNORE}), IGNORE; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"b": IGNORE}), UNSET; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
func TestLt(t *testing.T) {
	t.Run("TestLt with constants", func(t *testing.T) {
		e := Lt{Args: []Expression{Number(1), Number(2)}}
		if got, want := ToString(e), "lt(1, 2)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("TestLt with variable and constant", func(t *testing.T) {
		e := Lt{Args: []Expression{Variable{Name: "a"}, Number(2)}}
		if got, want := ToString(e), "lt(a, 2)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("TestLt with variables", func(t *testing.T) {
		e := Lt{Args: []Expression{Variable{Name: "a"}, Variable{Name: "b"}}}
		if got, want := ToString(e), "lt(a, b)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := ToString(e.Evaluate(Input{"a": Number(1), "b": UNKNOWN})), "lt(1, b)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1)}), UNSET; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1), "b": Number(2)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2), "b": Number(1)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2), "b": Number(2)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		if got, want := e.Evaluate(Input{"a": String("a"), "b": String("b")}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": String("b"), "b": String("a")}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": String("a"), "b": String("a")}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
func TestBetween(t *testing.T) {
	t.Run("TestBetween with constants", func(t *testing.T) {
		e := Between{Args: []Expression{Number(5), Number(1), Number(10)}}
		if got, want := ToString(e), "between(5, 1, 10)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("TestBetween with variable and constants", func(t *testing.T) {
		e := Between{Args: []Expression{Variable{Name: "a"}, Number(1), Number(10)}}
		if got, want := ToString(e), "between(a, 1, 10)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(5)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(0)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(11)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("TestBetween with variables", func(t *testing.T) {
		e := Between{Args: []Expression{Variable{Name: "a"}, Variable{Name: "b"}, Variable{Name: "c"}}}
		if got, want := ToString(e), "between(a, b, c)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := ToString(e.Evaluate(Input{"a": Number(5), "b": Number(1), "c": UNKNOWN})), "between(5, 1, c)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(5), "b": Number(1)}), UNSET; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(5), "b": Number(1), "c": Number(10)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(0), "b": Number(1), "c": Number(10)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(11), "b": Number(1), "c": Number(10)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		if got, want := e.Evaluate(Input{"a": String("m"), "b": String("a"), "c": String("z")}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": String("a"), "b": String("m"), "c": String("z")}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": String("z"), "b": String("a"), "c": String("m")}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
func TestNe(t *testing.T) {
	t.Run("TestNe with constants", func(t *testing.T) {
		e := Ne{Args: []Expression{Number(1), Number(2)}}
		if got, want := ToString(e), "ne(1, 2)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("TestNe with variable and constant", func(t *testing.T) {
		e := Ne{Args: []Expression{Variable{Name: "a"}, Number(1)}}
		if got, want := ToString(e), "ne(a, 1)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("TestNe with variables", func(t *testing.T) {
		e := Ne{Args: []Expression{Variable{Name: "a"}, Variable{Name: "b"}}}
		if got, want := ToString(e), "ne(a, b)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := ToString(e.Evaluate(Input{"a": Number(1), "b": UNKNOWN})), "ne(1, b)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1)}), UNSET; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1), "b": Number(1)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1), "b": Number(2)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2), "b": Number(1)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestLe(t *testing.T) {
	t.Run("TestLe with constants", func(t *testing.T) {
		e := Le{Args: []Expression{Number(1), Number(2)}}
		if got, want := ToString(e), "le(1, 2)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("TestLe with variable and constant", func(t *testing.T) {
		e := Le{Args: []Expression{Variable{Name: "a"}, Number(2)}}
		if got, want := ToString(e), "le(a, 2)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(3)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("TestLe with variables", func(t *testing.T) {
		e := Le{Args: []Expression{Variable{Name: "a"}, Variable{Name: "b"}}}
		if got, want := ToString(e), "le(a, b)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := ToString(e.Evaluate(Input{"a": Number(1), "b": UNKNOWN})), "le(1, b)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1)}), UNSET; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1), "b": Number(2)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2), "b": Number(1)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2), "b": Number(2)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
func TestGt(t *testing.T) {
	t.Run("TestGt with constants", func(t *testing.T) {
		e := Gt{Args: []Expression{Number(2), Number(1)}}
		if got, want := ToString(e), "gt(2, 1)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("TestGt with variable and constant", func(t *testing.T) {
		e := Gt{Args: []Expression{Variable{Name: "a"}, Number(1)}}
		if got, want := ToString(e), "gt(a, 1)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("TestGt with variables", func(t *testing.T) {
		e := Gt{Args: []Expression{Variable{Name: "a"}, Variable{Name: "b"}}}
		if got, want := ToString(e), "gt(a, b)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := ToString(e.Evaluate(Input{"a": Number(2), "b": UNKNOWN})), "gt(2, b)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2)}), UNSET; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2), "b": Number(1)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1), "b": Number(2)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2), "b": Number(2)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestGe(t *testing.T) {
	t.Run("TestGe with constants", func(t *testing.T) {
		e := Ge{Args: []Expression{Number(2), Number(1)}}
		if got, want := ToString(e), "ge(2, 1)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("TestGe with variable and constant", func(t *testing.T) {
		e := Ge{Args: []Expression{Variable{Name: "a"}, Number(2)}}
		if got, want := ToString(e), "ge(a, 2)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("TestGe with variables", func(t *testing.T) {
		e := Ge{Args: []Expression{Variable{Name: "a"}, Variable{Name: "b"}}}
		if got, want := ToString(e), "ge(a, b)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := ToString(e.Evaluate(Input{"a": Number(2), "b": UNKNOWN})), "ge(2, b)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2)}), UNSET; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2), "b": Number(1)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(1), "b": Number(2)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(2), "b": Number(2)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}

func TestNotBetween(t *testing.T) {
	t.Run("TestNotBetween with constants", func(t *testing.T) {
		e := NotBetween{Args: []Expression{Number(5), Number(1), Number(10)}}
		if got, want := ToString(e), "not_between(5, 1, 10)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("TestNotBetween with variable and constants", func(t *testing.T) {
		e := NotBetween{Args: []Expression{Variable{Name: "a"}, Number(1), Number(10)}}
		if got, want := ToString(e), "not_between(a, 1, 10)"; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		r := e.Evaluate(Input{"a": UNKNOWN})
		if !reflect.DeepEqual(r, e) {
			t.Errorf("got %v, want %v", r, e)
		}
		if got, want := e.Evaluate(Input{"a": UNSET}), UNSET; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(5)}), DCNFalse; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(0)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		if got, want := e.Evaluate(Input{"a": Number(11)}), DCNTrue; got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
