package runtime

import (
	"m.shebli.refaai/ht/lexer"
)

func (interpreter *Interpreter) EvaluateForLoop(openTag lexer.OpenTag, scope *Scope) *EvalValue {
	params := openTag.Params
	nodes := openTag.Children

	//Looking for 'var' param
	if len(params) == 0 || params[0].Key != "var" {
		interpreter.threwError("Expect 'var' param for 'For'")
	}
	//Looking for 'from' param
	if len(params) == 1 || params[1].Key != "from" {
		interpreter.threwError("Expect 'from' param for 'For'")
	}
	//Looking for 'to' param
	if len(params) == 2 || params[2].Key != "to" {
		interpreter.threwError("Expect 'to' param for 'For'")
	}
	newScope := &Scope{}
	result := interpreter.EvaluateParameters(params, newScope)
	varParmName := result["var"].Value.(string)
	initValue := result["from"]
	to := result["to"]
	//check 'from' param if it number
	if initValue.Type != VAR_TYPE_NUMBER {
		interpreter.threwError("Expect 'from' value to be number")
	}
	//check 'to' param if it number
	if to.Type != VAR_TYPE_NUMBER {
		interpreter.threwError("Expect 'to' value to be number")
	}
	//check 'from' param if it's greater than 'to' param
	if initValue.Value.(int) > to.Value.(int) {
		interpreter.threwError("Expect 'from' value to be less than 'to' value")
	}
	//save the var in scope variables for updating it after every iteration
	newScope.DefineVariable(Variable{
		Name:      varParmName,
		Value:     initValue.Value,
		ValueType: VAR_TYPE_NUMBER,
	})
	_, initiator := newScope.GetVariable(varParmName)
	hasBreak := false
	for {
		if initiator.Value == to.Value || hasBreak {
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
