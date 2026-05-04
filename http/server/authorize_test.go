package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	e "github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

type errorWriter struct {
	header http.Header
	status int
}

func (e *errorWriter) Header() http.Header {
	if e.header == nil {
		e.header = make(http.Header)
	}
	return e.header
}

func (e *errorWriter) Write(b []byte) (int, error) {
	return 0, errors.New("write error")
}

func (e *errorWriter) WriteHeader(statusCode int) {
	e.status = statusCode
}
func TestAuthzForPolicies(t *testing.T) {
	am := ams.NewAuthorizationManagerForFs("../../pkg/ams/test/scenarios/simple", nil)
	<-am.WhenReady()
	r := NewRouter(am, nopLogger{})
	var req AuthorizationRequest

	t.Run("Evaluate single policy without input", func(t *testing.T) {
		rr := httptest.NewRecorder()

		req = AuthorizationRequest{
			Policies: []string{"base.PublicStuff"},
			Action:   "read",
			Resource: "r1",
			Input: Input{
				"$env.$user.groups": e.StringArray{"g1", "g2"},
			},
		}
		r.Mux().ServeHTTP(rr, newAuthorizationRequest(req))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		got, err := expressionFromBody(rr.Body)
		if err != nil {
			t.Errorf("Error decoding expression from body: %v", err)
		}

		want := "or({$app.entity1.public}, in({$app.entity1.group}, [\"g1\" \"g2\"]))"
		if got.String() != want {
			t.Errorf("Expected expression %s, got %s", want, got.String())
		}
	})

	t.Run("Evaluate single policy with input", func(t *testing.T) {
		rr := httptest.NewRecorder()

		i := Input{
			"$env.$user.groups":  e.StringArray{"g1", "g2"},
			"$app.entity1.group": e.String("g1"),
		}
		req = AuthorizationRequest{
			Policies: []string{"base.PublicStuff"},
			Action:   "read",
			Resource: "r1",
			Input:    i,
		}
		r.Mux().ServeHTTP(rr, newAuthorizationRequest(req))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		got, err := expressionFromBody(rr.Body)
		if err != nil {
			t.Errorf("Error decoding expression from body: %v", err)
		}

		want := e.TRUE
		if got != want {
			t.Errorf("Expected expression %s, got %s", want.String(), got.String())
		}
	})

	t.Run("Evaluate single policy with empty array in input ", func(t *testing.T) {
		rr := httptest.NewRecorder()

		i := Input{
			"$env.$user.groups":   e.NumberArray{},
			"$app.entity1.group":  e.String("g1"),
			"$app.entity1.public": e.FALSE,
		}
		req = AuthorizationRequest{
			Policies: []string{"base.PublicStuff"},
			Action:   "read",
			Resource: "r1",
			Input:    i,
		}
		r.Mux().ServeHTTP(rr, newAuthorizationRequest(req))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		got, err := expressionFromBody(rr.Body)
		if err != nil {
			t.Errorf("Error decoding expression from body: %v", err)
		}

		want := e.FALSE
		if got != want {
			t.Errorf("Expected expression %s, got %s", want.String(), got.String())
		}
	})
	t.Run("Evaluate single policy nullify", func(t *testing.T) {
		rr := httptest.NewRecorder()

		req = AuthorizationRequest{
			Policies: []string{"base.PublicStuff"},
			Action:   "read",
			Resource: "r1",
			Input: Input{
				"$env.$user.groups": e.StringArray{"g1", "g2"},
			},
			NullifyExcept: []string{"$app.entity1.group"},
		}
		r.Mux().ServeHTTP(rr, newAuthorizationRequest(req))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		got, err := expressionFromBody(rr.Body)
		if err != nil {
			t.Errorf("Error decoding expression from body: %v", err)
		}
		want := "in({$app.entity1.group}, [\"g1\" \"g2\"])"
		if got.String() != want {
			t.Errorf("Expected expression %s, got %s", want, got.String())
		}
	})

	t.Run("Edge-case: no req body", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r.Mux().ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/v1/authorize", nil))
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})
}

type tokenClaim map[string]any

func TestEvaluateToken(t *testing.T) {
	am := ams.NewAuthorizationManagerForFs("../../pkg/ams/test/scenarios/simple", nil)
	<-am.WhenReady()
	r := NewRouter(am, nopLogger{})
	t.Run("Evaluate valid token", func(t *testing.T) {
		rr := httptest.NewRecorder()
		claims := tokenClaim{
			"groups":  []string{"g1", "g2"},
			"scim_id": "user1",
			"app_tid": "tenant1",
			"email":   "user@example.com",
		}

		req := AuthorizationRequest{
			Action:   "read",
			Resource: "r1",
			Token:    newToken(claims),
			Input: Input{
				"$app.entity1.group": e.String("g3"),
			},
		}
		r.Mux().ServeHTTP(rr, newAuthorizationRequest(req))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		got, err := expressionFromBody(rr.Body)
		if err != nil {
			t.Errorf("Error decoding expression from body: %v", err)
		}

		want := "{$app.entity1.public}"
		if got.Evaluate(e.Input{"$app.entity1.size": e.Number(99)}).String() != want {
			t.Errorf("Expected expression %s, got %s", want, got.String())
		}
		want = "ge({$app.entity1.size}, 100)"
		if got.Evaluate(e.Input{"$app.entity1.public": e.FALSE}).String() != want {
			t.Errorf("Expected expression %s, got %s", want, got.String())
		}
	})

	t.Run("Evaluate token with nullify", func(t *testing.T) {
		rr := httptest.NewRecorder()
		claims := tokenClaim{
			"groups":  []string{"g1", "g2"},
			"scim_id": "user1",
			"app_tid": "tenant1",
			"email":   "user@example.com",
		}

		req := AuthorizationRequest{
			Action:   "read",
			Resource: "r1",
			Token:    newToken(claims),
			Input: Input{
				"$app.entity1.group": e.String("g3"),
			},
			NullifyExcept: []string{"$app.entity1.size"},
		}
		r.Mux().ServeHTTP(rr, newAuthorizationRequest(req))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		got, err := expressionFromBody(rr.Body)
		if err != nil {
			t.Errorf("Error decoding expression from body: %v", err)
		}

		want := "ge({$app.entity1.size}, 100)"
		if got.String() != want {
			t.Errorf("Expected expression %s, got %s", want, got.String())
		}
	})

	t.Run("Edgecase: invalid token", func(t *testing.T) {
		rr := httptest.NewRecorder()
		req := AuthorizationRequest{
			Action:   "read",
			Resource: "r1",
			Token:    "invalidtoken",
		}
		r.Mux().ServeHTTP(rr, newAuthorizationRequest(req))
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("Edge-case: token with policies specified", func(t *testing.T) {
		rr := httptest.NewRecorder()
		claims := tokenClaim{
			"groups":  []string{"g1", "g2"},
			"scim_id": "user1",
			"app_tid": "tenant1",
			"email":   "user@example.com",
		}

		req := AuthorizationRequest{
			Action:   "read",
			Resource: "r1",
			Token:    newToken(claims),
			Policies: []string{"policy1"},
		}
		r.Mux().ServeHTTP(rr, newAuthorizationRequest(req))
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})

	t.Run("Edge-case: can't write response", func(t *testing.T) {
		rr := &errorWriter{}
		claims := tokenClaim{
			"groups":  []string{"g1", "g2"},
			"scim_id": "user1",
			"app_tid": "tenant1",
			"email":   "user@example.com",
		}
		req := AuthorizationRequest{
			Action:   "read",
			Resource: "r1",
			Token:    newToken(claims),
		}
		r.Mux().ServeHTTP(rr, newAuthorizationRequest(req))
		if rr.status != http.StatusInternalServerError {
			t.Errorf("Expected status 500, got %d", rr.status)
		}
	})
}

func TestHealth(t *testing.T) {
	t.Run("Health check ready when Authorization Manager is initialized", func(t *testing.T) {
		dcnChan := make(chan dcn.DcnContainer)
		assigmentChan := make(chan dcn.Assignments)
		am := ams.NewAuthorizationManager(dcnChan, assigmentChan, nil)
		r := NewRouter(am, nopLogger{})
		rr := httptest.NewRecorder()
		r.Mux().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/health", nil))
		if rr.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status 503, got %d", rr.Code)
		}
		dcnChan <- dcn.DcnContainer{}
		assigmentChan <- dcn.Assignments{}
		<-am.WhenReady()
		rr = httptest.NewRecorder()
		r.Mux().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/health", nil))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		// stays healthy after dcn update error
		dcnChan <- dcn.DcnContainer{
			Policies: []dcn.Policy{
				{
					QualifiedName: dcn.QualifiedName{"pkg", "p1"},
					Rules: []dcn.Rule{
						{
							Condition: &dcn.Expression{
								Call: []string{},
							},
						},
					},
				},
			},
		}
		rr = httptest.NewRecorder()
		r.Mux().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/v1/health", nil))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})
}

func newToken(claims tokenClaim) string {
	jwtRaw, err := json.Marshal(claims)
	if err != nil {
		panic(err)
	}
	return string(jwtRaw)
}

func expressionFromBody(body io.Reader) (e.Expression, error) {
	var resp AuthorizationResponse
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return nil, err
	}
	exprContainer, err := e.FromDCN(resp.Result, nil)
	if err != nil {
		return nil, err
	}
	return exprContainer.Expression, nil
}

func newAuthorizationRequest(body AuthorizationRequest) *http.Request {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}
	n := bytes.NewReader(bodyBytes)

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, "/v1/authorize", n)
	return req
}
