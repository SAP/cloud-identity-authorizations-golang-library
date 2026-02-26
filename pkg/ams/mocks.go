package ams

import (
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/internal"
)

func MockDecision(expr expression.Expression, am *AuthorizationManager) Decision {
	if am == nil {
		return Decision{
			condition: expr,
			schema:    internal.Schema{},
		}
	}
	return Decision{
		condition: expr,
		schema:    am.schema,
	}
}
