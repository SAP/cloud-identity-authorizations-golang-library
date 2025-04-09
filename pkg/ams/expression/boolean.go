package expression

type Not struct {
	Arg Expression
}

type And struct {
	Args []Expression
}

type Or struct {
	Args []Expression
}

func NewOr(args ...Expression) Expression {
	if len(args) == 1 {
		return args[0]
	}
	if len(args) == 0 {
		return Bool(false)
	}
	return Or{Args: args}
}

func (o Or) Evaluate(input Input) Expression {
	var resultArgs []Expression
	isIgnore := false
	hasUnset := false
	hasUnknown := false
	for _, arg := range o.Args {
		nextArg := arg.Evaluate(input)
		b, ok := nextArg.(Bool)
		if ok {
			if b == TRUE {
				return b
			}
			continue
		}
		if bp, ok := nextArg.(Wildcard); ok {
			if bp == IGNORE {
				isIgnore = true
			}
			if bp == UNSET {
				hasUnset = true
			}
			continue
		}
		hasUnknown = true
		resultArgs = append(resultArgs, nextArg)
	}
	if isIgnore {
		return IGNORE
	}
	if !hasUnknown && hasUnset {
		return UNSET
	}
	return NewOr(resultArgs...)
}

func NewAnd(args ...Expression) Expression {
	if len(args) == 1 {
		return args[0]
	}
	if len(args) == 0 {
		return Bool(true)
	}
	return And{Args: args}
}

func (a And) Evaluate(input Input) Expression {
	var resultArgs []Expression
	isUnset := false
	hasIgnore := false
	hasUnknown := false
	for _, arg := range a.Args {
		nextArg := arg.Evaluate(input)
		b, ok := nextArg.(Bool)
		if ok {
			if b == FALSE {
				return b
			}
			continue
		}
		bp, ok := nextArg.(Wildcard)
		if ok {
			if bp == UNSET {
				isUnset = true
			}
			if bp == IGNORE {
				hasIgnore = true
			}
			continue
		}
		hasUnknown = true
		resultArgs = append(resultArgs, nextArg)
	}
	if isUnset {
		return UNSET
	}

	if !hasUnknown && hasIgnore {
		return IGNORE
	}

	return NewAnd(resultArgs...)
}

func (n Not) Evaluate(input Input) Expression {
	r := n.Arg.Evaluate(input)
	b, ok := r.(Bool)
	if ok {
		return !b
	}
	v, ok := r.(Wildcard)
	if ok {
		if v == IGNORE || v == UNSET {
			return v
		}
	}
	return Not{Arg: r}
}
