package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"slices"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
	e "github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

type nopLogger struct{}

func (n nopLogger) Errorf(ctx context.Context, msg string, args ...interface{}) {}
func (n nopLogger) Warnf(ctx context.Context, msg string, args ...interface{})  {}
func (n nopLogger) Infof(ctx context.Context, msg string, args ...interface{})  {}
func (n nopLogger) Debugf(ctx context.Context, msg string, args ...interface{}) {}

func TestRouter(t *testing.T) {
	am := ams.NewAuthorizationManagerForFs("../../pkg/ams/test/scenarios/simple", nil)
	<-am.WhenReady()
	r := NewRouter(am, nopLogger{})

	t.Run("get resources for token", func(t *testing.T) {
		rr := httptest.NewRecorder()
		claims := tokenClaim{
			"groups":  []string{"g1", "g2"},
			"scim_id": "user1",
			"app_tid": "tenant1",
			"email":   "user@example.com",
		}
		req := ResourcesRequest{
			Token: newToken(claims),
		}

		r.HandleResources(rr, newResourcesRequest(req))

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", rr.Code)
		}
		var res ResourcesResponse
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}
		expectedResources := []string{"r1", "r2"}
		if len(res.Resources) != 2 {
			t.Errorf("Expected 2 resources, got %d", len(res.Resources))
		}
		if !slices.Contains(res.Resources, "r1") || !slices.Contains(res.Resources, "r2") {
			t.Errorf("Expected resources %v, got %v", expectedResources, res.Resources)
		}
	})

	t.Run("get actions for token and resource r2", func(t *testing.T) {
		rr := httptest.NewRecorder()
		claims := tokenClaim{
			"groups":  []string{"g1", "g2"},
			"scim_id": "user1",
			"app_tid": "tenant1",
			"email":   "user@example.com",
		}
		req := ActionsRequest{
			Token:    newToken(claims),
			Resource: "r2",
		}

		r.HandleActions(rr, newActionRequest(req))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", rr.Code)
		}
		var res ActionsResponse
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}
		expectedActions := []string{"read"}
		if !reflect.DeepEqual(res.Actions, expectedActions) {
			t.Errorf("Expected actions %v, got %v", expectedActions, res.Actions)
		}
	})
}

func TestInputEndpoint(t *testing.T) {
	am := ams.NewAuthorizationManagerForFs("../../pkg/ams/test/scenarios/simple", nil)
	<-am.WhenReady()
	r := NewRouter(am, nopLogger{})

	t.Run("removes undefined input fields with warnings", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := InputRequest{
			Input: Input{
				"$app.entity1.group": e.String("g3"),
				"undefined_field_1":  e.String("value"),
				"undefined_field_2":  e.String("value"),
			},
		}
		r.HandleInput(rr, newInputRequest(req))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", rr.Code)
		}
		var res InputResponse
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}
		expectedInput := Input{
			"$app.entity1.group": e.String("g3"),
		}
		if !reflect.DeepEqual(res.Input, expectedInput) {
			t.Errorf("Expected input %v, got %v", expectedInput, res.Input)
		}
		if len(res.Warnings) != 2 {
			t.Errorf("Expected 2 warnings, got %d", len(res.Warnings))
		}
		if len(res.Errors) != 0 {
			t.Errorf("Expected 0 errors, got %d", len(res.Errors))
		}
	})

	t.Run("reomves wrongly typed input fields with errors", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := InputRequest{
			Input: Input{
				"$app.entity1.group":  e.Number(1),
				"$app.entity1.public": e.Bool(true),
				"$app.entity1.size":   e.String("large"),
				"$app.entity1.name":   e.Bool(true),
			},
		}
		r.HandleInput(rr, newInputRequest(req))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status code 200, got %d", rr.Code)
		}
		var res InputResponse
		if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}
		expectedInput := Input{
			"$app.entity1.public": e.Bool(true),
		}
		if !reflect.DeepEqual(res.Input, expectedInput) {
			t.Errorf("Expected input %v, got %v", expectedInput, res.Input)
		}
		if len(res.Warnings) != 0 {
			t.Errorf("Expected 0 warnings, got %d", len(res.Warnings))
		}
		if len(res.Errors) != 3 {
			t.Errorf("Expected 3 errors, got %d", len(res.Errors))
		}
	})
}

func newInputRequest(body InputRequest) *http.Request {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	n := bytes.NewReader(bodyBytes)

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/input", n)
	return req
}

func newActionRequest(body ActionsRequest) *http.Request {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	n := bytes.NewReader(bodyBytes)

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/actions", n)
	return req
}

func newResourcesRequest(body ResourcesRequest) *http.Request {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	n := bytes.NewReader(bodyBytes)

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/resources", n)
	return req
}
