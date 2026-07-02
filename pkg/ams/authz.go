package ams

import (
	"reflect"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/internal/runtime"
)

type AuthorizationsNew struct {
	policyNames []string
	envInput    expression.Input
	r           *runtime.Runtime
	andJoined   []*AuthorizationsNew
}

func (a *AuthorizationsNew) Authorize(action string, resource string, input any) *Decision {
	internalInput := a.envInput.Copy()
	if input != nil {
		inputValue := reflect.ValueOf(input)
		runtime.InsertInput(internalInput, inputValue, []string{"$app"})
		if _, ok := a.r.Resources[resource]; ok {
			runtime.InsertInput(internalInput, inputValue, []string{"$resource"})
		}
		a.r.SanitizeInput(internalInput, resource)
	}
	expr := a.r.Evaluate(a.policyNames, action, resource, internalInput)
	return &Decision{
		condition:      expr,
		inputConverter: a.inputConverter(resource),
	}
}

func (a *AuthorizationsNew) AuthorizeRaw(action string, resource string, input expression.Input) *Decision {
	internalInput := a.envInput.Copy()
	for k, v := range input {
		internalInput[k] = v
	}
	a.r.SanitizeInput(internalInput, resource)
	expr := a.r.Evaluate(a.policyNames, action, resource, internalInput)
	return &Decision{
		condition:      expr,
		inputConverter: a.inputConverter(resource),
	}
}

func (a *AuthorizationsNew) GetResources() []string {
	return a.r.GetResources(a.policyNames)
}

func (a *AuthorizationsNew) GetActions(resource string) []string {
	return a.r.GetActions(a.policyNames, resource)
}

func (a *AuthorizationsNew) SetEnvInput(env any) {
	a.envInput = expression.Input{}
	runtime.InsertInput(a.envInput, reflect.ValueOf(env), []string{"$env"})
	a.r.SanitizeInput(a.envInput, "")
}

func (a *AuthorizationsNew) And(other *AuthorizationsNew) *AuthorizationsNew {
	return &AuthorizationsNew{
		andJoined:   []*AuthorizationsNew{a, other},
		r:           a.r,
		policyNames: a.policyNames,
		envInput:    a.envInput,
	}
}

func (a *AuthorizationsNew) inputConverter(resource string) func(any) expression.Input {
	return func(input any) expression.Input {
		result := make(expression.Input)
		if input == nil {
			return result
		}

		inputValue := reflect.ValueOf(input)
		runtime.InsertInput(result, inputValue, []string{"$app"})
		if _, ok := a.r.Resources[resource]; ok {
			runtime.InsertInput(result, inputValue, []string{"$resource"})
		}
		a.r.SanitizeInput(result, resource)
		return result
	}
}
