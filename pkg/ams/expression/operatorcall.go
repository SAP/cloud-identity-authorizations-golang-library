package expression

import (
	"fmt"
	"regexp"
	"strings"
)

type callOperator int

const (
	AND callOperator = iota
	OR
	NOT
	EQ
	NE
	LT
	LE
	GT
	GE
	BETWEEN
	LIKE
	IN
	IS_NULL
	IS_NOT_NULL
	NOT_BETWEEN
	NOT_LIKE
	NOT_IN
	NOT_RESTRICTED
	RESTRICTED
)

var operatorNames = map[callOperator]string{
	AND:            "and",
	OR:             "or",
	NOT:            "not",
	EQ:             "eq",
	NE:             "ne",
	LT:             "lt",
	LE:             "le",
	GT:             "gt",
	GE:             "ge",
	BETWEEN:        "between",
	LIKE:           "like",
	IN:             "in",
	IS_NULL:        "is_null",
	IS_NOT_NULL:    "is_not_null",
	NOT_BETWEEN:    "not_between",
	NOT_LIKE:       "not_like",
	NOT_IN:         "not_in",
	NOT_RESTRICTED: "not_restricted",
	RESTRICTED:     "restricted",
}

type OperatorCall struct {
	operator callOperator
	args     []Expression
	regex    *regexp.Regexp
}

func (o OperatorCall) String() string {
	args := make([]string, len(o.args))
	for i, arg := range o.args {
		args[i] = fmt.Sprintf("%v", arg)
	}

	return o.GetOperator() + "(" + strings.Join(args, ", ") + ")"
}

func (o OperatorCall) GetOperator() string {
	if name, ok := operatorNames[o.operator]; ok {
		return name
	}
	return ""
}

func (o OperatorCall) GetArgs() []Expression {
	return o.args
}

func (o OperatorCall) Evaluate(input Input) Expression {
	switch o.operator {
	case AND:
		return o.evaluateAnd(input)
	case OR:
		return o.evaluateOr(input)
	case NOT:
		return o.evaluateNot(input)
	case LIKE:
		newArg := o.args[0].Evaluate(input)
		str, ok := newArg.(String)
		if !ok {
			return OperatorCall{
				operator: LIKE,
				args:     o.args,
				regex:    o.regex,
			}
		}
		return Bool(o.regex.MatchString(string(str)))
	case NOT_LIKE:
		newArg := o.args[0].Evaluate(input)
		str, ok := newArg.(String)
		if !ok {
			return OperatorCall{
				operator: NOT_LIKE,
				args:     o.args,
				regex:    o.regex,
			}
		}
		return Bool(!o.regex.MatchString(string(str)))
	case IN:
		left := o.args[0].Evaluate(input)
		right := o.args[1].Evaluate(input)
		r, ok := right.(ArrayConstant)
		if !ok {
			return In(left, right)
		}
		if r.IsEmpty() {
			return FALSE
		}
		l, ok := left.(Constant)
		if !ok {
			return In(left, right)
		}
		if r.Contains(l) {
			return TRUE
		}
		return FALSE
	case NOT_IN:
		left := o.args[0].Evaluate(input)
		right := o.args[1].Evaluate(input)
		r, ok := right.(ArrayConstant)
		if !ok {
			return NotIn(left, right)
		}
		if r.IsEmpty() {
			return TRUE
		}
		l, ok := left.(Constant)
		if !ok {
			return NotIn(left, right)
		}
		if r.Contains(l) {
			return FALSE
		}
		return TRUE
	case IS_NULL:
		newArg := o.args[0].Evaluate(input)
		if _, ok := newArg.(Constant); ok {
			return FALSE
		}
		return OperatorCall{
			operator: IS_NULL,
			args:     o.args,
		}
	case IS_NOT_NULL:
		newArg := o.args[0].Evaluate(input)
		if _, ok := newArg.(Constant); ok {
			return TRUE
		}
		return OperatorCall{
			operator: IS_NOT_NULL,
			args:     o.args,
		}
	case RESTRICTED:
		return FALSE
	case NOT_RESTRICTED:
		return TRUE
	case EQ:
		c, newArgs := evaluateArgs(input, o.args)
		if len(c) == 2 {
			return Bool(c[0].equals((c[1])))
		}
		return Eq(newArgs...)

	case NE:
		c, newArgs := evaluateArgs(input, o.args)
		if len(c) == 2 {
			return Bool(!c[0].equals((c[1])))
		}
		return Ne(newArgs...)
	case LT:
		c, newArgs := evaluateArgs(input, o.args)
		if len(c) == 2 {
			return Bool(c[0].lessThan(c[1]))
		}
		return Lt(newArgs...)
	case LE:
		c, newArgs := evaluateArgs(input, o.args)
		if len(c) == 2 {
			return Bool(!c[1].lessThan(c[0]))
		}
		return Le(newArgs...)
	case GT:
		c, newArgs := evaluateArgs(input, o.args)
		if len(c) == 2 {
			return Bool(c[1].lessThan(c[0]))
		}
		return Gt(newArgs...)
	case GE:
		c, newArgs := evaluateArgs(input, o.args)
		if len(c) == 2 {
			return Bool(!c[0].lessThan(c[1]))
		}
		return Ge(newArgs...)
	case BETWEEN:
		c, newArgs := evaluateArgs(input, o.args)
		if len(c) == 3 {
			return Bool(!c[0].lessThan(c[1]) && !c[2].lessThan(c[0]))
		}
		return Between(newArgs...)
	case NOT_BETWEEN:
		c, newArgs := evaluateArgs(input, o.args)
		if len(c) == 3 {
			return Bool(c[0].lessThan(c[1]) || c[2].lessThan(c[0]))
		}
		return NotBetween(newArgs...)
	}
	return OperatorCall{
		operator: o.operator,
		args:     o.args,
	}
}

func (o OperatorCall) evaluateNot(input Input) Expression {
	newArg := o.args[0].Evaluate(input)
	if newArg == TRUE {
		return FALSE
	}
	if newArg == FALSE {
		return TRUE
	}
	return Not(newArg)
}

func (o OperatorCall) evaluateOr(input Input) Expression {
	newArgs := []Expression{}
	for _, arg := range o.args {
		nextArg := arg.Evaluate(input)
		b, ok := nextArg.(Bool)
		if ok {
			if b == TRUE {
				return b
			}
			continue
		}
		newArgs = append(newArgs, nextArg)
	}
	return Or(newArgs...)
}

func (o OperatorCall) evaluateAnd(input Input) Expression {
	newArgs := []Expression{}
	for _, arg := range o.args {
		nextArg := arg.Evaluate(input)
		b, ok := nextArg.(Bool)
		if ok {
			if b == FALSE {
				return b
			}
			continue
		}
		newArgs = append(newArgs, nextArg)
	}
	return And(newArgs...)
}

func CallOperator(name string, args ...Expression) Expression {
	if name == "like" {
		return Like(args...)
	}
	if name == "not_like" {
		return NotLike(args...)
	}
	for k, v := range operatorNames {
		if v == name {
			return OperatorCall{
				operator: k,
				args:     args,
			}
		}
	}
	return FunctionCall{
		name: name,
		args: args,
	}
}

func Or(args ...Expression) Expression {
	if len(args) == 1 {
		return args[0]
	}
	if len(args) == 0 {
		return Bool(false)
	}
	return OperatorCall{
		operator: OR,
		args:     args,
	}
}

func And(args ...Expression) Expression {
	if len(args) == 1 {
		return args[0]
	}
	if len(args) == 0 {
		return Bool(true)
	}
	return OperatorCall{
		operator: AND,
		args:     args,
	}
}

func Not(arg Expression) Expression {
	if op, ok := arg.(OperatorCall); ok {
		if op.operator == NOT {
			return op.args[0]
		}
	}
	return OperatorCall{
		operator: NOT,
		args:     []Expression{arg},
	}
}

func Function(name string, container *FunctionRegistry, args []Expression) Expression {
	return FunctionCall{
		name: name,
		args: args,
		fc:   container,
	}
}

func In(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: IN,
		args:     args,
	}
}

func NotIn(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: NOT_IN,
		args:     args,
	}
}

func IsNull(arg Expression) OperatorCall {
	return OperatorCall{
		operator: IS_NULL,
		args:     []Expression{arg},
	}
}

func IsNotNull(arg Expression) OperatorCall {
	return OperatorCall{
		operator: IS_NOT_NULL,
		args:     []Expression{arg},
	}
}

func Restricted(arg Expression) OperatorCall {
	return OperatorCall{
		operator: RESTRICTED,
		args:     []Expression{arg},
	}
}

func NotRestricted(arg Expression) OperatorCall {
	return OperatorCall{
		operator: NOT_RESTRICTED,
		args:     []Expression{arg},
	}
}

func Eq(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: EQ,
		args:     args,
	}
}

func Ne(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: NE,
		args:     args,
	}
}

func Lt(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: LT,
		args:     args,
	}
}

func Le(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: LE,
		args:     args,
	}
}

func Gt(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: GT,
		args:     args,
	}
}

func Ge(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: GE,
		args:     args,
	}
}

func Between(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: BETWEEN,
		args:     args,
	}
}

func NotBetween(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: NOT_BETWEEN,
		args:     args,
	}
}

func Like(args ...Expression) OperatorCall {
	escape := String("")
	if len(args) == 3 {
		escape, _ = args[2].(String)
	}

	pattern, _ := args[1].(String)

	regex := createLikeRegex(pattern, escape)
	return OperatorCall{
		operator: LIKE,
		args:     args,
		regex:    regex,
	}
}

func NotLike(args ...Expression) OperatorCall {
	r := Like(args...)
	r.operator = NOT_LIKE
	return r
}

func evaluateArgs(input Input, args []Expression) ([]Constant, []Expression) {
	var constants []Constant
	newArgs := make([]Expression, len(args))

	for i, arg := range args {
		newArg := arg.Evaluate(input)
		c, ok := newArg.(Constant)
		if ok {
			constants = append(constants, c)
			newArgs[i] = c
			continue
		}
		newArgs[i] = newArg
	}

	return constants, newArgs
}

func createLikeRegex(pattern, escape String) *regexp.Regexp {
	const (
		placeholder1 = "\x1c"
		placeholder2 = "\x1e"
		placeholder3 = "\x1f"
	)

	p := string(pattern)
	e := string(escape)
	if e != "" {
		p = strings.ReplaceAll(p, e+e, placeholder1)
		p = strings.ReplaceAll(p, e+"_", placeholder2)
		p = strings.ReplaceAll(p, e+"%", placeholder3)
	}
	// no we need to escape the regex characters
	p = regexp.QuoteMeta(p)
	p = strings.ReplaceAll(p, "%", ".*")
	p = strings.ReplaceAll(p, "_", ".")
	if escape != "" {
		p = strings.ReplaceAll(p, placeholder1, e)
		p = strings.ReplaceAll(p, placeholder2, "_")
		p = strings.ReplaceAll(p, placeholder3, "%")
	}
	return regexp.MustCompile(p)
}
