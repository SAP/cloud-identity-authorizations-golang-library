package ams

import (
	"reflect"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/internal"
)

type Decision struct {
	condition expression.Expression
	schema    internal.Schema
}

func (d Decision) Condition() expression.Expression {
	return d.condition
}

func (d Decision) IsGranted() bool {
	return d.condition == expression.TRUE
}

func (d Decision) IsDenied() bool {
	return d.condition == expression.FALSE
}

func (d Decision) Inquire(app any) Decision {
	i := expression.Input{}

	d.schema.InsertCustomInput(i, reflect.ValueOf(app), []string{"$app"})
	return d.Evaluate(i)
}

func (d Decision) Evaluate(input expression.Input) Decision {
	return Decision{
		schema:    d.schema,
		condition: d.condition.Evaluate(input),
	}
}
