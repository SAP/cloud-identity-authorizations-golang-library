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

func (b Wildcard) equals(Constant) bool {
	return false
}

func (b Wildcard) lessThan(Constant) bool {
	return false
}
