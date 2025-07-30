package api

import (
	"context"
	"net/http"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

type API struct {
	am          *ams.AuthorizationManager
	getIdentity func(context.Context) ams.Identity
}

func NewAPI(am *ams.AuthorizationManager, getIdentity func(context.Context) ams.Identity) *API {
	return &API{
		am:          am,
		getIdentity: getIdentity,
	}
}

func (a *API) Middleware(resource, action string, input any) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authz := authFromContext(r.Context())
			nextR := r
			identity := a.getIdentity(r.Context())
			if identity == nil {
				http.Error(w, "Missing identity in context", http.StatusNotExtended)
				return
			}
			if authz == nil {
				authz = a.am.AuthorizationsForIndentiy(identity)
				nextR = r.WithContext(context.WithValue(r.Context(), "ams_auth", authz))
			}

			decision := authz.Inquire(action, resource, input)
			if decision == expression.FALSE {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
			nextR = r.WithContext(context.WithValue(r.Context(), "ams_decision", authz))
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
