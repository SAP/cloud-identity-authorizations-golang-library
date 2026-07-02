package runtime

import (
	_ "embed"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

//go:embed testdata/example.json
var exampleJSON []byte

func TestExampleJSON(t *testing.T) {
	var r Runtime
	err := json.Unmarshal(exampleJSON, &r)
	if err != nil {
		t.Fatalf("Failed to unmarshal example JSON: %v", err)
	}
	t.Run("has default policy", func(t *testing.T) {
		want := []string{"my.package.EmployeeDefault"}
		got := r.GetDefaultPolicyNames("")
		if !reflect.DeepEqual(got, want) {
			t.Errorf("GetDefaultPolicyNames() = %v, want %v", got, want)
		}
	})

	t.Run("has conditions", func(t *testing.T) {
		p, ok := r.Policies["my.package.GlobalReadPolicy"]
		if !ok {
			t.Fatalf("Policy my.package.GlobalReadPolicy not found")
		}
		res, ok := p.Resources["*"]
		if !ok {
			t.Fatalf("Resource * not found in policy my.package.GlobalReadPolicy")
		}
		action, ok := res.Conditions["read"]
		if !ok {
			t.Fatalf("Condition read not found in resource * of policy my.package.GlobalReadPolicy")
		}
		want := expression.TRUE
		if !reflect.DeepEqual(action, want) {
			t.Errorf("Condition read in resource * of policy my.package.GlobalReadPolicy = %v, want %v", action, want)
		}
	})
}
