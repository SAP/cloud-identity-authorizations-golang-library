package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/logging"
)

type tokenClaim map[string]any
type AuthorizationManager struct {
	c   *http.Client
	url string
	l   logging.Logger
}

func NewAuthorizationManager(url string, client *http.Client, logger logging.Logger) *AuthorizationManager {
	return &AuthorizationManager{
		c:   client,
		url: url,
		l:   logger,
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
	if i == nil {
		return &Authorizations{
			ctx:      ctx,
			identity: nil,
			client:   a,
			andJoin:  []*Authorizations{},
			envInput: reqInput{},
			l:        a.l,
		}
	}
	return &Authorizations{
		ctx:      ctx,
		identity: i,
		client:   a,
		andJoin:  []*Authorizations{},
		l:        a.l,
		envInput: reqInput{
			"$env.$user.email":     expression.String(i.Email()),
			"$env.$user.user_uuid": expression.String(i.UserUUID()),
			"$env.$user.groups":    expression.ArrayFrom(i.Groups()),
		},
	}
}

func (a *AuthorizationManager) AuthorizationsForPolicies(ctx context.Context, policyNames []string) *Authorizations {
	return &Authorizations{
		ctx:      ctx,
		policies: policyNames,
		client:   a,
		andJoin:  []*Authorizations{},
		envInput: reqInput{},
		l:        a.l,
	}
}

func (a *AuthorizationManager) GetDefaultPolicyNames(ctx context.Context, tenant string) ([]string, error) {
	var response DefaultPoliciesResponse
	err := a.get(ctx, PATH_DEFAULT_POLICIES+"/"+tenant, &response)
	if err != nil {
		a.l.Errorf(ctx, "Error getting default policies for tenant %s: %v", tenant, err)
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
		a.l.Errorf(ctx, "Error getting assigned policies for tenant %s and user %s: %v", tenant, user, err)
		return nil, err
	}
	return response.Policies, nil
}

func (a *AuthorizationManager) CreateInput(
	ctx context.Context,
	action,
	resource string,
	input any,
	env any,
) (expression.Input, error) {
	reqInput := reqInput{}

	insertCustomInput(reqInput, reflect.ValueOf(input), []string{"$app"})
	insertCustomInput(reqInput, reflect.ValueOf(env), []string{"$env"})

	req := InputRequest{
		Action:   action,
		Resource: resource,
		Input:    reqInput,
	}
	var response InputResponse
	err := a.post(ctx, PATH_CREATE_INPUT, req, &response)
	if err != nil {
		a.l.Errorf(ctx, "Error creating input for action %s and resource %s: %v", action, resource, err)
		return nil, err
	}
	return expression.Input(response.Input), nil
}

func (a *AuthorizationManager) get(ctx context.Context, path string, responseBody any) error {
	result := make(chan error, 1)
	go func() {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.url+path, nil)
		if err != nil {
			a.l.Errorf(ctx, "Error creating GET request for path %s: %v", path, err)
			result <- err
			return
		}
		resp, err := a.c.Do(req)
		if err != nil {
			a.l.Errorf(ctx, "Error executing GET request for path %s: %v", path, err)
			result <- err
			return
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			a.l.Errorf(ctx, "Unexpected status code for GET request to path %s: %d", path, resp.StatusCode)
			result <- fmt.Errorf("unexpected on GET %s status code: %d", a.url+path, resp.StatusCode)
			return
		}
		if responseBody == nil {
			result <- nil
			return
		}
		result <- json.NewDecoder(resp.Body).Decode(responseBody)
	}()

	select {
	case err := <-result:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (a *AuthorizationManager) post(ctx context.Context, path string, requestBody any, responseBody any) error {
	result := make(chan error, 1)
	go func() {
		reqBodyBytes, err := json.Marshal(requestBody)
		if err != nil {
			a.l.Errorf(ctx, "Error marshalling request body for POST request to path %s: %v", path, err)
			result <- err
			return
		}
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			a.url+path,
			bytes.NewReader(reqBodyBytes),
		)
		if err != nil {
			a.l.Errorf(ctx, "Error creating POST request for path %s: %v", path, err)
			result <- err
			return
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := a.c.Do(req)
		if err != nil {
			a.l.Errorf(ctx, "Error executing POST request for path %s: %v", path, err)
			result <- err
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			a.l.Errorf(ctx, "Unexpected status code for POST request to path %s: %d", path, resp.StatusCode)
			result <- fmt.Errorf("unexpected on POST %s status code: %d", a.url+path, resp.StatusCode)
			return
		}
		if responseBody == nil {
			result <- nil
			return
		}
		result <- json.NewDecoder(resp.Body).Decode(responseBody)
	}()

	select {
	case err := <-result:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
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
