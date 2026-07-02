package runtime

import (
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
)

func (r *Runtime) GetDefaultPolicyNames(_ string) []string {
	return r.defaultPolicyNames
}

func (r *Runtime) GetResources(policyNames []string) []string {
	resourceSet := make(map[string]struct{})
	for _, p := range policyNames {
		policy, ok := r.Policies[p]
		if !ok {
			continue
		}
		for resource := range policy.Resources {
			resourceSet[resource] = struct{}{}
		}
	}
	resources := make([]string, 0, len(resourceSet))
	for resource := range resourceSet {
		resources = append(resources, resource)
	}
	return resources
}

func (r *Runtime) GetActions(policyNames []string, resource string) []string {
	actionSet := make(map[string]struct{})
	for _, p := range policyNames {
		policy, ok := r.Policies[p]
		if !ok {
			continue
		}
		if res, ok := policy.Resources[resource]; ok {
			for action := range res.Conditions {
				actionSet[action] = struct{}{}
			}
		}
		if res, ok := policy.Resources[WildCard]; ok {
			for action := range res.Conditions {
				actionSet[action] = struct{}{}
			}
		}
	}
	actions := make([]string, 0, len(actionSet))
	for action := range actionSet {
		actions = append(actions, action)
	}
	return actions
}

func (r *Runtime) Evaluate(
	policyNames []string,
	action string,
	resource string,
	input expression.Input) expression.Expression {
	args := make([]expression.Expression, 0)
	var arg expression.Expression
	for _, p := range policyNames {
		policy, ok := r.Policies[p]
		if !ok {
			continue
		}
		if resource, ok := policy.Resources[resource]; ok {
			arg = resource.Conditions.evaluate(action, input)
			if arg == expression.TRUE {
				return arg
			}
			if arg != expression.FALSE {
				args = append(args, arg)
			}
		}
		if resource, ok := policy.Resources[WildCard]; ok {
			arg = resource.Conditions.evaluate(action, input)
			if arg == expression.TRUE {
				return arg
			}
			if arg != expression.FALSE {
				args = append(args, arg)
			}
		}
	}
	return expression.Or(args...)
}

func (a *ActionsConditions) evaluate(action string, input expression.Input) expression.Expression {
	args := make([]expression.Expression, 0)
	var arg expression.Expression
	if expr, ok := (*a)[action]; ok {
		arg = expr
		if len(input) > 0 {
			arg = expr.Evaluate(input)
		}
		if arg == expression.TRUE {
			return arg
		}
		if arg != expression.FALSE {
			args = append(args, arg)
		}
	}
	if expr, ok := (*a)[WildCard]; ok {
		arg = expr
		if len(input) > 0 {
			arg = expr.Evaluate(input)
		}
		if arg == expression.TRUE {
			return arg
		}
		if arg != expression.FALSE {
			args = append(args, arg)
		}
	}
	return expression.Or(args...)
}
