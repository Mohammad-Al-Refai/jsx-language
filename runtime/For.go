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
	forScope := &Scope{Prev: scope}
	varParmName := params[0].Value.Body.(lexer.Statement).Body.(lexer.Expression).Statements[0].Body.(string)
	initValue := interpreter.EvaluateExpression(params[1].Value.Body.(lexer.Statement).Body.(lexer.Expression), forScope)
	to := interpreter.EvaluateExpression(params[2].Value.Body.(lexer.Statement).Body.(lexer.Expression), forScope)

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
	forScope.DefineVariable(Variable{
		Name:      varParmName,
		Value:     initValue.Value,
		ValueType: VAR_TYPE_NUMBER,
	})
	_, initiator := forScope.GetVariable(varParmName)
	hasBreak := false

	for {
		nodesScope := &Scope{
			Prev:      forScope,
			Variables: []*Variable{initiator}}
		for _, node := range nodes {

			r := interpreter.Evaluate(node, nodesScope)

			if r.Value == "continue" {
				if initiator.Value.(int) < to.Value.(int) || hasBreak {
					forScope.UpdateVariable(initiator.Name, initiator.Value.(int)+1)
					nodesScope.Free()
				}
				continue
			}
			if r.Value == "break" {
				hasBreak = true
				break
			}
		}
		if initiator.Value.(int) == to.Value.(int) || hasBreak {
			break
		}
		forScope.UpdateVariable(initiator.Name, initiator.Value.(int)+1)
		nodesScope.Free()
	}

	return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}
