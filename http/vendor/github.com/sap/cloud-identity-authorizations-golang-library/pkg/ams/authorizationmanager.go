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
	"time"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/internal"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/util"
)

type AuthorizationManager struct {
	ready              chan bool
	policies           internal.PolicySet
	Assignments        dcn.Assignments
	m                  sync.RWMutex
	schema             internal.Schema
	dcnChannel         chan dcn.DcnContainer
	assignmentsChannel chan dcn.Assignments
	// Tests              []dcn.Test
	hasDCN            bool
	hasAssignments    bool
	functionContainer *expression.FunctionRegistry

	bundleLoader *dcn.BundleLoader
	errHandlers  []func(error)
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
		policies:           internal.PolicySet{},
		dcnChannel:         dcnC,
		assignmentsChannel: assignmentsC,
		m:                  sync.RWMutex{},
		hasDCN:             false,
		hasAssignments:     false,
		functionContainer:  expression.NewFunctionRegistry(),
		bundleLoader:       nil,
		errHandlers:        []func(error){},
	}
	if errorCallback != nil {
		result.errHandlers = append(result.errHandlers, errorCallback)
	}

	return &result
}

func NewOfflineAuthorizationManager(dcn dcn.DcnContainer, assignments dcn.Assignments) (*AuthorizationManager, error) {
	var err error
	result := &AuthorizationManager{
		ready:             make(chan bool),
		m:                 sync.RWMutex{},
		hasDCN:            false,
		hasAssignments:    false,
		functionContainer: expression.NewFunctionRegistry(),
		bundleLoader:      nil,
		errHandlers:       []func(error){},
	}
	result.Assignments = assignments
	result.schema = internal.SchemaFromDCN(dcn.Schemas)
	result.policies, err = internal.PoliciesFromDCN(dcn.Policies, result.schema, result.functionContainer)
	if err != nil {
		return nil, err
	}
	return result, result.Run(context.Background())
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
				a.m.Lock()
				a.Assignments = assignments
				a.hasAssignments = true
				if !a.isReady() && a.hasDCN {
					close(a.ready)
				}
				a.m.Unlock()
				continue
			case dcn := <-a.dcnChannel:
				a.m.Lock()
				a.schema = internal.SchemaFromDCN(dcn.Schemas)
				for _, f := range dcn.Functions {
					expr, err := expression.FromDCN(f.Result, a.functionContainer)
					if err != nil {
						a.notifyError(err)
						continue
					}
					name := util.StringifyQualifiedName(f.QualifiedName)
					a.functionContainer.RegisterExpressionFunction(name, expr.Expression)
				}
				var err error
				a.policies, err = internal.PoliciesFromDCN(dcn.Policies, a.schema, a.functionContainer)
				if err != nil {
					a.notifyError(err)
				} else {
					a.hasDCN = true
				}

				a.m.Unlock()
				if !a.isReady() {
					if a.hasDCN && a.hasAssignments {
						close(a.ready)
					}
				}
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
	a.m.Lock()
	defer a.m.Unlock()
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
	a.m.RLock()
	defer a.m.RUnlock()
	if i == nil {
		return &Authorizations{
			policies: a.policies.GetSubset([]string{}, "", false),
			a:        a,
		}
	}

	defaultPolicyNames := a.policies.GetDefaultPolicyNames(i.AppTID())

	assignmentPolicyNames := a.GetAssignments(i.AppTID(), i.ScimID())
	policyNames := make([]string, 0, len(defaultPolicyNames)+len(assignmentPolicyNames))
	policyNames = append(policyNames, defaultPolicyNames...)
	policyNames = append(policyNames, assignmentPolicyNames...)

	return &Authorizations{
		policies: a.policies.GetSubset(policyNames, i.AppTID(), true),
		a:        a,
		envInput: expression.Input{
			"$env.$user.email":     expression.String(i.Email()),
			"$env.$user.user_uuid": expression.String(i.UserUUID()),
			"$env.$user.groups":    expression.ArrayFrom(i.Groups()),
		},
	}
}

func (a *AuthorizationManager) AuthorizationsForToken(t Token) *Authorizations {
	a.m.RLock()
	defer a.m.RUnlock()
	if t == nil {
		return &Authorizations{
			policies: a.policies.GetSubset([]string{}, "", false),
			a:        a,
		}
	}

	defaultPolicyNames := a.policies.GetDefaultPolicyNames(t.AppTID())
	assignmentPolicyNames := a.GetAssignments(t.AppTID(), t.ScimID())
	policyNames := make([]string, 0, len(defaultPolicyNames)+len(assignmentPolicyNames))
	policyNames = append(policyNames, defaultPolicyNames...)
	policyNames = append(policyNames, assignmentPolicyNames...)

	envInput := expression.Input{}
	claims := t.GetAllClaimsAsMap()
	v := reflect.ValueOf(claims)
	a.schema.InsertCustomInput(envInput, v, []string{"$env", "$user"})

	return &Authorizations{
		policies: a.policies.GetSubset(policyNames, t.AppTID(), true),
		a:        a,
		envInput: envInput,
	}
}

// Returns Authorizations, based on the provided policy names and optionally the default policies
// and filtered filtering out admin policies from tenants other than the provided tenant.
// for tenant-independent queries, use "" as tenant.
func (a *AuthorizationManager) AuthorizationsForPolicies(policyNames []string, tenant string) *Authorizations {
	a.m.RLock()
	defer a.m.RUnlock()
	return &Authorizations{
		policies: a.policies.GetSubset(policyNames, tenant, false),
		a:        a,
	}
}

func (a *AuthorizationManager) GetUserFields() map[string]expression.Type {
	a.m.RLock()
	defer a.m.RUnlock()
	allFields := a.schema.GetAllInputFields()
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
	a.m.RLock()
	defer a.m.RUnlock()
	return a.policies.GetDefaultPolicyNames(tenant)
}

// Returns the policies that are assigned to the user in the given tenant.
func (a *AuthorizationManager) GetAssignments(tenant, user string) []string {
	a.m.RLock()
	defer a.m.RUnlock()
	t, ok := a.Assignments[tenant]
	if !ok {
		return []string{}
	}
	assignment, ok := t[user]
	if !ok {
		return []string{}
	}
	return assignment
}

func (a *AuthorizationManager) CreateInput(action, resource string, input any, env any) expression.Input {
	a.m.RLock()
	defer a.m.RUnlock()
	return a.schema.CustomInput(action, resource, input, env)
}

func (a *AuthorizationManager) ValidateInput(input expression.Input) ([]string, []string) {
	a.m.RLock()
	defer a.m.RUnlock()
	return a.schema.PurgeInvalidInput(input)
}

func (a *AuthorizationManager) notifyError(err error) {
	for _, handler := range a.errHandlers {
		handler(err)
	}
}
