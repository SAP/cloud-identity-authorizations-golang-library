package expression

type wildcard byte

const (
	unset wildcard = iota
	ignore
	null
)

func (w wildcard) Evaluate(input Input) Expression {
	return nil
}

func NullifyExcept(e Expression, unknowns map[string]bool) Expression {
	result := nullifyExcept(e, unknowns, false)
	if result == null {
		return FALSE
	}
	return result
}
func nullifyExcept(e Expression, unknowns map[string]bool, inv bool) Expression {
	switch e := e.(type) {
	case Constant:
		return e
	case Reference:
		if _, ok := unknowns[e.GetName()]; ok {
			return null
		}
		return e
	case OperatorCall:
		switch e.operator { //nolint:exhaustive
		case IS_NULL:
			arg := nullifyExcept(e.args[0], unknowns, inv)
			if arg == null {
				return TRUE
			}
			return IsNull(arg)
		case IS_NOT_NULL:
			arg := nullifyExcept(e.args[0], unknowns, inv)
			if arg == null {
				return FALSE
			}
			return IsNotNull(arg)
		case NOT:
			arg := nullifyExcept(e.args[0], unknowns, inv)
			if arg == null {
				return null
			}
			if arg == TRUE {
				return FALSE
			}
			if arg == FALSE {
				return TRUE
			}
			return Not(arg)
		case AND:
			newArgs := []Expression{}
			hasNull := false
			for _, arg := range e.args {
				arg := nullifyExcept(arg, unknowns, inv)
				if arg == FALSE {
					return FALSE
				}
				if arg == TRUE {
					continue
				}
				if arg == null {
					hasNull = true
					continue
				}
				newArgs = append(newArgs, arg)
			}
			if !inv {
				if hasNull {
					return null
				}
				return And(newArgs...)
			} else {
				if len(newArgs) > 0 {
					return And(newArgs...)
				}
				if hasNull {
					return null
				}
				return TRUE
			}
		case OR:
			newArgs := []Expression{}
			hasNull := false
			for _, arg := range e.args {
				arg := nullifyExcept(arg, unknowns, inv)
				if arg == TRUE {
					return TRUE
				}
				if arg == FALSE {
					continue
				}
				if arg == null {
					hasNull = true
					continue
				}
				newArgs = append(newArgs, arg)
			}
			if !inv {
				if len(newArgs) > 0 {
					return Or(newArgs...)
				}
				if hasNull {
					return null
				}
				return FALSE
			} else {
				if hasNull {
					return null
				}
				if len(newArgs) > 0 {
					return Or(newArgs...)
				}
				return FALSE
			}
		default:
			newArgs := []Expression{}
			for _, arg := range e.args {
				arg := nullifyExcept(arg, unknowns, inv)
				if arg == null {
					return null
				}
				newArgs = append(newArgs, arg)
			}
			return OperatorCall{
				operator: e.operator,
				args:     newArgs,
				regex:    e.regex,
			}
		}
	}
	return e
}
func UnknownIgnore(e Expression, unknowns, ignores map[string]bool) Expression {
	res := unkownIgnore(e, unknowns, ignores, false)
	if res == unset {
		return FALSE
	}
	if res == ignore {
		return TRUE
	}
	return res
}

func unkownIgnore(e Expression, unknowns, ignores map[string]bool, inv bool) Expression {
	switch e := e.(type) {
	case Constant:
		return e
	case Reference:
		if _, ok := ignores[e.GetName()]; ok {
			return ignore
		}
		if _, ok := unknowns[e.GetName()]; !ok {
			return unset
		}
		return e
	case OperatorCall:
		switch e.operator { //nolint:exhaustive
		case IS_NULL:
			arg := unkownIgnore(e.args[0], unknowns, ignores, inv)
			if arg == unset {
				return TRUE
			}
			if arg == ignore {
				return ignore
			}
			return IsNull(arg)
		case IS_NOT_NULL:
			arg := unkownIgnore(e.args[0], unknowns, ignores, inv)
			if arg == unset {
				return FALSE
			}
			if arg == ignore {
				return ignore
			}
			return IsNotNull(arg)
		case NOT:
			arg := unkownIgnore(e.args[0], unknowns, ignores, !inv)
			if arg == unset {
				return unset
			}
			if arg == ignore {
				return ignore
			}
			if arg == TRUE {
				return FALSE
			}
			if arg == FALSE {
				return TRUE
			}
			return Not(arg)
		case AND:
			newArgs := []Expression{}
			hasUnset := false
			hasIgnore := false
			for _, arg := range e.args {
				arg := unkownIgnore(arg, unknowns, ignores, inv)
				if arg == FALSE {
					return FALSE
				}
				if arg == TRUE {
					continue
				}
				if arg == unset {
					hasUnset = true
					continue
				}
				if arg == ignore {
					hasIgnore = true
					continue
				}
				newArgs = append(newArgs, arg)
			}
			if !inv {
				if hasUnset {
					return unset
				}
				if len(newArgs) > 0 {
					return And(newArgs...)
				}
				if hasIgnore {
					return ignore
				}
				return TRUE
			} else {
				if hasIgnore {
					return ignore
				}
				if len(newArgs) > 0 {
					return And(newArgs...)
				}
				if hasUnset {
					return unset
				}
				return TRUE
			}
		case OR:
			newArgs := []Expression{}
			hasUnset := false
			hasIgnore := false
			for _, arg := range e.args {
				arg := unkownIgnore(arg, unknowns, ignores, inv)
				if arg == TRUE {
					return TRUE
				}
				if arg == FALSE {
					continue
				}
				if arg == unset {
					hasUnset = true
					continue
				}
				if arg == ignore {
					hasIgnore = true
					continue
				}
				newArgs = append(newArgs, arg)
			}
			if !inv {
				if hasIgnore {
					return ignore
				}
				if len(newArgs) > 0 {
					return Or(newArgs...)
				}
				if hasUnset {
					return unset
				}
				return FALSE
			} else {
				if hasUnset {
					return unset
				}
				if len(newArgs) > 0 {
					return Or(newArgs...)
				}
				if hasIgnore {
					return ignore
				}
				return FALSE
			}
		default:
			newArgs := []Expression{}
			hasIgnore := false
			for _, arg := range e.args {
				arg := unkownIgnore(arg, unknowns, ignores, inv)
				if arg == unset {
					return unset
				}
				if arg == ignore {
					hasIgnore = true
					continue
				}
				newArgs = append(newArgs, arg)
			}
			if hasIgnore {
				return ignore
			}
			return OperatorCall{
				operator: e.operator,
				args:     newArgs,
				regex:    e.regex,
			}
		}
	}
	return e
}
