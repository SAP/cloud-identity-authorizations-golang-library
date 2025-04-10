package expression

import (
	"fmt"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/dcn"
	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/util"
)

type Function struct {
	body Expression
	// name string
}

type Functions map[string]Function

func (f Function) Evaluate(input Input) Expression {
	return f.body.Evaluate(input)
}

func FunctionsFromDCN(dcn []dcn.Function) (Functions, error) {
	functions := make(Functions)
	dcnF, err := topologicalSort(dcn)
	if err != nil {
		return nil, err
	}
	for _, f := range dcnF {
		name := util.StringifyQualifiedName(f.QualifiedName)
		expContainer, err := FromDCN(f.Result, functions)
		if err != nil {
			return nil, err
		}
		functions[name] = Function{
			body: expContainer.Expression,
		}
	}

	return functions, nil
}

func topologicalSort(functions []dcn.Function) ([]dcn.Function, error) {
	graph := make(map[string][]string)
	inDegree := make(map[string]int)
	functionsMap := make(map[string]dcn.Function)

	for _, function := range functions {
		name := util.StringifyQualifiedName(function.QualifiedName)
		functionsMap[name] = function
		inDegree[name] = 0
	}

	for _, function := range functions {
		name := util.StringifyQualifiedName(function.QualifiedName)
		for _, call := range getFunctionCalls(function.Result) {
			graph[call] = append(graph[call], name)
			inDegree[name]++
		}
	}

	queue := []string{}
	for name, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, name)
		}
	}

	result := []dcn.Function{}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		result = append(result, functionsMap[current])
		for _, child := range graph[current] {
			inDegree[child]--
			if inDegree[child] == 0 {
				queue = append(queue, child)
			}
		}
	}

	if len(result) != len(functions) {
		return nil, fmt.Errorf("cyclic dependency detected")
	}

	return result, nil
}

func getFunctionCalls(exp dcn.Expression) []string {
	if exp.Call == nil {
		return nil
	}

	if len(exp.Call) == 1 {
		result := []string{}
		for _, arg := range exp.Args {
			result = append(result, getFunctionCalls(arg)...)
		}
		return result
	}

	return []string{util.StringifyQualifiedName(exp.Call)}
}
