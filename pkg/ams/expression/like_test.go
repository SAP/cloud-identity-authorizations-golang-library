package expression

import (
	"reflect"
	"testing"
)

func TestLike(t *testing.T) {
	t.Run("string variable like constant String", func(t *testing.T) { //nolint:dupl
		like := Like(Ref("x"), String("a"))
		result := like.Evaluate(Input{"x": String("a")})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = like.Evaluate(Input{"x": String("c")})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = like.Evaluate(Input{})
		if !reflect.DeepEqual(result, like) {
			t.Errorf("Expected %v, got %v", like, result)
		}

		if ToString(like) != "like(x, \"a\")" {
			t.Errorf("Expected like(x, a), got %v", ToString(like))
		}
	})

	t.Run("evaluate like with _ as escape character", func(t *testing.T) { //nolint:dupl
		like := Like(Ref("x"), String("a"), String("_"))
		result := like.Evaluate(Input{"x": String("a")})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = like.Evaluate(Input{"x": String("c")})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = like.Evaluate(Input{})
		if !reflect.DeepEqual(result, like) {
			t.Errorf("Expected %v, got %v", like, result)
		}

		if ToString(like) != "like(x, \"a\", \"_\")" {
			t.Errorf("Expected like(x, a, _), got %v", ToString(like))
		}
	})

	t.Run("evalutate Pattern _TEST_", func(t *testing.T) {
		like := Like(Ref("x"), String("_TEST_"), String(""))
		result := like.Evaluate(Input{"x": String("TEST")})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = like.Evaluate(Input{"x": String("_TEST_")})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = like.Evaluate(Input{"x": String("1TESTx")})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
	})

	t.Run("usage of regex characters in pattern", func(t *testing.T) {
		like := Like(Ref("x"), String("a.*b"), String(""))
		result := like.Evaluate(Input{"x": String("a.*b")})
		if result != Bool(true) {
			t.Errorf("Expected true, got %v", result)
		}
		result = like.Evaluate(Input{"x": String("ac")})
		if result != Bool(false) {
			t.Errorf("Expected false, got %v", result)
		}
		result = like.Evaluate(Input{"x": String("a.b")})
		if result != Bool(false) {
			t.Errorf("Expected true, got %v", result)
		}
	})
}

func TestNotLike(t *testing.T) {
	t.Run("string variable like constant String", func(t *testing.T) { //nolint:dupl
		notLike := NotLike(Ref("x"), String("a"))
		result := notLike.Evaluate(Input{"x": String("a")})
		if result != Bool(false) {
			t.Errorf("Expected true, got %v", result)
		}
		result = notLike.Evaluate(Input{"x": String("c")})
		if result != Bool(true) {
			t.Errorf("Expected false, got %v", result)
		}
		result = notLike.Evaluate(Input{})
		if !reflect.DeepEqual(result, notLike) {
			t.Errorf("Expected %v, got %v", notLike, result)
		}

		if ToString(notLike) != "not_like(x, \"a\")" {
			t.Errorf("Expected not_like(x, a), got %v", ToString(notLike))
		}
	})

	t.Run("evaluate like with _ as escape character", func(t *testing.T) { //nolint:dupl
		notLike := NotLike(Ref("x"), String("a"), String("_"))
		result := notLike.Evaluate(Input{"x": String("a")})
		if result != Bool(false) {
			t.Errorf("Expected true, got %v", result)
		}
		result = notLike.Evaluate(Input{"x": String("c")})
		if result != Bool(true) {
			t.Errorf("Expected false, got %v", result)
		}
		result = notLike.Evaluate(Input{})
		if !reflect.DeepEqual(result, notLike) {
			t.Errorf("Expected %v, got %v", notLike, result)
		}

		if ToString(notLike) != "not_like(x, \"a\", \"_\")" {
			t.Errorf("Expected not_like(x, a, _), got %v", ToString(notLike))
		}
	})
}
