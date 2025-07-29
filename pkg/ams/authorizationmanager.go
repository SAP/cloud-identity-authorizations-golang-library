package ams

import (
	"crypto/tls"
	"log"
	"net/http"
	"net/url"
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
	Tests              []dcn.Test
	l                  log.Logger
	hasDCN             bool
	hasAssignments     bool
	functionContainer  *expression.FunctionRegistry
}

// Returns a new AuthorizationManager that listens to the provided DCN and Assignments channels,
// to update its policies and assignments during runtime.
// the instance must receive (possibly empty) data on both channels to be ready.
func NewAuthorizationManager(dcnC chan dcn.DcnContainer, assignmentsC chan dcn.Assignments) *AuthorizationManager {
	result := AuthorizationManager{
		ready:              make(chan bool),
		policies:           internal.PolicySet{},
		dcnChannel:         dcnC,
		assignmentsChannel: assignmentsC,
		m:                  sync.RWMutex{},
		hasDCN:             false,
		hasAssignments:     false,
		functionContainer:  expression.NewFunctionRegistry(),
	}

	go result.start()

	return &result
}

// Returns a new AuthorizationManager that loads the DCN and Assignments for the given AMS instance
// the provided data should be taken from the identity binding.
func NewAuthorizationManagerForIAS(bundleUrl, amsInstanceID, cert, key string) (*AuthorizationManager, error) {
	// parse the cert and key
	certificate, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return nil, err
	}

	stringURL, err := url.JoinPath(bundleUrl, amsInstanceID+".dcn.tar.gz")
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(stringURL)
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

	loader := dcn.NewBundleLoader(
		parsedURL,
		client,
		*time.NewTicker(time.Second * 20),
	)

	result := NewAuthorizationManager(loader.DCNChannel, loader.AssignmentsChannel)
	loader.RegisterErrorHandler(result.notifyError)
	return result, nil
}

// Returns a new AuthorizationManager that loads the DCN and Assignments from the local file system
// the provided path should contain the schema.dcn and the data.json files and subdirectories
// containing the other dcn files// the data.json file should contain the assignments, if needed
// and could be omitted.
func NewAuthorizationManagerForFs(path string) *AuthorizationManager {
	loader := dcn.NewLocalLoader(path)
	result := NewAuthorizationManager(loader.DCNChannel, loader.AssignmentsChannel)
	loader.RegisterErrorHandler(result.notifyError)
	return result
}

// Register a new log callback that will be called when a log message is generated.
func (a *AuthorizationManager) SetLogger(l log.Logger) {
	a.m.Lock()
	defer a.m.Unlock()
	a.l = l
}

func (a *AuthorizationManager) start() {
	for {
		select {
		case assignments := <-a.assignmentsChannel:
			a.m.Lock()
			a.Assignments = assignments
			a.hasAssignments = true
			if a.hasDCN {
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
			a.Tests = dcn.Tests
			a.m.Unlock()
			if !a.IsReady() {
				if a.hasDCN && a.hasAssignments {
					close(a.ready)
				}
			}
		}
	}
}

// Returns a channel that will be closed when the AuthorizationManager is ready to be used.
func (a *AuthorizationManager) WhenReady() <-chan bool {
	return a.ready
}

// Returns true if the AuthorizationManager is ready to be used
// This is the case when both the DCN and Assignments have been loaded.
func (a *AuthorizationManager) IsReady() bool {
	select {
	case <-a.ready:
		return true
	default:
		return false
	}
}

// Returns Schema that can be used for input creation/validation based on the DCL schema.
func (a *AuthorizationManager) GetSchema() internal.Schema {
	a.m.RLock()
	defer a.m.RUnlock()
	return a.schema
}

// Returns Authorizations, based on the provided identity and the default policies
func (a *AuthorizationManager) AuthorizationsForIndentiy(i Identity) *Authorizations {
	a.m.RLock()
	defer a.m.RUnlock()
	if i == nil {
		return &Authorizations{
			policies: a.policies.GetSubset([]string{}, "", true),
			schema:   a.schema,
		}
	}

	policyNames := a.policies.GetDefaultPolicyNames(i.AppTID())

	policyNames = append(policyNames, a.GetAssignments(i.AppTID(), i.ScimID())...)

	return &Authorizations{
		policies: a.policies.GetSubset(policyNames, i.AppTID(), true),
		schema:   a.schema,
		envInput: expression.Input{
			"$env.$user.email":     expression.String(i.Email()),
			"$env.$user.user_uuid": expression.String(i.UserUUID()),
			"$env.$user.groups":    expression.StringArrayFrom(i.Groups()),
		},
	}
}

// Returns Authorizations, based on the provided policy names and optionally the default policies
// and filtered filtering out admin policies from tenants other than the provided tenant.
// for tenant-independent queries, use "" as tenant.
func (a *AuthorizationManager) AuthorizationsForPolicies(policyNames []string) *Authorizations {
	a.m.RLock()
	defer a.m.RUnlock()
	return &Authorizations{
		policies: a.policies.GetSubset(policyNames, "-", false),
		schema:   a.schema,
	}
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
