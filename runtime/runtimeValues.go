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
	IsNative bool
	Name     string
	Scope    *Scope
	Nodes    []lexer.Statement
	Call     func(Parameters) *EvalValue
}
