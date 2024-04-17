package runtime

import (
	"fmt"

	"m.shebli.refaai/ht/lexer"
)

func (interpreter *Interpreter) EvaluateLetDeclaration(closeTag lexer.CloseTag, scope *Scope) *EvalValue {
	params := closeTag.Params
	evaluatedParams := interpreter.EvaluateParameters(params, scope)
	id, isId := evaluatedParams["id"]
	value, isValue := evaluatedParams["value"]

	if !isValue {
		interpreter.threwError("Expect 'value' param")
		return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
	}
	if !isId {
		interpreter.threwError("Expect 'id' param")
		return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
	}

	isOk := interpreter.Scope.DefineVariable(Variable{Name: id.Value.(string), Value: value.Value, ValueType: value.Type})
	if !isOk {
		interpreter.threwError(fmt.Sprintf("%v is already declared", id))
	}
	return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}
