package expression

import "fmt"

func Implies(a, b Expression) bool {
	a_m := minimizeOperatorSet(a)
	b_m := minimizeOperatorSet(b)
	return implies(a_m, b_m)
}

func implies(a, b Expression) bool {
	b_oc, ok := b.(OperatorCall)
	if !ok {
		if b == TRUE {
			return true
		}
		if b == FALSE {
			// todo it is true if a is always false
			return false
		}
		panic(fmt.Sprintf("not a boolean expression: %v", b))
	}
	switch b_oc.operator {
	case AND:
		return every(b_oc.args, func(arg Expression) bool {
			return implies(a, arg)
		})
	case OR:
		return exists(b_oc.args, func(arg Expression) bool {
			return implies(a, arg)
		})
	case NOT:
		return impliesNot(a, b_oc.args[0])
	}

	a_oc, ok := a.(OperatorCall)
	if !ok {
		if a == FALSE {
			return true
		}
		if a == TRUE {
			// todo it is true if b is always true
			return false
		}
		panic(fmt.Sprintf("not a boolean expression: %v", a))
	}
	switch a_oc.operator {
	case AND:
		return exists(a_oc.args, func(arg Expression) bool {
			return implies(arg, b)
		})
	case OR:
		return every(a_oc.args, func(arg Expression) bool {
			return implies(arg, b)
		})
	case NOT:
		return notImplies(a_oc.args[0], b)
	}

	switch a_oc.operator {
	case LT:
		switch b_oc.operator {
		case LT:
			a_ref, a_const := resolveArgs(a)
			b_ref, b_const := resolveArgs(b)
			if len(a_ref) == 0 || len(b_ref) == 0 {
				return implies(a.Evaluate(nil), b.Evaluate(nil))
			}
			if len(a_ref) != len(b_ref) {
				return false
			}
			if len(a_ref) > 1 {
				return equals(a, b)
			}
			for aName, aIndex := range a_ref {
				bIndex, ok := b_ref[aName]
				if !ok {
					return false
				}
				if aIndex != bIndex {
					return false
				}
				return !b_const.lessThan(a_const)
			}
		}
	case IS_NULL:
		switch b_oc.operator {
		case IS_NULL:
			a_ref, _ := resolveArgs(a)
			b_ref, _ := resolveArgs(b)
			for aName := range a_ref {
				_, ok := b_ref[aName]
				if !ok {
					return true
				}
			}
		}
	case IN:
		switch b_oc.operator {
		case IN:
			return equals(a, b)
		}
	case LIKE:
		switch b_oc.operator {
		case LIKE:
			return equals(a, b)
		}
	}

	return false
}

func impliesNot(a, b Expression) bool {
	a_oc, ok := a.(OperatorCall)
	if !ok {
		if a == FALSE {
			return true
		}
		if a == TRUE {
			// todo it is true if b is always true
			return false
		}
		panic(fmt.Sprintf("not a boolean expression: %v", a))
	}
	if a_oc.operator == NOT {
		return implies(b, a_oc.args[0])
	}
	b_oc, ok := b.(OperatorCall)
	if !ok {
		if b == TRUE {
			return true
		}
		if b == FALSE {
			// todo it is true if a is always false
			return false
		}
		panic(fmt.Sprintf("not a boolean expression: %v", b))
	}
	switch b_oc.operator {
	case AND:
		return exists(b_oc.args, func(arg Expression) bool {
			return impliesNot(a, arg)
		})
	case OR:
		return every(b_oc.args, func(arg Expression) bool {
			return impliesNot(a, arg)
		})
	case NOT:
		return implies(a, b_oc.args[0])
	}

	switch a_oc.operator {
	case AND:
		return exists(a_oc.args, func(arg Expression) bool {
			return impliesNot(arg, b)
		})
	case OR:
		return every(a_oc.args, func(arg Expression) bool {
			return impliesNot(arg, b)
		})
	}

	switch a_oc.operator {
	case LT:
		switch b_oc.operator {
		case LT:
			a_ref, a_const := resolveArgs(a)
			b_ref, b_const := resolveArgs(b)
			if len(a_ref) == 0 || len(b_ref) == 0 {
				return impliesNot(a.Evaluate(nil), b.Evaluate(nil))
			}
			if len(a_ref) != len(b_ref) {
				return false
			}
			if len(a_ref) > 1 {
				return equals(a, Lt(b_oc.args[0], b_oc.args[1]))
			}
			for aName, aIndex := range a_ref {
				bIndex, ok := b_ref[aName]
				if !ok {
					return false
				}
				if aIndex == bIndex {
					return false
				}
				return !b_const.lessThan(a_const)
			}
		case IS_NULL:
			a_ref, _ := resolveArgs(a)
			b_ref, _ := resolveArgs(b)
			for aName := range a_ref {
				_, ok := b_ref[aName]
				if ok {
					// a < 1 => a is not null
					return true
				}
			}
		}
	}
	return false
}

func notImplies(a, b Expression) bool {
	a_oc, ok := a.(OperatorCall)
	if !ok {
		if a == TRUE {
			return true
		}
		if a == FALSE {
			// todo it is true if b is always true
			return false
		}
		panic(fmt.Sprintf("not a boolean expression: %v", a))
	}
	b_oc, ok := b.(OperatorCall)
	if !ok {
		if b == FALSE {
			return true
		}
		if b == TRUE {
			// todo it is true if a is always false
			return false
		}
		panic(fmt.Sprintf("not a boolean expression: %v", b))
	}
	switch b_oc.operator {
	case AND:
		return every(b_oc.args, func(arg Expression) bool {
			return notImplies(a, arg)
		})
	case OR:
		return exists(b_oc.args, func(arg Expression) bool {
			return notImplies(a, arg)
		})
	case NOT:
		return implies(b_oc.args[0], a)
	}
	switch a_oc.operator {
	case AND:
		return exists(a_oc.args, func(arg Expression) bool {
			return notImplies(arg, b)
		})
	case OR:
		return every(a_oc.args, func(arg Expression) bool {
			return notImplies(arg, b)
		})
	case NOT:
		return implies(a_oc.args[0], b)
	}
	switch a_oc.operator {
	case LT:
		switch b_oc.operator {
		case LT:
			a_ref, a_const := resolveArgs(a)
			b_ref, b_const := resolveArgs(b)
			if len(a_ref) == 0 || len(b_ref) == 0 {
				return notImplies(a.Evaluate(nil), b.Evaluate(nil))
			}
			if len(a_ref) != len(b_ref) {
				return false
			}
			if len(a_ref) > 1 {
				return false
			}
			for aName, aIndex := range a_ref {
				bIndex, ok := b_ref[aName]
				if !ok {
					return false
				}
				if aIndex == bIndex {
					return b_const.lessThan(a_const)
				}
				return a_const.lessThan(b_const)
			}
		}
	}

	return false
}

func equals(a, b Expression) bool {
	switch b := b.(type) {
	case OperatorCall:
		a, ok := a.(OperatorCall)
		if !ok {
			return false
		}
		if len(a.args) != len(b.args) {
			return false
		}
		for i, arg := range a.args {
			if !equals(arg, b.args[i]) {
				return false
			}
		}
		return true
	case Reference:
		if a, ok := a.(Reference); ok {
			return a.Name == b.Name
		}
		return false
	case Bool:
		if a, ok := a.(Bool); ok {
			return a == b
		}
		return false
	case String:
		if a, ok := a.(String); ok {
			return a == b
		}
		return false
	case Number:
		if a, ok := a.(Number); ok {
			return a == b
		}
		return false
	}
	return false
}

func resolveArgs(a Expression) (map[string]int, Constant) {
	refs := make(map[string]int)
	var consts Constant
	lt := a.(OperatorCall)

	if c, ok := lt.args[0].(Constant); ok {
		consts = c
	}

	if ref, ok := lt.args[0].(Reference); ok {
		refs[ref.Name] = 0
	}
	if c, ok := lt.args[1].(Constant); ok {
		consts = c
	}

	if ref, ok := lt.args[1].(Reference); ok {
		refs[ref.Name] = 1
	}
	return refs, consts
}

func map_[from Expression, to Expression](e []from, f func(from) to) []to {
	if e == nil {
		return nil
	}
	out := make([]to, len(e))
	for i, v := range e {
		out[i] = f(v)
	}
	return out
}

func every(e []Expression, f func(Expression) bool) bool {
	if e == nil {
		return false
	}
	for _, v := range e {
		if !f(v) {
			return false
		}
	}
	return true
}

func exists(e []Expression, f func(Expression) bool) bool {
	if e == nil {
		return false
	}
	for _, v := range e {
		if f(v) {
			return true
		}
	}
	return false
}
