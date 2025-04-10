package internal

import (
	_ "embed"
	"encoding/json"
	"testing"

	dcn "github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

//go:embed testfiles/rules/grant_all.dcn
var grantAll []byte

//go:embed testfiles/rules/grant_actions.dcn
var grantActions []byte

//go:embed testfiles/rules/grant_conditional.dcn
var grantConditional []byte

//go:embed testfiles/rules/grant_resources.dcn
var grantResources []byte

func TestRule(t *testing.T) {
	t.Run("Evaluate Grant All Rule", func(t *testing.T) {
		dcnRule := dcn.Rule{}
		err := json.Unmarshal(grantAll, &dcnRule)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		rule, err := RuleFromDCN(dcnRule, nil)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		result := rule.Evaluate(nil)
		if result != expression.Bool(true) {
			t.Errorf("Unexpected result: %s", result)
		}
	})

	t.Run("Evaluate Grant Actions Rule", func(t *testing.T) {
		dcnRule := dcn.Rule{}
		err := json.Unmarshal(grantActions, &dcnRule)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		rule, err := RuleFromDCN(dcnRule, nil)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		result := rule.Evaluate(expression.Input{
			"$dcl.action": expression.String("read"),
		})

		if result != expression.Bool(true) {
			t.Errorf("Unexpected result: %s", result)
		}

		result = rule.Evaluate(expression.Input{
			"$dcl.action": expression.String("write"),
		})

		if result != expression.Bool(true) {
			t.Errorf("Unexpected result: %s", result)
		}

		result = rule.Evaluate(expression.Input{
			"$dcl.action": expression.String("forbidden"),
		})
		if result != expression.Bool(false) {
			t.Errorf("Unexpected result: %s", result)
		}
	})

	t.Run("Evaluate Grant Resources Rule", func(t *testing.T) {
		dcnRule := dcn.Rule{}
		err := json.Unmarshal(grantResources, &dcnRule)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		rule, err := RuleFromDCN(dcnRule, nil)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		result := rule.Evaluate(expression.Input{
			"$dcl.resource": expression.String("resource1"),
		})
		if result != expression.Bool(true) {
			t.Errorf("Unexpected result: %s", result)
		}

		result = rule.Evaluate(expression.Input{
			"$dcl.action": expression.String("asdf"),
		})
		if result != expression.Bool(false) {
			t.Errorf("Unexpected result: %s", result)
		}

		result = rule.Evaluate(expression.Input{
			"$dcl.resource": expression.IGNORE,
		})
		if result != expression.Bool(true) {
			t.Errorf("Unexpected result: %s", result)
		}

		result = rule.Evaluate(expression.Input{
			"$dcl.resource": expression.String("database"),
		})
		if result != expression.Bool(false) {
			t.Errorf("Unexpected result: %s", result)
		}
	})

	t.Run("Evaluate Grant Conditional Rule", func(t *testing.T) {
		dcnRule := dcn.Rule{}
		err := json.Unmarshal(grantConditional, &dcnRule)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		rule, err := RuleFromDCN(dcnRule, nil)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		result := rule.Evaluate(expression.Input{
			"$dcl.action":   expression.String("read"),
			"$dcl.resource": expression.String("resource1"),
			"x":             expression.Number(1),
		})
		if result != expression.Bool(true) {
			t.Errorf("Unexpected result: %s", result)
		}
		result = rule.Evaluate(expression.Input{
			"$dcl.action":   expression.String("read"),
			"$dcl.resource": expression.String("resource1"),
			"x":             expression.Number(2),
		})
		if result != expression.Bool(false) {
			t.Errorf("Unexpected result: %s", result)
		}
		result = rule.Evaluate(expression.Input{
			"$dcl.action":   expression.String("write"),
			"$dcl.resource": expression.String("resource1"),
		})
		if result != expression.Bool(false) {
			t.Errorf("Unexpected result: %s", result)
		}

		result = rule.Evaluate(expression.Input{
			"$dcl.action":   expression.String("read"),
			"$dcl.resource": expression.UNKNOWN,
			"x":             expression.Number(1),
		})
		if expression.ToString(result) != "in($dcl.resource, [resource1 resource2 resource3])" {
			t.Errorf("Unexpected result: %s", expression.ToString(result))
		}
	})
}
