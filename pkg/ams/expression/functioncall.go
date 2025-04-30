package expression

import "sync"

type FunctionCall struct {
	name string
	args []Expression
	fc   *FunctionRegistry
}

type FunctionRegistry struct {
	m    *sync.RWMutex
	impl map[string]func(Input, ...Constant) Expression
}

func NewFunctionRegistry() *FunctionRegistry {
	return &FunctionRegistry{
		m:    &sync.RWMutex{},
		impl: make(map[string]func(Input, ...Constant) Expression),
	}
}

func (o FunctionCall) Evaluate(input Input) Expression {
	c, newArgs := evaluateArgs(input, o.args)
	if len(c) < len(o.args) || o.fc == nil {
		return &FunctionCall{
			name: o.name,
			args: newArgs,
			fc:   o.fc,
		}
	}
	r := o.fc.Call(o.name, input, c)
	if r == nil {
		return &FunctionCall{
			name: o.name,
			args: newArgs,
			fc:   o.fc,
		}
	}
	return r
}

// something like this in future for user privide functions
// func (fc *FunctionRegistry) Register(name string, impl func(Input, ...Constant) Expression) {
// 	fc.m.Lock()
// 	defer fc.m.Unlock()
// 	fc.impl[name] = impl
// }

func (fc *FunctionRegistry) RegisterExpressionFunction(name string, expr Expression) {
	fc.m.Lock()
	defer fc.m.Unlock()
	fc.impl[name] = func(input Input, args ...Constant) Expression { //nolint:unparam
		return expr.Evaluate(input)
	}
}

func (fc *FunctionRegistry) Call(name string, input Input, args []Constant) Expression {
	fc.m.RLock()
	defer fc.m.RUnlock()
	if impl, ok := fc.impl[name]; ok {
		return impl(input, args...)
	}
	return nil
}
