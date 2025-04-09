package expression

type IsNull struct {
	Arg Expression
}

type IsNotNull struct {
	Arg Expression
}

func (i IsNotNull) Evaluate(input Input) Expression {
	arg := i.Arg.Evaluate(input)
	if arg == UNSET {
		return Bool(false)
	}
	if arg == IGNORE {
		return arg
	}
	_, ok := arg.(Constant)
	if ok {
		return Bool(true)
	}
	return i
}

func (i IsNull) Evaluate(input Input) Expression {
	arg := i.Arg.Evaluate(input)
	if arg == UNSET {
		return Bool(true)
	}
	if arg == IGNORE {
		return arg
	}
	_, ok := arg.(Constant)
	if ok {
		return Bool(false)
	}
	return i
}
