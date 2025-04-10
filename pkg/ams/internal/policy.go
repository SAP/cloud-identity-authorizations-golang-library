package internal

import (
	"fmt"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/util"
)

type Policy struct {
	rules         []Rule
	tenant        string
	defaultPolicy bool
	name          string
}

type PolicySet struct {
	allPolicies     map[string]Policy
	defaultPolicies map[string][]Policy
}

func (p PolicySet) Evaluate(input expression.Input) expression.Expression {
	results := []expression.Expression{}
	for _, policy := range p.allPolicies {
		r := policy.Evaluate(input)
		if r == expression.Bool(true) {
			return r
		}
		if r != expression.Bool(false) {
			results = append(results, r)
		}
	}
	return expression.NewOr(results...)
}

func (p Policy) Evaluate(input expression.Input) expression.Expression {
	results := []expression.Expression{}
	for _, rule := range p.rules {
		r := rule.Evaluate(input)
		if r == expression.Bool(true) {
			return r
		}
		if r != expression.Bool(false) {
			results = append(results, r)
		}
	}
	return expression.NewOr(results...)
}

func (p PolicySet) GetSubset(names []string, tenant string, includeDefault bool) PolicySet {
	result := PolicySet{
		allPolicies:     make(map[string]Policy),
		defaultPolicies: make(map[string][]Policy),
	}
	for _, name := range names {
		if policy, ok := p.allPolicies[name]; ok {
			if policy.tenant == tenant || policy.tenant == "" {
				result.allPolicies[name] = policy
			}
		}
	}
	if includeDefault {
		result.defaultPolicies[""] = p.defaultPolicies[""]
		if tenant != "" {
			result.defaultPolicies[tenant] = p.defaultPolicies[tenant]
		}
	}
	for _, t := range result.defaultPolicies {
		for _, policy := range t {
			result.allPolicies[policy.name] = policy
		}
	}

	return result
}

func PoliciesFromDCN(policies []dcn.Policy, schema Schema, f expression.Functions) (PolicySet, error) {
	result := PolicySet{
		allPolicies:     make(map[string]Policy),
		defaultPolicies: make(map[string][]Policy),
	}

	policies, err := topologicalSort(policies)
	if err != nil {
		return result, err
	}

	for _, policy := range policies {
		policyName := util.StringifyQualifiedName(policy.QualifiedName)
		p := Policy{
			tenant:        schema.GetTenantForQualifiedName(policy.QualifiedName),
			defaultPolicy: policy.Default,
			name:          policyName,
		}

		for _, rule := range policy.Rules {
			r, err := RuleFromDCN(rule, f)
			if err != nil {
				return result, err
			}
			p.rules = append(p.rules, r)
		}
		for _, use := range policy.Uses {
			usedPolicyName := util.StringifyQualifiedName(use.QualifiedPolicyName)
			// the used policy is always present and initialized, because of the topological sort
			usedPolicy := result.allPolicies[usedPolicyName]

			for _, usedRule := range usedPolicy.rules {
				if expression.IsRestrictable(usedRule.asExpression) {
					if len(use.Restrictions) == 0 {
						newRule := Rule{
							asExpression: usedRule.asExpression,
						}
						p.rules = append(p.rules, newRule)
					}
					for _, restriction := range use.Restrictions {
						restrictions := []expression.ExpressionContainer{}
						for _, r := range restriction {
							container, err := expression.FromDCN(r, f)
							if err != nil {
								return result, err
							}
							restrictions = append(restrictions, container)
						}

						newRule := Rule{
							asExpression: expression.ApplyRestriction(usedRule.asExpression, restrictions),
						}
						p.rules = append(p.rules, newRule)
					}
				} else {
					p.rules = append(p.rules, usedRule)
				}
			}
		}
		result.allPolicies[policyName] = p
		if policy.Default {
			if _, ok := result.defaultPolicies[p.tenant]; !ok {
				result.defaultPolicies[p.tenant] = []Policy{}
			}
			result.defaultPolicies[p.tenant] = append(result.defaultPolicies[p.tenant], p)
		}
	}
	return result, nil
}

func topologicalSort(policies []dcn.Policy) ([]dcn.Policy, error) {
	graph := make(map[string][]string)
	inDegree := make(map[string]int)
	policyMap := make(map[string]dcn.Policy)

	for _, policy := range policies {
		name := util.StringifyQualifiedName(policy.QualifiedName)
		policyMap[name] = policy
		inDegree[name] = 0
	}
	for _, policy := range policies {
		name := util.StringifyQualifiedName(policy.QualifiedName)
		for _, use := range policy.Uses {
			useName := util.StringifyQualifiedName(use.QualifiedPolicyName)
			graph[useName] = append(graph[useName], name)
			inDegree[name]++
		}
	}
	queue := []string{}
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}
	result := []dcn.Policy{}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		result = append(result, policyMap[node])
		for _, child := range graph[node] {
			inDegree[child]--
			if inDegree[child] == 0 {
				queue = append(queue, child)
			}
		}
	}
	if len(result) != len(policies) {
		return nil, fmt.Errorf("cyclic dependency detected")
	}
	return result, nil
}
