package runtime

import "fmt"

func NativePrint(params Parameters) EvalValue {
	values := []any{}
	for _, param := range params {
		switch param.Type {
		case VAR_TYPE_STRING:
			values = append(values, param.Value)
		case VAR_TYPE_IDENTIFIER:
			values = append(values, param.Value.(EvalValue).Value)
		}
	}
	fmt.Println(values...)
	return EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}

func NativeIfStatement(params Parameters) EvalValue {
	if params["condition"].Value.(bool) {
		return EvalValue{Type: VAR_TYPE_BOOLEAN, Value: "true"}
	}
	return EvalValue{Type: VAR_TYPE_BOOLEAN, Value: "false"}
}
