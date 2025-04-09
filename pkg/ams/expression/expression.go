package expression

import (
	"fmt"
	"strings"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/util"
)

type Input map[string]Constant

type Variable struct {
	Name string
}
type variableSet map[string]bool
type ExpressionContainer struct {
	Expression    Expression `json:"expression"`
	VariableNames variableSet
}

type Expression interface {
	Evaluate(Input) Expression
}

func FromDCN(e dcn.Expression, f Functions) (ExpressionContainer, error) {
	result := ExpressionContainer{
		VariableNames: make(variableSet),
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
			for name := range container.VariableNames {
				result.VariableNames[name] = true
			}
		}
		if len(e.Call) == 1 {
			switch e.Call[0] {
			case "and":
				result.Expression = And{Args: args}
			case "or":
				result.Expression = Or{Args: args}
			case "not":
				result.Expression = Not{Arg: args[0]}
			case "is_null":
				result.Expression = IsNull{Arg: args[0]}
			case "is_not_null":
				result.Expression = IsNotNull{Arg: args[0]}
			case "like":
				pattern := args[1].(String)
				var escape String
				if len(args) == 3 {
					escape = args[2].(String)
				}
				result.Expression = NewLike(args[0], pattern, escape)
			case "not_like":
				pattern := args[1].(String)
				var escape String
				if len(args) == 3 {
					escape = args[2].(String)
				}
				result.Expression = NewNotLike(args[0], pattern, escape)
			case "between":
				result.Expression = Between{Args: args}
			case "not_between":
				result.Expression = NotBetween{Args: args}
			case "in":
				result.Expression = In{Args: args}
			case "not_in":
				result.Expression = NotIn{Args: args}
			case "eq":
				result.Expression = Eq{Args: args}
			case "ne":
				result.Expression = Ne{Args: args}
			case "lt":
				result.Expression = Lt{Args: args}
			case "le":
				result.Expression = Le{Args: args}
			case "gt":
				result.Expression = Gt{Args: args}
			case "ge":
				result.Expression = Ge{Args: args}
			case "restricted":
				variable := args[0].(Variable)
				result.Expression = IsRestricted{
					Not:          Bool(false),
					VariableName: variable.Name,
				}
			case "not_restricted":
				variable := args[0].(Variable)
				result.Expression = IsRestricted{
					Not:          Bool(true),
					VariableName: variable.Name,
				}
			default:
				return result, fmt.Errorf("unknown call: %s", e.Call[0])
			}
			return result, nil
		}

		if len(e.Call) > 1 {
			name := util.StringifyReference(e.Call)
			function, ok := f[name]
			if !ok {
				return result, fmt.Errorf("unknown function %s", name)
			}
			result.Expression = function
		}

	}

	if e.Ref != nil {
		name := util.StringifyReference(e.Ref)
		result.Expression = Variable{Name: name}
		result.VariableNames[name] = true
	}

	if e.Constant != nil {
		result.Expression = ConstantFrom(e.Constant)
		if result.Expression == UNSET {
			return result, fmt.Errorf("unexpected constant %v", e.Constant)
		}
	}
	return result, nil
}

func (v Variable) Evaluate(input Input) Expression {
	val, ok := input[v.Name]
	if !ok {
		return UNSET
	}
	if val == UNKNOWN {
		return v
	}

	return val
}
func ToString(e Expression) string {
	return Visit(e,
		func(name string, args []string) string {
			return name + "(" + strings.Join(args, ", ") + ")"
		},
		func(v Variable) string {
			return v.Name
		},
		func(c Constant) string {
			switch c := c.(type) {
			case String:
				return fmt.Sprintf("\"%v\"", c)
			case ArrayConstant:
				return fmt.Sprintf("%v", c)
			default:
				return fmt.Sprintf("%v", c)
			}
		},
	)
}

func IsRestrictable(e Expression) bool {

	switch e := e.(type) {
	case And:
		for _, arg := range e.Args {
			if IsRestrictable(arg) {
				return true
			}
		}
		return false
	case Or:
		for _, arg := range e.Args {
			if IsRestrictable(arg) {
				return true
			}
		}
		return false
	case Not:
		return IsRestrictable(e.Arg)
	case IsRestricted:
		return true
	default:
		return false
	}
}

func ApplyRestriction(e Expression, restriction []ExpressionContainer) Expression {
	switch e := e.(type) {
	case And:
		args := make([]Expression, len(e.Args))
		for i, arg := range e.Args {
			args[i] = ApplyRestriction(arg, restriction)
		}
		return And{args}
	case Or:
		args := make([]Expression, len(e.Args))
		for i, arg := range e.Args {
			args[i] = ApplyRestriction(arg, restriction)
		}
		return Or{args}
	case Not:
		return Not{ApplyRestriction(e.Arg, restriction)}
	case IsRestricted:
		for _, r := range restriction {
			if _, ok := r.VariableNames[e.VariableName]; ok {
				return r.Expression
			}
		}
		return e
	default:
		return e
	}
}

func VisitExpression[T any](e Expression, f func(Expression, []T) T) T {
	switch e := e.(type) {
	case And:
		args := make([]T, len(e.Args))
		for i, arg := range e.Args {
			args[i] = VisitExpression(arg, f)
		}
		return f(e, args)
	case Or:
		args := make([]T, len(e.Args))
		for i, arg := range e.Args {
			args[i] = VisitExpression(arg, f)
		}
		return f(e, args)
	case Not:
		return f(e, []T{VisitExpression(e.Arg, f)})
	default:
		return f(e, []T{})
	}
}

func Visit[T any](e Expression, fCall func(string, []T) T, fRef func(Variable) T, fConst func(Constant) T) T {
	switch e := e.(type) {
	case Variable:
		return fRef(e)
	case Constant:
		return fConst(e)
	case And:
		args := make([]T, len(e.Args))
		for i, arg := range e.Args {
			args[i] = Visit(arg, fCall, fRef, fConst)
		}
		return fCall("and", args)
	case Or:
		args := make([]T, len(e.Args))
		for i, arg := range e.Args {
			args[i] = Visit(arg, fCall, fRef, fConst)
		}
		return fCall("or", args)
	case Not:
		return fCall("not", []T{Visit(e.Arg, fCall, fRef, fConst)})
	case Eq:
		args := make([]T, len(e.Args))
		for i, arg := range e.Args {
			args[i] = Visit(arg, fCall, fRef, fConst)
		}
		return fCall("eq", args)
	case Ne:
		args := make([]T, len(e.Args))
		for i, arg := range e.Args {
			args[i] = Visit(arg, fCall, fRef, fConst)
		}
		return fCall("ne", args)
	case Lt:
		args := make([]T, len(e.Args))
		for i, arg := range e.Args {
			args[i] = Visit(arg, fCall, fRef, fConst)
		}
		return fCall("lt", args)
	case Le:
		args := make([]T, len(e.Args))
		for i, arg := range e.Args {
			args[i] = Visit(arg, fCall, fRef, fConst)
		}
		return fCall("le", args)
	case Gt:
		args := make([]T, len(e.Args))
		for i, arg := range e.Args {
			args[i] = Visit(arg, fCall, fRef, fConst)
		}
		return fCall("gt", args)
	case Ge:
		args := make([]T, len(e.Args))
		for i, arg := range e.Args {
			args[i] = Visit(arg, fCall, fRef, fConst)
		}
		return fCall("ge", args)
	case Between:
		args := make([]T, len(e.Args))
		for i, arg := range e.Args {
			args[i] = Visit(arg, fCall, fRef, fConst)
		}
		return fCall("between", args)
	case NotBetween:
		args := make([]T, len(e.Args))
		for i, arg := range e.Args {
			args[i] = Visit(arg, fCall, fRef, fConst)
		}
		return fCall("not_between", args)
	case In:
		args := make([]T, len(e.Args))
		for i, arg := range e.Args {
			args[i] = Visit(arg, fCall, fRef, fConst)
		}
		return fCall("in", args)
	case NotIn:
		args := make([]T, len(e.Args))
		for i, arg := range e.Args {
			args[i] = Visit(arg, fCall, fRef, fConst)
		}
		return fCall("not_in", args)
	case Like:
		args := []T{
			Visit(e.Arg, fCall, fRef, fConst),
			fConst(String(e.Pattern)),
		}
		if e.Escape != "" {
			args = append(args, fConst(String(e.Escape)))
		}
		return fCall("like", args)
	case NotLike:
		args := []T{
			Visit(e.Arg, fCall, fRef, fConst),
			fConst(String(e.Pattern)),
		}
		if e.Escape != "" {
			args = append(args, fConst(String(e.Escape)))
		}
		return fCall("not_like", args)
	case IsNull:
		return fCall("is_null", []T{Visit(e.Arg, fCall, fRef, fConst)})
	case IsNotNull:
		return fCall("is_not_null", []T{Visit(e.Arg, fCall, fRef, fConst)})
	case IsRestricted:
		if e.Not {
			return fCall("is_not_restricted", []T{fRef(Variable{Name: e.VariableName})})
		} else {
			return fCall("is_restricted", []T{fRef(Variable{Name: e.VariableName})})
		}
	}
	return fCall("unexpected_expression", []T{})
}
