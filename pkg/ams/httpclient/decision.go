package httpclient

import (
	"context"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

type Decision struct {
	condition      expression.Expression
	inputConverter func(any) (expression.Input, error)
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

func (d Decision) Inquire(ctx context.Context, app any) (Decision, error) {
	input, err := d.inputConverter(app)
	if err != nil {
		return Decision{}, err
	}
	return d.Evaluate(input), nil
}

func (d Decision) Evaluate(input expression.Input) Decision {
	return Decision{
		condition:      d.condition.Evaluate(input),
		inputConverter: d.inputConverter,
	}
}
