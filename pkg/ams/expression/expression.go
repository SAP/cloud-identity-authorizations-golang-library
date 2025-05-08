package expression

import (
	"fmt"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/util"
)

type Input map[string]Constant

type Reference struct {
	Name string
}
type (
	referenceSet        map[string]bool
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
	return Reference{Name: name}
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
				result.References[name] = true
			}
		}
		if len(e.Call) == 1 {
			switch e.Call[0] {
			case "and":
				result.Expression = And(args...)
			case "or":
				result.Expression = Or(args...)
			case "not":
				result.Expression = Not(args[0])
			case "is_null":
				result.Expression = IsNull(args[0])
			case "is_not_null":
				result.Expression = IsNotNull(args[0])
			case "like":
				result.Expression = Like(args...)
			case "not_like":
				result.Expression = NotLike(args...)
			case "between":
				result.Expression = Between(args...)
			case "not_between":
				result.Expression = NotBetween(args...)
			case "in":
				result.Expression = In(args...)
			case "not_in":
				result.Expression = NotIn(args...)
			case "eq":
				result.Expression = Eq(args...)
			case "ne":
				result.Expression = Ne(args...)
			case "lt":
				result.Expression = Lt(args...)
			case "le":
				result.Expression = Le(args...)
			case "gt":
				result.Expression = Gt(args...)
			case "ge":
				result.Expression = Ge(args...)
			case "restricted":

				result.Expression = Restricted(args[0])
			case "not_restricted":
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
		result.Expression = Reference{Name: name}
		result.References[name] = true
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
	val, ok := input[v.Name]
	if !ok {
		return v
	}
	return val
}

func (v Reference) GetName() string {
	return v.Name
}

func ToString(e Expression) string {
	return fmt.Sprintf("%v", e)
}

func IsRestrictable(e Expression) bool {
	oc, ok := e.(OperatorCall)
	if !ok {
		return false
	}

	if oc.operator == RESTRICTED || oc.operator == NOT_RESTRICTED {
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
	if oc.operator == RESTRICTED || oc.operator == NOT_RESTRICTED {
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
	if oc.operator == AND || oc.operator == OR || oc.operator == NOT {
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

func Visit[T any](e Expression, fCall func(string, []T) T, fRef func(string) T, fConst func(Constant) T) T {
	switch e := e.(type) {
	case Reference:
		return fRef(e.Name)
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
