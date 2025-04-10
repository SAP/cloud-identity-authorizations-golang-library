package expression

type IsRestricted struct {
	Not       Bool
	Reference string
}

func (e IsRestricted) Evaluate(input Input) Expression {
	val, ok := input[e.Reference]
	if ok && val == UNKNOWN {
		return e
	}
	return e.Not
}
