package expression

import (
	"fmt"
	"regexp"
	"strings"
)

type callOperator int

const (
	and callOperator = iota
	or
	not
	eq
	ne
	lt
	le
	gt
	ge
	between
	like
	in
	is_null
	is_not_null
	not_between
	not_like
	not_in
	not_restricted
	restricted
)

const (
	AND            = "and"
	OR             = "or"
	NOT            = "not"
	EQ             = "eq"
	NE             = "ne"
	LT             = "lt"
	LE             = "le"
	GT             = "gt"
	GE             = "ge"
	BETWEEN        = "between"
	LIKE           = "like"
	IN             = "in"
	IS_NULL        = "is_null"
	IS_NOT_NULL    = "is_not_null"
	NOT_BETWEEN    = "not_between"
	NOT_LIKE       = "not_like"
	NOT_IN         = "not_in"
	NOT_RESTRICTED = "not_restricted"
	RESTRICTED     = "restricted"
)

var operatorNames = map[callOperator]string{
	and:            AND,
	or:             OR,
	not:            NOT,
	eq:             EQ,
	ne:             NE,
	lt:             LT,
	le:             LE,
	gt:             GT,
	ge:             GE,
	between:        BETWEEN,
	like:           LIKE,
	in:             IN,
	is_null:        IS_NULL,
	is_not_null:    IS_NOT_NULL,
	not_between:    NOT_BETWEEN,
	not_like:       NOT_LIKE,
	not_in:         NOT_IN,
	not_restricted: NOT_RESTRICTED,
	restricted:     RESTRICTED,
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
	case and:
		return o.evaluateAnd(input)
	case or:
		return o.evaluateOr(input)
	case not:
		return o.evaluateNot(input)
	case like:
		newArg := o.args[0].Evaluate(input)
		str, ok := newArg.(String)
		if !ok {
			return OperatorCall{
				operator: like,
				args:     o.args,
				regex:    o.regex,
			}
		}
		return Bool(o.regex.MatchString(string(str)))
	case not_like:
		newArg := o.args[0].Evaluate(input)
		str, ok := newArg.(String)
		if !ok {
			return OperatorCall{
				operator: not_like,
				args:     o.args,
				regex:    o.regex,
			}
		}
		return Bool(!o.regex.MatchString(string(str)))
	case in:
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
	case not_in:
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
	case is_null:
		newArg := o.args[0].Evaluate(input)
		if _, ok := newArg.(Constant); ok {
			return FALSE
		}
		return OperatorCall{
			operator: is_null,
			args:     o.args,
		}
	case is_not_null:
		newArg := o.args[0].Evaluate(input)
		if _, ok := newArg.(Constant); ok {
			return TRUE
		}
		return OperatorCall{
			operator: is_not_null,
			args:     o.args,
		}
	case restricted:
		return FALSE
	case not_restricted:
		return TRUE
	case eq:
		c, newArgs := evaluateArgs(input, o.args)
		if len(c) == 2 {
			return Bool(c[0].equals((c[1])))
		}
		return Eq(newArgs...)

	case ne:
		c, newArgs := evaluateArgs(input, o.args)
		if len(c) == 2 {
			return Bool(!c[0].equals((c[1])))
		}
		return Ne(newArgs...)
	case lt:
		c, newArgs := evaluateArgs(input, o.args)
		if len(c) == 2 {
			return Bool(c[0].LessThan(c[1]))
		}
		return Lt(newArgs...)
	case le:
		c, newArgs := evaluateArgs(input, o.args)
		if len(c) == 2 {
			return Bool(!c[1].LessThan(c[0]))
		}
		return Le(newArgs...)
	case gt:
		c, newArgs := evaluateArgs(input, o.args)
		if len(c) == 2 {
			return Bool(c[1].LessThan(c[0]))
		}
		return Gt(newArgs...)
	case ge:
		c, newArgs := evaluateArgs(input, o.args)
		if len(c) == 2 {
			return Bool(!c[0].LessThan(c[1]))
		}
		return Ge(newArgs...)
	case between:
		c, newArgs := evaluateArgs(input, o.args)
		if len(c) == 3 {
			return Bool(!c[0].LessThan(c[1]) && !c[2].LessThan(c[0]))
		}
		return Between(newArgs...)
	case not_between:
		c, newArgs := evaluateArgs(input, o.args)
		if len(c) == 3 {
			return Bool(c[0].LessThan(c[1]) || c[2].LessThan(c[0]))
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
		operator: or,
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
		operator: and,
		args:     args,
	}
}

func Not(arg Expression) Expression {
	if op, ok := arg.(OperatorCall); ok {
		if op.operator == not {
			return op.args[0]
		}
	}

	if arg == TRUE {
		return FALSE
	}
	if arg == FALSE {
		return TRUE
	}

	return OperatorCall{
		operator: not,
		args:     []Expression{arg},
	}
}

func In(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: in,
		args:     args,
	}
}

func NotIn(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: not_in,
		args:     args,
	}
}

func IsNull(arg Expression) OperatorCall {
	return OperatorCall{
		operator: is_null,
		args:     []Expression{arg},
	}
}

func IsNotNull(arg Expression) OperatorCall {
	return OperatorCall{
		operator: is_not_null,
		args:     []Expression{arg},
	}
}

func Restricted(arg Expression) OperatorCall {
	return OperatorCall{
		operator: restricted,
		args:     []Expression{arg},
	}
}

func NotRestricted(arg Expression) OperatorCall {
	return OperatorCall{
		operator: not_restricted,
		args:     []Expression{arg},
	}
}

func Eq(args ...Expression) Expression {
	return OperatorCall{
		operator: eq,
		args:     args,
	}
}

func Ne(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: ne,
		args:     args,
	}
}

func Lt(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: lt,
		args:     args,
	}
}

func Le(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: le,
		args:     args,
	}
}

func Gt(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: gt,
		args:     args,
	}
}

func Ge(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: ge,
		args:     args,
	}
}

func Between(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: between,
		args:     args,
	}
}

func NotBetween(args ...Expression) OperatorCall {
	return OperatorCall{
		operator: not_between,
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
		operator: like,
		args:     args,
		regex:    regex,
	}
}

func NotLike(args ...Expression) OperatorCall {
	r := Like(args...)
	r.operator = not_like
	return r
}

func evaluateArgs(input Input, args []Expression) ([]Constant, []Expression) {
	if args == nil {
		return nil, nil
	}
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
