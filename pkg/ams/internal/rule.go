package internal

import (
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

type Rule struct {
	asExpression expression.Expression
}

func RuleFromDCN(rawRule dcn.Rule, f expression.Functions) (Rule, error) {
	var rule Rule
	args := []expression.Expression{}

	if rawRule.Condition != nil {
		cond, err := expression.FromDCN(*rawRule.Condition, f)
		if err != nil {
			return rule, err
		}
		args = append(args, cond.Expression)
	}
	if len(rawRule.Actions) > 0 {
		args = append(args, expression.In{
			Args: []expression.Expression{
				expression.Reference{Name: "$dcl.action"},
				expression.ConstantFrom(rawRule.Actions),
			},
		})
	}
	if len(rawRule.Resources) > 0 {
		args = append(args, expression.In{
			Args: []expression.Expression{
				expression.Reference{Name: "$dcl.resource"},
				expression.ConstantFrom(rawRule.Resources),
			},
		})
	}
	rule.asExpression = expression.NewAnd(args...)
	return rule, nil
}

func (r *Rule) Evaluate(input expression.Input) expression.Expression {
	result := r.asExpression.Evaluate(input)
	if result == expression.UNSET {
		return expression.Bool(false)
	}
	if result == expression.IGNORE {
		return expression.Bool(true)
	}
	return result
}
