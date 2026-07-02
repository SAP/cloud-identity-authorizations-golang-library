package runtime

import (
	_ "embed"
	"reflect"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

//go:embed testdata/example.json
var evaluateTestJSON []byte

func TestEvaluate(t *testing.T) {
	r := &Runtime{}
	err := r.UnmarshalJSON(evaluateTestJSON)
	if err != nil {
		t.Fatalf("Failed to unmarshal evaluate test JSON: %v", err)
	}
	t.Run("get resources of policies", func(t *testing.T) {
		policyNames := []string{"my.package.EmployeePolicy"}
		resources := r.GetResources(policyNames)
		want := []string{"my.package.Employee"}
		if !reflect.DeepEqual(resources, want) {
			t.Errorf("GetResourcesOfPolicies() = %v, want %v", resources, want)
		}
		policyNames = []string{"my.package.EmployeePolicy", "my.package.GlobalReadPolicy", "non.existent.policy"}
		resources = r.GetResources(policyNames)
		want = []string{"my.package.Employee", "*"}
		if !compareStringSlices(resources, want) {
			t.Errorf("GetResourcesOfPolicies() = %v, want %v", resources, want)
		}
	})

	t.Run("get actions of policies for resource", func(t *testing.T) {
		policyNames := []string{"my.package.EmployeePolicy", "my.package.GlobalReadPolicy", "non.existent.policy"}
		resource := "my.package.Employee"
		actions := r.GetActions(policyNames, resource)
		want := []string{"read", "write"}
		if !compareStringSlices(actions, want) {
			t.Errorf("GetActionsOfPoliciesForResource() = %v, want %v", actions, want)
		}
		resource = "anonymous.resource"
		actions = r.GetActions(policyNames, resource)
		want = []string{"read"}
		if !compareStringSlices(actions, want) {
			t.Errorf("GetActionsOfPoliciesForResource() = %v, want %v", actions, want)
		}
	})

	t.Run("evaluate policies for action and resource", func(t *testing.T) {
		policyNames := []string{"my.package.EmployeePolicy", "non.existent.policy"}
		input := expression.Input{
			"$env.$user.user_uuid":           expression.String("id1"),
			"$resource.department.head.uuid": expression.String("id2"),
		}
		result := r.Evaluate(policyNames, "write", "my.package.Employee", input)
		want := (expression.Expression)(expression.FALSE)
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Evaluate() = %v, want %v", result, want)
		}

		result = r.Evaluate(policyNames, "read", "my.package.Employee", input)
		want = expression.Ref("$resource.is_public")
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Evaluate() = %v, want %v", result, want)
		}
	})

	t.Run("evaluate policies with wildcards", func(t *testing.T) {
		policyNames := []string{"my.package.AnyActionPolicy", "my.package.GlobalReadPolicy"}
		input := expression.Input{
			"$resource.department.name": expression.String("Engineering"),
		}
		result := r.Evaluate(policyNames, "read", "anonymous.resource", input)
		want := expression.TRUE
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Evaluate() = %v, want %v", result, want)
		}

		result = r.Evaluate(policyNames, "anonymous.action", "my.package.Employee", input)
		want = expression.TRUE
		if !reflect.DeepEqual(result, want) {
			t.Errorf("Evaluate() = %v, want %v", result, want)
		}
	})

}
