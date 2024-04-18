package runtime

import (
	"m.shebli.refaai/ht/lexer"
)

func (interpreter *Interpreter) EvaluateFunctionDeclaration(openTag lexer.OpenTag, scope *Scope) *EvalValue {
	if len(openTag.Params) == 0 || openTag.Params[0].Key != "id" {
		interpreter.threwError("Missing 'id' param for function deceleration")
		return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
	}
	params := openTag.Params
	functionName := interpreter.Evaluate(params[0].Value, scope)
	if functionName.Type != VAR_TYPE_STRING {
		interpreter.threwError("Expect 'id' param value to be string")
	}
	if len(params) > 1 {
		if params[1].Key == "args" {
			for _, expr := range params[1].Value.Body.(lexer.Statement).Body.(lexer.Expression).Statements {
				p := interpreter.Evaluate(expr, scope)
				scope.DefineVariable(Variable{
					Name:      p.Value.(string),
					ValueType: p.Type,
					Value:     "undefined"})
			}
		}
	}
	function := &RuntimeFunctionCall{
		Name:  functionName.Value.(string),
		Scope: scope,
		Nodes: openTag.Children,
	}
	interpreter.Scope.DefineFunction(function)
	return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}
