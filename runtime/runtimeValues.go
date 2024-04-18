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
func (ev *EvalValue) IsArray() bool {
	return ev.Type == VAR_TYPE_ARRAY
}
func (ev *EvalValue) IsObject() bool {
	return ev.Type == VAR_TYPE_OBJECT
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
type RuntimeObject struct {
	IsNative bool
	Name     string
	Members  []RuntimeObjectMember
}
type RuntimeObjectMember struct {
	Name string
	Call func(Parameters) *EvalValue
}

type ArrayRuntime struct {
	Size  int
	Items []*EvalValue
}

func (a *ArrayRuntime) Push(value *EvalValue) {
	a.Items = append(a.Items, value)
}

func (a *ArrayRuntime) Pop() *EvalValue {
	last := a.Items[len(a.Items)-1]
	a.Items = a.Items[:len(a.Items)-1]
	return last
}
