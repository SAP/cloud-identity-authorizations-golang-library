package expression

// expresses the same with just using AND/OR/NOT/LT/IN/LIKE/IS_NULL
func minimizeOperatorSet(e Expression) Expression {

	oc, ok := e.(OperatorCall)
	if !ok {
		return e
	}
	switch oc.operator {
	case AND:
		newArgs := []Expression{}
		for _, arg := range oc.args {
			newArg := minimizeOperatorSet(arg)
			if newArg == TRUE {
				continue
			}
			if newArg == FALSE {
				return FALSE
			}
			if and, ok := newArg.(OperatorCall); ok && and.operator == AND {
				newArgs = append(newArgs, and.args...)
			} else {
				newArgs = append(newArgs, newArg)
			}
		}
		return And(newArgs...)
	case OR:
		newArgs := []Expression{}
		for _, arg := range oc.args {
			newArg := minimizeOperatorSet(arg)
			if newArg == FALSE {
				continue
			}
			if newArg == TRUE {
				return TRUE
			}
			if or, ok := newArg.(OperatorCall); ok && or.operator == OR {
				newArgs = append(newArgs, or.args...)
			} else {
				newArgs = append(newArgs, newArg)
			}
		}
		return Or(newArgs...)
	case NOT:
		return Not(minimizeOperatorSet(oc.args[0]))
	case EQ:
		return And(Not(Lt(oc.args[0], oc.args[1])), Not(Lt(oc.args[1], oc.args[0])))
	case NE:
		return Or(Lt(oc.args[0], oc.args[1]), Lt(oc.args[1], oc.args[0]))
	case LT:
		return e
	case GT:
		return Lt(oc.args[1], oc.args[0])
	case LE:
		return Not(Lt(oc.args[1], oc.args[0]))
	case GE:
		return Not(Lt(oc.args[0], oc.args[1]))
	case IN:
		array, ok := oc.args[1].(ArrayConstant)
		if !ok {
			return e
		}
		if array.IsEmpty() {
			return FALSE
		}
		newArgs := []Expression{}
		for _, v := range array.Elements() {
			newArgs = append(newArgs,
				And(Not(Lt(oc.args[0], v)),
					Not(Lt(v, oc.args[0])),
				),
			)
		}
		return Or(newArgs...)
	case NOT_IN:
		array, ok := oc.args[1].(ArrayConstant)
		if !ok {
			return Not(In(oc.args...))
		}
		if array.IsEmpty() {
			return TRUE
		}
		newArgs := []Expression{}
		for _, v := range array.Elements() {
			newArgs = append(newArgs,
				Or(Lt(oc.args[0], v),
					Lt(v, oc.args[0]),
				),
			)
		}
		return And(newArgs...)
	case LIKE:
		return e
	case NOT_LIKE:
		return Not(Like(oc.args...))
	case IS_NULL:
		return e
	case IS_NOT_NULL:
		return IsNotNull(oc.args[0])
	case BETWEEN:
		return And(Not(Lt(oc.args[0], oc.args[1])), Not(Lt(oc.args[2], oc.args[0])))
	case NOT_BETWEEN:
		return Or(Lt(oc.args[0], oc.args[1]), Lt(oc.args[2], oc.args[0]))

	}
	return e
}
