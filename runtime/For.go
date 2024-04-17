package runtime

import "m.shebli.refaai/ht/lexer"

func (interpreter *Interpreter) EvaluateForLoop(openTag lexer.OpenTag, scope *Scope) *EvalValue {
	params := openTag.Params
	nodes := openTag.Children

	if len(params) == 0 || params[0].Key != "var" {
		interpreter.threwError("Expect 'var' param for 'For'")
	}
	if len(params) == 1 || params[1].Key != "from" {
		interpreter.threwError("Expect 'from' param for 'For'")
	}
	if len(params) == 2 || params[2].Key != "to" {
		interpreter.threwError("Expect 'to' param for 'For'")
	}

	newScope := &Scope{}
	result := interpreter.EvaluateParameters(params, newScope)

	varParmName := result["var"].Value.(string)
	initValue := result["from"].Value
	to := result["to"].Value
	newScope.DefineVariable(Variable{
		Name:      varParmName,
		Value:     initValue,
		ValueType: VAR_TYPE_NUMBER,
	})
	_, initiator := newScope.GetVariable(varParmName)
	hasBreak := false
	for {
		if initiator.Value == to || hasBreak {
			break
		}

		for _, node := range nodes {
			r := interpreter.Evaluate(node, newScope)
			if r.Value == "break" {
				hasBreak = true
				break
			}
		}
		newScope.UpdateVariable(initiator.Name, initiator.Value.(int)+1)

	}

	return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}
