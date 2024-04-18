package runtime

import (
	"fmt"

	"m.shebli.refaai/ht/lexer"
)

func (interpreter *Interpreter) EvaluateSet(closeTag lexer.CloseTag, scope *Scope) *EvalValue {
	params := closeTag.Params
	evaluatedParams := interpreter.EvaluateParameters(params, scope)
	id, isId := evaluatedParams["id"]
	to, isTo := evaluatedParams["to"]
	if !isTo {
		interpreter.threwError("Expect 'to' param")
		return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
	}
	if !isId {
		interpreter.threwError("Expect 'id' param")
		return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
	}
	isOk, _ := scope.UpdateVariable(id.Value.(string), to.Value)
	if !isOk {
		isOk, _ := interpreter.Scope.UpdateVariable(id.Value.(string), to.Value)
		if !isOk {
			interpreter.threwError(fmt.Sprintf("%v is undeclared", id))
		}
	}
	return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}
