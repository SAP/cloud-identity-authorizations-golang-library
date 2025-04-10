package ams

import (
	"reflect"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func TestAuthorizationManager(t *testing.T) { //nolint:maintidx
	t.Run("has schema", func(t *testing.T) {
		dcnChannel := make(chan dcn.DcnContainer)
		assignmentsChannel := make(chan dcn.Assignments)
		am := NewAuthorizationManager(dcnChannel, assignmentsChannel)

		dcnChannel <- dcn.DcnContainer{
			Policies: []dcn.Policy{},
			Schemas: []dcn.Schema{
				{
					QualifiedName: []string{"pkg", "schema1"},
					Tenant:        "tenant1",
				},
			},
			Functions: []dcn.Function{},
		}
		assignmentsChannel <- dcn.Assignments{}

		<-am.WhenReady()

		tenant := am.GetSchema().GetTenantForQualifiedName([]string{"pkg", "p1"})
		if tenant != "tenant1" {
			t.Errorf("expected tenant1, got %v", tenant)
		}
	})
	t.Run("is ready after receiving DCN", func(t *testing.T) {
		dcnChannel := make(chan dcn.DcnContainer)
		assignmentsChannel := make(chan dcn.Assignments)
		am := NewAuthorizationManager(dcnChannel, assignmentsChannel)
		assignmentsChannel <- dcn.Assignments{}

		if am.IsReady() {
			t.Error("is ready before receiving DCN")
		}
		dcnChannel <- dcn.DcnContainer{
			Policies:  []dcn.Policy{},
			Schemas:   []dcn.Schema{},
			Functions: []dcn.Function{},
		}

		<-am.WhenReady()

		// update again
		dcnChannel <- dcn.DcnContainer{
			Policies:  []dcn.Policy{},
			Schemas:   []dcn.Schema{},
			Functions: []dcn.Function{},
		}

		if !am.IsReady() {
			t.Error("is not ready after receiving DCN")
		}

		// update again
		dcnChannel <- dcn.DcnContainer{
			Policies:  []dcn.Policy{},
			Schemas:   []dcn.Schema{},
			Functions: []dcn.Function{},
		}

		if !am.IsReady() {
			t.Error("is not ready after receiving DCN")
		}
	})

	t.Run("error in functions", func(t *testing.T) {
		dcnChannel := make(chan dcn.DcnContainer)
		assignmentsChannel := make(chan dcn.Assignments)
		am := NewAuthorizationManager(dcnChannel, assignmentsChannel)
		assignmentsChannel <- dcn.Assignments{}

		errors := []error{}

		done := make(chan struct{})

		am.RegisterErrorHandler(func(err error) {
			errors = append(errors, err)
			done <- struct{}{}
		})

		if len(errors) != 0 {
			t.Error("errors before receiving DCN")
		}
		dcnChannel <- dcn.DcnContainer{
			Policies: []dcn.Policy{},
			Schemas:  []dcn.Schema{},
			Functions: []dcn.Function{
				{
					QualifiedName: []string{"func1"},
					Result: dcn.Expression{
						Call: []string{"func2"},
					},
				},
			},
		}
		<-done
		if len(errors) != 1 {
			t.Errorf("expected 1 error, got %v", errors)
		}
	})

	t.Run("error in policies", func(t *testing.T) {
		dcnChannel := make(chan dcn.DcnContainer)
		assignmentsChannel := make(chan dcn.Assignments)
		am := NewAuthorizationManager(dcnChannel, assignmentsChannel)
		assignmentsChannel <- dcn.Assignments{}

		errors := []error{}
		done := make(chan struct{})

		am.RegisterErrorHandler(func(err error) {
			errors = append(errors, err)
			done <- struct{}{}
		})

		if len(errors) != 0 {
			t.Error("errors before receiving DCN")
		}
		dcnChannel <- dcn.DcnContainer{
			Policies: []dcn.Policy{
				{
					QualifiedName: []string{"policy1"},
					Rules: []dcn.Rule{
						{
							Condition: &dcn.Expression{
								Call: []string{"func1"},
							},
						},
					},
				},
			},
			Schemas:   []dcn.Schema{},
			Functions: []dcn.Function{},
		}

		<-done
		if len(errors) != 1 {
			t.Errorf("expected 1 error, got %v", errors)
		}
	})

	t.Run("get Authorizations", func(t *testing.T) {
		dcnChannel := make(chan dcn.DcnContainer)
		assignmentsChannel := make(chan dcn.Assignments)
		am := NewAuthorizationManager(dcnChannel, assignmentsChannel)
		assignmentsChannel <- dcn.Assignments{}

		dcnChannel <- dcn.DcnContainer{
			Policies: []dcn.Policy{
				{
					QualifiedName: []string{"pkg", "policy1"},
					Rules: []dcn.Rule{
						{
							Actions:   []string{"action1"},
							Resources: []string{"resource1"},
						},
					},
				},
				{
					QualifiedName: []string{"pkg", "policy2"},
					Rules: []dcn.Rule{
						{
							Actions:   []string{"action2"},
							Resources: []string{"resource2"},
						},
					},
				},
				{
					QualifiedName: []string{"pkg", "policy3"},
					Rules: []dcn.Rule{
						{
							Actions:   []string{"action3"},
							Resources: []string{"resource2"},
						},
					},
				},
			},
			Schemas: []dcn.Schema{
				{
					QualifiedName: []string{"pkg", "schema1"},
					Tenant:        "tenant1",
				},
			},
			Functions: []dcn.Function{},
		}

		<-am.WhenReady()

		auths := am.GetAuthorizations([]string{"pkg.policy1"}, "tenant1", false)

		r := auths.Evaluate(expression.Input{
			"$dcl.resource": expression.String("resource1"),
			"$dcl.action":   expression.String("action1"),
		})
		if r != expression.TRUE {
			t.Errorf("expected true, got %v", r)
		}
		r = auths.Evaluate(expression.Input{
			"$dcl.resource": expression.String("resource2"),
			"$dcl.action":   expression.String("action2"),
		})
		if r != expression.FALSE {
			t.Errorf("expected false, got %v", r)
		}

		auth2 := am.GetAuthorizations([]string{"pkg.policy2"}, "tenant1", false)

		r = auth2.Evaluate(expression.Input{
			"$dcl.resource": expression.String("resource1"),
			"$dcl.action":   expression.String("action1"),
		})
		if r != expression.FALSE {
			t.Errorf("expected false, got %v", r)
		}
		r = auth2.Evaluate(expression.Input{
			"$dcl.resource": expression.String("resource2"),
			"$dcl.action":   expression.String("action2"),
		})
		if r != expression.TRUE {
			t.Errorf("expected true, got %v", r)
		}

		andJoined := auths.AndJoin(auth2)

		r = andJoined.Evaluate(expression.Input{
			"$dcl.resource": expression.String("resource1"),
			"$dcl.action":   expression.String("action1"),
		})
		if r != expression.FALSE {
			t.Errorf("expected false, got %v", r)
		}
		r = andJoined.Evaluate(expression.Input{
			"$dcl.resource": expression.String("resource2"),
			"$dcl.action":   expression.String("action2"),
		})
		if r != expression.FALSE {
			t.Errorf("expected false, got %v", r)
		}

		r = andJoined.Evaluate(expression.Input{
			"$dcl.resource": expression.String("resource2"),
			"$dcl.action":   expression.UNKNOWN,
		})
		if r != expression.FALSE {
			t.Errorf("expected false, got %v", r)
		}

		auth3 := am.GetAuthorizations([]string{"pkg.policy3"}, "tenant1", false)

		andJoined = auth2.AndJoin(auth3)
		r = andJoined.Evaluate(expression.Input{
			"$dcl.resource": expression.String("resource2"),
			"$dcl.action":   expression.UNKNOWN,
		})

		in1 := expression.In{
			Args: []expression.Expression{
				expression.Reference{Name: "$dcl.action"},
				expression.StringArray{"action2"},
			},
		}
		in2 := expression.In{
			Args: []expression.Expression{
				expression.Reference{Name: "$dcl.action"},
				expression.StringArray{"action3"},
			},
		}

		expected := expression.NewAnd(in1, in2)
		if !reflect.DeepEqual(r, expected) {
			t.Errorf("expected %+v, got %+v", expected, r)
		}
	})

	t.Run("get assignments", func(t *testing.T) {
		dcnChannel := make(chan dcn.DcnContainer)
		assignmentsChannel := make(chan dcn.Assignments)
		am := NewAuthorizationManager(dcnChannel, assignmentsChannel)

		dcnChannel <- dcn.DcnContainer{
			Policies:  []dcn.Policy{},
			Schemas:   []dcn.Schema{},
			Functions: []dcn.Function{},
		}

		if am.IsReady() {
			t.Error("is ready before receiving DCN")
		}

		assignmentsChannel <- dcn.Assignments{
			"tenant1": dcn.UserAssignments{
				"user1": []string{"pkg.policy1"},
			},
		}

		<-am.WhenReady()

		r := am.GetAssignments("tenant1", "user1")
		expected := []string{"pkg.policy1"}
		if !reflect.DeepEqual(r, expected) {
			t.Errorf("expected %v, got %v", expected, r)
		}
		r = am.GetAssignments("tenant1", "user2")
		expected = []string{}
		if !reflect.DeepEqual(r, expected) {
			t.Errorf("expected %v, got %v", expected, r)
		}
		r = am.GetAssignments("tenant2", "user1")
		expected = []string{}
		if !reflect.DeepEqual(r, expected) {
			t.Errorf("expected %v, got %v", expected, r)
		}
	})
}
