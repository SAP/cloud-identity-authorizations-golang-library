package ams

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/internal"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/util"
)

type authorizationState struct {
	policies    internal.PolicySet
	assignments dcn.Assignments
	schema      internal.Schema
}

type AuthorizationManager struct {
	ready              chan bool
	readyOnce          sync.Once
	state              atomic.Pointer[authorizationState]
	dcnChannel         chan dcn.DcnContainer
	assignmentsChannel chan dcn.Assignments
	// Tests              []dcn.Test
	hasDCN            atomic.Bool
	hasAssignments    atomic.Bool
	functionContainer *expression.FunctionRegistry

	bundleLoader  *dcn.BundleLoader
	errHandlersMu sync.RWMutex
	errHandlers   []func(error)
}

// Returns a new AuthorizationManager that listens to the provided DCN and Assignments channels,
// to update its policies and assignments during runtime.
// the instance must receive (possibly empty) data on both channels to be ready.
func NewAuthorizationManager(
	dcnC chan dcn.DcnContainer,
	assignmentsC chan dcn.Assignments,
	errorCallback func(error),
) *AuthorizationManager {
	result := AuthorizationManager{
		ready:              make(chan bool),
		dcnChannel:         dcnC,
		assignmentsChannel: assignmentsC,
		functionContainer:  expression.NewFunctionRegistry(),
		bundleLoader:       nil,
		errHandlers:        []func(error){},
	}
	result.state.Store(&authorizationState{
		policies:    internal.PolicySet{},
		assignments: dcn.Assignments{},
		schema:      internal.Schema{},
	})
	if errorCallback != nil {
		result.errHandlers = append(result.errHandlers, errorCallback)
	}

	return &result
}

func NewOfflineAuthorizationManager(dcn dcn.DcnContainer, assignments dcn.Assignments) (*AuthorizationManager, error) {
	var err error
	result := &AuthorizationManager{
		ready:             make(chan bool),
		functionContainer: expression.NewFunctionRegistry(),
		bundleLoader:      nil,
		errHandlers:       []func(error){},
	}
	state := &authorizationState{
		assignments: assignments,
		schema:      internal.SchemaFromDCN(dcn.Schemas),
	}
	state.policies, err = internal.PoliciesFromDCN(dcn.Policies, state.schema, result.functionContainer)
	if err != nil {
		return nil, err
	}
	result.state.Store(state)
	result.hasDCN.Store(true)
	result.hasAssignments.Store(true)
	result.tryMarkReady()
	return result, nil
}

// Returns a new AuthorizationManager that loads the DCN and Assignments for the given AMS instance
// the provided data should be taken from the identity binding.
func NewAuthorizationManagerForIASConfig(
	config IASConfig,
	errorCallback func(error),
) (*AuthorizationManager, error) {
	return NewAuthorizationManagerForIAS(
		config.GetAuthorizationBundleURL(),
		config.GetAuthorizationInstanceID(),
		config.GetCertificate(),
		config.GetKey(),
		errorCallback,
	)
}

// Returns a new AuthorizationManager that loads the DCN and Assignments for the given AMS instance
// the provided data should be taken from the identity binding.
func NewAuthorizationManagerForIAS(
	bundleUrl,
	amsInstanceID,
	cert,
	key string,
	errorCallback func(error),
) (*AuthorizationManager, error) {
	// parse the cert and key
	certificate, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{certificate},
				MinVersion:   tls.VersionTLS12,
			},
		},
	}
	return newAuthorizationManagerForIAS(bundleUrl, amsInstanceID, client, errorCallback)
}

func newAuthorizationManagerForIAS(
	bundleUrl,
	amsInstanceID string,
	client *http.Client,
	errorCallback func(error),
) (*AuthorizationManager, error) {
	// parse the cert and key

	stringURL, err := url.JoinPath(bundleUrl, amsInstanceID+".dcn.tar.gz")
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(stringURL)
	if err != nil {
		return nil, err
	}

	loader := dcn.NewBundleLoader(
		parsedURL,
		client,
		*time.NewTicker(time.Second * 20),
		errorCallback,
	)

	result := NewAuthorizationManager(loader.DCNChannel, loader.AssignmentsChannel, errorCallback)
	result.bundleLoader = loader
	return result, nil
}

// Returns a new AuthorizationManager that loads the DCN and Assignments from the local file system
// the provided path should contain the schema.dcn and the data.json files and subdirectories
// containing the other dcn files// the data.json file should contain the assignments, if needed
// and could be omitted.
func NewAuthorizationManagerForFs(path string, errorCallback func(error)) *AuthorizationManager {
	loader := dcn.NewLocalLoader(path, errorCallback)
	result := NewAuthorizationManager(loader.DCNChannel, loader.AssignmentsChannel, errorCallback)
	return result
}

func (a *AuthorizationManager) UpdateIASX509Certificate(ctx context.Context, certificate tls.Certificate) error {
	if a.bundleLoader == nil {
		return fmt.Errorf("bundleLoader not initialized, cannot update certificate")
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{certificate},
				MinVersion:   tls.VersionTLS12,
			},
		},
	}
	return a.bundleLoader.SetHttpClient(ctx, client)
}

func (a *AuthorizationManager) Run(ctx context.Context) error {
	if a.bundleLoader != nil {
		a.bundleLoader.Run(ctx)
	}
	readinessContext, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case assignments := <-a.assignmentsChannel:
				current := a.getState()
				a.state.Store(&authorizationState{
					policies:    current.policies,
					assignments: assignments,
					schema:      current.schema,
				})
				a.hasAssignments.Store(true)
				a.tryMarkReady()
				continue
			case dcn := <-a.dcnChannel:
				schema := internal.SchemaFromDCN(dcn.Schemas)
				for _, f := range dcn.Functions {
					expr, err := expression.FromDCN(f.Result, a.functionContainer)
					if err != nil {
						a.notifyError(err)
						continue
					}
					name := util.StringifyQualifiedName(f.QualifiedName)
					a.functionContainer.RegisterExpressionFunction(name, expr.Expression)
				}
				current := a.getState()
				next := &authorizationState{
					policies:    current.policies,
					assignments: current.assignments,
					schema:      schema,
				}

				policies, err := internal.PoliciesFromDCN(dcn.Policies, schema, a.functionContainer)
				if err != nil {
					a.notifyError(err)
				} else {
					next.policies = policies
					a.hasDCN.Store(true)
				}
				a.state.Store(next)
				a.tryMarkReady()
			}
		}
	}()
	select {
	case <-readinessContext.Done():
		return readinessContext.Err()
	case <-a.ready:
		return nil
	}
}

func (a *AuthorizationManager) RegisterErrorHandler(handler func(error)) {
	if handler == nil {
		return
	}
	a.errHandlersMu.Lock()
	defer a.errHandlersMu.Unlock()
	a.errHandlers = append(a.errHandlers, handler)
}

func (a *AuthorizationManager) isReady() bool {
	select {
	case <-a.ready:
		return true
	default:
		return false
	}
}

// Returns Authorizations, based on the provided identity and the default policies.
func (a *AuthorizationManager) AuthorizationsForIdentity(i Identity) *Authorizations {
	state := a.getState()
	if i == nil {
		return &Authorizations{
			policies: state.policies.GetSubset([]string{}, "", false),
			a:        a,
		}
	}

	defaultPolicyNames := state.policies.GetDefaultPolicyNames(i.AppTID())
	assignmentPolicyNames := getAssignments(state.assignments, i.AppTID(), i.ScimID())
	policyNames := make([]string, 0, len(defaultPolicyNames)+len(assignmentPolicyNames))
	policyNames = append(policyNames, defaultPolicyNames...)
	policyNames = append(policyNames, assignmentPolicyNames...)

	return &Authorizations{
		policies: state.policies.GetSubset(policyNames, i.AppTID(), true),
		a:        a,
		envInput: expression.Input{
			"$env.$user.email":     expression.String(i.Email()),
			"$env.$user.user_uuid": expression.String(i.UserUUID()),
			"$env.$user.groups":    expression.ArrayFrom(i.Groups()),
		},
	}
}

func (a *AuthorizationManager) AuthorizationsForToken(t Token) *Authorizations {
	state := a.getState()
	if t == nil {
		return &Authorizations{
			policies: state.policies.GetSubset([]string{}, "", false),
			a:        a,
		}
	}

	defaultPolicyNames := state.policies.GetDefaultPolicyNames(t.AppTID())
	assignmentPolicyNames := getAssignments(state.assignments, t.AppTID(), t.ScimID())
	policyNames := make([]string, 0, len(defaultPolicyNames)+len(assignmentPolicyNames))
	policyNames = append(policyNames, defaultPolicyNames...)
	policyNames = append(policyNames, assignmentPolicyNames...)

	envInput := expression.Input{}
	claims := t.GetAllClaimsAsMap()
	v := reflect.ValueOf(claims)
	state.schema.InsertCustomInput(envInput, v, []string{"$env", "$user"})

	return &Authorizations{
		policies: state.policies.GetSubset(policyNames, t.AppTID(), true),
		a:        a,
		envInput: envInput,
	}
}

// Returns Authorizations, based on the provided policy names and optionally the default policies
// and filtered filtering out admin policies from tenants other than the provided tenant.
// for tenant-independent queries, use "" as tenant.
func (a *AuthorizationManager) AuthorizationsForPolicies(policyNames []string, tenant string) *Authorizations {
	state := a.getState()
	return &Authorizations{
		policies: state.policies.GetSubset(policyNames, tenant, false),
		a:        a,
	}
}

func (a *AuthorizationManager) GetUserFields() map[string]expression.Type {
	allFields := a.getState().schema.GetAllInputFields()
	result := make(map[string]expression.Type)
	for k, t := range allFields {
		if !strings.HasPrefix(k, "$env.$user") {
			continue
		}
		switch t {
		case internal.STRING:
			result[k] = expression.TypeString
		case internal.BOOLEAN:
			result[k] = expression.TypeBool
		case internal.NUMBER:
			result[k] = expression.TypeNumber
		case internal.STRING_ARRAY:
			result[k] = expression.TypeStringArray
		case internal.BOOLEAN_ARRAY:
			result[k] = expression.TypeBoolArray
		case internal.NUMBER_ARRAY:
			result[k] = expression.TypeNumberArray
		case internal.STRUCTURE, internal.UNDEFINED:
			// ignore structures and undefined types,
			// as they cannot be set directly by the user and are not relevant for the user input validation
			continue
		}
	}
	return result
}

func (a *AuthorizationManager) GetDefaultPolicyNames(tenant string) []string {
	return a.getState().policies.GetDefaultPolicyNames(tenant)
}

// Returns the policies that are assigned to the user in the given tenant.
func (a *AuthorizationManager) GetAssignments(tenant, user string) []string {
	return getAssignments(a.getState().assignments, tenant, user)
}

func (a *AuthorizationManager) CreateInput(action, resource string, input any, env any) expression.Input {
	return a.getState().schema.CustomInput(action, resource, input, env)
}

func (a *AuthorizationManager) ValidateInput(input expression.Input) ([]string, []string) {
	return a.getState().schema.PurgeInvalidInput(input)
}

func (a *AuthorizationManager) notifyError(err error) {
	a.errHandlersMu.RLock()
	handlers := make([]func(error), len(a.errHandlers))
	copy(handlers, a.errHandlers)
	a.errHandlersMu.RUnlock()

	for _, handler := range handlers {
		handler(err)
	}
}

func (a *AuthorizationManager) getState() *authorizationState {
	return a.state.Load()
}

func (a *AuthorizationManager) tryMarkReady() {
	if !a.hasDCN.Load() || !a.hasAssignments.Load() {
		return
	}
	a.readyOnce.Do(func() {
		close(a.ready)
	})
}

func getAssignments(assignments dcn.Assignments, tenant, user string) []string {
	t, ok := assignments[tenant]
	if !ok {
		return []string{}
	}
	assignment, ok := t[user]
	if !ok {
		return []string{}
	}
	return assignment
}
