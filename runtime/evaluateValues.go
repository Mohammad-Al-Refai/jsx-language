package runtime

import "m.shebli.refaai/ht/lexer"

type EvalValue struct {
	Value interface{}
	Type  VarType
}

func (ev *EvalValue) Is(t VarType) bool {
	return ev.Type == t
}
func (ev *EvalValue) ExpectAnyOf(t []VarType) bool {
	for _, ty := range t {
		if ev.Type == ty {
			return true
		}
	}
	return false
}

type RuntimeFunction struct {
	Name  string
	Scope Scope
	Nodes []EvalValue
}
type RuntimeFunctionCall struct {
	Name     string
	Params   Parameters
	IsNative bool
	Call     func(Parameters) EvalValue
	Children []lexer.Statement
}
type RuntimeIfStatement struct {
	Condition EvalValue
	Scope     Scope
	Execute   func()
}
