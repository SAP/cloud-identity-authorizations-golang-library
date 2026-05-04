package ams

import (
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func MockDecision(expr expression.Expression, am *AuthorizationManager) Decision {
	if am == nil {
		return Decision{
			condition: expr,
			inputConverter: func(_ any) expression.Input {
				return expression.Input{}
			},
		}
	}
	return Decision{
		condition: expr,
		inputConverter: func(_ any) expression.Input {
			return expression.Input{}
		},
	}
}
