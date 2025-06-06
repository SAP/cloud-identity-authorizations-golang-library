package test

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/internal"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/util"
)

func TestRun(t *testing.T) {
	// tmp, _ := os.ReadDir("./")
	// t.Fatalf("tmp: %v", tmp)
	testDirs, err := os.ReadDir("scenarios")
	if err != nil {
		t.Fatalf("can't open test directories: %v", err)
	}
	for _, testDir := range testDirs {
		t.Run(testDir.Name(), func(t *testing.T) {
			ams := ams.NewAuthorizationManagerForLocal(path.Join("scenarios", testDir.Name()))

			ams.RegisterErrorHandler(func(err error) {
				t.Errorf("error in authorization manager: %v", err)
				panic(err)
			})

			<-ams.WhenReady()

			for _, test := range ams.Tests {
				t.Run(util.StringifyQualifiedName(test.Test), func(t *testing.T) {
					for _, assertion := range test.Assertions {
						actions := assertion.Actions
						resources := assertion.Resources
						inputs := assertion.Inputs
						if len(actions) == 0 {
							actions = []string{""}
						}
						if len(resources) == 0 {
							resources = []string{""}
						}
						if len(inputs) == 0 {
							inputs = []dcn.Input{{
								Input:    make(map[string]any),
								Unknowns: []dcn.Reference{},
								Ignores:  []dcn.Reference{},
							}}
						}
						policies := ams.GetDefaultPolicyNames("")
						for _, policy := range assertion.Policies {
							policies = append(policies, util.StringifyQualifiedName(policy))
						}
						scopeFilter := []string{}
						for _, filter := range assertion.ScopeFilter {
							scopeFilter = append(scopeFilter, util.StringifyQualifiedName(filter))
						}
						authz := ams.GetAuthorizations(policies)
						if len(scopeFilter) > 0 {
							scopeFilter := ams.GetAuthorizations(scopeFilter)
							authz = authz.AndJoin(scopeFilter)
						}
						t.Run(fmt.Sprintf("policies: %v, scopeFilter: %v", policies, scopeFilter), func(t *testing.T) {
							for _, action := range actions {
								for _, resource := range resources {
									for _, tInput := range inputs {
										t.Run(assertionCaption(action, resource, tInput), func(t *testing.T) {
											input := createInput(ams.GetSchema(), tInput, action, resource)

											result := authz.Evaluate(input)
											result = unsetIgnore(result, tInput)
											result = NormalizeExpression(result)
											expectedContainer, err := expression.FromDCN(assertion.Expect, &expression.FunctionRegistry{})
											expected := NormalizeExpression(expectedContainer.Expression)
											if err != nil {
												t.Fatalf("error in expected expression: %v", err)
											}
											if !reflect.DeepEqual(result, expected) {
												createInput(ams.GetSchema(), tInput, action, resource)
												authz.Evaluate(input)
												t.Errorf("expected %v, got %v", expected, result)
											}
										})
									}
								}
							}
						})
					}
				})
			}
		})
	}
}

func assertionCaption(action string, resource string, input dcn.Input) string {
	return fmt.Sprintf("action: %s, resource: %s, input: %+v", action, resource, input)
}

func unsetIgnore(e expression.Expression, input dcn.Input) expression.Expression {
	u := map[string]bool{}
	i := map[string]bool{}
	for _, ref := range input.Unknowns {
		u[util.StringifyQualifiedName(ref.Ref)] = true
	}
	for _, ref := range input.Ignores {
		i[util.StringifyQualifiedName(ref.Ref)] = true
	}
	return expression.UnknownIgnore(e, u, i) //nolint:staticcheck
}

func createInput(schema internal.Schema, input dcn.Input, action, resource string) expression.Input {
	app, ok := input.Input["$app"]
	if !ok {
		app = nil
	}
	env, ok := input.Input["$env"]
	if !ok {
		env = nil
	}
	result := schema.CustomInput(action, resource, app, env)

	return result
}
