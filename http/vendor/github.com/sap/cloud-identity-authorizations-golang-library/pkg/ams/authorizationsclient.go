package ams

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/httpclient"
)

type tokenClaim map[string]any
type authorizationsClient struct {
	c           http.Client
	url         string
	errHandlers []func(error)
}

type authorizationsObj struct {
	identity Identity
	policies []string
	andJoin  []Authorizations
	envInput any
	client   *authorizationsClient
}

func NewAuthorizationManagerClient(url string, client http.Client) AuthorizationManager {
	return &authorizationsClient{
		c:   client,
		url: url,
	}
}

func (a *authorizationsClient) IsReady() bool {
	return a.get(httpclient.PATH_HEALTH, nil) == nil
}

func (a *authorizationsClient) WhenReady() <-chan bool {
	ch := make(chan bool)
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			if a.IsReady() {
				ch <- true
				return
			}
			<-ticker.C
		}
	}()
	return ch
}

func (a *authorizationsClient) AuthorizationsForIdentity(i Identity) Authorizations {
	return &authorizationsObj{
		identity: i,
		client:   a,
		andJoin:  []Authorizations{},
	}
}

func (a *authorizationsClient) AuthorizationsForPolicies(policyNames []string) Authorizations {
	return &authorizationsObj{
		policies: policyNames,
		client:   a,
		andJoin:  []Authorizations{},
	}
}

func (a *authorizationsClient) GetDefaultPolicyNames(tenant string) []string {
	var response httpclient.DefaultPoliciesResponse
	err := a.get(httpclient.PATH_DEFAULT_POLICIES+"/"+tenant, &response)
	if err != nil {
		a.notifyError(err)
		return nil
	}
	return response.DefaultPolicies
}

func (a *authorizationsClient) GetAssignments(tenant, user string) []string {
	req := httpclient.AssignedPoliciesRequest{
		Token: newToken(tokenClaim{
			"scim_id": user,
			"app_tid": tenant,
		}),
	}
	var response httpclient.AssignedPoliciesResponse
	err := a.post(httpclient.PATH_ASSIGNED_POLICIES, req, &response)
	if err != nil {
		a.notifyError(err)
		return nil
	}
	return response.Policies
}

func (a *authorizationsClient) RegisterErrorHandler(handler func(error)) {
	a.errHandlers = append(a.errHandlers, handler)
}

func (a *authorizationsClient) CreateInput(action, resource string, input any, env any) expression.Input {
	req := httpclient.CreateInputRequest{
		Action:   action,
		Resource: resource,
		Input:    httpclient.ConvertInput(input),
		Env:      httpclient.ConvertInput(env),
	}
	var response httpclient.CreateInputResponse
	err := a.post(httpclient.PATH_CREATE_INPUT, req, &response)
	if err != nil {
		a.notifyError(err)
		return nil
	}
	return response.RawInput
}

func (a *authorizationsClient) get(path string, responseBody any) error {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, a.url+path, nil)
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

func (a *authorizationsClient) post(path string, requestBody any, responseBody any) error {
	reqBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(
		context.Background(),
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

func (a *authorizationsClient) notifyError(err error) {
	for _, handler := range a.errHandlers {
		handler(err)
	}
}

func (a *authorizationsObj) decisionForDCN(dcnExpression dcn.Expression) Decision {
	condition, err := expression.FromDCN(dcnExpression, nil)
	inputConverter := func(app any) expression.Input {
		req := httpclient.CreateInputRequest{
			Action:   "", // action and resource are not relevant for the input conversion, as the condition is already evaluated
			Resource: "",
			Input:    httpclient.ConvertInput(app),
		}
		var response httpclient.CreateInputResponse
		err := a.client.post(httpclient.PATH_CREATE_INPUT, req, &response)
		if err != nil {
			a.client.notifyError(err)
			return nil
		}
		return response.RawInput
	}

	if err != nil {
		a.client.notifyError(err)
		return Decision{
			condition:      expression.FALSE,
			inputConverter: inputConverter,
		}
	}
	return Decision{
		condition:      condition.Expression,
		inputConverter: inputConverter,
	}
}

func (a *authorizationsObj) Evaluate(input expression.Input) Decision {
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
	req := httpclient.UnvalidatedAuthorizationRequest{
		Policies: a.policies,
		Token:    token,
		Input:    input,
	}
	res := httpclient.AuthorizationResponse{}
	err := a.client.post(httpclient.PATH_AUTHORIZE_UNVALIDATED, req, &res)
	if err != nil {
		a.client.notifyError(err)
		return Decision{condition: expression.FALSE}
	}
	result := a.decisionForDCN(res.Result)
	if result.Condition() == expression.FALSE {
		return result
	}
	for _, aa := range a.andJoin {
		r := aa.Evaluate(input).Condition()
		if r == expression.Bool(false) {
			return Decision{
				condition:      r,
				inputConverter: result.inputConverter,
			}
		}
		if r != expression.Bool(true) {
			result.condition = expression.And(result.condition, r)
		}
	}
	return result
}

func (a *authorizationsObj) AndJoin(other Authorizations) Authorizations {
	a.andJoin = append(a.andJoin, other)
	return a
}

func (a *authorizationsObj) GetActions(resource string) []string {
	token := ""
	if a.identity != nil {
		token = newToken(tokenClaim{
			"scim_id": a.identity.ScimID(),
			"app_tid": a.identity.AppTID(),
		})
	}
	req := httpclient.ActionsRequest{
		Policies: a.policies,
		Token:    token,
		Resource: resource,
	}
	var response httpclient.ActionsResponse
	err := a.client.post(httpclient.PATH_ACTIONS, req, &response)
	if err != nil {
		a.client.notifyError(err)
		return nil
	}
	return response.Actions
}

func (a *authorizationsObj) GetResources() []string {
	token := ""
	if a.identity != nil {
		token = newToken(tokenClaim{
			"scim_id": a.identity.ScimID(),
			"app_tid": a.identity.AppTID(),
		})
	}
	req := httpclient.ResourcesRequest{
		Policies: a.policies,
		Token:    token,
	}
	var response httpclient.ResourcesResponse
	err := a.client.post(httpclient.PATH_RESOURCES, req, &response)
	if err != nil {
		a.client.notifyError(err)
		return nil
	}
	return response.Resources
}

func (a *authorizationsObj) Inquire(action, resource string, app any) Decision {
	token := ""
	if a.identity != nil {
		token = newToken(tokenClaim{
			"scim_id": a.identity.ScimID(),
			"app_tid": a.identity.AppTID(),
		})
	}
	req := httpclient.AuthorizationRequest{
		Action:   action,
		Resource: resource,
		Policies: a.policies,
		Token:    token,
		Input:    httpclient.ConvertInput(app),
	}
	var response httpclient.AuthorizationResponse
	err := a.client.post(httpclient.PATH_AUTHORIZE, req, &response)
	if err != nil {
		a.client.notifyError(err)
		return Decision{condition: expression.FALSE}
	}
	result := a.decisionForDCN(response.Result)
	if result.Condition() == expression.FALSE {
		return result
	}
	for _, aa := range a.andJoin {
		r := aa.Inquire(action, resource, app).Condition()
		if r == expression.Bool(false) {
			return Decision{
				condition:      r,
				inputConverter: result.inputConverter,
			}
		}
		if r != expression.Bool(true) {
			result.condition = expression.And(result.condition, r)
		}
	}
	return result
}

func (a *authorizationsObj) SetEnvInput(env any) {
	a.envInput = env
}

func (a *authorizationsClient) ValidateInput(input expression.Input) ([]string, []string) {
	panic("not implemented")
}

func newToken(claims tokenClaim) string {
	jwtRaw, err := json.Marshal(claims)
	if err != nil {
		panic(err)
	}
	return string(jwtRaw)
}
