package runtime

import "fmt"

func GlobalScope() *Scope {
	globalScope := Scope{}
	globalScope.DefineFunction(Print(Parameters{"value": &EvalValue{}}))
	globalScope.DefineVariable(Variable{Name: "true", ValueType: VAR_TYPE_BOOLEAN, Value: true})
	globalScope.DefineVariable(Variable{Name: "false", ValueType: VAR_TYPE_BOOLEAN, Value: false})
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
