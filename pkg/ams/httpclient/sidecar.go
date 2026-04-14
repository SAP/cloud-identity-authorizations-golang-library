package httpclient

import (
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

const (
	PATH_AUTHORIZE             = "/v1/authorize"
	PATH_AUTHORIZE_UNVALIDATED = "/v1/authorize_unvalidated"
	PATH_RESOURCES             = "/v1/resources"
	PATH_ACTIONS               = "/v1/actions"
	PATH_HEALTH                = "/v1/health"
	PATH_DEFAULT_POLICIES      = "/v1/policies/default"
	PATH_ASSIGNED_POLICIES     = "/v1/policies/assigned"
	PATH_CREATE_INPUT          = "/v1/input"
)

type AuthorizationRequest struct {
	Action        string         `json:"action"`
	Resource      string         `json:"resource"`
	Policies      []string       `json:"policies,omitempty"`
	Token         string         `json:"token,omitempty"`
	Input         map[string]any `json:"input"`
	Env           map[string]any `json:"env,omitempty"`
	NullifyExcept []string       `json:"nullify_except,omitempty"`
}
type UnvalidatedAuthorizationRequest struct {
	Policies      []string         `json:"policies,omitempty"`
	Token         string           `json:"token,omitempty"`
	Input         expression.Input `json:"input,omitempty"`
	NullifyExcept []string         `json:"nullify_except,omitempty"`
}

type AuthorizationResponse struct {
	Result dcn.Expression `json:"result"`
}

type ResourcesRequest struct {
	Policies []string `json:"policies,omitempty"`
	Token    string   `json:"token,omitempty"`
}

type ResourcesResponse struct {
	Resources []string `json:"resources"`
}

type ActionsRequest struct {
	Policies []string `json:"policies,omitempty"`
	Token    string   `json:"token,omitempty"`
	Resource string   `json:"resource"`
}

type ActionsResponse struct {
	Actions []string `json:"actions"`
}

type DefaultPoliciesResponse struct {
	DefaultPolicies []string `json:"default_policies"`
}

type AssignedPoliciesRequest struct {
	Token string `json:"token,omitempty"`
}

type AssignedPoliciesResponse struct {
	Policies []string `json:"policies"`
}
type CreateInputRequest struct {
	Action   string         `json:"action"`
	Resource string         `json:"resource"`
	Input    map[string]any `json:"input"`
	Env      map[string]any `json:"env,omitempty"`
}

type CreateInputResponse struct {
	RawInput expression.Input `json:"raw_input"`
}
