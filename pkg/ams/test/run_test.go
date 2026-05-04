package test

import (
	"context"
	"fmt"
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/util"
)

type crashLogger struct{}

func (l crashLogger) Debugf(ctx context.Context, format string, args ...interface{}) {}
func (l crashLogger) Infof(ctx context.Context, format string, args ...interface{})  {}
func (l crashLogger) Warnf(ctx context.Context, format string, args ...interface{})  {}
func (l crashLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	panic(fmt.Sprintf(format, args...))
}
func TestRun(t *testing.T) {
	// tmp, _ := os.ReadDir("./")
	// t.Fatalf("tmp: %v", tmp)
	testDirs, err := os.ReadDir("scenarios")
	if err != nil {
		t.Fatalf("can't open test directories: %v", err)
	}
	for _, testDir := range testDirs {
		t.Run(testDir.Name(), func(t *testing.T) {
			loader := dcn.NewLocalLoader(path.Join("scenarios", testDir.Name()), crashLogger{})
			tests := []dcn.Test{}
			dcnChannel := make(chan dcn.DcnContainer)
			go func() {
				for {
					dcnContainer := <-loader.DCNChannel
					tests = dcnContainer.Tests
					dcnChannel <- dcnContainer
				}
			}()
			ams := ams.NewAuthorizationManager(context.Background(), dcnChannel, loader.AssignmentsChannel, crashLogger{})

			<-ams.WhenReady()

			for _, test := range tests {
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
						authz := ams.AuthorizationsForPolicies(context.Background(), policies)
						if len(scopeFilter) > 0 {
							scopeFilter := ams.AuthorizationsForPolicies(context.Background(), scopeFilter)
							authz = authz.AndJoin(scopeFilter)
						}
						t.Run(fmt.Sprintf("policies: %v, scopeFilter: %v", policies, scopeFilter), func(t *testing.T) {
							for _, action := range actions {
								for _, resource := range resources {
									for _, tInput := range inputs {
										t.Run(assertionCaption(action, resource, tInput), func(t *testing.T) {
											input := createInput(ams, tInput, action, resource)

											result := authz.Evaluate(input).Condition()
											result = unsetIgnore(result, tInput)
											result = NormalizeExpression(result)
											expectedContainer, err := expression.FromDCN(assertion.Expect, &expression.FunctionRegistry{})
											expected := NormalizeExpression(expectedContainer.Expression)
											if err != nil {
												t.Fatalf("error in expected expression: %v", err)
											}
											if !reflect.DeepEqual(result, expected) {
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

func createInput(am *ams.AuthorizationManager, input dcn.Input, action, resource string) expression.Input {
	app, ok := input.Input["$app"]
	if !ok {
		app = nil
	}
	env, ok := input.Input["$env"]
	if !ok {
		env = nil
	}
	result := am.CreateInput(action, resource, app, env)

	return result
}
