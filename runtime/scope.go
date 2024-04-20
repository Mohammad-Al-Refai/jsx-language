package runtime

import "fmt"

type ScopeStack struct {
	Stack []*EvalValue
}
type Scope struct {
	Variables []*Variable
	Functions []*RuntimeFunctionCall
	Objects   []*RuntimeObject
	Stack     ScopeStack
	Prev      *Scope
}

func (scope *ScopeStack) Push(value *EvalValue) {
	// fmt.Printf("Push: %v\n", value)
	scope.Stack = append(scope.Stack, value)
}
func (scope *ScopeStack) Pop() *EvalValue {
	if len(scope.Stack) == 0 {
		return &EvalValue{Type: VAR_TYPE_UNDEFINED}
	}
	last := scope.Stack[len(scope.Stack)-1]
	scope.Stack = scope.Stack[:len(scope.Stack)-1]
	// fmt.Printf("Pop: %v Length: %v\n", last, len(scope.Stack))
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
func (scope *Scope) DefineObject(obj *RuntimeObject) bool {
	for _, declaration := range scope.Objects {
		if declaration.Name == obj.Name {
			return false
		}
	}
	scope.Objects = append(scope.Objects, obj)
	return true
}
func (scope *Scope) DefineFunction(function *RuntimeFunctionCall) bool {
	for _, declaration := range scope.Functions {
		if declaration.Name == function.Name {
			return false
		}
	}
	scope.Functions = append(scope.Functions, function)
	return true
}
func (scope *Scope) Debug() {
	vars := []*EvalValue{}
	vars = append(vars, scope.Stack.Stack...)
	fmt.Printf("Stack: %v\n--------------\n", vars)
}
func (scope *Scope) GetVariable(name string) (bool, *Variable) {
	for _, declaration := range scope.Variables {
		if declaration.Name == name {
			return true, declaration
		}
	}
	return false, &Variable{}
}
func (scope *Scope) GetFunction(name string) (bool, *RuntimeFunctionCall) {
	for _, declaration := range scope.Functions {
		if declaration.Name == name {
			return true, declaration
		}
	}
	return false, &RuntimeFunctionCall{}
}
func (scope *Scope) GetObject(name string) (bool, *RuntimeObject) {
	for _, declaration := range scope.Objects {
		if declaration.Name == name {
			return true, declaration
		}
	}
	return false, &RuntimeObject{}
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
	scope.Stack = ScopeStack{}
	scope.Variables = []*Variable{}
}
