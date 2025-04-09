package expression

type Eq struct {
	Args []Expression
}

type Ne struct {
	Args []Expression
}

type Lt struct {
	Args []Expression
}

type Le struct {
	Args []Expression
}

type Gt struct {
	Args []Expression
}

type Ge struct {
	Args []Expression
}

type Between struct {
	Args []Expression
}

type NotBetween struct {
	Args []Expression
}

func (e Eq) Evaluate(input Input) Expression {
	constants, nextArgs, bp := evaluateArgs(input, e.Args)
	if bp != UNKNOWN {
		return bp
	}
	if len(constants) == 2 {
		return Bool(constants[0].Equals(constants[1]))
	}
	return Eq{Args: nextArgs}
}

func (e Ne) Evaluate(input Input) Expression {
	constants, nextArgs, bp := evaluateArgs(input, e.Args)
	if bp != UNKNOWN {
		return bp
	}
	if len(constants) == 2 {
		return Bool(!constants[0].Equals(constants[1]))
	}
	return Ne{Args: nextArgs}
}

func (e Lt) Evaluate(input Input) Expression {
	constants, nextArgs, bp := evaluateArgs(input, e.Args)
	if bp != UNKNOWN {
		return bp
	}
	if len(constants) == 2 {
		return Bool(constants[0].LessThan(constants[1]))
	}
	return Lt{Args: nextArgs}
}

func (e Le) Evaluate(input Input) Expression {
	constants, nextArgs, bp := evaluateArgs(input, e.Args)
	if bp != UNKNOWN {
		return bp
	}
	if len(constants) == 2 {
		return Bool(!constants[1].LessThan(constants[0]))
	}
	return Le{Args: nextArgs}
}

func (e Gt) Evaluate(input Input) Expression {
	constants, nextArgs, bp := evaluateArgs(input, e.Args)
	if bp != UNKNOWN {
		return bp
	}
	if len(constants) == 2 {
		return Bool(constants[1].LessThan(constants[0]))
	}
	return Gt{Args: nextArgs}
}

func (e Ge) Evaluate(input Input) Expression {
	constants, nextArgs, bp := evaluateArgs(input, e.Args)
	if bp != UNKNOWN {
		return bp
	}
	if len(constants) == 2 {
		return Bool(!constants[0].LessThan(constants[1]))
	}
	return Ge{Args: nextArgs}
}

func (e Between) Evaluate(input Input) Expression {
	constants, nextArgs, bp := evaluateArgs(input, e.Args)
	if bp != UNKNOWN {
		return bp
	}
	if len(constants) == 3 {
		return Bool(!constants[0].LessThan(constants[1]) && !constants[2].LessThan(constants[0]))
	}
	return Between{Args: nextArgs}
}

func (e NotBetween) Evaluate(input Input) Expression {
	constants, nextArgs, bp := evaluateArgs(input, e.Args)
	if bp != UNKNOWN {
		return bp
	}
	if len(constants) == 3 {
		return Bool(constants[0].LessThan(constants[1]) || constants[2].LessThan(constants[0]))
	}
	return NotBetween{Args: nextArgs}
}

func evaluateArgs(input Input, args []Expression) ([]Constant, []Expression, Wildcard) {
	var constants []Constant
	nextArgs := make([]Expression, len(args))
	ignore := false
	for i, arg := range args {
		nextArg := arg.Evaluate(input)
		switch nextArg := nextArg.(type) {
		case Wildcard:
			if nextArg == UNSET {
				return nil, nil, nextArg
			}
			if nextArg == IGNORE {
				ignore = true
				continue
			}
		case Constant:
			constants = append(constants, nextArg)
			nextArgs[i] = nextArg
		default:
			nextArgs[i] = nextArg
		}
	}
	if ignore {
		return nil, nil, IGNORE
	}
	return constants, nextArgs, UNKNOWN
}
