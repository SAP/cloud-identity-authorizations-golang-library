package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/logging"
)

type tokenClaim map[string]any
type AuthorizationManager struct {
	c           http.Client
	url         string
	errHandlers []func(error)
}

type Authorizations struct {
	ctx      context.Context
	identity ams.Identity
	policies []string
	andJoin  []*Authorizations
	envInput any
	client   *AuthorizationManager
}

func NewAuthorizationManager(url string, client http.Client, logger logging.Logger) *AuthorizationManager {
	return &AuthorizationManager{
		c:   client,
		url: url,
	}
}

func (a *AuthorizationManager) IsReady(ctx context.Context) bool {
	return a.get(ctx, PATH_HEALTH, nil) == nil
}

func (a *AuthorizationManager) WhenReady(ctx context.Context) <-chan bool {
	ch := make(chan bool)
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			if a.IsReady(ctx) {
				ch <- true
				return
			}
			<-ticker.C
		}
	}()
	return ch
}

func (a *AuthorizationManager) AuthorizationsForIdentity(ctx context.Context, i ams.Identity) *Authorizations {
	return &Authorizations{
		ctx:      ctx,
		identity: i,
		client:   a,
		andJoin:  []*Authorizations{},
	}
}

func (a *AuthorizationManager) AuthorizationsForPolicies(ctx context.Context, policyNames []string) *Authorizations {
	return &Authorizations{
		ctx:      ctx,
		policies: policyNames,
		client:   a,
		andJoin:  []*Authorizations{},
	}
}

func (a *AuthorizationManager) GetDefaultPolicyNames(ctx context.Context, tenant string) ([]string, error) {
	var response DefaultPoliciesResponse
	err := a.get(ctx, PATH_DEFAULT_POLICIES+"/"+tenant, &response)
	if err != nil {
		return nil, err
	}
	return response.DefaultPolicies, nil
}

func (a *AuthorizationManager) GetAssignments(ctx context.Context, tenant, user string) ([]string, error) {
	req := AssignedPoliciesRequest{
		Token: newToken(tokenClaim{
			"scim_id": user,
			"app_tid": tenant,
		}),
	}
	var response AssignedPoliciesResponse
	err := a.post(ctx, PATH_ASSIGNED_POLICIES, req, &response)
	if err != nil {
		return nil, err
	}
	return response.Policies, nil
}

func (a *AuthorizationManager) CreateInput(ctx context.Context, action, resource string, input any, env any) (expression.Input, error) {
	req := CreateInputRequest{
		Action:   action,
		Resource: resource,
		Input:    ConvertInput(input),
		Env:      ConvertInput(env),
	}
	var response CreateInputResponse
	err := a.post(ctx, PATH_CREATE_INPUT, req, &response)
	if err != nil {
		return nil, err
	}
	return response.RawInput, nil
}

func (a *AuthorizationManager) get(ctx context.Context, path string, responseBody any) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.url+path, nil)
	if err != nil {
		return err
	}
	resp, err := a.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected on GET %s status code: %d", a.url+path, resp.StatusCode)
	}
	if responseBody == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(responseBody)
}

func (a *AuthorizationManager) post(ctx context.Context, path string, requestBody any, responseBody any) error {
	reqBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		a.url+path,
		bytes.NewReader(reqBodyBytes),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected on POST %s status code: %d", a.url+path, resp.StatusCode)
	}
	if responseBody == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(responseBody)
}

func (a *Authorizations) decisionForDCN(ctx context.Context, dcnExpression dcn.Expression) (Decision, error) {
	condition, err := expression.FromDCN(dcnExpression, nil)
	inputConverter := func(app any) (expression.Input, error) {
		req := CreateInputRequest{
			Action:   "", // action and resource are not relevant for the input conversion, as the condition is already evaluated
			Resource: "",
			Input:    ConvertInput(app),
		}
		var response CreateInputResponse
		err := a.client.post(ctx, PATH_CREATE_INPUT, req, &response)
		if err != nil {
			return nil, err
		}
		return response.RawInput, nil
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

func (a *Authorizations) Evaluate(ctx context.Context, input expression.Input) (Decision, error) {
	token := ""
	if a.identity != nil {
		if input == nil {
			input = expression.Input{}
		}
		input["$env.$user.email"] = expression.String(a.identity.Email())
		input["$env.$user.user_uuid"] = expression.String(a.identity.UserUUID())
		input["$env.$user.groups"] = expression.ArrayFrom(a.identity.Groups())
		token = newToken(tokenClaim{
			"scim_id": a.identity.ScimID(),
			"app_tid": a.identity.AppTID(),
		})
	}
	req := UnvalidatedAuthorizationRequest{
		Policies: a.policies,
		Token:    token,
		Input:    input,
	}
	res := AuthorizationResponse{}
	err := a.client.post(ctx, PATH_AUTHORIZE_UNVALIDATED, req, &res)
	if err != nil {
		return Decision{condition: expression.FALSE}, err
	}
	result, err := a.decisionForDCN(ctx, res.Result)
	if err != nil {
		return Decision{condition: expression.FALSE}, err
	}
	if result.Condition() == expression.FALSE {
		return result, nil
	}
	for _, aa := range a.andJoin {
		r, err := aa.Evaluate(ctx, input)
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
	token := ""
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
		Input:    ConvertInput(app),
	}
	var response AuthorizationResponse
	err := a.client.post(ctx, PATH_AUTHORIZE, req, &response)
	if err != nil {
		return Decision{condition: expression.FALSE}, err
	}
	result, err := a.decisionForDCN(ctx, response.Result)
	if err != nil {
		return Decision{condition: expression.FALSE}, err
	}
	if result.Condition() == expression.FALSE {
		return result, nil
	}
	for _, aa := range a.andJoin {
		r, err := aa.Inquire(ctx, action, resource, app)
		if err != nil {
			return Decision{condition: expression.FALSE}, err
		}
		if r.Condition() == expression.FALSE {
			return Decision{
				condition:      r.Condition(),
				inputConverter: result.inputConverter,
			}, nil
		}
		if r.Condition() != expression.Bool(true) {
			result.condition = expression.And(result.condition, r.Condition())
		}
	}
	return result, nil
}

func (a *Authorizations) SetEnvInput(env any) {
	a.envInput = env
}

func (a *AuthorizationManager) ValidateInput(input expression.Input) ([]string, []string) {
	panic("not implemented")
}

func newToken(claims tokenClaim) string {
	jwtRaw, err := json.Marshal(claims)
	if err != nil {
		panic(err)
	}
	return string(jwtRaw)
}
