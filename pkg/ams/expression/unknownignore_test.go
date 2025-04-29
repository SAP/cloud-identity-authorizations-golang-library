package expression

import (
	"reflect"
	"testing"
)

func TestUnknownIgnore(t *testing.T) { //nolint:dupl

	empty := map[string]bool{}
	xl := map[string]bool{"x": true}
	yl := map[string]bool{"y": true}
	t.Run("is_null(unset) => true", func(t *testing.T) {
		e := IsNull(Ref("x"))
		result := UnknownIgnore(e, empty, empty)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})
	t.Run("is_not_null(unset) => false", func(t *testing.T) {
		e := IsNotNull(Ref("x"))
		result := UnknownIgnore(e, empty, empty)
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})
	t.Run("is_null(ignore) => ignore => true", func(t *testing.T) {
		e := IsNull(Ref("x"))
		result := UnknownIgnore(e, empty, xl)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("false and unset =>  false", func(t *testing.T) {
		e := And(Bool(false), Ref("x"))
		result := UnknownIgnore(e, empty, empty)
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})
	t.Run("not(false and unset) => true", func(t *testing.T) {
		e := Not(And(Bool(false), Ref("x")))
		result := UnknownIgnore(e, empty, empty)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})
	t.Run("true and unset => unset => false", func(t *testing.T) {
		e := And(Bool(true), Ref("x"))
		result := UnknownIgnore(e, empty, empty)
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})
	t.Run("not(true and unset) => unset => false", func(t *testing.T) {
		e := Not(And(Bool(true), Ref("x")))
		result := UnknownIgnore(e, empty, empty)
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})
	t.Run("ignore and unset => unset => false", func(t *testing.T) {
		e := And(Ref("x"), Ref("y"))
		result := UnknownIgnore(e, empty, xl)
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})
	t.Run("not(ignore and unset) => not(ignore) or not(unset) => ignore or unset => ignore => true", func(t *testing.T) {
		e := Not(And(Ref("x"), Ref("y")))
		result := UnknownIgnore(e, empty, xl)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}

	})
	t.Run("ignore and true => ignore => true", func(t *testing.T) {
		e := And(Ref("x"), Bool(true))
		result := UnknownIgnore(e, empty, xl)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})
	t.Run("not(ignore and true) => not(ignore) or not(true) => ignore or false => ignore => true", func(t *testing.T) {
		e := Not(And(Ref("x"), Bool(true)))
		result := UnknownIgnore(e, empty, xl)
		if result != Bool(true) {
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
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("unset and y => unset => false", func(t *testing.T) {
		e := And(Ref("x"), Ref("y"))
		result := UnknownIgnore(e, yl, empty)
		if result != Bool(false) {
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
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
	})

	t.Run("true or unset => true", func(t *testing.T) {
		e := Or(Bool(true), Ref("x"))
		result := UnknownIgnore(e, empty, empty)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("ignore or unset => ignore => true", func(t *testing.T) {
		e := Or(Ref("x"), Ref("y"))
		result := UnknownIgnore(e, empty, xl)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("ignore or y => ignore => true", func(t *testing.T) {
		e := Or(Ref("x"), Ref("y"))
		result := UnknownIgnore(e, yl, xl)
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("unset or y => y", func(t *testing.T) {
		e := Or(Ref("x"), Ref("y"))
		result := UnknownIgnore(e, yl, empty)
		if result != Ref("y") {
			t.Errorf("Expected y, got %v", result)
		}
	})

}
