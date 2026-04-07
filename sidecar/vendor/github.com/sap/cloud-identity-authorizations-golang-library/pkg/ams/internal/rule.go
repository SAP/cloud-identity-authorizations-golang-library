package internal

import (
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

type Rule struct {
	asExpression expression.Expression
	actions      []string
	resources    []string
}

func RuleFromDCN(rawRule dcn.Rule, f *expression.FunctionRegistry) (Rule, error) {
	var rule Rule
	args := []expression.Expression{}
	rule.actions = []string{}
	rule.resources = []string{}

	if rawRule.Condition != nil {
		cond, err := expression.FromDCN(*rawRule.Condition, f)
		if err != nil {
			return rule, err
		}
		args = append(args, cond.Expression)
	}
	if len(rawRule.Actions) > 0 {
		rule.actions = rawRule.Actions
		args = append(args, expression.In(
			expression.Ref("$dcl.action"),
			expression.ConstantFrom(rawRule.Actions),
		))
	}
	if len(rawRule.Resources) > 0 {
		rule.resources = rawRule.Resources
		args = append(args, expression.In(
			expression.Ref("$dcl.resource"),
			expression.ConstantFrom(rawRule.Resources),
		))
	}
	rule.asExpression = expression.And(args...)

	return rule, nil
}

func (r *Rule) Evaluate(input expression.Input) expression.Expression {
	result := r.asExpression.Evaluate(input)
	return result
}
