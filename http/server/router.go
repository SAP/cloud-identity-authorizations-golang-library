package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/logging"
	"github.com/sap/cloud-security-client-go/auth"
)

type Router struct {
	am        *ams.AuthorizationManager
	l         logging.Logger
	lastError error
}

func NewRouter(am *ams.AuthorizationManager, l logging.Logger) *Router {
	s := &Router{am: am, l: l}
	return s
}

func (s *Router) Mux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v1/authorize", s.handleAuthorize)
	mux.HandleFunc("POST /v1/resources", s.handleResources)
	mux.HandleFunc("POST /v1/actions", s.handleActions)
	mux.HandleFunc("GET /v1/health", s.handleHealth)
	mux.HandleFunc("GET /v1/policies/default", s.handleDefaultPolicies)
	mux.HandleFunc("GET /v1/policies/default/{tenant_id}", s.handleDefaultPolicies)
	mux.HandleFunc("POST /v1/policies/assigned", s.handleAssignedPolicies)
	mux.HandleFunc("POST /v1/input", s.handleInput)
	return mux
}

func (s *Router) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req AuthorizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// s.l.Info(ctx, "Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	a, err := s.authzForRequest(ctx, req.Token, req.Policies)
	if err != nil {
		// s.l.Info(ctx, "Error authorizing request: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if req.Input == nil {
		req.Input = Input{}
	}
	input := expression.Input(req.Input)
	input["$dcl.action"] = expression.String(req.Action)
	input["$dcl.resource"] = expression.String(req.Resource)

	undefinedFields, wrongTypedFields := s.am.ValidateInput(input)

	result := a.Evaluate(input)
	condition := result.Condition()
	if len(req.NullifyExcept) > 0 {
		keepRefs := make(map[string]bool)
		for _, ref := range req.NullifyExcept {
			keepRefs[ref] = true
		}
		condition = expression.NullifyExcept(condition, keepRefs)
	}
	resp := AuthorizationResponse{Result: expression.ToDCN(condition)}
	for _, field := range undefinedFields {
		resp.Warnings = append(resp.Warnings, fmt.Sprintf("Input field '%s' is not defined in schema", field))
	}
	for _, field := range wrongTypedFields {
		resp.Errors = append(resp.Errors, fmt.Sprintf("Type input field '%s' is incompatible to schema definition", field))
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Router) handleResources(w http.ResponseWriter, r *http.Request) {
	var req ResourcesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	a, err := s.authzForRequest(r.Context(), req.Token, req.Policies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resources := a.GetResources()
	resp := ResourcesResponse{Resources: resources}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Router) handleActions(w http.ResponseWriter, r *http.Request) {
	var req ActionsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	a, err := s.authzForRequest(r.Context(), req.Token, req.Policies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	actions := a.GetActions(req.Resource)
	resp := ActionsResponse{Actions: actions}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Router) handleHealth(w http.ResponseWriter, r *http.Request) {
	if s.am.IsReady() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}
}

func (s *Router) authzForRequest(ctx context.Context, tokenStr string, policies []string) (*ams.Authorizations, error) {
	if tokenStr != "" {
		if policies != nil {
			return nil, fmt.Errorf("cannot specify both policies and token")
		}
		token, err := auth.NewToken(tokenStr)
		if err != nil {
			return nil, fmt.Errorf("error decoding token: %w", err)
		}
		return s.am.AuthorizationsForIdentity(ctx, token), nil
	} else {
		return s.am.AuthorizationsForPolicies(ctx, policies), nil
	}
}

func (s *Router) handleDefaultPolicies(w http.ResponseWriter, r *http.Request) {
	tenant := r.PathValue("tenant_id")
	defaultPolicies := s.am.GetDefaultPolicyNames(tenant)
	resp := DefaultPoliciesResponse{DefaultPolicies: defaultPolicies}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Router) handleAssignedPolicies(w http.ResponseWriter, r *http.Request) {
	req := AssignedPoliciesRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	token, err := auth.NewToken(req.Token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusBadRequest)
		return
	}
	assignedPolicies := s.am.GetAssignments(token.AppTID(), token.ScimID())
	resp := AssignedPoliciesResponse{Policies: assignedPolicies}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (s *Router) handleInput(w http.ResponseWriter, r *http.Request) {
	req := InputRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	input := expression.Input(req.Input)

	undefinedFields, wrongTypedFields := s.am.ValidateInput(input)
	resp := InputResponse{Input: input}
	for _, field := range undefinedFields {
		resp.Warnings = append(resp.Warnings, fmt.Sprintf("Input field '%s' is not defined in schema", field))
	}
	for _, field := range wrongTypedFields {
		resp.Errors = append(resp.Errors, fmt.Sprintf("Type input field '%s' is incompatible to schema definition", field))
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
