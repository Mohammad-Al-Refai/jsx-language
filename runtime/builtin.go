package runtime

import (
	"fmt"

	scopename "m.shebli.refaai/ht/runtime/scopeName"
)

func GlobalScope() *Scope {
	globalScope := Scope{Name: scopename.APP}
	globalScope.DefineVariable(Variable{Name: "Print", ValueType: VAR_TYPE_NATIVE_FUNCTION, Value: RuntimeNativeFunctionCall{
		IsNative: true,
		Name:     "Print",
		Call:     NativePrint,
	}})
	globalScope.DefineVariable(Variable{Name: "true", ValueType: VAR_TYPE_BOOLEAN, Value: true})
	globalScope.DefineVariable(Variable{Name: "false", ValueType: VAR_TYPE_BOOLEAN, Value: false})
	return &globalScope
}
func NativePrint(param Parameters) *EvalValue {
	fmt.Println(param["value"].Value)
	return NewUndefinedValue()
}
