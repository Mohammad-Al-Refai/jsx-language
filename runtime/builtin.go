package runtime

import (
	"fmt"
)

func GlobalScope() *Scope {
	globalScope := Scope{}
	globalScope.DefineFunction(Print(Parameters{"value": &EvalValue{}}))
	globalScope.DefineVariable(Variable{Name: "true", ValueType: VAR_TYPE_BOOLEAN, Value: true})
	globalScope.DefineVariable(Variable{Name: "false", ValueType: VAR_TYPE_BOOLEAN, Value: false})
	ApplyArray(globalScope)
	return &globalScope
}
func Print(param Parameters) *RuntimeFunctionCall {
	return &RuntimeFunctionCall{
		IsNative: true,
		Name:     "Print",
		Scope:    &Scope{},
		Call: func(p Parameters) *EvalValue {
			value := p["value"]
			fmt.Println(value.Value)
			return &EvalValue{}
		},
	}
}

func ApplyArray(scope Scope) {
	members := []RuntimeObjectMember{{
		Name: "length",
		Call: func(p Parameters) *EvalValue {
			arg := p["value"]
			if arg.Type == VAR_TYPE_ARRAY {
				return &EvalValue{Value: arg.Value.(ArrayRuntime).Size, Type: VAR_TYPE_NUMBER}
			}
			return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
		},
	}}
	scope.DefineObject(&RuntimeObject{
		Name:    "array",
		Members: members,
	})
}
