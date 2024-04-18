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
	Members  []*RuntimeObjectMember
}

func (obj *RuntimeObject) GetObjectMember(name string) (bool, *RuntimeObjectMember) {
	for _, member := range obj.Members {
		if member.Name == name {
			return true, member
		}
	}
	return false, &RuntimeObjectMember{}
}

type RuntimeObjectMember struct {
	Name string
	Call func(*ScopeStack) *EvalValue
}

type ArrayRuntime struct {
	Size  int
	Items []*EvalValue
}

func (a *ArrayRuntime) Push(value *EvalValue) {
	a.Items = append(a.Items, value)
	a.Size = len(a.Items)
}

func (a *ArrayRuntime) Pop() *EvalValue {
	last := a.Items[len(a.Items)-1]
	a.Items = a.Items[:len(a.Items)-1]
	a.Size = len(a.Items)
	return last
}
func (a *ArrayRuntime) At(index int) *EvalValue {
	if index > a.Size-1 {
		return &EvalValue{Value: "undefined", Type: VAR_TYPE_UNDEFINED}
	}
	return a.Items[index]
}
