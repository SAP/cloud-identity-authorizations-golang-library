package runtime

import (
	"encoding/json"
	"fmt"

	"github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression"
	. "github.com/sap/cloud-identity-authorizations-golang-library/pkg/ams/expression" // nolint:staticcheck
)

func expressionFromJSON(data json.RawMessage) (Expression, error) {
	var rawOperator map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawOperator); err != nil {
		return constantFromJSON(data)
	}
	if v, ok := rawOperator["var"]; ok {
		var varName string
		if err := json.Unmarshal(v, &varName); err != nil {
			return nil, fmt.Errorf("invalid var expression: %s", string(data))
		}
		return expression.Ref(varName), nil
	}

	for op, args := range rawOperator {
		args, err := argsFromJSON(args)
		if err != nil {
			return nil, err
		}
		switch op {
		case "and":
			return expression.And(args...), nil
		case "or":
			return expression.Or(args...), nil
		case "eq":
			return expression.Eq(args[0], args[1]), nil
		case "ne":
			return expression.Ne(args[0], args[1]), nil
		case "in":
			return expression.In(args[0], args[1]), nil
		case "not_in":
			return expression.NotIn(args[0], args[1]), nil
		case "like":
			return expression.Like(args[0], args[1]), nil
		case "not_like":
			return expression.NotLike(args[0], args[1]), nil
		case "lt":
			return expression.Lt(args[0], args[1]), nil
		case "le":
			return expression.Le(args[0], args[1]), nil
		case "is_null":
			return expression.IsNull(args[0]), nil
		case "is_not_null":
			return expression.IsNotNull(args[0]), nil

		}
	}
	return nil, fmt.Errorf("unknown operator in expression: %s", string(data))
}

func argsFromJSON(data json.RawMessage) ([]Expression, error) {
	var rawArgs []json.RawMessage
	if err := json.Unmarshal(data, &rawArgs); err != nil {
		return nil, err
	}
	args := []Expression{}
	for _, rawArg := range rawArgs {
		arg, err := expressionFromJSON(rawArg)
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
	}
	return args, nil
}

func constantFromJSON(data []byte) (Expression, error) {
	var rawConstant any
	if err := json.Unmarshal(data, &rawConstant); err != nil {
		return nil, fmt.Errorf("failed to unmarshal expression %s: %w", string(data), err)
	}

	return expression.ConstantFrom(rawConstant)
}
