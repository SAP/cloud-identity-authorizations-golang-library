package testserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	x "github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func TestScenarioAllowAction(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("Could not get current filename")
	}
	dcnPath := filepath.Join(filepath.Dir(filename), "../../pkg/ams/test/scenarios/test_001")

	absPath, err := filepath.Abs(dcnPath)
	if err != nil {
		t.Fatalf("Failed to get absolute path: %v", err)
	}

	router := Router{}
	mux := router.Mux()

	testserver := httptest.NewServer(router.Mux())
	defer testserver.Close()

	post := func(t *testing.T, path string, body any) *httptest.ResponseRecorder {
		result := httptest.NewRecorder()

		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("Failed to marshal request body: %v", err)
		}
		req := httptest.NewRequest("POST", testserver.URL+path, bytes.NewReader(data))
		req.Header.Set("Content-Type", "application/json")
		mux.ServeHTTP(result, req)
		return result
	}

	loader := dcn.NewLocalLoader(absPath, nil)
	dcn := <-loader.DCNChannel
	<-loader.AssignmentsChannel

	resp := post(t, PATH_LOAD_DCN, dcn)
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.Code)
	}

	t.Run("evaluate policies", func(t *testing.T) {
		resp := post(t, PATH_EVALUATE_POLICIES, PolicyEvaluationRequest{
			Policies: []string{"cas.StarOnSingleResource", "cas.SingleActionOnStar"},
		})
		if resp.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.Code)
		}
		var result EvaluationResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}

		got, err := x.FromDCN(result.Expression, nil)
		if err != nil {
			t.Fatalf("Failed to convert response expression: %v", err)
		}
		want := x.Or(
			x.In(
				x.Ref("$dcl.resource"),
				x.StringArray{x.String("Resource1")},
			),
			x.In(
				x.Ref("$dcl.action"),
				x.StringArray{x.String("Action1")},
			),
		)

		if !reflect.DeepEqual(got.Expression, want) {
			t.Errorf("Expected expression %v, got %v", want, got)
		}

		resp = post(t, PATH_NULLIFY_EXCEPT, NullifyExceptRequest{
			Expression: result.Expression,
			KeepRefs: [][]string{
				{"$dcl", "action"},
			}})
		if resp.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.Code)
		}
		var nullifyResult EvaluationResponse
		if err := json.NewDecoder(resp.Body).Decode(&nullifyResult); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}
		nullified, err := x.FromDCN(nullifyResult.Expression, nil)
		if err != nil {
			t.Fatalf("Failed to convert nullified expression: %v", err)
		}
		expectedNullified := x.Or(
			x.In(
				x.Ref("$dcl.action"),
				x.StringArray{x.String("Action1")},
			),
		)
		if !reflect.DeepEqual(nullified.Expression, expectedNullified) {
			t.Errorf("Expected nullified expression %v, got %v", expectedNullified, nullified)
		}
	})

	t.Run("get default policies", func(t *testing.T) {
		resp := post(t, PATH_GET_DEFAULT_POLICIES, DefaultPoliciesRequest{
			Tenant: "tenant1",
		})
		if resp.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, resp.Code)
		}
		var result DefaultPoliciesResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			t.Fatalf("Failed to decode response body: %v", err)
		}
		expectedPolicies := []string{}
		if !reflect.DeepEqual(result.Policies, expectedPolicies) {
			t.Errorf("Expected policies %v, got %v", expectedPolicies, result.Policies)
		}
	})

}
