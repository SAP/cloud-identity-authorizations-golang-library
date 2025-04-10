package test

import (
	"reflect"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func TestNormalizeExpression(t *testing.T) {

	t.Run(" In => Or equals", func(t *testing.T) {
		in := expression.In{Args: []expression.Expression{expression.Reference{Name: "x"}, expression.StringArray{expression.String("a"), expression.String("b")}}}
		got := NormalizeExpression(in)
		want := expression.Or{Args: []expression.Expression{
			expression.Eq{Args: []expression.Expression{expression.Reference{Name: "x"}, expression.String("a")}},
			expression.Eq{Args: []expression.Expression{expression.Reference{Name: "x"}, expression.String("b")}},
		}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run(" Not In => And not equals", func(t *testing.T) {
		notIt := expression.NotIn{Args: []expression.Expression{expression.Reference{Name: "x"}, expression.StringArray{expression.String("a"), expression.String("b")}}}
		got := NormalizeExpression(notIt)
		want := expression.And{Args: []expression.Expression{
			expression.Ne{Args: []expression.Expression{expression.Reference{Name: "x"}, expression.String("a")}},
			expression.Ne{Args: []expression.Expression{expression.Reference{Name: "x"}, expression.String("b")}},
		}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})
}
