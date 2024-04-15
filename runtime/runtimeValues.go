package runtime

import "m.shebli.refaai/ht/lexer"

type EvalValue struct {
	Value interface{}
	Type  VarType
}

func (ev *EvalValue) IsNumber() bool {
	return ev.Type == VAR_TYPE_NUMBER
}
func (ev *EvalValue) IsString() bool {
	return ev.Type == VAR_TYPE_STRING
}
func (ev *EvalValue) IsBoolean() bool {
	return ev.Type == VAR_TYPE_BOOLEAN
}
func (ev *EvalValue) ExpectAnyOf(t []VarType) bool {
	for _, ty := range t {
		if ev.Type == ty {
			return true
		}
	}
	return false
}

type RuntimeFunctionCall struct {
	Name       string
	Parameters Parameters
	Scope      *Scope
	Nodes      []lexer.Statement
}

func NewRuntimeFunctionCall() *RuntimeFunctionCall {
	return &RuntimeFunctionCall{
		Scope:      &Scope{},
		Nodes:      []lexer.Statement{},
		Name:       "",
		Parameters: make(Parameters),
	}
}

type RuntimeNativeFunctionCall struct {
	Name     string
	Params   Parameters
	IsNative bool
	Call     func(Parameters) *EvalValue
	Children []lexer.Statement
}
