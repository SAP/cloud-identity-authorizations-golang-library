package test

import (
	"reflect"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func TestNormalizeExpression(t *testing.T) {
	t.Run(" TRUE and FALSE => FALSE", func(t *testing.T) {
		and := expression.And(
			expression.TRUE,
			expression.FALSE,
		)
		got := NormalizeExpression(and)
		want := expression.FALSE
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run(" resolves And in And and removes duplicates of eq", func(t *testing.T) {
		and := expression.And(
			expression.Eq(
				expression.Ref("x"),
				expression.String("a"),
			),
			expression.And(
				expression.Eq(
					expression.Ref("x"),
					expression.String("a"),
				),
				expression.Eq(
					expression.Ref("y"),
					expression.String("b"),
				),
			),
		)
		got := NormalizeExpression(and)
		want := expression.And(
			expression.Eq(
				expression.Ref("x"),
				expression.String("a"),
			),
			expression.Eq(
				expression.Ref("y"),
				expression.String("b"),
			),
		)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run("resolves Or in Or", func(t *testing.T) {
		or := expression.Or(
			expression.Or(
				expression.Eq(
					expression.Ref("x"),
					expression.String("a"),
				),
				expression.Eq(
					expression.Ref("y"),
					expression.String("b"),
				),
			),
			expression.Eq(
				expression.Ref("z"),
				expression.String("c"),
			),
		)
		got := NormalizeExpression(or)
		want := expression.Or(
			expression.Eq(
				expression.Ref("x"),
				expression.String("a"),
			),
			expression.Eq(
				expression.Ref("y"),
				expression.String("b"),
			),
			expression.Eq(
				expression.Ref("z"),
				expression.String("c"),
			),
		)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})
	t.Run(" In => Or equals", func(t *testing.T) {
		in := expression.In(
			expression.Ref("x"),
			expression.StringArray{expression.String("a"), expression.String("b")},
		)
		got := NormalizeExpression(in)
		want := expression.Or(
			expression.Eq(expression.Ref("x"), expression.String("a")),
			expression.Eq(expression.Ref("x"), expression.String("b")),
		)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run(" In x => In x", func(t *testing.T) {
		in := expression.In(
			expression.Ref("x"),
			expression.Ref("y"),
		)
		got := NormalizeExpression(in)
		want := in
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run(" Not In => And not equals", func(t *testing.T) {
		notIt := expression.NotIn(
			expression.Ref("x"),
			expression.StringArray{
				expression.String("a"),
				expression.String("b"),
			},
		)
		got := NormalizeExpression(notIt)
		want := expression.And(
			expression.Ne(expression.Ref("x"), expression.String("a")),
			expression.Ne(expression.Ref("x"), expression.String("b")),
		)
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run(" Not In x => Not In x", func(t *testing.T) {
		notIt := expression.NotIn(
			expression.Ref("x"),
			expression.Ref("y"),
		)
		got := NormalizeExpression(notIt)
		want := notIt
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run(" x=true and x=false => FALSE", func(t *testing.T) {
		and := expression.And(
			expression.Eq(
				expression.Ref("x"),
				expression.Bool(true),
			),
			expression.Eq(
				expression.Ref("x"),
				expression.Bool(false),
			),
		)
		got := NormalizeExpression(and)
		want := expression.FALSE
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run("x=x => x is not null", func(t *testing.T) {
		eq := expression.Eq(
			expression.Ref("x"),
			expression.Ref("x"),
		)
		got := NormalizeExpression(eq)
		want := expression.IsNotNull(expression.Ref("x"))
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run(" is_restricted => FALSE", func(t *testing.T) {
		restricted := expression.Restricted(
			expression.Ref("x"),
		)
		got := NormalizeExpression(restricted)
		want := expression.FALSE
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})
}
