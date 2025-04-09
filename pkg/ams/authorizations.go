package ams

import (
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/internal"
)

type Authorizations struct {
	policies  internal.PolicySet
	andJoined []*Authorizations
	schema    internal.Schema
}

func (a Authorizations) Inquire(action, resource string, input any, env any) expression.Expression {
	i := a.schema.CustomInput(action, resource, input, env)
	return a.Evaluate(i)
}

func (a Authorizations) Evaluate(input expression.Input) expression.Expression {
	r := a.policies.Evaluate(input)
	if r == expression.FALSE {
		return r
	}
	results := []expression.Expression{
		r,
	}

	for _, aa := range a.andJoined {
		r := aa.Evaluate(input)
		if r == expression.Bool(false) {
			return r
		}
		if r != expression.Bool(true) {
			results = append(results, r)
		}
	}
	return expression.NewAnd(results...)
}

func (a Authorizations) AndJoin(aa *Authorizations) *Authorizations {
	return &Authorizations{
		policies:  a.policies,
		andJoined: append(a.andJoined, aa),
	}
}
