package runtime

import "fmt"

type Scope struct {
	Variables []*Variable
	Functions []*RuntimeFunctionCall
	Stack     []*EvalValue
	Prev      *Scope
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
	vars := []string{}
	for _, x := range scope.Variables {
		vars = append(vars, x.Name)
	}
	fmt.Printf("Variables: %v\n--------------\n", vars)
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
