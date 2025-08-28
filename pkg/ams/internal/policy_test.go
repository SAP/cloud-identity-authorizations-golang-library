package internal

import (
	_ "embed"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

//go:embed testfiles/policies/defaults.dcn
var defaultsPolicy []byte

//go:embed testfiles/policies/simple.dcn
var simplePolicy []byte

//go:embed testfiles/policies/simple_use.dcn
var simpleUsePolicy []byte

//go:embed testfiles/policies/use_with_restrictions.dcn
var useWithRestrictionsPolicy []byte

//go:embed testfiles/policies/use_without_restrictions.dcn
var useWithoutRestrictionsPolicy []byte

//go:embed testfiles/policies/use_non_existent.dcn
var useNonExistentPolicy []byte

//go:embed testfiles/policies/with_tenant.dcn
var withTenantPolicy []byte

//go:embed testfiles/policies/use_with_broken_restriction.dcn
var useWithBrokenRestrictionPolicy []byte

func TestPolicy(t *testing.T) { //nolint:maintidx
	schema := Schema{
		tenantSchemas: map[string]string{
			"tenant.package.name": "mytenantid",
		},
	}
	t.Run("simple policy from DCN", func(t *testing.T) {
		var sp []dcn.Policy
		err := json.Unmarshal(simplePolicy, &sp)
		if err != nil {
			t.Errorf("Error parsing policy: %v", err)
		}
		p, err := PoliciesFromDCN(sp, schema, nil)
		if err != nil {
			t.Errorf("Error parsing policy: %v", err)
		}
		if len(p.allPolicies) != 1 {
			t.Errorf("Expected 1 policy, got %d", len(p.allPolicies))
		}
		policy, ok := p.allPolicies["simple.name"]
		if !ok {
			t.Errorf("Policy not found")
		}
		if policy.tenant != "" {
			t.Errorf("Expected empty tenant, got %s", policy.tenant)
		}
		if len(policy.rules) != 1 {
			t.Errorf("Expected 1 rule, got %d", len(policy.rules))
		}
	})

	t.Run("simple use policy from DCN", func(t *testing.T) {
		var sp []dcn.Policy
		err := json.Unmarshal(simpleUsePolicy, &sp)
		if err != nil {
			t.Errorf("Error parsing policy: %v", err)
		}
		p, err := PoliciesFromDCN(sp, schema, nil)
		if err != nil {
			t.Errorf("Error parsing policy: %v", err)
		}
		if len(p.allPolicies) != 2 {
			t.Errorf("Expected 2 policies, got %d", len(p.allPolicies))
		}
		policy, ok := p.allPolicies["simple.name"]
		if !ok {
			t.Errorf("Policy not found")
		}
		if policy.tenant != "" {
			t.Errorf("Expected empty tenant, got %s", policy.tenant)
		}
		if len(policy.rules) != 2 {
			t.Errorf("Expected 1 rule, got %d", len(policy.rules))
		}
		policy, ok = p.allPolicies["simple.use"]
		if !ok {
			t.Errorf("Policy not found")
		}
		if policy.tenant != "" {
			t.Errorf("Expected empty tenant, got %s", policy.tenant)
		}
		if len(policy.rules) != 2 {
			t.Errorf("Expected 2 rules, got %d", len(policy.rules))
		}

		r := p.Evaluate(expression.Input{
			"$dcl.resource": expression.String("data"),
			"$dcl.action":   expression.String("read"),
		})
		if r != expression.Bool(true) {
			t.Errorf("Expected true, got %v", r)
		}

		r = p.Evaluate(expression.Input{
			"$dcl.resource": expression.String("data"),
			"$dcl.action":   expression.String("delete"),
		})
		if r != expression.Bool(false) {
			t.Errorf("Expected false, got %v", r)
		}

		p = p.GetSubset([]string{"simple.use"}, "mytenantid", false)

		r = p.Evaluate(expression.Input{
			"$dcl.resource": expression.String("data"),
		})
		expected := expression.Or(
			expression.In(
				expression.Ref("$dcl.action"),
				expression.StringArray{
					"read",
				},
			),
			expression.In(
				expression.Ref("$dcl.action"),
				expression.StringArray{
					"write",
				},
			),
		)

		if !reflect.DeepEqual(r, expected) {
			t.Errorf("Expected %+v, got %+v", expected, r)
		}
	})

	t.Run("use with restrictions policy from DCN", func(t *testing.T) { //nolint:dupl
		var sp []dcn.Policy
		err := json.Unmarshal(useWithRestrictionsPolicy, &sp)
		if err != nil {
			t.Errorf("Error parsing policy: %v", err)
		}
		p, err := PoliciesFromDCN(sp, schema, nil)
		if err != nil {
			t.Errorf("Error parsing policy: %v", err)
		}
		if len(p.allPolicies) != 2 {
			t.Errorf("Expected 2 policies, got %d", len(p.allPolicies))
		}
		policy, ok := p.allPolicies["simple.name"]
		if !ok {
			t.Errorf("Policy not found")
		}
		if policy.tenant != "" {
			t.Errorf("Expected empty tenant, got %s", policy.tenant)
		}
		if len(policy.rules) != 2 {
			t.Errorf("Expected 2 rules, got %d", len(policy.rules))
		}
		policy, ok = p.allPolicies["simple.use"]
		if !ok {
			t.Errorf("Policy not found")
		}
		if policy.tenant != "" {
			t.Errorf("Expected empty tenant, got %s", policy.tenant)
		}
		if len(policy.rules) != 3 {
			t.Errorf("Expected 3 rules, got %d", len(policy.rules))
		}
	})

	t.Run("use non existent policy from DCN", func(t *testing.T) {
		var sp []dcn.Policy
		err := json.Unmarshal(useNonExistentPolicy, &sp)
		if err != nil {
			t.Errorf("Error parsing policy: %v", err)
		}
		_, err = PoliciesFromDCN(sp, schema, nil)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("use without restrictions policy from DCN", func(t *testing.T) { //nolint:dupl
		var sp []dcn.Policy
		err := json.Unmarshal(useWithoutRestrictionsPolicy, &sp)
		if err != nil {
			t.Errorf("Error parsing policy: %v", err)
		}
		p, err := PoliciesFromDCN(sp, schema, nil)
		if err != nil {
			t.Errorf("Error parsing policy: %v", err)
		}
		if len(p.allPolicies) != 2 {
			t.Errorf("Expected 2 policies, got %d", len(p.allPolicies))
		}
		policy, ok := p.allPolicies["simple.name"]
		if !ok {
			t.Errorf("Policy not found")
		}
		if policy.tenant != "" {
			t.Errorf("Expected empty tenant, got %s", policy.tenant)
		}
		if len(policy.rules) != 2 {
			t.Errorf("Expected 2 rules, got %d", len(policy.rules))
		}
		policy, ok = p.allPolicies["simple.use"]
		if !ok {
			t.Errorf("Policy not found")
		}
		if policy.tenant != "" {
			t.Errorf("Expected empty tenant, got %s", policy.tenant)
		}
		if len(policy.rules) != 2 {
			t.Errorf("Expected 2 rules, got %d", len(policy.rules))
		}
	})

	t.Run("with tenant policy from DCN", func(t *testing.T) {
		var sp []dcn.Policy
		err := json.Unmarshal(withTenantPolicy, &sp)
		if err != nil {
			t.Errorf("Error parsing policy: %v", err)
		}
		p, err := PoliciesFromDCN(sp, schema, nil)
		if err != nil {
			t.Errorf("Error parsing policy: %v", err)
		}
		if len(p.allPolicies) != 1 {
			t.Errorf("Expected 1 policy, got %d", len(p.allPolicies))
		}
		policy, ok := p.allPolicies["tenant.package.name.p"]
		if !ok {
			t.Errorf("Policy not found")
		}
		if policy.tenant != "mytenantid" {
			t.Errorf("Expected mytenantid, got %s", policy.tenant)
		}
		if len(policy.rules) != 1 {
			t.Errorf("Expected 1 rule, got %d", len(policy.rules))
		}
	})

	t.Run("2 default policies", func(t *testing.T) {
		var sp []dcn.Policy
		err := json.Unmarshal(defaultsPolicy, &sp)
		if err != nil {
			t.Errorf("Error parsing policy: %v", err)
		}
		p, err := PoliciesFromDCN(sp, schema, nil)
		if err != nil {
			t.Errorf("Error parsing policy: %v", err)
		}
		if len(p.allPolicies) != 3 {
			t.Errorf("Expected 3 policies, got %d", len(p.allPolicies))
		}
		sub := p.GetSubset([]string{"base.simple"}, "mytenantid", true)

		if len(sub.allPolicies) != 3 {
			t.Errorf("Expected 3 policies, got %d", len(sub.allPolicies))
		}

		sub = p.GetSubset([]string{"non-existent"}, "mytenantid", true)

		if len(sub.allPolicies) != 2 {
			t.Errorf("Expected 2 policies, got %d", len(sub.allPolicies))
		}

		sub = p.GetSubset([]string{"non-existent"}, "mytenantid", false)

		if len(sub.allPolicies) != 0 {
			t.Errorf("Expected 2 policies, got %d", len(sub.allPolicies))
		}

		sub = p.GetSubset([]string{"base.simple"}, "non-existent-tenant", true)

		if len(sub.allPolicies) != 2 {
			t.Errorf("Expected 2 policies, got %d", len(sub.allPolicies))
		}

		defaultNames := p.GetDefaultPolicyNames("mytenantid")
		want := []string{"base.default", "tenant.package.name.p"}
		if !reflect.DeepEqual(defaultNames, want) {
			t.Errorf("Expected %v, got %v", want, defaultNames)
		}

		defaultNames = p.GetDefaultPolicyNames("non-existent-tenant")
		want = []string{"base.default"}
		if !reflect.DeepEqual(defaultNames, want) {
			t.Errorf("Expected %v, got %v", want, defaultNames)
		}
	})

	t.Run("use with broken restriction policy from DCN", func(t *testing.T) {
		var sp []dcn.Policy
		err := json.Unmarshal(useWithBrokenRestrictionPolicy, &sp)
		if err != nil {
			t.Errorf("Error parsing policy: %v", err)
		}
		_, err = PoliciesFromDCN(sp, schema, nil)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
	})

	t.Run("get meta info", func(t *testing.T) {
		var sp []dcn.Policy
		err := json.Unmarshal(simplePolicy, &sp)
		if err != nil {
			t.Errorf("Error parsing policy: %v", err)
		}

		ps, err := PoliciesFromDCN(sp, schema, nil)
		if err != nil {
			t.Errorf("Error creating policy set: %v", err)
		}

		t.Run("get resources", func(t *testing.T) {
			expected := []string{"data"}
			if !reflect.DeepEqual(ps.GetResources(), expected) {
				t.Errorf("Expected %v, got %v", expected, ps.GetResources())
			}
		})

		t.Run("get actions", func(t *testing.T) {
			t.Run("returns read for resource data", func(t *testing.T) {
				expected := []string{"read"}
				if !reflect.DeepEqual(ps.GetActions("data"), expected) {
					t.Errorf("Expected %v, got %v", expected, ps.GetResources())
				}
			})

			t.Run("returns empty array for non-existing resource", func(t *testing.T) {
				expected := []string{}
				if !reflect.DeepEqual(ps.GetActions("non-existent"), expected) {
					t.Errorf("Expected %v, got %v", expected, ps.GetResources())
				}
			})

		})
	})

	// t.Run("error on not parseable policy", func(t *testing.T) {

	// 	_, err := PoliciesFromDCN([]byte("not a policy"), schema)
	// 	if err == nil {
	// 		t.Errorf("Expected error, got nil")
	// 	}
	// })
}
