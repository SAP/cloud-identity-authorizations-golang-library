package ams

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/internal"
)

type AuthorizationManager struct {
	ready              chan bool
	policies           internal.PolicySet
	Assignments        dcn.Assignments
	m                  sync.RWMutex
	schema             internal.Schema
	dcnChannel         chan dcn.DcnContainer
	assignmentsChannel chan dcn.Assignments
	errHandlers        []func(error)
	Tests              []dcn.Test
	hasDCN             bool
	hasAssignments     bool
}

func NewAuthorizationManager(dcnChannel chan dcn.DcnContainer, assignmentsChannel chan dcn.Assignments) *AuthorizationManager {
	result := AuthorizationManager{
		ready:              make(chan bool),
		policies:           internal.PolicySet{},
		dcnChannel:         dcnChannel,
		assignmentsChannel: assignmentsChannel,
		errHandlers:        []func(error){},
		m:                  sync.RWMutex{},
		hasDCN:             false,
		hasAssignments:     false,
	}

	go result.start()

	return &result
}

func AuthorizationManagerForAMS(bundleURL, amsInstanceID, cert, key string) (*AuthorizationManager, error) {

	//parse the cert and key
	certificate, err := tls.X509KeyPair([]byte(cert), []byte(key))
	if err != nil {
		return nil, err
	}

	parsedURL, err := url.Parse(bundleURL)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{certificate},
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

func AuthorizationManagerForLocal(path string) *AuthorizationManager {
	loader := dcn.NewLocalLoader(path)
	result := NewAuthorizationManager(loader.DCNChannel, loader.AssignmentsChannel)
	loader.RegisterErrorHandler(result.notifyError)
	return result
}

func (a *AuthorizationManager) RegisterErrorHandler(handler func(error)) {
	a.m.Lock()
	defer a.m.Unlock()
	a.errHandlers = append(a.errHandlers, handler)
}

func (a *AuthorizationManager) notifyError(err error) {
	for _, handler := range a.errHandlers {
		handler(err)
	}
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
			functions, err := expression.FunctionsFromDCN(dcn.Functions)
			if err != nil {
				a.notifyError(err)
				a.m.Unlock()
				continue
			}
			a.policies, err = internal.PoliciesFromDCN(dcn.Policies, a.schema, functions)
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

func (a *AuthorizationManager) WhenReady() <-chan bool {
	return a.ready
}

func (a *AuthorizationManager) IsReady() bool {
	select {
	case <-a.ready:
		return true
	default:
		return false
	}
}

func (a *AuthorizationManager) GetSchema() internal.Schema {
	a.m.RLock()
	defer a.m.RUnlock()
	return a.schema
}
func (a *AuthorizationManager) UserAuthorizations(tenant, user string) *Authorizations {
	pNames := a.GetAssignments(tenant, user)
	return a.GetAuthorizations(pNames, tenant, true)
}

func (a *AuthorizationManager) GetAuthorizations(names []string, tenant string, includeDefault bool) *Authorizations {
	a.m.RLock()
	defer a.m.RUnlock()
	return &Authorizations{
		policies: a.policies.GetSubset(names, tenant, includeDefault),
		schema:   a.schema,
	}
}

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
