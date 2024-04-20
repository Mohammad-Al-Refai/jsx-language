package runtime

import (
	"m.shebli.refaai/ht/lexer"
)

func (interpreter *Interpreter) EvaluateIfStatement(openTag lexer.OpenTag, scope *Scope) *EvalValue {
	params := openTag.Params
	nodes := openTag.Children
	newScope := &Scope{
		Prev: scope,
	}
	if len(params) == 0 || params[0].Key != "condition" {
		interpreter.threwError("Expect 'condition' param for if statement")
	}
	result := interpreter.EvaluateCondition(params[0], newScope)
	hasBreak := false
	hasContinue := false
	if result.Value == true {
		for _, node := range nodes {
			r := interpreter.Evaluate(node, newScope)
			if r.Value == "break" {
				hasBreak = true
			}
			if r.Value == "continue" {
				hasContinue = true
			}
		}
	}
	if hasBreak {
		return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "break"}
	}
	if hasContinue {
		return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "continue"}
	}
	return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}
