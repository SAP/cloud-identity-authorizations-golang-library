package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-security-client-go/auth"
)

type Router struct {
	am        *ams.AuthorizationManager
	lastError error
}

const (
	PATH_AUTHORIZE = "/v1/authorize"
	PATH_RESOURCES = "/v1/resources"
	PATH_ACTIONS   = "/v1/actions"
	PATH_HEALTH    = "/v1/health"
)

func NewRouter(am *ams.AuthorizationManager) *Router {
	s := &Router{am: am}
	am.RegisterErrorHandler(func(err error) {
		// Store the last error encountered by the AuthorizationManager
		s.lastError = err
	})
	return s
}

func (s *Router) Mux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(PATH_AUTHORIZE, s.handleAuthorize)
	mux.HandleFunc(PATH_RESOURCES, s.handleResources)
	mux.HandleFunc(PATH_ACTIONS, s.handleActions)
	mux.HandleFunc(PATH_HEALTH, s.handleHealth)
	return mux
}

func (s *Router) handleAuthorize(w http.ResponseWriter, r *http.Request) {
	var req AuthorizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	a, err := s.authzForRequest(req.Token, req.Policies)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Env != nil {
		a.SetEnvInput(req.Env)
	}
	result := a.Inquire(req.Action, req.Resource, req.Input)
	condition := result.Condition()
	if len(req.NullifyExcept) > 0 {
		keepRefs := make(map[string]bool)
		for _, ref := range req.NullifyExcept {
			keepRefs[ref] = true
		}
		condition = expression.NullifyExcept(result.Condition(), keepRefs)
	}
	resp := AuthorizationResponse{Result: expression.ToDCN(condition)}
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
	a, err := s.authzForRequest(req.Token, req.Policies)
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
	a, err := s.authzForRequest(req.Token, req.Policies)
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

func (s *Router) authzForRequest(tokenStr string, policies []string) (*ams.Authorizations, error) {
	if tokenStr != "" {
		if policies != nil {
			return nil, fmt.Errorf("Cannot specify both policies and token")
		}
		token, err := auth.NewToken(tokenStr)
		if err != nil {
			return nil, fmt.Errorf("Error decoding token: %w", err)
		}
		return s.am.AuthorizationsForIdentity(token), nil
	} else {
		return s.am.AuthorizationsForPolicies(policies), nil
	}
}
