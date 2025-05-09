package expression

// // EXPERIMENTAL returns the disjunctive normal form of the expression
// func DNF(expr Expression) Expression {
// 	return dnf(expr, false)

// }

// func dnf(expr Expression, inv bool) Expression {
// 	oc, ok := expr.(OperatorCall)
// 	if ok {
// 		if oc.operator == NOT {
// 			return dnf(oc.args[0], !inv)
// 		}
// 		if oc.operator == OR && inv || oc.operator == AND && !inv {
// 			return expand(oc.args, inv)
// 		}
// 		if oc.operator == OR && !inv || oc.operator == AND && inv {
// 			newArgs := []Expression{}
// 			for _, arg := range oc.args {
// 				newArg := dnf(arg, inv)
// 				if arg, ok := newArg.(OperatorCall); ok && arg.operator == OR {
// 					newArgs = append(newArgs, arg.args...)
// 				} else {
// 					newArgs = append(newArgs, newArg)
// 				}
// 			}
// 			return Or(newArgs...)
// 		}
// 	}
// 	if inv {
// 		return Not(expr)
// 	}

// 	return expr

// }

// func expand(args []Expression, inv bool) Expression {
// 	result := [][]Expression{{}}
// 	for _, arg := range args {
// 		r := dnf(arg, inv)
// 		oc, ok := r.(OperatorCall)
// 		if !ok {
// 			for i := 0; i < len(result); i++ {
// 				result[i] = append(result[i], r)
// 			}
// 			continue
// 		}
// 		if oc.operator == OR {
// 			newResult := [][]Expression{}
// 			for _, arg := range oc.args {
// 				for _, r := range result {
// 					if andArg, ok := arg.(OperatorCall); ok && andArg.operator == AND {
// 						newResult = append(newResult, append(r, andArg.args...))
// 					} else {
// 						newResult = append(newResult, append(r, arg))
// 					}
// 				}
// 			}
// 			result = newResult
// 		} else if oc.operator == AND {
// 			for i := 0; i < len(result); i++ {
// 				result[i] = append(result[i], oc.args...)
// 			}
// 		} else {
// 			for i := 0; i < len(result); i++ {
// 				result[i] = append(result[i], oc)
// 			}
// 		}
// 	}
// 	ands := []Expression{}
// 	for _, r := range result {
// 		ands = append(ands, And(r...))
// 	}
// 	return Or(ands...)
// }
