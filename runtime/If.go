package runtime

import "m.shebli.refaai/ht/lexer"

func (interpreter *Interpreter) EvaluateIfStatement(openTag lexer.OpenTag, scope *Scope) *EvalValue {
	params := openTag.Params
	nodes := openTag.Children
	if len(params) == 0 || params[0].Key != "condition" {
		interpreter.threwError("Expect 'condition' param for if statement")
	}
	result := interpreter.EvaluateCondition(params[0], scope)
	hasBreak := false
	if result.Value == true {
		for _, node := range nodes {
			r := interpreter.Evaluate(node, scope)
			if r.Value == "break" {
				hasBreak = true
			}
		}
	}
	if hasBreak {
		return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "break"}
	}
	return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}
