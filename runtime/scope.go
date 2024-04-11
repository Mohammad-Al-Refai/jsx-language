package runtime

type Scope struct {
	Variables []Variable
}

func (scope *Scope) DefineVariable(variable Variable) bool {
	for _, declaration := range scope.Variables {
		if declaration.Name == variable.Name {
			return false
		}
	}
	scope.Variables = append(scope.Variables, variable)
	return true
}

func (scope *Scope) GetVariable(name string) (bool, Variable) {
	for _, declaration := range scope.Variables {
		if declaration.Name == name {
			return true, declaration
		}
	}
	return false, Variable{}
}
