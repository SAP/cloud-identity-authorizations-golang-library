package expression

// import (
// 	"reflect"
// 	"testing"
// )

// func TestSimpleDNF(t *testing.T) {

// 	t.Run("a or b => a or b", func(t *testing.T) {
// 		e := Or(Ref("a"), Ref("b"))
// 		got := DNF(e)
// 		want := Or(Ref("a"), Ref("b"))
// 		if !reflect.DeepEqual(got, want) {
// 			t.Errorf("got %v, want %v", got, want)
// 		}
// 	})

// 	t.Run("a and b => a and b", func(t *testing.T) {
// 		e := And(Ref("a"), Ref("b"))
// 		got := DNF(e)
// 		want := And(Ref("a"), Ref("b"))
// 		if !reflect.DeepEqual(got, want) {
// 			t.Errorf("got %v, want %v", got, want)
// 		}
// 	})

// 	t.Run("not(a and b) => not(a) or not(b)", func(t *testing.T) {
// 		e := Not(And(Ref("a"), Ref("b")))
// 		got := DNF(e)
// 		want := Or(Not(Ref("a")), Not(Ref("b")))
// 		if !reflect.DeepEqual(got, want) {
// 			t.Errorf("got %v, want %v", got, want)
// 		}
// 	})

// 	t.Run("not(a or b) => not(a) and not(b)", func(t *testing.T) {
// 		e := Not(Or(Ref("a"), Ref("b")))
// 		got := DNF(e)
// 		want := And(Not(Ref("a")), Not(Ref("b")))
// 		if !reflect.DeepEqual(got, want) {
// 			t.Errorf("got %v, want %v", got, want)
// 		}
// 	})

// 	t.Run("(a or b) and (c or d) => (a and c) or (b and c) or (a and d) or (b and d)", func(t *testing.T) {
// 		e := And(Or(Ref("a"), Ref("b")), Or(Ref("c"), Ref("d")))
// 		got := DNF(e)
// 		want := Or(
// 			And(Ref("a"), Ref("c")),
// 			And(Ref("b"), Ref("c")),
// 			And(Ref("a"), Ref("d")),
// 			And(Ref("b"), Ref("d")),
// 		)
// 		if !reflect.DeepEqual(got, want) {
// 			t.Errorf("got %v, want %v", got, want)
// 		}
// 	})

// 	t.Run("(a or (b and c)) and d => (a and d) or (b and c and d)", func(t *testing.T) {
// 		e := And(Or(Ref("a"), And(Ref("b"), Ref("c"))), Ref("d"))
// 		got := DNF(e)
// 		want := Or(
// 			And(Ref("a"), Ref("d")),
// 			And(Ref("b"), Ref("c"), Ref("d")),
// 		)
// 		if !reflect.DeepEqual(got, want) {
// 			t.Errorf("got %v, want %v", got, want)
// 		}
// 	})

// }
