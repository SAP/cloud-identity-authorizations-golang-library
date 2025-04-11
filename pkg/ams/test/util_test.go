package test

import (
	"reflect"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func TestNormalizeExpression(t *testing.T) {
	t.Run(" TRUE and FALSE => FALSE", func(t *testing.T) {
		and := expression.And{
			Args: []expression.Expression{
				expression.TRUE,
				expression.FALSE,
			},
		}
		got := NormalizeExpression(and)
		want := expression.FALSE
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run(" resolves And in And and removes duplicates of eq", func(t *testing.T) {
		and := expression.And{
			Args: []expression.Expression{
				expression.Eq{
					Args: []expression.Expression{
						expression.Reference{Name: "x"},
						expression.String("a"),
					},
				},
				expression.And{
					Args: []expression.Expression{
						expression.Eq{
							Args: []expression.Expression{
								expression.Reference{Name: "x"},
								expression.String("a"),
							},
						},
						expression.Eq{
							Args: []expression.Expression{
								expression.Reference{Name: "y"},
								expression.String("b"),
							},
						},
					},
				},
			},
		}
		got := NormalizeExpression(and)
		want := expression.And{
			Args: []expression.Expression{
				expression.Eq{
					Args: []expression.Expression{
						expression.Reference{Name: "x"},
						expression.String("a"),
					},
				},
				expression.Eq{
					Args: []expression.Expression{
						expression.Reference{Name: "y"},
						expression.String("b"),
					},
				},
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run("resolves Or in Or", func(t *testing.T) {
		or := expression.Or{
			Args: []expression.Expression{
				expression.Or{
					Args: []expression.Expression{
						expression.Eq{
							Args: []expression.Expression{
								expression.Reference{Name: "x"},
								expression.String("a"),
							},
						},
						expression.Eq{
							Args: []expression.Expression{
								expression.Reference{Name: "y"},
								expression.String("b"),
							},
						},
					},
				},
				expression.Eq{
					Args: []expression.Expression{
						expression.Reference{Name: "z"},
						expression.String("c"),
					},
				},
			},
		}
		got := NormalizeExpression(or)
		want := expression.Or{
			Args: []expression.Expression{
				expression.Eq{
					Args: []expression.Expression{
						expression.Reference{Name: "x"},
						expression.String("a"),
					},
				},
				expression.Eq{
					Args: []expression.Expression{
						expression.Reference{Name: "y"},
						expression.String("b"),
					},
				},
				expression.Eq{
					Args: []expression.Expression{
						expression.Reference{Name: "z"},
						expression.String("c"),
					},
				},
			},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})
	t.Run(" In => Or equals", func(t *testing.T) {
		in := expression.In{
			Args: []expression.Expression{
				expression.Reference{Name: "x"},
				expression.StringArray{expression.String("a"), expression.String("b")},
			},
		}
		got := NormalizeExpression(in)
		want := expression.Or{Args: []expression.Expression{
			expression.Eq{Args: []expression.Expression{expression.Reference{Name: "x"}, expression.String("a")}},
			expression.Eq{Args: []expression.Expression{expression.Reference{Name: "x"}, expression.String("b")}},
		}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run(" In x => In x", func(t *testing.T) {
		in := expression.In{
			Args: []expression.Expression{
				expression.Reference{Name: "x"},
				expression.Reference{Name: "y"},
			},
		}
		got := NormalizeExpression(in)
		want := in
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run(" Not In => And not equals", func(t *testing.T) {
		notIt := expression.NotIn{
			Args: []expression.Expression{
				expression.Reference{Name: "x"},
				expression.StringArray{
					expression.String("a"),
					expression.String("b"),
				},
			},
		}
		got := NormalizeExpression(notIt)
		want := expression.And{Args: []expression.Expression{
			expression.Ne{Args: []expression.Expression{expression.Reference{Name: "x"}, expression.String("a")}},
			expression.Ne{Args: []expression.Expression{expression.Reference{Name: "x"}, expression.String("b")}},
		}}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run(" Not In x => Not In x", func(t *testing.T) {
		notIt := expression.NotIn{
			Args: []expression.Expression{
				expression.Reference{Name: "x"},
				expression.Reference{Name: "y"},
			},
		}
		got := NormalizeExpression(notIt)
		want := notIt
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run(" x=true and x=false => FALSE", func(t *testing.T) {
		and := expression.And{
			Args: []expression.Expression{
				expression.Eq{
					Args: []expression.Expression{
						expression.Reference{Name: "x"},
						expression.Bool(true),
					},
				},
				expression.Eq{
					Args: []expression.Expression{
						expression.Reference{Name: "x"},
						expression.Bool(false),
					},
				},
			},
		}
		got := NormalizeExpression(and)
		want := expression.FALSE
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run("x=x => x is not null", func(t *testing.T) {
		eq := expression.Eq{
			Args: []expression.Expression{
				expression.Reference{Name: "x"},
				expression.Reference{Name: "x"},
			},
		}
		got := NormalizeExpression(eq)
		want := expression.IsNotNull{
			Arg: expression.Reference{Name: "x"},
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})

	t.Run(" is_restricted => FALSE", func(t *testing.T) {
		restricted := expression.IsRestricted{
			Not:       expression.Bool(false),
			Reference: "x",
		}
		got := NormalizeExpression(restricted)
		want := expression.FALSE
		if !reflect.DeepEqual(got, want) {
			t.Errorf("NormalizeExpression() = %v, want %v", got, want)
		}
	})
}
