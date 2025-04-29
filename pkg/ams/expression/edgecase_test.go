package expression

import (
	"testing"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
)

type UnexpectedExpression struct{}

func (e UnexpectedExpression) Evaluate(input Input) Expression {
	return FALSE
}

func TestEdgeCases(t *testing.T) {
	t.Run("Visit unexpected expression", func(t *testing.T) {
		e := UnexpectedExpression{}
		r := Visit(e, func(c string, args []string) string {
			return c
		}, func(Reference) string {
			return ""
		}, func(Constant) string {
			return ""
		})
		if r != "unexpected_expression" {
			t.Errorf("Unexpected result: %s", r)
		}
	})

	t.Run("unsupported constant as expression Arguemnt", func(t *testing.T) {
		cExp := dcn.Expression{
			Call: []string{"and"},
			Args: []dcn.Expression{{Constant: struct{}{}}},
		}
		_, err := FromDCN(cExp, nil)
		if err == nil {
			t.Errorf("Expected error")
		}
	})

	t.Run("empty call", func(t *testing.T) {
		cExp := dcn.Expression{
			Call: []string{},
		}
		_, err := FromDCN(cExp, nil)
		if err == nil {
			t.Errorf("Expected error")
		}
	})
}
