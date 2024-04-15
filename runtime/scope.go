package runtime

type Scope struct {
	Variables []*Variable
	Stack     []*EvalValue
}

func (scope *Scope) Push(value *EvalValue) {
	// fmt.Printf("PUSH %+v\n", value)
	scope.Stack = append(scope.Stack, value)
	// fmt.Printf("NEW STACK %+v\n", scope.Stack)

}
func (scope *Scope) Pop() *EvalValue {
	if len(scope.Stack) == 0 {
		return &EvalValue{Type: VAR_TYPE_UNDEFINED}
	}
	last := scope.Stack[len(scope.Stack)-1]
	scope.Stack = scope.Stack[:len(scope.Stack)-1]
	// fmt.Printf("POP %+v\n", last)

	// fmt.Printf("NEW STACK %+v\n", scope.Stack)
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
