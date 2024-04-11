package runtime

import "fmt"

func NativePrint(params Parameters) EvalValue {
	values := []any{}
	for _, param := range params {
		values = append(values, param.Value.(EvalValue).Value)
	}

	fmt.Println(values...)
	return EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}
