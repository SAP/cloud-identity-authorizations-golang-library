package expression

// // expresses the same with just using AND/OR/NOT/LT/IN/LIKE/IS_NULL
// // comparisons to bool or bool refs are maped to Eq(x,true)
// func minimizeOperatorSet(e Expression) Expression {

// 	oc, ok := e.(OperatorCall)
// 	if !ok {
// 		return e
// 	}
// 	switch oc.operator {
// 	case AND:
// 		newArgs := []Expression{}
// 		for _, arg := range oc.args {
// 			newArg := minimizeOperatorSet(arg)
// 			if newArg == TRUE {
// 				continue
// 			}
// 			if newArg == FALSE {
// 				return FALSE
// 			}
// 			if and, ok := newArg.(OperatorCall); ok && and.operator == AND {
// 				newArgs = append(newArgs, and.args...)
// 			} else {
// 				newArgs = append(newArgs, newArg)
// 			}
// 		}
// 		return And(newArgs...)
// 	case OR:
// 		newArgs := []Expression{}
// 		for _, arg := range oc.args {
// 			newArg := minimizeOperatorSet(arg)
// 			if newArg == FALSE {
// 				continue
// 			}
// 			if newArg == TRUE {
// 				return TRUE
// 			}
// 			if or, ok := newArg.(OperatorCall); ok && or.operator == OR {
// 				newArgs = append(newArgs, or.args...)
// 			} else {
// 				newArgs = append(newArgs, newArg)
// 			}
// 		}
// 		return Or(newArgs...)
// 	case NOT:
// 		return Not(minimizeOperatorSet(oc.args[0]))
// 	case EQ:
// 		if oc.args[1] == TRUE {
// 			return Eq(oc.args[0], TRUE)
// 		}
// 		if oc.args[0] == TRUE {
// 			return Eq(oc.args[1], TRUE)
// 		}
// 		if oc.args[1] == FALSE {
// 			return Not(Eq(oc.args[0], TRUE))
// 		}
// 		if oc.args[0] == FALSE {
// 			return Not(Eq(oc.args[1], TRUE))
// 		}
// 		return minimizeOperatorSet(And(Not(lt(oc.args[0], oc.args[1])), Not(lt(oc.args[1], oc.args[0]))))
// 	case NE:
// 		return minimizeOperatorSet(Or(lt(oc.args[0], oc.args[1]), lt(oc.args[1], oc.args[0])))
// 	case LT:
// 		return lt(oc.args[0], oc.args[1])
// 	case GT:
// 		return lt(oc.args[1], oc.args[0])
// 	case LE:
// 		return Not(lt(oc.args[1], oc.args[0]))
// 	case GE:
// 		return Not(lt(oc.args[0], oc.args[1]))
// 	case IN:
// 		array, ok := oc.args[1].(ArrayConstant)
// 		if !ok {
// 			return e
// 		}
// 		if array.IsEmpty() {
// 			return FALSE
// 		}
// 		newArgs := []Expression{}
// 		for _, v := range array.Elements() {
// 			newArgs = append(newArgs,
// 				And(Not(lt(oc.args[0], v)),
// 					Not(lt(v, oc.args[0])),
// 				),
// 			)
// 		}
// 		return minimizeOperatorSet(Or(newArgs...))
// 	case NOT_IN:
// 		array, ok := oc.args[1].(ArrayConstant)
// 		if !ok {
// 			return Not(In(oc.args...))
// 		}
// 		if array.IsEmpty() {
// 			return TRUE
// 		}
// 		newArgs := []Expression{}
// 		for _, v := range array.Elements() {
// 			newArgs = append(newArgs,
// 				Or(lt(oc.args[0], v),
// 					lt(v, oc.args[0]),
// 				),
// 			)
// 		}
// 		return minimizeOperatorSet(And(newArgs...))
// 	case LIKE:
// 		return e
// 	case NOT_LIKE:
// 		return Not(Like(oc.args...))
// 	case IS_NULL:
// 		return e
// 	case IS_NOT_NULL:
// 		return Not(IsNull(oc.args[0]))
// 	case BETWEEN:
// 		return minimizeOperatorSet(And(Not(lt(oc.args[0], oc.args[1])), Not(lt(oc.args[2], oc.args[0]))))
// 	case NOT_BETWEEN:
// 		return minimizeOperatorSet(Or(lt(oc.args[0], oc.args[1]), lt(oc.args[2], oc.args[0])))

// 	}
// 	return e
// }

// func lt(a, b Expression) Expression {
// 	if a == TRUE {
// 		return FALSE
// 	}
// 	if b == TRUE {
// 		return Not(Eq(a, TRUE))
// 	}
// 	if a == FALSE {
// 		return Eq(b, TRUE)
// 	}
// 	if b == FALSE {
// 		return FALSE
// 	}

// 	return Lt(a, b)
// }
