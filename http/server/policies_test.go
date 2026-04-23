package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
)

func TestPolicyEndpoints(t *testing.T) {
	am := ams.NewAuthorizationManagerForFs("../../pkg/ams/test/scenarios/simple", func(err error) {
		t.Errorf("Error in AuthorizationManager: %v", err)
		t.Fail()
	})
	<-am.WhenReady()
	r := NewRouter(am)
	t.Run("Get assigned policies of user1", func(t *testing.T) {
		req := AssignedPoliciesRequest{
			Token: newToken(tokenClaim{
				"scim_id": "user1",
				"app_tid": "tenant1",
			}),
		}
		rr := httptest.NewRecorder()
		r.Mux().ServeHTTP(rr, newAssignedPoliciesRequest(req))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		var resp AssignedPoliciesResponse
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Errorf("Error decoding response body: %v", err)
		}
		want := []string{"_dcltenant_._tenant1.R1BigSized"}

		if !reflect.DeepEqual(resp.Policies, want) {
			t.Errorf("Expected policies %v, got %v", want, resp.Policies)
		}
	})

	t.Run("Get default policies for tenant1", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r.Mux().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/policies/default/tenant1", nil))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		var resp DefaultPoliciesResponse
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Errorf("Error decoding response body: %v", err)
		}
		want := []string{"base.PublicStuff"}

		if !reflect.DeepEqual(resp.DefaultPolicies, want) {
			t.Errorf("Expected default policies %v, got %v", want, resp.DefaultPolicies)
		}
	})

	t.Run("Get not-tenant specific default policies", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r.Mux().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/policies/default", nil))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		var resp DefaultPoliciesResponse
		if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
			t.Errorf("Error decoding response body: %v", err)
		}
		want := []string{"base.PublicStuff"}

		if !reflect.DeepEqual(resp.DefaultPolicies, want) {
			t.Errorf("Expected default policies %v, got %v", want, resp.DefaultPolicies)
		}
	})
}

func newAssignedPoliciesRequest(req AssignedPoliciesRequest) *http.Request {
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	n := bytes.NewReader(bodyBytes)
	r := httptest.NewRequest(http.MethodPost, "/v1/policies/assigned", n)
	r.Header.Set("Content-Type", "application/json")
	return r
}
