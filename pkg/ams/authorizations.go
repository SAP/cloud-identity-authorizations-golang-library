package ams

import (
	"reflect"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/internal"
)

type Authorizations struct {
	policies  internal.PolicySet
	andJoined []*Authorizations
	schema    internal.Schema
	envInput  expression.Input
}

// Retrieve a access decision for a given action and resource and possibly some custom input
// the app input should correspond to the DCL schema definition and will be mapped into $app fields.
// This can be achieved by providing either:
//   - deeply nested map[string] where the keys are the schema names and the values can translated to the schema types
//   - a struct, thats fields are tagged with 'ams:"<fieldname>"' where the field name corresponds to the schema
//     name or the fields name is EXACTLY the same as the schema name
func (a Authorizations) Inquire(action, resource string, app any) expression.Expression {

	i := expression.Input{
		"$dcl.action":   expression.String(action),
		"$dcl.resource": expression.String(resource),
	}
	if action == "" {
		delete(i, "$dcl.action")
	}
	if resource == "" {
		delete(i, "$dcl.resource")
	}
	for k, v := range a.envInput {
		i[k] = v
	}
	a.schema.InsertCustomInput(i, reflect.ValueOf(app), []string{"$app"})
	return a.Evaluate(i)
}

func (a Authorizations) SetEnvInput(env any) {
	a.envInput = expression.Input{}
	a.schema.InsertCustomInput(a.envInput, reflect.ValueOf(env), []string{"$env"})
}

// Retrieve a access decision for a given action and resource and possibly some custom input
// this function is ment to provide generic quick access to the authorizations and is dangerous to use
// the provided input must be a map[string]expression.Constant where:
//
//   - the keys are the stringified qualified names from the schema (see util.StringifyQualifiedName)
//   - the values are the expression constants that match exactly the schema types
//   - the evaluation will panic if the input is wrongly typed
//
// the input can savely created/purged by the Schema.
func (a Authorizations) Evaluate(input expression.Input) expression.Expression {
	for k, v := range a.envInput {
		input[k] = v
	}
	r := a.policies.Evaluate(input)
	for k := range a.envInput {
		delete(input, k)
	}
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
	return expression.And(results...)
}

// Restrict an authorizations object by another one
// a possible scenario would be to restrict a users authorizations by other technical authorizations.
func (a Authorizations) AndJoin(aa *Authorizations) *Authorizations {
	return &Authorizations{
		policies:  a.policies,
		andJoined: append(a.andJoined, aa),
	}
}
