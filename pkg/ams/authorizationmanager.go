package ams

import (
	"context"
	"crypto/tls"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/internal"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/logging"
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
	l                 logging.Logger
	ctx               context.Context
	cancel            context.CancelFunc
	closed            chan bool
	closeBundleLoader func(context.Context) error
}

// Returns a new AuthorizationManager that listens to the provided DCN and Assignments channels,
// to update its policies and assignments during runtime.
// the instance must receive (possibly empty) data on both channels to be ready.
func NewAuthorizationManager(
	ctx context.Context,
	dcnC chan dcn.DcnContainer,
	assignmentsC chan dcn.Assignments,
	log logging.Logger,
) *AuthorizationManager {
	ctx, cancel := context.WithCancel(ctx)
	result := AuthorizationManager{
		ready:              make(chan bool),
		policies:           internal.PolicySet{},
		dcnChannel:         dcnC,
		assignmentsChannel: assignmentsC,
		m:                  sync.RWMutex{},
		hasDCN:             false,
		hasAssignments:     false,
		functionContainer:  expression.NewFunctionRegistry(),
		l:                  log,
		ctx:                ctx,
		cancel:             cancel,
		closed:             make(chan bool),
	}
	if result.l == nil {
		result.l = logging.Default()
	}

	go result.start()

	return &result
}

// Returns a new AuthorizationManager that loads the DCN and Assignments for the given AMS instance
// the provided data should be taken from the identity binding.
func NewAuthorizationManagerForIASConfig(
	ctx context.Context,
	config IASConfig,
	log logging.Logger,
) (*AuthorizationManager, error) {
	return NewAuthorizationManagerForIAS(
		ctx,
		config.GetAuthorizationBundleURL(),
		config.GetAuthorizationInstanceID(),
		config.GetCertificate(),
		config.GetKey(),
		log,
	)
}

// Returns a new AuthorizationManager that loads the DCN and Assignments for the given AMS instance
// the provided data should be taken from the identity binding.
func NewAuthorizationManagerForIAS(
	ctx context.Context,
	bundleUrl,
	amsInstanceID,
	cert,
	key string,
	log logging.Logger,
) (*AuthorizationManager, error) {
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
		ctx,
		parsedURL,
		client,
		*time.NewTicker(time.Second * 20),
		log,
	)

	result := NewAuthorizationManager(ctx, loader.DCNChannel, loader.AssignmentsChannel, log)
	result.closeBundleLoader = loader.Close
	return result, nil
}

// Returns a new AuthorizationManager that loads the DCN and Assignments from the local file system
// the provided path should contain the schema.dcn and the data.json files and subdirectories
// containing the other dcn files// the data.json file should contain the assignments, if needed
// and could be omitted.
func NewAuthorizationManagerForFs(path string, log logging.Logger) *AuthorizationManager {
	loader := dcn.NewLocalLoader(path, log)
	result := NewAuthorizationManager(context.Background(), loader.DCNChannel, loader.AssignmentsChannel, log)

	return result
}

func (a *AuthorizationManager) start() {
	for {
		select {
		case <-a.ctx.Done():
			close(a.closed)
			return
		case assignments := <-a.assignmentsChannel:
			a.m.Lock()
			a.Assignments = assignments
			a.hasAssignments = true
			if !a.IsReady() && a.hasDCN {
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

func (a *AuthorizationManager) Close(ctx context.Context) error {
	a.cancel()
	err := a.closeBundleLoader(ctx)
	if err != nil {
		return err
	}
	select {
	case <-a.closed:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Returns Authorizations, based on the provided identity and the default policies.
func (a *AuthorizationManager) AuthorizationsForIdentity(ctx context.Context, i Identity) *Authorizations {
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
	policyNames := append(defaultPolicyNames, assignmentPolicyNames...)
	a.l.Infof(ctx, "AuthorizationsForIdentity: for user %s in tenant %s, default policies: %v, assignment policies: %v", i.ScimID(), i.AppTID(), len(defaultPolicyNames), len(assignmentPolicyNames))

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

// Returns Authorizations, based on the provided policy names and optionally the default policies
// and filtered filtering out admin policies from tenants other than the provided tenant.
// for tenant-independent queries, use "" as tenant.
func (a *AuthorizationManager) AuthorizationsForPolicies(ctx context.Context, policyNames []string) *Authorizations {
	a.m.RLock()
	defer a.m.RUnlock()
	return &Authorizations{
		policies: a.policies.GetSubset(policyNames, "-", false),
		a:        a,
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
	if a.l != nil {
		a.l.Errorf(a.ctx, err.Error())
	}
}
