package test

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func NormalizeExpression(expr expression.Expression) expression.Expression {
	return expression.Visit(expr,
		func(s string, e []expression.Expression) expression.Expression {
			switch s {
			case "and":
				return normalizeAnd(e)
			case "or":
				return normalizeOr(e)
			case "in":
				array, ok := e[1].(expression.ArrayConstant)
				if !ok {
					return expression.In(e[0], e[1])
				}
				newArgs := []expression.Expression{}
				for _, arg := range array.Elements() {
					newArgs = append(newArgs, expression.Eq(e[0], arg))
				}
				return expression.Or(newArgs...)
			case "not_in":
				array, ok := e[1].(expression.ArrayConstant)
				if !ok {
					return expression.NotIn(e[0], e[1])
				}
				newArgs := []expression.Expression{}
				for _, arg := range array.Elements() {
					newArgs = append(newArgs, expression.Ne(e[0], arg))
				}
				return expression.And(newArgs...)
			case "eq":
				l, ok := e[1].(expression.Reference)
				if ok {
					r, ok := e[0].(expression.Reference)
					if !ok {
						return expression.Eq(e[1], e[0])
					} else if l.GetName() == r.GetName() {
						return expression.IsNotNull(l)
					}
				}
			case "ne":
				_, ok := e[1].(expression.Reference)
				if ok {
					_, ok := e[0].(expression.Reference)
					if !ok {
						return expression.Ne(e[1], e[0])
					}
				}
			case "restricted":
				return expression.FALSE
			}
			return expression.CallOperator(s, e...)
		},
		func(ref expression.Reference) expression.Expression {
			return expression.Ref(ref.GetName())
		},
		func(c expression.Constant) expression.Expression {
			return c
		},
	)
}

func normalizeAnd(args []expression.Expression) expression.Expression {
	newArgs := []expression.Expression{}
	for _, arg := range args {
		if arg == expression.FALSE {
			return arg
		}
		if arg == expression.TRUE {
			continue
		}

		if and, ok := castOp(arg, "and"); ok {
			for _, andArg := range and.GetArgs() {
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
		if eq1, ok := castOp(newArgs[0], "eq"); ok {
			if eq2, ok := castOp(newArgs[1], "eq"); ok {
				if var1, ok := eq1.GetArgs()[0].(expression.Reference); ok {
					if var2, ok := eq2.GetArgs()[0].(expression.Reference); ok {
						if var1.GetName() == var2.GetName() {
							return expression.FALSE
						}
					}
				}
			}
		}
	}

	sort.Slice(newArgs, func(i, j int) bool {
		return fmt.Sprintf("%v", newArgs[i]) < fmt.Sprintf("%v", newArgs[j])
	})
	return expression.And(newArgs...)
}

func normalizeOr(args []expression.Expression) expression.Expression {
	newArgs := []expression.Expression{}
	for _, arg := range args {
		alreadyExists := false
		if or, ok := castOp(arg, "or"); ok {
			newArgs = append(newArgs, or.GetArgs()...)
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
	return expression.Or(newArgs...)
}

func castOp(e expression.Expression, operator string) (expression.OperatorCall, bool) {
	oc, ok := e.(expression.OperatorCall)
	if !ok {
		return oc, false
	}
	if oc.GetOperator() == operator {
		return oc, true
	}
	return oc, false
}

// func visitExpression(e expression.Expression, args []expression.Expression) expression.Expression {
// 	switch e := e.(type) {

// 	case expression.In:
// 		array, ok := e.Args[1].(expression.ArrayConstant)
// 		if !ok {
// 			return e
// 		}
// 		newArgs := []expression.Expression{}
// 		for _, arg := range array.Elements() {
// 			newArgs = append(newArgs, expression.Eq(e.Args[0], arg))
// 		}
// 		return expression.NewOr(newArgs...)
// 	case expression.NotIn:
// 		array, ok := e.Args[1].(expression.ArrayConstant)
// 		if !ok {
// 			return e
// 		}
// 		newArgs := []expression.Expression{}
// 		for _, arg := range array.Elements() {
// 			newArgs = append(newArgs, expression.Ne(e.Args[0], arg))
// 		}
// 		return expression.NewAnd(newArgs...)
// 	case expression.Eq:

// 		l, ok := e.Args[1].(expression.Reference)
// 		if ok {
// 			r, ok := e.Args[0].(expression.Reference)
// 			if !ok {
// 				return expression.Eq(e.Args[1], e.Args[0])
// 			} else if l.Name == r.Name {
// 				return expression.IsNotNull{Arg: e.Args[0]}
// 			}
// 		}
// 		return e
// 	case expression.Ne:
// 		_, ok := e.Args[1].(expression.Reference)
// 		if ok {
// 			_, ok := e.Args[0].(expression.Reference)
// 			if !ok {
// 				return expression.Ne(e.Args[1], e.Args[0])
// 			}
// 		}
// 		return e

// }
