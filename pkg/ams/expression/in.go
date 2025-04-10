package expression

type In struct {
	Args []Expression
}

type NotIn struct {
	Args []Expression
}

func (e In) Evaluate(input Input) Expression {
	left := e.Args[0].Evaluate(input)
	right := e.Args[1].Evaluate(input)
	if left == UNSET || right == UNSET {
		return UNSET
	}
	if left == IGNORE || right == IGNORE {
		return IGNORE
	}

	if _, ok := right.(Reference); ok {
		return In{Args: []Expression{left, right}}
	}
	r := right.(ArrayConstant)
	if r.IsEmpty() {
		return FALSE
	}

	if _, ok := left.(Reference); ok {
		return In{Args: []Expression{left, right}}
	}
	l := left.(Constant)

	if r.Contains(l) {
		return TRUE
	}
	return FALSE
}

func (e NotIn) Evaluate(input Input) Expression {
	left := e.Args[0].Evaluate(input)
	right := e.Args[1].Evaluate(input)
	if left == UNSET || right == UNSET {
		return UNSET
	}
	if left == IGNORE || right == IGNORE {
		return IGNORE
	}

	if _, ok := right.(Reference); ok {
		return NotIn{Args: []Expression{left, right}}
	}
	r := right.(ArrayConstant)
	if r.IsEmpty() {
		return TRUE
	}
	if _, ok := left.(Reference); ok {
		return NotIn{Args: []Expression{left, right}}
	}
	l := left.(Constant)

	if r.Contains(l) {
		return Bool(false)
	}
	return Bool(true)
}
