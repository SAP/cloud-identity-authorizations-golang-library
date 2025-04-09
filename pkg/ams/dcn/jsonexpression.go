package dcn

import (
	"encoding/json"
	"fmt"
)

type CallObject struct {
	Call []string     `json:"call"`
	Args []Expression `json:"args"`
}
type Variable struct {
	Ref []string `json:"ref"`
}

func (e Expression) MarshalJSON() ([]byte, error) {
	if e.Constant != nil {
		result, err := json.Marshal(e.Constant)
		return result, err
	}
	if len(e.Call) > 0 {
		return json.Marshal(CallObject{
			Call: e.Call,
			Args: e.Args,
		})
	}
	if len(e.Ref) > 0 {
		return json.Marshal(Variable{
			Ref: e.Ref,
		})
	}
	return nil, fmt.Errorf("unexpected expression %+v", e)
}

func (e *Expression) UnmarshalJSON(data []byte) error {

	var c CallObject
	if err := json.Unmarshal(data, &c); err == nil && len(c.Call) > 0 {
		e.Call = c.Call
		e.Args = c.Args
		return nil
	}

	var v Variable
	if err := json.Unmarshal(data, &v); err == nil && len(v.Ref) > 0 {
		e.Ref = v.Ref
		return nil
	}

	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		e.Constant = b
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		e.Constant = s
		return nil
	}

	var n float64
	if err := json.Unmarshal(data, &n); err == nil {
		e.Constant = n
		return nil
	}

	var sa []string
	if err := json.Unmarshal(data, &sa); err == nil {
		e.Constant = sa
		return nil
	}

	var na []float64
	if err := json.Unmarshal(data, &na); err == nil {
		e.Constant = na
		return nil
	}

	var ba []bool
	if err := json.Unmarshal(data, &ba); err == nil {
		e.Constant = ba
		return nil
	}

	return fmt.Errorf("not parseable expression %s", data)
}
