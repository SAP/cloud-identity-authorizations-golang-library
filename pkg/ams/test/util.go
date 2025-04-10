package test

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func NormalizeExpression(expr expression.Expression) expression.Expression {
	return expression.VisitExpression(expr, visitExpression)
}

func visitExpression(e expression.Expression, args []expression.Expression) expression.Expression {
	switch e := e.(type) {
	case expression.And:
		newArgs := []expression.Expression{}
		for _, arg := range args {
			if arg == expression.FALSE {
				return arg
			}
			if arg == expression.TRUE {
				continue
			}

			if and, ok := arg.(expression.And); ok {
				for _, andArg := range and.Args {
					alreadyExists := false
					for _, newArg := range newArgs {
						if reflect.DeepEqual(newArg, andArg) {
							alreadyExists = true
							break
						}
					}
					if !alreadyExists {
						newArgs = append(newArgs, andArg)
					}
				}
				continue
			}
			alreadyExists := false
			for _, newArg := range newArgs {
				if reflect.DeepEqual(newArg, arg) {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				newArgs = append(newArgs, arg)
			}
		}
		if len(newArgs) == 2 {
			if eq1, ok := newArgs[0].(expression.Eq); ok {
				if eq2, ok := newArgs[1].(expression.Eq); ok {
					if var1, ok := eq1.Args[0].(expression.Reference); ok {
						if var2, ok := eq2.Args[0].(expression.Reference); ok {
							if var1.Name == var2.Name {
								if eq1.Args[1] != eq2.Args[1] {
									return expression.FALSE
								}
							}
						}
					}
				}
			}
		}

		sort.Slice(newArgs, func(i, j int) bool {
			return fmt.Sprintf("%v", newArgs[i]) < fmt.Sprintf("%v", newArgs[j])
		})
		return expression.NewAnd(newArgs...)
	case expression.Or:
		newArgs := []expression.Expression{}
		for _, arg := range args {
			alreadyExists := false
			if or, ok := arg.(expression.Or); ok {
				newArgs = append(newArgs, or.Args...)
				continue
			}
			for _, newArg := range newArgs {
				if reflect.DeepEqual(newArg, arg) {
					alreadyExists = true
					break
				}
			}
			if !alreadyExists {
				newArgs = append(newArgs, arg)
			}
		}
		sort.Slice(newArgs, func(i, j int) bool {
			return fmt.Sprintf("%v", newArgs[i]) < fmt.Sprintf("%v", newArgs[j])
		})
		return expression.NewOr(newArgs...)
	case expression.Not:
		return expression.Not{Arg: e.Arg}
	case expression.In:
		array, ok := e.Args[1].(expression.ArrayConstant)
		if !ok {
			return e
		}
		newArgs := []expression.Expression{}
		for _, arg := range array.Elements() {
			newArgs = append(newArgs, expression.Eq{Args: []expression.Expression{e.Args[0], arg}})
		}
		return expression.NewOr(newArgs...)
	case expression.NotIn:
		array, ok := e.Args[1].(expression.ArrayConstant)
		if !ok {
			return e
		}
		newArgs := []expression.Expression{}
		for _, arg := range array.Elements() {
			newArgs = append(newArgs, expression.Ne{Args: []expression.Expression{e.Args[0], arg}})
		}
		return expression.NewAnd(newArgs...)
	case expression.Eq:

		l, ok := e.Args[1].(expression.Reference)
		if ok {
			r, ok := e.Args[0].(expression.Reference)
			if !ok {
				return expression.Eq{Args: []expression.Expression{e.Args[1], e.Args[0]}}
			} else if l.Name == r.Name {
				return expression.IsNotNull{Arg: e.Args[0]}
			}
		}
		return e
	case expression.Ne:
		_, ok := e.Args[1].(expression.Reference)
		if ok {
			_, ok := e.Args[0].(expression.Reference)
			if !ok {
				return expression.Ne{Args: []expression.Expression{e.Args[1], e.Args[0]}}
			}
		}
		return e
	case expression.IsRestricted:
		return e.Not

	default:
		return e
	}
}
