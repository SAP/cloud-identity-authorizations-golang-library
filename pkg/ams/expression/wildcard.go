package expression

type Wildcard byte

const (
	UNKNOWN Wildcard = iota
	IGNORE
	UNSET
)

func (b Wildcard) Evaluate(input Input) Expression {
	return b
}

func (b Wildcard) Equals(Constant) bool {
	return false
}

func (b Wildcard) LessThan(Constant) bool {
	return false
}
