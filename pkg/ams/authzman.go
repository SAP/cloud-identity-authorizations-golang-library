package ams

import (
	"context"
	"reflect"
	"sync"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/internal/runtime"
)

type AuthorizationManagerNew struct {
	r                  *runtime.Runtime
	assignments        dcn.Assignments
	runtimeChannel     chan runtime.Runtime
	assignmentsChannel chan dcn.Assignments
	ready              chan bool
	hasRuntime         bool
	hasAssignments     bool
	ctx                context.Context
	cancel             context.CancelFunc
	closed             chan bool
	closeBundleLoader  func(context.Context) error
	m                  sync.RWMutex
}

func NewAuthorizationManagerNew(
	ctx context.Context,
	runtimeC chan runtime.Runtime,
	assignmentsC chan dcn.Assignments,
) *AuthorizationManagerNew {
	ctx, cancel := context.WithCancel(ctx)
	result := AuthorizationManagerNew{
		runtimeChannel:     runtimeC,
		assignmentsChannel: assignmentsC,
		ready:              make(chan bool),
		hasRuntime:         false,
		hasAssignments:     false,
		ctx:                ctx,
		cancel:             cancel,
		closed:             make(chan bool),
		closeBundleLoader:  func(ctx context.Context) error { return nil },
	}

	go result.listen()

	return &result
}

func (am *AuthorizationManagerNew) AuthorizationsForPolicies(policyNames []string) *AuthorizationsNew {
	return &AuthorizationsNew{
		policyNames: policyNames,
		envInput:    expression.Input{},
		r:           am.r,
	}
}

func (am *AuthorizationManagerNew) AuthorizationsForToken(token Token) *AuthorizationsNew {
	envInput := expression.Input{}
	value := reflect.ValueOf(token.GetAllClaimsAsMap())
	runtime.InsertInput(envInput, value, []string{"$env", "$user"})
	am.r.SanitizeInput(envInput, "")
	policyNames := am.GetAssignments(token.AppTID(), token.ScimID())
	policyNames = append(policyNames, am.r.GetDefaultPolicyNames(token.AppTID())...)
	return &AuthorizationsNew{
		policyNames: policyNames,
		envInput:    envInput,
		r:           am.r,
	}
}

func (a *AuthorizationManagerNew) GetAssignments(tenant, user string) []string {
	a.m.RLock()
	defer a.m.RUnlock()
	t, ok := a.assignments[tenant]
	if !ok {
		return []string{}
	}
	assignment, ok := t[user]
	if !ok {
		return []string{}
	}
	return assignment
}

func (a *AuthorizationManagerNew) WhenReady() chan bool {
	return a.ready
}

func (a *AuthorizationManagerNew) IsReady() bool {
	select {
	case <-a.ready:
		return true
	default:
		return false
	}
}

func (am *AuthorizationManagerNew) listen() {
	for {
		select {
		case r := <-am.runtimeChannel:
			am.m.Lock()
			am.r = &r
			am.hasRuntime = true
			if am.hasAssignments && !am.IsReady() {
				close(am.ready)
			}
			am.m.Unlock()
		case a := <-am.assignmentsChannel:
			am.m.Lock()
			am.assignments = a
			am.hasAssignments = true
			if am.hasRuntime && !am.IsReady() {
				close(am.ready)
			}
			am.m.Unlock()
		case <-am.ctx.Done():
			return
		}
	}
}

func (a *AuthorizationManagerNew) Close(ctx context.Context) error {
	a.cancel()
	if a.closeBundleLoader != nil {
		err := a.closeBundleLoader(ctx)
		if err != nil {
			return err
		}
	}
	select {
	case <-a.closed:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
