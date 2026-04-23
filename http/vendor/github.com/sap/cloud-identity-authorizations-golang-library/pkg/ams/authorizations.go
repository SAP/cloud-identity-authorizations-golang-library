package ams

import (
	"reflect"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/internal"
)

type Authorizations interface {
	Inquire(action, resource string, app any) Decision

	SetEnvInput(env any)

	GetResources() []string

	GetActions(resource string) []string

	Evaluate(input expression.Input) Decision

	AndJoin(aa Authorizations) Authorizations
}

type authorizations struct {
	policies  internal.PolicySet
	andJoined []Authorizations
	schema    internal.Schema
	envInput  expression.Input
}

const (
	DCL_ACTION   = "$dcl.action"
	DCL_RESOURCE = "$dcl.resource"
)

// Retrieve a access decision for a given action and resource and possibly some custom input
// the app input should correspond to the DCL schema definition and will be mapped into $app fields.
// This can be achieved by providing either:
//   - deeply nested map[string] where the keys are the schema names and the values can translated to the schema types
//   - a struct, thats fields are tagged with 'ams:"<fieldname>"' where the field name corresponds to the schema
//     name or the fields name is EXACTLY the same as the schema name
func (a *authorizations) Inquire(action, resource string, app any) Decision {
	i := a.schema.CustomInput(action, resource, app, nil)
	for k, v := range a.envInput {
		i[k] = v
	}

	return a.Evaluate(i)
}

func (a *authorizations) SetEnvInput(env any) {
	a.envInput = expression.Input{}
	a.schema.InsertCustomInput(a.envInput, reflect.ValueOf(env), []string{"$env"})
}

func (a *authorizations) GetResources() []string {
	return a.policies.GetResources()
}

func (a *authorizations) GetActions(resource string) []string {
	return a.policies.GetActions(resource)
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
func (a *authorizations) Evaluate(input expression.Input) Decision {
	for k, v := range a.envInput {
		input[k] = v
	}
	r := a.policies.Evaluate(input)
	for k := range a.envInput {
		delete(input, k)
	}
	if r == expression.FALSE {
		return Decision{
			condition: r,
			inputConverter: func(app any) expression.Input {
				result := expression.Input{}
				a.schema.InsertCustomInput(result, reflect.ValueOf(app), []string{"$app"})
				return result
			},
		}
	}
	results := []expression.Expression{
		r,
	}

	for _, aa := range a.andJoined {
		r := aa.Evaluate(input).Condition()
		if r == expression.Bool(false) {
			return Decision{
				condition: r,
				inputConverter: func(app any) expression.Input {
					result := expression.Input{}
					a.schema.InsertCustomInput(result, reflect.ValueOf(app), []string{"$app"})
					return result
				},
			}
		}
		if r != expression.Bool(true) {
			results = append(results, r)
		}
	}
	return Decision{
		condition: expression.And(results...),
		inputConverter: func(app any) expression.Input {
			result := expression.Input{}
			a.schema.InsertCustomInput(result, reflect.ValueOf(app), []string{"$app"})
			return result
		},
	}
}

// Restrict an authorizations object by another one
// a possible scenario would be to restrict a users authorizations by other technical authorizations.
func (a *authorizations) AndJoin(aa Authorizations) Authorizations {
	return &authorizations{
		policies:  a.policies,
		andJoined: append(a.andJoined, aa),
	}
}
