package ams

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

type concurrencyIdentity struct {
	appTID   string
	scimID   string
	userUUID string
	groups   []string
	email    string
}

func (i concurrencyIdentity) AppTID() string {
	return i.appTID
}

func (i concurrencyIdentity) ScimID() string {
	return i.scimID
}

func (i concurrencyIdentity) UserUUID() string {
	return i.userUUID
}

func (i concurrencyIdentity) Groups() []string {
	return i.groups
}

func (i concurrencyIdentity) Email() string {
	return i.email
}

type concurrencyToken struct {
	appTID string
	scimID string
	claims map[string]interface{}
}

func (t concurrencyToken) AppTID() string {
	return t.appTID
}

func (t concurrencyToken) ScimID() string {
	return t.scimID
}

func (t concurrencyToken) GetAllClaimsAsMap() map[string]interface{} {
	return t.claims
}

func makeValidDCN(version int) dcn.DcnContainer {
	return dcn.DcnContainer{
		Version: version,
		Policies: []dcn.Policy{
			{
				QualifiedName: []string{"pkg", fmt.Sprintf("policy_%d", version)},
				Default:       true,
				Rules: []dcn.Rule{
					{
						Actions:   []string{"read", "write"},
						Resources: []string{"resource1", "resource2"},
					},
				},
			},
		},
		Schemas: []dcn.Schema{
			{
				QualifiedName: []string{"pkg", "tenant", "schema1"},
				Tenant:        "tenant1",
				Definition: dcn.SchemaAttribute{
					Nested: map[string]dcn.SchemaAttribute{
						"$env": {
							Nested: map[string]dcn.SchemaAttribute{
								"$user": {
									Nested: map[string]dcn.SchemaAttribute{
										"email":  {Type: "String"},
										"groups": {Type: "String[]"},
									},
								},
							},
						},
						"$app": {
							Nested: map[string]dcn.SchemaAttribute{
								"field": {Type: "String"},
							},
						},
					},
				},
			},
		},
		Functions: []dcn.Function{
			{
				QualifiedName: []string{"pkg", "alwaysTrue"},
				Result: dcn.Expression{
					Constant: true,
				},
			},
		},
		Tests: []dcn.Test{},
	}
}

func makeBrokenFunctionDCN(version int) dcn.DcnContainer {
	return dcn.DcnContainer{
		Version:  version,
		Policies: []dcn.Policy{},
		Functions: []dcn.Function{{
			QualifiedName: []string{fmt.Sprintf("func_%d", version)},
			Result:        dcn.Expression{Call: []string{"func_missing"}}}},
		Schemas: []dcn.Schema{},
		Tests:   []dcn.Test{},
	}
}

func TestAuthorizationManagerConcurrency_ReadsWithStreamingUpdates(t *testing.T) {
	dcnChannel := make(chan dcn.DcnContainer, 128)
	assignmentsChannel := make(chan dcn.Assignments, 128)
	am := NewAuthorizationManager(dcnChannel, assignmentsChannel, nil)

	initialAssignments := dcn.Assignments{
		"tenant1": dcn.UserAssignments{
			"user1": []string{"pkg.policy_0"},
		},
	}
	assignmentsChannel <- initialAssignments
	dcnChannel <- makeValidDCN(0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := am.Run(ctx); err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}

	const readerCount = 8
	const iterationsPerReader = 200

	var readersWG sync.WaitGroup
	readersWG.Add(readerCount)
	for range readerCount {
		go func() {
			defer readersWG.Done()
			for range iterationsPerReader {
				authzPolicies := am.AuthorizationsForPolicies([]string{"pkg.policy_0"}, "tenant1")
				_ = authzPolicies.GetResources()

				_ = am.GetAssignments("tenant1", "user1")
				_ = am.GetDefaultPolicyNames("tenant1")
				_ = am.GetUserFields()

				input := am.CreateInput("read", "resource1", map[string]any{"field": "x"}, DefaultEnvironmentInput{
					UserInfo: UserInfo{
						Email:  "user1@example.com",
						Groups: []string{"g1"},
					},
				})
				am.ValidateInput(input)
			}
		}()
	}

	const updateIterations = 250
	updatesDone := make(chan struct{})
	go func() {
		defer close(updatesDone)
		for i := 1; i <= updateIterations; i++ {
			select {
			case <-ctx.Done():
				return
			case assignmentsChannel <- dcn.Assignments{
				"tenant1": dcn.UserAssignments{
					"user1": []string{fmt.Sprintf("pkg.policy_%d", i)},
				},
			}:
			}

			select {
			case <-ctx.Done():
				return
			case dcnChannel <- makeValidDCN(i):
			}
		}
	}()

	done := make(chan struct{})
	go func() {
		defer close(done)
		readersWG.Wait()
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("reader goroutines timed out")
	}

	select {
	case <-updatesDone:
	case <-time.After(10 * time.Second):
		t.Fatal("update goroutine timed out")
	}
}

func TestAuthorizationManagerConcurrency_RegisterErrorHandlersDuringErrors(t *testing.T) {
	var callbackErrors atomic.Int64
	dcnChannel := make(chan dcn.DcnContainer, 128)
	assignmentsChannel := make(chan dcn.Assignments, 128)
	am := NewAuthorizationManager(dcnChannel, assignmentsChannel, func(err error) {
		if err != nil {
			callbackErrors.Add(1)
		}
	})

	assignmentsChannel <- dcn.Assignments{
		"tenant1": dcn.UserAssignments{
			"user1": []string{"pkg.policy_0"},
		},
	}
	dcnChannel <- makeValidDCN(0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := am.Run(ctx); err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}

	var registeredHandlersCalled atomic.Int64
	const registerIterations = 200
	const errorIterations = 200

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for range registerIterations {
			am.RegisterErrorHandler(func(err error) {
				if err != nil {
					registeredHandlersCalled.Add(1)
				}
			})
		}
	}()

	go func() {
		defer wg.Done()
		for i := 1; i <= errorIterations; i++ {
			select {
			case <-ctx.Done():
				return
			case dcnChannel <- makeBrokenFunctionDCN(i):
			}
		}
	}()

	waitDone := make(chan struct{})
	go func() {
		defer close(waitDone)
		wg.Wait()
	}()

	select {
	case <-waitDone:
	case <-time.After(10 * time.Second):
		t.Fatal("concurrency scenario timed out")
	}

	deadline := time.Now().Add(2 * time.Second)
	for callbackErrors.Load() == 0 && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}

	if callbackErrors.Load() == 0 {
		t.Fatal("expected at least one error callback to be called")
	}
	if registeredHandlersCalled.Load() == 0 {
		t.Fatal("expected at least one registered error handler to be called")
	}
}

func TestAuthorizationManagerConcurrency_ReadAPIsWithNilInputs(t *testing.T) {
	dcnChannel := make(chan dcn.DcnContainer, 8)
	assignmentsChannel := make(chan dcn.Assignments, 8)
	am := NewAuthorizationManager(dcnChannel, assignmentsChannel, nil)

	assignmentsChannel <- dcn.Assignments{}
	dcnChannel <- makeValidDCN(0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := am.Run(ctx); err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}

	const goroutines = 12
	const iterations = 150

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()
			for range iterations {
				_ = am.AuthorizationsForIdentity(nil)
				_ = am.AuthorizationsForToken(nil)
				_ = am.AuthorizationsForPolicies([]string{}, "")
				_ = am.GetAssignments("unknown", "unknown")
				am.ValidateInput(expression.Input{})
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("nil input concurrency scenario timed out")
	}
}

func TestAuthorizationManagerConcurrency_IdentityAndTokenReads(t *testing.T) {
	dcnChannel := make(chan dcn.DcnContainer, 8)
	assignmentsChannel := make(chan dcn.Assignments, 8)
	am := NewAuthorizationManager(dcnChannel, assignmentsChannel, nil)

	assignmentsChannel <- dcn.Assignments{
		"tenant1": dcn.UserAssignments{
			"user1": []string{"pkg.policy_0"},
		},
	}
	dcnChannel <- makeValidDCN(0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := am.Run(ctx); err != nil {
		t.Fatalf("unexpected run error: %v", err)
	}

	identity := concurrencyIdentity{
		appTID:   "tenant1",
		scimID:   "user1",
		userUUID: "uuid-1",
		groups:   []string{"g1", "g2"},
		email:    "user1@example.com",
	}
	token := concurrencyToken{
		appTID: "tenant1",
		scimID: "user1",
		claims: map[string]interface{}{
			"email":  "user1@example.com",
			"groups": []string{"g1", "g2"},
		},
	}

	const goroutines = 8
	const iterations = 200

	var wg sync.WaitGroup
	wg.Add(goroutines)

	for range goroutines {
		go func() {
			defer wg.Done()
			for range iterations {
				authzIdentity := am.AuthorizationsForIdentity(identity)
				_ = authzIdentity.GetResources()

				authzToken := am.AuthorizationsForToken(token)
				_ = authzToken.GetResources()
			}
		}()
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		wg.Wait()
	}()

	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatal("identity/token read concurrency scenario timed out")
	}
}
