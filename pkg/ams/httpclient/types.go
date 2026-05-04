package httpclient

import (
	"encoding/json"
	"fmt"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

const (
	PATH_AUTHORIZE         = "/v1/authorize"
	PATH_RESOURCES         = "/v1/resources"
	PATH_ACTIONS           = "/v1/actions"
	PATH_HEALTH            = "/v1/health"
	PATH_DEFAULT_POLICIES  = "/v1/policies/default"
	PATH_ASSIGNED_POLICIES = "/v1/policies/assigned"
	PATH_CREATE_INPUT      = "/v1/input"
)

type reqInput map[string]any

type Input map[string]expression.Constant

type AuthorizationRequest struct {
	Action        string   `json:"action"`
	Resource      string   `json:"resource"`
	Policies      []string `json:"policies,omitempty"`
	Token         string   `json:"token,omitempty"`
	Input         reqInput `json:"input,omitempty"`
	NullifyExcept []string `json:"nullify_except,omitempty"`
}

type AuthorizationResponse struct {
	Result   dcn.Expression `json:"result"`
	Errors   []string       `json:"errors,omitempty"`
	Warnings []string       `json:"warnings,omitempty"`
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

type InputRequest struct {
	Action   string   `json:"action"`
	Resource string   `json:"resource"`
	Input    reqInput `json:"input"`
}

type InputResponse struct {
	Input    Input    `json:"input"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

func (i *Input) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if *i == nil {
		*i = make(Input, len(raw))
	}

	for k, v := range raw {
		switch val := v.(type) {
		case string:
			(*i)[k] = expression.String(val)
		case bool:
			(*i)[k] = expression.Bool(val)
		case float64:
			(*i)[k] = expression.Number(val)
		case []interface{}:
			if len(val) == 0 {
				(*i)[k] = expression.EmptyArray{}
				continue
			}
			switch val[0].(type) {
			case string:
				strs := make([]string, len(val))
				for i, s := range val {
					var ok bool
					strs[i], ok = s.(string)
					if !ok {
						return fmt.Errorf("first element of array field %s has type string but field[%d] has type %T", k, i, s)
					}
				}
				(*i)[k] = expression.ArrayFrom(strs)
			case bool:
				bools := make([]bool, len(val))
				for i, b := range val {
					var ok bool
					bools[i], ok = b.(bool)
					if !ok {
						return fmt.Errorf("first element of array field %s has type boolean but field[%d] has type %T", k, i, b)
					}
				}
				(*i)[k] = expression.ArrayFrom(bools)
			case float64:
				nums := make([]float64, len(val))
				for i, n := range val {
					var ok bool
					nums[i], ok = n.(float64)
					if !ok {
						return fmt.Errorf("first element of array field %s has type number but field[%d] has type %T", k, i, n)
					}
				}
				(*i)[k] = expression.ArrayFrom(nums)
			default:
				return fmt.Errorf("unable to interpret elements of array field %s as string, number or boolean", k)
			}
		default:
			return fmt.Errorf("unable to interpret value of %s as string, number or boolean or array thereof", k)
		}
	}
	return nil
}
