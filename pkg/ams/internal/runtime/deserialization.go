package runtime

import (
	"encoding/json"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

const WildCard = "*"

type AttributeType int

const (
	attrTypeString AttributeType = iota
	attrTypeNumber
	attrTypeBool
	attrTypeArray
)
const (
	AttrTypeString = "String"
	AttrTypeNumber = "Number"
	AttrTypeBool   = "Boolean"
	AttrTypeArray  = "Array"
)

type Runtime struct {
	Resources          map[string]ResourceDefinition `json:"resources,omitempty"`
	Policies           map[string]PolicyDefinition   `json:"policies,omitempty"`
	defaultPolicyNames []string
}

type Attribute struct {
	Type        AttributeType `json:"type,omitempty"`
	ElementType AttributeType `json:"elementType,omitempty"`
}

type Attributes map[string]Attribute

type PolicyDefinition struct {
	Default   bool                      `json:"default,omitempty"`
	AppTID    string                    `json:"appTid,omitempty"`
	Resources map[string]PolicyResource `json:"resources,omitempty"`
}

type PolicyResource struct {
	Conditions ActionsConditions `json:"conditions,omitempty"`
}

type ResourceDefinition struct {
	Attributes Attributes `json:"attributes,omitempty"`
}

type ActionsConditions map[string]expression.Expression

func (r *Runtime) UnmarshalJSON(data []byte) error {
	type Alias Runtime
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	r.defaultPolicyNames = make([]string, 0)

	for policyName, policy := range r.Policies {
		if policy.Default {
			r.defaultPolicyNames = append(r.defaultPolicyNames, policyName)
		}
	}
	return nil
}

func (em *ActionsConditions) UnmarshalJSON(data []byte) error {
	result := make(ActionsConditions)

	var rawActionConditions map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawActionConditions); err != nil {
		return err
	}
	for k, v := range rawActionConditions {
		var err error
		result[k], err = expressionFromJSON(v)
		if err != nil {
			return err
		}
	}
	*em = result
	return nil
}

func (a *AttributeType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case AttrTypeString:
		*a = attrTypeString
	case AttrTypeNumber:
		*a = attrTypeNumber
	case AttrTypeBool:
		*a = attrTypeBool
	case AttrTypeArray:
		*a = attrTypeArray
	default:
		return &json.UnmarshalTypeError{
			Value: "unknown attribute type",
			Type:  nil,
		}
	}
	return nil
}
