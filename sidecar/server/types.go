package server

import (
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
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
