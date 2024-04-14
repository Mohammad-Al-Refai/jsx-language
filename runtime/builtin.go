package runtime

import "fmt"

func GlobalScope() *Scope {
	globalScope := Scope{}
	globalScope.DefineVariable(Variable{Name: "Print", ValueType: VAR_TYPE_NATIVE_FUNCTION, Value: RuntimeFunctionCall{
		IsNative: true,
		Name:     "Print",
		Call:     NativePrint,
	}})
	// globalScope.DefineVariable(Variable{Name: "If", ValueType: VAR_TYPE_NATIVE_FUNCTION, Value: RuntimeFunctionCall{
	// 	IsNative: true,
	// 	Name:     "If",
	// 	Call:     NativeIfStatement,
	// }})
	return &globalScope
}
func NativePrint(param Parameters) EvalValue {
	fmt.Println(param["value"].Value)
	return EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}

// func NativeIfStatement(condition EvalValue) EvalValue {
// 	if condition.Value.(bool) {
// 		return EvalValue{Type: VAR_TYPE_BOOLEAN, Value: "true"}
// 	}
// 	return EvalValue{Type: VAR_TYPE_BOOLEAN, Value: "false"}
// }
