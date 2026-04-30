package httpclient

import (
	"context"
	"reflect"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

type Authorizations struct {
	ctx      context.Context
	identity ams.Identity
	policies []string
	andJoin  []*Authorizations
	envInput reqInput
	client   *AuthorizationManager
}

func (a *Authorizations) Evaluate(ctx context.Context, input expression.Input) (Decision, error) {
	reqInput := make(reqInput)
	for k, v := range input {
		reqInput[k] = v
	}
	action := ""
	resource := ""
	if a, ok := input["$dcl.action"]; ok {
		action = a.String()
	}
	if r, ok := input["$dcl.resource"]; ok {
		resource = r.String()
	}

	return a.evaluate(ctx, action, resource, reqInput)
}

func (a *Authorizations) evaluate(ctx context.Context, action, resource string, input reqInput) (Decision, error) {
	token := ""
	reqInput := make(reqInput)
	for k, v := range input {
		reqInput[k] = v
	}
	for k, v := range a.envInput {
		reqInput[k] = v
	}
	if a.identity != nil {
		token = newToken(tokenClaim{
			"scim_id": a.identity.ScimID(),
			"app_tid": a.identity.AppTID(),
		})
	}
	req := AuthorizationRequest{
		Action:   action,
		Resource: resource,
		Policies: a.policies,
		Token:    token,
		Input:    reqInput,
	}
	res := AuthorizationResponse{}
	err := a.client.post(ctx, PATH_AUTHORIZE, req, &res)
	if err != nil {
		return Decision{condition: expression.FALSE}, err
	}
	result, err := a.decisionForDCN(ctx, res.Result)
	for _, aa := range a.andJoin {
		r, err := aa.evaluate(ctx, action, resource, input)
		if err != nil {
			return Decision{condition: expression.FALSE}, err
		}
		if r.Condition() == expression.FALSE {
			return Decision{
				condition:      expression.FALSE,
				inputConverter: result.inputConverter,
			}, nil
		}
		if r.Condition() != expression.TRUE {
			result.condition = expression.And(result.condition, r.Condition())
		}
	}
	return result, nil
}

func (a *Authorizations) AndJoin(other *Authorizations) *Authorizations {
	a.andJoin = append(a.andJoin, other)
	return a
}

func (a *Authorizations) GetActions(ctx context.Context, resource string) ([]string, error) {
	token := ""
	if a.identity != nil {
		token = newToken(tokenClaim{
			"scim_id": a.identity.ScimID(),
			"app_tid": a.identity.AppTID(),
		})
	}
	req := ActionsRequest{
		Policies: a.policies,
		Token:    token,
		Resource: resource,
	}
	var response ActionsResponse
	err := a.client.post(ctx, PATH_ACTIONS, req, &response)
	if err != nil {
		return nil, err
	}
	return response.Actions, nil
}

func (a *Authorizations) GetResources(ctx context.Context) ([]string, error) {
	token := ""
	if a.identity != nil {
		token = newToken(tokenClaim{
			"scim_id": a.identity.ScimID(),
			"app_tid": a.identity.AppTID(),
		})
	}
	req := ResourcesRequest{
		Policies: a.policies,
		Token:    token,
	}
	var response ResourcesResponse
	err := a.client.post(ctx, PATH_RESOURCES, req, &response)
	if err != nil {
		return nil, err
	}
	return response.Resources, nil
}

func (a *Authorizations) Inquire(ctx context.Context, action, resource string, app any) (Decision, error) {
	input := reqInput{
		"$dcl.action":   expression.String(action),
		"$dcl.resource": expression.String(resource),
	}
	insertCustomInput(input, reflect.ValueOf(app), []string{"$app"})
	return a.evaluate(ctx, action, resource, input)
}

func (a *Authorizations) SetEnvInput(env any) {
	insertCustomInput(a.envInput, reflect.ValueOf(env), []string{"$env"})
}

func (a *Authorizations) decisionForDCN(ctx context.Context, dcnExpression dcn.Expression) (Decision, error) {
	condition, err := expression.FromDCN(dcnExpression, nil)
	inputConverter := func(app any) (expression.Input, error) {
		return a.client.CreateInput(ctx, "", "", app, nil)
	}

	if err != nil {
		return Decision{
			condition:      expression.FALSE,
			inputConverter: inputConverter,
		}, err
	}
	return Decision{
		condition:      condition.Expression,
		inputConverter: inputConverter,
	}, nil
}
