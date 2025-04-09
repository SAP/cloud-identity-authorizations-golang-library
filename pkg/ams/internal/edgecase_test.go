package internal

import (
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
)

func TestEdgeCases(t *testing.T) {
	t.Run("Rule with error in condition", func(t *testing.T) {
		rule := dcn.Rule{
			Condition: &dcn.Expression{
				Call: []string{"and"},
				Args: []dcn.Expression{{Constant: struct{}{}}},
			},
		}
		_, err := RuleFromDCN(rule, nil)
		if err == nil {
			t.Errorf("Expected error")
		}
	})

	t.Run("Policy with error in rule condition", func(t *testing.T) {
		policies := []dcn.Policy{
			{
				QualifiedName: dcn.QualifiedName{"test"},
				Rules: []dcn.Rule{
					{
						Condition: &dcn.Expression{
							Call: []string{"and"},
							Args: []dcn.Expression{{Constant: struct{}{}}},
						},
					},
				},
			},
		}
		_, err := PoliciesFromDCN(policies, Schema{}, nil)
		if err == nil {
			t.Errorf("Expected error")
		}
	})

	t.Run("Policy with error in restriction condition", func(t *testing.T) {
		policies := []dcn.Policy{
			{
				QualifiedName: dcn.QualifiedName{"base"},
				Rules: []dcn.Rule{
					{
						Effect: "grant",
						Condition: &dcn.Expression{
							Call: []string{"is_restricted"},
							Args: []dcn.Expression{{Ref: []string{"x"}}},
						},
					},
				},
			},
			{
				QualifiedName: dcn.QualifiedName{"test"},
				Uses: []dcn.Use{{
					QualifiedPolicyName: dcn.QualifiedName{"base"},
					Restrictions: [][]dcn.Expression{
						{
							{
								Call: []string{"and"},
								Args: []dcn.Expression{{Constant: struct{}{}}},
							},
						},
					}}},
			},
		}
		_, err := PoliciesFromDCN(policies, Schema{}, nil)
		if err == nil {
			t.Errorf("Expected error")
		}
	})
}
