package expression

import "sync"

type FunctionCall struct {
	name string
	args []Expression
	fc   *FunctionContainer
}

type FunctionContainer struct {
	m    *sync.RWMutex
	impl map[string]func(Input, ...Constant) Expression
}

func NewFunctionContainer() *FunctionContainer {
	return &FunctionContainer{
		m:    &sync.RWMutex{},
		impl: make(map[string]func(Input, ...Constant) Expression),
	}
}

func (o FunctionCall) Evaluate(input Input) Expression {
	c, newArgs := evaluateArgs(input, o.args)
	if len(c) < len(o.args) {
		return &FunctionCall{
			name: o.name,
			args: newArgs,
			fc:   o.fc,
		}
	}
	return o.fc.Call(o.name, input, c)
}

func (fc *FunctionContainer) Register(name string, impl func(Input, ...Constant) Expression) {
	fc.m.Lock()
	defer fc.m.Unlock()
	fc.impl[name] = impl
}

func (fc *FunctionContainer) RegisterExpressionFunction(name string, expr Expression) {
	fc.m.Lock()
	defer fc.m.Unlock()
	fc.impl[name] = func(input Input, args ...Constant) Expression { //nolint:unparam
		return expr.Evaluate(input)
	}
}

func (fc *FunctionContainer) Call(name string, input Input, args []Constant) Expression {
	fc.m.RLock()
	defer fc.m.RUnlock()
	if impl, ok := fc.impl[name]; ok {
		return impl(input, args...)
	}
	return FALSE
}
