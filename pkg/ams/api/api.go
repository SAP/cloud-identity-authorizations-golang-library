package api

import (
	"context"
	"net/http"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

type API struct {
	am          *ams.AuthorizationManager
	getIdentity func(context.Context) (ams.Identity, error)
}

func NewAPI(am *ams.AuthorizationManager, getIdentity func(context.Context) (ams.Identity, error)) *API {
	return &API{
		am:          am,
		getIdentity: getIdentity,
	}
}

type AmsCtxKey string

const (
	AMSDecisionCtxKey AmsCtxKey = "ams_decision"
	AMSAuthzCtxKey    AmsCtxKey = "ams_authz"
)

func (a *API) Middleware(resource, action string, inputFunc func(*http.Request) any) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authz := authFromContext(r.Context())
			nextR := r
			identity, err := a.getIdentity(r.Context())
			if err != nil {
				http.Error(w, err.Error(), http.StatusNotExtended)
				return
			}
			if identity == nil {
				http.Error(w, "Missing identity in context", http.StatusNotExtended)
				return
			}
			if authz == nil {
				authz = a.am.AuthorizationsForIndentiy(identity)
				nextR = r.WithContext(context.WithValue(r.Context(), AMSAuthzCtxKey, authz))
			}
			var input any
			if inputFunc != nil {
				input = inputFunc(r)
			}
			decision := authz.Inquire(action, resource, input)
			if decision == expression.FALSE {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			nextR = r.WithContext(context.WithValue(r.Context(), AMSDecisionCtxKey, authz))
			next.ServeHTTP(w, nextR)
		})
	}
}

func authFromContext(c context.Context) *ams.Authorizations {
	authz, ok := c.Value("ams_auth").(*ams.Authorizations)
	if !ok {
		return nil
	}
	return authz
}
