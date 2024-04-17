package runtime

import "m.shebli.refaai/ht/lexer"

func (interpreter *Interpreter) EvaluateFunctionDeclaration(openTag lexer.OpenTag, scope *Scope) *EvalValue {
	if len(openTag.Params) == 0 || openTag.Params[0].Key != "id" {
		interpreter.threwError("Missing 'id' param for function deceleration")
		return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
	}
	params := openTag.Params
	functionName := interpreter.Evaluate(params[0].Value, scope).Value.(string)
	args := Parameters{}
	if len(params) > 1 {
		if params[1].Key == "args" {
			for _, expr := range params[1].Value.Body.(lexer.Statement).Body.(lexer.Expression).Statements {
				p := interpreter.Evaluate(expr, scope)
				args[p.Value.(string)] = p
				scope.DefineVariable(Variable{
					Name:      p.Value.(string),
					ValueType: p.Type,
					Value:     "undefined"})
			}
		}
	}
	function := NewRuntimeFunctionCall()
	function.Name = functionName
	function.Scope = scope
	function.Nodes = openTag.Children
	interpreter.Scope.DefineVariable(Variable{
		Name:      functionName,
		ValueType: VAR_TYPE_FUNCTION,
		Value:     function,
	})
	return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}
