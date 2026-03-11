package expression

import (
	"fmt"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/util"
)

type Input map[string]Constant

type Reference struct {
	name string
}
type (
	// referenceSet is a map of reference names to the last index of its occurrence.
	referenceSet        map[string]int
	ExpressionContainer struct {
		Expression Expression `json:"expression"`
		References referenceSet
	}
)

// represents a logic expression. Is comparable to expression.TRUE and expression.FALSE.
//
// And can be processed using the Visit function.
type Expression interface {
	// uses the input to resolve references and evaluate to a new expression. Possibly TRUE or FALSE.
	//
	// this function may panic if the type of the input does not match the schema definition.
	Evaluate(input Input) Expression
}

func Ref(name string) Reference {
	return Reference{name: name}
}

func ToDCN(e Expression) dcn.Expression {
	switch e := e.(type) {
	case Reference:
		return dcn.Expression{Ref: util.ParseQualifiedName(e.name)}
	case Constant:
		return dcn.Expression{Constant: e}
	case OperatorCall:
		args := make([]dcn.Expression, len(e.args))
		for i, arg := range e.args {
			args[i] = ToDCN(arg)
		}
		return dcn.Expression{
			Call: []string{operatorNames[e.operator]},
			Args: args,
		}
	case FunctionCall:
		args := make([]dcn.Expression, len(e.args))
		for i, arg := range e.args {
			args[i] = ToDCN(arg)
		}
		return dcn.Expression{
			Call: util.ParseQualifiedName(e.name),
			Args: args,
		}
	default:
		panic(fmt.Sprintf("unexpected expression type %T", e))
	}
}

func FromDCN(e dcn.Expression, f *FunctionRegistry) (ExpressionContainer, error) {
	result := ExpressionContainer{
		References: make(referenceSet),
	}
	if e.Call != nil {
		if len(e.Call) == 0 {
			return result, fmt.Errorf("empty call")
		}
		args := make([]Expression, len(e.Args))
		for i, arg := range e.Args {
			container, err := FromDCN(arg, f)
			if err != nil {
				return result, err
			}
			args[i] = container.Expression
			for name := range container.References {
				result.References[name] = i
			}
		}
		if len(e.Call) == 1 {
			switch e.Call[0] {
			case AND:
				result.Expression = And(args...)
			case OR:
				result.Expression = Or(args...)
			case NOT:
				result.Expression = Not(args[0])
			case IS_NULL:
				result.Expression = IsNull(args[0])
			case IS_NOT_NULL:
				result.Expression = IsNotNull(args[0])
			case LIKE:
				result.Expression = Like(args...)
			case NOT_LIKE:
				result.Expression = NotLike(args...)
			case BETWEEN:
				result.Expression = Between(args...)
			case NOT_BETWEEN:
				result.Expression = NotBetween(args...)
			case IN:
				result.Expression = In(args...)
			case NOT_IN:
				result.Expression = NotIn(args...)
			case EQ:
				result.Expression = Eq(args...)
			case NE:
				result.Expression = Ne(args...)
			case LT:
				result.Expression = Lt(args...)
			case LE:
				result.Expression = Le(args...)
			case GT:
				result.Expression = Gt(args...)
			case GE:
				result.Expression = Ge(args...)
			case RESTRICTED:
				result.Expression = Restricted(args[0])
			case NOT_RESTRICTED:
				result.Expression = NotRestricted(args[0])
			default:
				return result, fmt.Errorf("unknown call: %s", e.Call[0])
			}
			return result, nil
		}

		if len(e.Call) > 1 {
			name := util.StringifyQualifiedName(e.Call)
			result.Expression = Function(name, f, args)
			return result, nil
		}
	}
	if e.Ref != nil {
		name := util.StringifyQualifiedName(e.Ref)
		result.Expression = Reference{name: name}
		result.References[name] = 0
	}
	if e.Constant != nil {
		result.Expression = ConstantFrom(e.Constant)
		if result.Expression == nil {
			return result, fmt.Errorf("unexpected constant %v", e.Constant)
		}
	}
	return result, nil
}

func (v Reference) Evaluate(input Input) Expression {
	val, ok := input[v.name]
	if !ok {
		return v
	}
	return val
}

func (v Reference) GetName() string {
	return v.name
}

func ToString(e Expression) string {
	return fmt.Sprintf("%v", e)
}

func IsRestrictable(e Expression) bool {
	oc, ok := e.(OperatorCall)
	if !ok {
		return false
	}

	if oc.operator == restricted || oc.operator == not_restricted {
		return true
	}

	for _, arg := range oc.args {
		if IsRestrictable(arg) {
			return true
		}
	}
	return false
}

func ApplyRestriction(e Expression, restriction []ExpressionContainer) Expression {
	oc, ok := e.(OperatorCall)
	if !ok {
		return e
	}
	if oc.operator == restricted || oc.operator == not_restricted {
		for _, r := range restriction {
			ref, ok := oc.args[0].(Reference)
			if !ok {
				continue
			}
			if _, ok := r.References[ref.GetName()]; ok {
				return r.Expression
			}
		}
	}
	if oc.operator == and || oc.operator == or || oc.operator == not {
		newArgs := make([]Expression, len(oc.args))
		for i, arg := range oc.args {
			newArgs[i] = ApplyRestriction(arg, restriction)
		}
		return OperatorCall{
			operator: oc.operator,
			args:     newArgs,
		}
	}
	return e
}

func Visit[T any](e Expression, fCall func(string, []T) T, fRef func(Reference) T, fConst func(Constant) T) T {
	switch e := e.(type) {
	case Reference:
		return fRef(e)
	case Constant:
		return fConst(e)
	case OperatorCall:
		args := make([]T, len(e.args))
		for i, arg := range e.args {
			args[i] = Visit(arg, fCall, fRef, fConst)
		}
		c, ok := operatorNames[e.operator]
		if ok {
			return fCall(c, args)
		}
	case FunctionCall:
		args := make([]T, len(e.args))
		for i, arg := range e.args {
			args[i] = Visit(arg, fCall, fRef, fConst)
		}
		return fCall(e.name, args)
	}
	return fCall("unexpected_expression", []T{})
}
