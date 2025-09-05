package testserver

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/util"
)

type Router struct {
	am *ams.AuthorizationManager
}

type DCLInput struct {
	Action   string `json:"action,omitempty"`
	Resource string `json:"resource,omitempty"`
}

type Input struct {
	DCL DCLInput               `json:"$dcl"`
	App map[string]interface{} `json:"$app"`
	Env map[string]interface{} `json:"$env"`
}

type PolicyEvaluationRequest struct {
	Policies []string `json:"policies"`
	Input    Input    `json:"input"`
}

type ScopedPolicyEvaluationRequest struct {
	Policies [][]string `json:"policies"`
	Input    Input      `json:"input"`
}

type NullifyExceptRequest struct {
	Expression dcn.Expression `json:"expression"`
	KeepRefs   [][]string     `json:"keep_refs"`
}

type DefaultPoliciesRequest struct {
	Tenant string `json:"tenant"`
}
type DefaultPoliciesResponse struct {
	Policies []string `json:"policies"`
}

type EvaluationResponse struct {
	Expression dcn.Expression `json:"expression"`
}

const (
	PATH_LOAD_DCN                 = "/v1/load_dcn"
	PATH_EVALUATE_POLICIES        = "/v1/evaluate_policies"
	PATH_EVALUATE_POLICIES_SCOPED = "/v1/evaluate_policies_scoped"
	PATH_NULLIFY_EXCEPT           = "/v1/nullify_except"
	PATH_GET_DEFAULT_POLICIES     = "/v1/get_default_policies"
)

func (s *Router) Mux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc(PATH_LOAD_DCN, s.handleLoadDCN)
	mux.HandleFunc(PATH_EVALUATE_POLICIES, s.handleEvaluatePolicies)
	mux.HandleFunc(PATH_EVALUATE_POLICIES_SCOPED, s.handleEvaluatePoliciesScoped)
	mux.HandleFunc(PATH_NULLIFY_EXCEPT, s.handleNullifyExcept)
	mux.HandleFunc(PATH_GET_DEFAULT_POLICIES, s.handleGetDefaultPolicies)
	return mux
}

func (s *Router) handleLoadDCN(w http.ResponseWriter, r *http.Request) {
	var rb dcn.DcnContainer
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse request body %v", err), http.StatusBadRequest)
		return
	}

	dcnChannel := make(chan dcn.DcnContainer, 1)
	assignmentsChannel := make(chan dcn.Assignments, 1)

	s.am = ams.NewAuthorizationManager(dcnChannel, assignmentsChannel, func(err error) {
		log.Printf("error in authorization manager: %v\n", err)
	})
	done := false
	s.am.RegisterErrorHandler(func(err error) {
		log.Printf("error in authorization manager: %v\n", err)
		if !done {
			http.Error(w, fmt.Sprintf("error in authorization manager: %v", err), http.StatusInternalServerError)
		}
	})
	dcnChannel <- rb
	assignmentsChannel <- dcn.Assignments{}

	<-s.am.WhenReady()
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(`"OK"`))
	if err != nil {
		log.Printf("could not write response %v\n", err)
	}
	done = true
}

func (s *Router) handleEvaluatePolicies(w http.ResponseWriter, r *http.Request) {
	if s.am == nil {
		http.Error(w, "Authorization manager not initialized", http.StatusInternalServerError)
		return
	}
	if !s.am.IsReady() {
		http.Error(w, "Authorization manager not ready", http.StatusServiceUnavailable)
		return
	}
	var rb PolicyEvaluationRequest
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse request body %v", err), http.StatusBadRequest)
		return
	}
	auth := s.am.AuthorizationsForPolicies(rb.Policies)
	input := s.am.GetSchema().CustomInput(
		rb.Input.DCL.Action,
		rb.Input.DCL.Resource,
		rb.Input.App,
		rb.Input.Env,
	)
	result := auth.Evaluate(input)

	response := EvaluationResponse{
		Expression: expression.ToDCN(result.Condition()),
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode response %v", err), http.StatusInternalServerError)
		return
	}
}

func (s *Router) handleEvaluatePoliciesScoped(w http.ResponseWriter, r *http.Request) {
	if s.am == nil {
		http.Error(w, "Authorization manager not initialized", http.StatusInternalServerError)
		return
	}
	if !s.am.IsReady() {
		http.Error(w, "Authorization manager not ready", http.StatusServiceUnavailable)
		return
	}
	var rb ScopedPolicyEvaluationRequest
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse request body %v", err), http.StatusBadRequest)
		return
	}

	if len(rb.Policies) == 0 {
		http.Error(w, "no policies provided", http.StatusBadRequest)
		return
	}
	auth := s.am.AuthorizationsForPolicies(rb.Policies[0])
	for _, policies := range rb.Policies[1:] {
		auth = auth.AndJoin(s.am.AuthorizationsForPolicies(policies))
	}

	input := s.am.GetSchema().CustomInput(
		rb.Input.DCL.Action,
		rb.Input.DCL.Resource,
		rb.Input.App,
		rb.Input.Env,
	)

	result := auth.Evaluate(input)

	response := EvaluationResponse{
		Expression: expression.ToDCN(result.Condition()),
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode response %v", err), http.StatusInternalServerError)
		return
	}
}

// func (s *Router) handleAuthorizeUser(w http.ResponseWriter, r *http.Request) {
// 	if s.am == nil {
// 		http.Error(w, "Authorization manager not initialized", http.StatusInternalServerError)
// 		return
// 	}
// 	if !s.am.IsReady() {
// 		http.Error(w, "Authorization manager not ready", http.StatusServiceUnavailable)
// 		return
// 	}
// 	var rb UserEvaluationRequest
// 	err := json.NewDecoder(r.Body).Decode(&rb)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("could not parse request body %v", err), http.StatusBadRequest)
// 		return
// 	}
// 	auth := s.am.UserAuthorizations(rb.Tenant, rb.User)
// 	input := s.am.GetSchema().CustomInput(
// 		rb.Input.DCL.Action,
// 		rb.Input.DCL.Resource,
// 		rb.Input.App,
// 		rb.Input.Env,
// 	)
// 	result := auth.Evaluate(input)

// 	response := EvaluationResponse{
// 		Expression: expression.ToDCN(result),
// 	}
// 	w.WriteHeader(http.StatusOK)
// 	w.Header().Set("Content-Type", "application/json")
// 	err = json.NewEncoder(w).Encode(response)
// 	if err != nil {
// 		http.Error(w, fmt.Sprintf("could not encode response %v", err), http.StatusInternalServerError)
// 		return
// 	}
// }

func (s *Router) handleNullifyExcept(w http.ResponseWriter, r *http.Request) {
	if s.am == nil {
		http.Error(w, "Authorization manager not initialized", http.StatusInternalServerError)
		return
	}
	if !s.am.IsReady() {
		http.Error(w, "Authorization manager not ready", http.StatusServiceUnavailable)
		return
	}
	var rb NullifyExceptRequest
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse request body %v", err), http.StatusBadRequest)
		return
	}
	e, err := expression.FromDCN(rb.Expression, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse expression %v", err), http.StatusBadRequest)
		return
	}
	keepRefs := make(map[string]bool)
	for _, ref := range rb.KeepRefs {
		keepRefs[util.StringifyQualifiedName(ref)] = true
	}
	result := expression.NullifyExcept(e.Expression, keepRefs)

	response := EvaluationResponse{
		Expression: expression.ToDCN(result),
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode response %v", err), http.StatusInternalServerError)
		return
	}
}

func (s *Router) handleGetDefaultPolicies(w http.ResponseWriter, r *http.Request) {
	if s.am == nil {
		http.Error(w, "Authorization manager not initialized", http.StatusInternalServerError)
		return
	}
	if !s.am.IsReady() {
		http.Error(w, "Authorization manager not ready", http.StatusServiceUnavailable)
		return
	}

	var rb DefaultPoliciesRequest
	err := json.NewDecoder(r.Body).Decode(&rb)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not parse request body %v", err), http.StatusBadRequest)
		return
	}

	policies := DefaultPoliciesResponse{
		Policies: s.am.GetDefaultPolicyNames(rb.Tenant),
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(policies)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not encode response %v", err), http.StatusInternalServerError)
		return
	}
}
