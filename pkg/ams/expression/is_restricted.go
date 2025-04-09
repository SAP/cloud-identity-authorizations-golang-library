package expression

type IsRestricted struct {
	Not          Bool
	VariableName string
}

func (e IsRestricted) Evaluate(input Input) Expression {
	val, ok := input[e.VariableName]
	if ok && val == UNKNOWN {
		return e
	}
	return e.Not
}
