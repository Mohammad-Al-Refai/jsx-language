package runtime

import (
	scopename "m.shebli.refaai/ht/runtime/scopeName"
)

type Scope struct {
	Name      scopename.ScopeName
	Next      *Scope
	Previous  *Scope
	Variables []*Variable
	Stack     []*EvalValue
}

func (scope *Scope) SetNext(new *Scope) *Scope {
	scope.Next = new
	new.Previous = scope
	return scope.Next
}
func (scope *Scope) Push(value *EvalValue) {
	scope.Stack = append(scope.Stack, value)
}
func (scope *Scope) Pop() *EvalValue {
	if len(scope.Stack) == 0 {
		return &EvalValue{Type: VAR_TYPE_UNDEFINED}
	}
	last := scope.Stack[len(scope.Stack)-1]
	scope.Stack = scope.Stack[:len(scope.Stack)-1]
	return last
}
func (scope *Scope) DefineVariable(variable Variable) bool {
	for _, declaration := range scope.Variables {
		if declaration.Name == variable.Name {
			return false
		}
	}
	scope.Variables = append(scope.Variables, &variable)
	return true
}

func (scope *Scope) GetVariable(name string) (bool, *Variable) {
	for _, declaration := range scope.Variables {
		if declaration.Name == name {
			return true, declaration
		}
	}
	return false, &Variable{}
}

func (scope *Scope) UpdateVariable(name string, value interface{}) (bool, *Variable) {
	for _, declaration := range scope.Variables {
		if declaration.Name == name {
			declaration.Value = value
			return true, declaration
		}
	}
	return false, &Variable{}
}
func (scope *Scope) Free() {
	scope.Stack = []*EvalValue{}
	scope.Variables = []*Variable{}
}
