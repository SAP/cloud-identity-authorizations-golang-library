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

	"github.com/lestrrat-go/jwx/jwt"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
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
	am := ams.NewAuthorizationManagerForFs("../../pkg/ams/test/scenarios/simple", func(err error) {
		t.Errorf("Error in AuthorizationManager: %v", err)
		t.Fail()
	})
	<-am.WhenReady()
	r := NewRouter(am)
	var req AuthorizationRequest

	t.Run("Evaluate single policy without input", func(t *testing.T) {
		rr := httptest.NewRecorder()
		env := map[string]any{
			"$user": map[string]any{
				"groups": []string{"g1", "g2"},
			},
		}
		req = AuthorizationRequest{
			Policies: []string{"base.PublicStuff"},
			Action:   "read",
			Resource: "r1",
			Env:      env,
		}
		r.Mux().ServeHTTP(rr, newAuthorizationRequest(req))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		got, err := expressionFromBody(rr.Body)
		if err != nil {
			t.Errorf("Error decoding expression from body: %v", err)
		}

		want := "or(Ref($app.entity1.public), in(Ref($app.entity1.group), [g1 g2]))"
		if got.String() != want {
			t.Errorf("Expected expression %s, got %s", want, got.String())
		}
	})

	t.Run("Evaluate single policy with input", func(t *testing.T) {
		rr := httptest.NewRecorder()
		env := map[string]any{
			"$user": map[string]any{
				"groups": []string{"g1", "g2"},
			},
		}
		input := map[string]any{
			"entity1": map[string]any{
				"group": "g1",
			},
		}
		req = AuthorizationRequest{
			Policies: []string{"base.PublicStuff"},
			Action:   "read",
			Resource: "r1",
			Env:      env,
			Input:    input,
		}
		r.Mux().ServeHTTP(rr, newAuthorizationRequest(req))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		got, err := expressionFromBody(rr.Body)
		if err != nil {
			t.Errorf("Error decoding expression from body: %v", err)
		}

		want := expression.TRUE
		if got != want {
			t.Errorf("Expected expression %s, got %s", want.String(), got.String())
		}
	})
	t.Run("Evaluate single policy nullify", func(t *testing.T) {
		rr := httptest.NewRecorder()
		env := map[string]any{
			"$user": map[string]any{
				"groups": []string{"g1", "g2"},
			},
		}
		req = AuthorizationRequest{
			Policies:      []string{"base.PublicStuff"},
			Action:        "read",
			Resource:      "r1",
			Env:           env,
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
		want := "in(Ref($app.entity1.group), [g1 g2])"
		if got.String() != want {
			t.Errorf("Expected expression %s, got %s", want, got.String())
		}
	})

	t.Run("Edge-case: no req body", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r.handleAuthorize(rr, httptest.NewRequest(http.MethodPost, "/v1/evaluate_policies", nil))
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})
}

type tokenClaim map[string]any

func TestEvaluateToken(t *testing.T) {
	am := ams.NewAuthorizationManagerForFs("../../pkg/ams/test/scenarios/simple", func(err error) {
		t.Errorf("Error in AuthorizationManager: %v", err)
		t.Fail()
	})
	<-am.WhenReady()
	r := NewRouter(am)
	t.Run("Evaluate valid token", func(t *testing.T) {
		rr := httptest.NewRecorder()
		claims := tokenClaim{
			"groups":  []string{"g1", "g2"},
			"scim_id": "user1",
			"app_tid": "tenant1",
			"email":   "user@example.com",
		}
		input := map[string]any{
			"entity1": map[string]any{
				"group": "g3",
			},
		}

		req := AuthorizationRequest{
			Action:   "read",
			Resource: "r1",
			Token:    newToken(claims),
			Input:    input,
		}
		r.Mux().ServeHTTP(rr, newAuthorizationRequest(req))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		got, err := expressionFromBody(rr.Body)
		if err != nil {
			t.Errorf("Error decoding expression from body: %v", err)
		}

		want := "Ref($app.entity1.public)"
		if got.Evaluate(expression.Input{"$app.entity1.size": expression.Number(99)}).String() != want {
			t.Errorf("Expected expression %s, got %s", want, got.String())
		}
		want = "ge(Ref($app.entity1.size), 100)"
		if got.Evaluate(expression.Input{"$app.entity1.public": expression.FALSE}).String() != want {
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
		input := map[string]any{
			"entity1": map[string]any{
				"group": "g3",
			},
		}

		req := AuthorizationRequest{
			Action:        "read",
			Resource:      "r1",
			Token:         newToken(claims),
			Input:         input,
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

		want := "ge(Ref($app.entity1.size), 100)"
		if got.String() != want {
			t.Errorf("Expected expression %s, got %s", want, got.String())
		}
	})

	t.Run("Evaluate token override env", func(t *testing.T) {
		rr := httptest.NewRecorder()
		claims := tokenClaim{
			"groups":  []string{"g1", "g2"},
			"scim_id": "user1",
			"app_tid": "tenant1",
			"email":   "user@example.com",
		}
		env := map[string]any{
			"$user": map[string]any{
				"groups": []string{"g3"},
			},
		}

		input := map[string]any{
			"entity1": map[string]any{
				"group": "g3",
			},
		}

		req := AuthorizationRequest{
			Action:   "read",
			Resource: "r1",
			Token:    newToken(claims),
			Input:    input,
			Env:      env,
		}
		r.Mux().ServeHTTP(rr, newAuthorizationRequest(req))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		got, err := expressionFromBody(rr.Body)
		if err != nil {
			t.Errorf("Error decoding expression from body: %v", err)
		}
		want := expression.TRUE
		if got != want {
			t.Errorf("Expected expression %s, got %s", want.String(), got.String())
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

	t.Run("Edge-case: no req body", func(t *testing.T) {
		rr := httptest.NewRecorder()
		r.handleAuthorize(rr, httptest.NewRequest(http.MethodPost, "/v1/evaluate_token", nil))
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
		r := NewRouter(am)
		rr := httptest.NewRecorder()
		r.Mux().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, PATH_HEALTH, nil))
		if rr.Code != http.StatusServiceUnavailable {
			t.Errorf("Expected status 503, got %d", rr.Code)
		}
		dcnChan <- dcn.DcnContainer{}
		assigmentChan <- dcn.Assignments{}
		<-am.WhenReady()
		rr = httptest.NewRecorder()
		r.Mux().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, PATH_HEALTH, nil))
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
		r.Mux().ServeHTTP(rr, httptest.NewRequest(http.MethodGet, PATH_HEALTH, nil))
		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
	})
}

func newToken(claims tokenClaim) string {
	token := jwt.New()
	for k, v := range claims {
		err := token.Set(k, v)
		if err != nil {
			panic(err)
		}
	}
	jwtRaw, err := jwt.NewSerializer().Serialize(token)
	if err != nil {
		panic(err)
	}
	return string(jwtRaw)
}

func expressionFromBody(body io.Reader) (expression.Expression, error) {
	var resp AuthorizationResponse
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return nil, err
	}
	exprContainer, err := expression.FromDCN(resp.Result, nil)
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

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost, PATH_AUTHORIZE, n)
	return req
}
