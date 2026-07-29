package ams

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

type errorHandler struct {
	errors         []error
	errorsReceived chan bool
}

func (l *errorHandler) Callback(err error) {
	l.errors = append(l.errors, err)
	l.errorsReceived <- true
}

func createErrorHandler() *errorHandler {
	return &errorHandler{
		errors:         []error{},
		errorsReceived: make(chan bool),
	}
}

type TestIdentity struct {
	appTID   string
	scimID   string
	userUUID string
	groups   []string
	email    string
}

func (i TestIdentity) AppTID() string {
	return i.appTID
}
func (i TestIdentity) ScimID() string {
	return i.scimID
}
func (i TestIdentity) UserUUID() string {
	return i.userUUID
}
func (i TestIdentity) Groups() []string {
	return i.groups
}
func (i TestIdentity) Email() string {
	return i.email
}

func TestAuthorizationManager(t *testing.T) { //nolint:maintidx
	t.Run("is ready after receiving DCN", func(t *testing.T) {
		dcnChannel := make(chan dcn.DcnContainer, 1)
		assignmentsChannel := make(chan dcn.Assignments, 1)
		am := NewAuthorizationManager(dcnChannel, assignmentsChannel, nil)
		assignmentsChannel <- dcn.Assignments{}

		if am.isReady() {
			t.Error("is ready before receiving DCN")
		}
		dcnChannel <- dcn.DcnContainer{
			Policies:  []dcn.Policy{},
			Schemas:   []dcn.Schema{},
			Functions: []dcn.Function{},
		}

		err := am.Run(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// update again
		dcnChannel <- dcn.DcnContainer{
			Policies:  []dcn.Policy{},
			Schemas:   []dcn.Schema{},
			Functions: []dcn.Function{},
		}

		if !am.isReady() {
			t.Error("is not ready after receiving DCN")
		}

		// update again
		dcnChannel <- dcn.DcnContainer{
			Policies:  []dcn.Policy{},
			Schemas:   []dcn.Schema{},
			Functions: []dcn.Function{},
		}

		if !am.isReady() {
			t.Error("is not ready after receiving DCN")
		}
	})

	t.Run("with functions", func(t *testing.T) {
		dcnChannel := make(chan dcn.DcnContainer, 1)
		assignmentsChannel := make(chan dcn.Assignments, 1)
		am := NewAuthorizationManager(dcnChannel, assignmentsChannel, nil)
		assignmentsChannel <- dcn.Assignments{}
		dcnChannel <- dcn.DcnContainer{
			Policies: []dcn.Policy{
				{
					QualifiedName: []string{"pkg", "policy1"},
					Rules: []dcn.Rule{
						{
							Condition: &dcn.Expression{
								Call: []string{"pkg", "func1"},
							},
						},
					},
					Default: true,
				},
			},
			Schemas: []dcn.Schema{},
			Functions: []dcn.Function{
				{
					QualifiedName: []string{"pkg", "func1"},
					Result: dcn.Expression{
						Ref: []string{"x"},
					},
				},
			},
		}
		err := am.Run(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		a := am.AuthorizationsForPolicies([]string{"pkg.policy1"})
		got := a.Evaluate("action1", "resource1", expression.Input{}).Condition()
		want := expression.Ref("x")
		if !reflect.DeepEqual(got, want) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("error in functions", func(t *testing.T) {
		dcnChannel := make(chan dcn.DcnContainer, 1)
		assignmentsChannel := make(chan dcn.Assignments, 1)
		ml := createErrorHandler()

		am := NewAuthorizationManager(dcnChannel, assignmentsChannel, ml.Callback)
		go func() {
			err := am.Run(context.Background())
			if err != nil {
				panic(fmt.Sprintf("unexpected error: %v", err))
			}
		}()
		assignmentsChannel <- dcn.Assignments{}

		if len(ml.errors) != 0 {
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
		<-ml.errorsReceived
		if len(ml.errors) != 1 {
			t.Errorf("expected 1 error, got %v", ml.errors)
		}
	})

	t.Run("error in policies", func(t *testing.T) {
		dcnChannel := make(chan dcn.DcnContainer, 1)
		assignmentsChannel := make(chan dcn.Assignments, 1)
		ml := createErrorHandler()
		am := NewAuthorizationManager(dcnChannel, assignmentsChannel, ml.Callback)
		go func() {
			err := am.Run(context.Background())
			if err != nil {
				panic(fmt.Sprintf("unexpected error: %v", err))
			}
		}()
		assignmentsChannel <- dcn.Assignments{}

		if len(ml.errors) != 0 {
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

		<-ml.errorsReceived
		if len(ml.errors) != 1 {
			t.Errorf("expected 1 error, got %v", ml.errors)
		}
	})

	t.Run("get Authorizations", func(t *testing.T) {
		dcnChannel := make(chan dcn.DcnContainer, 1)
		assignmentsChannel := make(chan dcn.Assignments, 1)
		am := NewAuthorizationManager(dcnChannel, assignmentsChannel, nil)
		assignmentsChannel <- dcn.Assignments{}

		dcnChannel <- dcn.DcnContainer{
			Policies: []dcn.Policy{
				{
					QualifiedName: []string{"pkg", "policy1"},
					Rules: []dcn.Rule{
						{
							Actions:   []string{"action1"},
							Resources: []string{"resource1"},
							Condition: &dcn.Expression{
								Call: []string{"eq"},
								Args: []dcn.Expression{
									{Ref: []string{"field1"}},
									{Constant: "value1"},
								},
							},
						},
					},
				},
				{
					QualifiedName: []string{"pkg", "policy2"},
					Rules: []dcn.Rule{
						{
							Actions:   []string{"action2"},
							Resources: []string{"resource2"},
							Condition: &dcn.Expression{
								Call: []string{"eq"},
								Args: []dcn.Expression{
									{Ref: []string{"field2"}},
									{Constant: "value2"},
								},
							},
						},
					},
				},
				{
					QualifiedName: []string{"pkg", "policy3"},
					Rules: []dcn.Rule{
						{
							Actions:   []string{"action3", "action2"},
							Resources: []string{"resource2"},
							Condition: &dcn.Expression{
								Call: []string{"eq"},
								Args: []dcn.Expression{
									{Ref: []string{"field3"}},
									{Constant: "value3"},
								},
							},
						},
					},
				},
			},
			Schemas: []dcn.Schema{
				{
					QualifiedName: []string{"pkg", "tenant", "schema1"},
					Tenant:        "tenant1",
				},
			},
			Functions: []dcn.Function{},
		}

		err := am.Run(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		auths := am.AuthorizationsForPolicies([]string{"pkg.policy1"})

		r := auths.Evaluate("action1", "resource1", expression.Input{"field1": expression.String("value1")})
		if !r.IsGranted() {
			t.Errorf("expected true, got %v", r)
		}
		r = auths.Evaluate("action2", "resource2", expression.Input{})
		if !r.IsDenied() {
			t.Errorf("expected false, got %v", r)
		}

		auth2 := am.AuthorizationsForPolicies([]string{"pkg.policy2"})

		r = auth2.Evaluate("action1", "resource1", expression.Input{})
		if !r.IsDenied() {
			t.Errorf("expected false, got %v", r)
		}
		r = auth2.Evaluate("action2", "resource2", expression.Input{"field2": expression.String("value2")})
		if !r.IsGranted() {
			t.Errorf("expected true, got %v", r)
		}

		andJoined := auths.AndJoin(auth2)

		r = andJoined.Evaluate("action1", "resource1", expression.Input{})
		if !r.IsDenied() {
			t.Errorf("expected false, got %v", r)
		}
		r = andJoined.Evaluate("action2", "resource2", expression.Input{})
		if !r.IsDenied() {
			t.Errorf("expected false, got %v", r)
		}

		r = andJoined.Evaluate("action2", "resource2", expression.Input{})
		if !r.IsDenied() {
			t.Errorf("expected false, got %v", r)
		}

		auth3 := am.AuthorizationsForPolicies([]string{"pkg.policy3"})

		andJoined = auth2.AndJoin(auth3)
		r = andJoined.Evaluate("action2", "resource2", expression.Input{})

		expected := expression.And(
			expression.Eq(
				expression.Ref("field2"),
				expression.String("value2"),
			),
			expression.Eq(
				expression.Ref("field3"),
				expression.String("value3"),
			),
		)
		if !reflect.DeepEqual(r.Condition(), expected) {
			t.Errorf("expected %+v, got %+v", expected, r.condition)
		}
	})

	t.Run("get assignments", func(t *testing.T) {
		dcnChannel := make(chan dcn.DcnContainer, 1)
		assignmentsChannel := make(chan dcn.Assignments, 1)
		am := NewAuthorizationManager(dcnChannel, assignmentsChannel, nil)

		dcnChannel <- dcn.DcnContainer{
			Policies: []dcn.Policy{
				{
					QualifiedName: []string{"pkg", "policy1"},
					Rules: []dcn.Rule{
						{
							Resources: []string{"resource1"},
						},
					},
				},
			},
			Schemas:   []dcn.Schema{},
			Functions: []dcn.Function{},
		}

		if am.isReady() {
			t.Error("is ready before receiving DCN")
		}

		assignmentsChannel <- dcn.Assignments{
			"tenant1": dcn.UserAssignments{
				"user1": []string{"pkg.policy1"},
			},
		}

		err := am.Run(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

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

	t.Run("Authorizations for identity with user attribues", func(t *testing.T) {
		dcnChannel := make(chan dcn.DcnContainer, 1)
		assignmentsChannel := make(chan dcn.Assignments, 1)
		am := NewAuthorizationManager(dcnChannel, assignmentsChannel, nil)

		dcnChannel <- dcn.DcnContainer{
			Policies: []dcn.Policy{
				{
					QualifiedName: []string{"pkg", "policy1"},
					Rules: []dcn.Rule{
						{
							Resources: []string{"r1"},
							Condition: &dcn.Expression{
								Call: []string{"eq"},
								Args: []dcn.Expression{
									{Ref: []string{"$env", "$user", "email"}},
									{Constant: "user1@example.com"},
								},
							},
						},
					},
				},
			},
			Schemas:   []dcn.Schema{},
			Functions: []dcn.Function{},
		}
		assignmentsChannel <- dcn.Assignments{
			"tenant1": dcn.UserAssignments{
				"user1": []string{"pkg.policy1"},
			},
		}

		err := am.Run(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		authz := am.AuthorizationsForIdentity(TestIdentity{
			email:    "user1@example.com",
			appTID:   "tenant1",
			scimID:   "user1",
			userUUID: "user1",
			groups:   []string{"group1"},
		})

		r := authz.Inquire("*", "r1", nil)
		if !r.IsGranted() {
			t.Errorf("expected true, got %v", r.condition)
		}
	})

	t.Run("get default policy names", func(t *testing.T) {
		dcnChannel := make(chan dcn.DcnContainer, 1)
		assignmentsChannel := make(chan dcn.Assignments, 1)
		am := NewAuthorizationManager(dcnChannel, assignmentsChannel, nil)

		dcnChannel <- dcn.DcnContainer{
			Policies: []dcn.Policy{
				{
					QualifiedName: []string{"pkg", "policy1"},
					Rules:         []dcn.Rule{},
					Default:       true,
				},
				{
					QualifiedName: []string{"tenant1", "policy2"},
					Rules:         []dcn.Rule{},
					Default:       true,
				},
			},
			Schemas: []dcn.Schema{
				{QualifiedName: []string{"tenant1", "schema1"}, Tenant: "tenant1"},
			},
			Functions: []dcn.Function{},
		}

		if am.isReady() {
			t.Error("is ready before receiving DCN")
		}

		assignmentsChannel <- dcn.Assignments{}

		err := am.Run(context.Background())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		r := am.GetDefaultPolicyNames("")
		expected := []string{"pkg.policy1"}
		if !reflect.DeepEqual(r, expected) {
			t.Errorf("expected %v, got %v", expected, r)
		}

		r = am.GetDefaultPolicyNames("tenant1")
		expected = []string{"pkg.policy1", "tenant1.policy2"}
		if !reflect.DeepEqual(r, expected) {
			t.Errorf("expected %v, got %v", expected, r)
		}
	})

	t.Run("error on load dcn", func(t *testing.T) {
		dcnChannel := make(chan dcn.DcnContainer, 1)
		assignmentsChannel := make(chan dcn.Assignments, 1)
		ml := createErrorHandler()
		am := NewAuthorizationManager(dcnChannel, assignmentsChannel, ml.Callback)
		go func() {
			err := am.Run(context.Background())
			if err != nil {
				panic(fmt.Sprintf("unexpected error: %v", err))
			}
		}()

		assignmentsChannel <- dcn.Assignments{}
		dcnChannel <- dcn.DcnContainer{
			Policies: []dcn.Policy{
				{
					QualifiedName: []string{"pkg", "policy1"},
					Rules: []dcn.Rule{
						{
							Condition: &dcn.Expression{
								Call: []string{"invalid"},
							},
						},
					},
				},
			},
		}

		<-ml.errorsReceived
		if len(ml.errors) != 1 {
			t.Errorf("expected 1 error, got %d", len(ml.errors))
		}
	})
}
