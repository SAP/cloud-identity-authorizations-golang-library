package expression

import (
	"reflect"
	"testing"
)

func TestVisitExpression(t *testing.T) {
	// TestVisitExpression tests the VisitExpression function.
	t.Run("Visit without changing", func(t *testing.T) {
		e := Or{Args: []Expression{
			And{Args: []Expression{
				Not{Arg: Reference{Name: "a"}},
				Reference{Name: "b"},
			}}, Reference{Name: "c"}}}
		result := VisitExpression(e, func(e Expression, args []Expression) Expression {
			return e
		})
		if !reflect.DeepEqual(result, e) {
			t.Errorf("Expected %v, got %v", e, result)
		}
	})
}
