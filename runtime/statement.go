package runtime

import (
	"fmt"

	"m.shebli.refaai/ht/lexer"
)

func (interpreter *Interpreter) EvaluateCondition(param lexer.Parameter, scope *Scope) *EvalValue {
	if param.Key != "condition" {
		interpreter.threwError(fmt.Sprintf("Expect 'condition' param for if statement found '%v'", param.Key))
	}
	return interpreter.EvaluateExpression(param.Value.Body.(lexer.Statement).Body.(lexer.Expression), scope)
}

func (interpreter *Interpreter) EvaluateExpression(expr lexer.Expression, scope *Scope) *EvalValue {
	for _, ex := range expr.Statements {
		if ex.Kind == lexer.K_OPERATOR {
			scope.Push(interpreter.EvaluateOperator(ex, scope))
			continue
		}
		scope.Push(interpreter.Evaluate(ex, scope))
	}
	return scope.Pop()
}

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
