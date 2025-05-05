package expression

import (
	"testing"
)

func TestSimpleImplications(t *testing.T) {
	x := Ref("x")
	// y := Ref("y")

	t.Run("x<1 => x<=1", func(t *testing.T) {
		l := Lt(x, Number(1))
		r := Le(x, Number(1))
		got := Implies(l, r)
		want := true
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		got = Implies(r, l)
		want = false
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("x>1 => x>=1", func(t *testing.T) {
		l := Gt(x, Number(1))
		r := Ge(x, Number(1))
		got := Implies(l, r)
		want := true
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		got = Implies(r, l)
		want = false
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("x>1 != x<1", func(t *testing.T) {
		l := Gt(x, Number(1))
		r := Lt(x, Number(1))
		got := Implies(l, r)
		want := false
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		got = Implies(r, l)
		want = false
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("not(1<x) => x<2", func(t *testing.T) {
		l := Not(Lt(Number(1), x))
		r := Lt(x, Number(2))
		got := Implies(l, r)
		want := true
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		got = Implies(r, l)
		want = false
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("not(1<x) => x>0", func(t *testing.T) {
		l := Not(Lt(Number(1), x))
		r := Gt(x, Number(0))
		got := Implies(l, r)
		want := true
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		got = Implies(r, l)
		want = false
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("x=1 => x<2", func(t *testing.T) {

		l := Eq(x, Number(1))
		r := Lt(x, Number(2))

		got := Implies(l, r)
		want := true
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		got = Implies(r, l)
		want = false
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("x=1 => x>0", func(t *testing.T) {
		l := Eq(x, Number(1))
		r := Gt(x, Number(0))

		// t.Errorf("l %v", minimizeOperatorSet(l))
		// t.Errorf("r %v", minimizeOperatorSet(r))

		got := Implies(l, r)
		want := true
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}

		got = Implies(r, l)
		want = false
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
	t.Run("x=1 => x=1", func(t *testing.T) {
		l := Eq(x, Number(1))
		r := Eq(x, Number(1))
		got := Implies(l, r)
		want := true
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
		got = Implies(r, l)
		want = true
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
