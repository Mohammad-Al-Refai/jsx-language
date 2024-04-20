package runtime

import (
	"fmt"
	"os"

	"m.shebli.refaai/ht/lexer"
)

type Parameters map[string]*EvalValue

type Interpreter struct {
	Scope            *Scope
	AST              lexer.Program
	IsFinish         bool
	CurrentStatement lexer.Statement
	CurrentIndex     int
	CallStack        *CallStack
}

func (interpreter *Interpreter) next() {
	interpreter.CurrentIndex++
	if interpreter.CurrentIndex < len(interpreter.AST.Statements) {
		interpreter.CurrentStatement = interpreter.AST.Statements[interpreter.CurrentIndex]
	} else {
		interpreter.IsFinish = true
	}
}
func NewInterpreter(ast lexer.Program) *Interpreter {
	return &Interpreter{
		AST:              ast,
		Scope:            GlobalScope(),
		CurrentStatement: ast.Statements[0],
		CurrentIndex:     0,
		CallStack:        NewCallStack(),
	}
}
func (interpreter *Interpreter) threwError(message string) {
	fmt.Println(fmt.Errorf(fmt.Sprintf("[RuntimeError] %v", message)))
	os.Exit(1)
}
func (interpreter *Interpreter) Setup() {
	first := interpreter.AST.Statements[0]
	if first.Kind != lexer.K_OPEN_TAG ||
		(first.Kind == lexer.K_OPEN_TAG &&
			first.Body.(lexer.OpenTag).Name != "App") {
		interpreter.threwError("Missing <App>")
		return
	}
	if interpreter.AST.Declarations != nil || len(interpreter.AST.Declarations) != 0 {
		declares := interpreter.AST.Declarations
		statements := interpreter.AST.Statements
		newStatements := []lexer.Statement{}
		newStatements = append(newStatements, declares...)
		newStatements = append(newStatements, statements...)
		interpreter.AST.Statements = newStatements
		interpreter.CurrentStatement = interpreter.AST.Statements[0]
	}
}
func (interpreter *Interpreter) Run() {
	interpreter.Setup()
	for {
		interpreter.Evaluate(interpreter.CurrentStatement, interpreter.Scope)
		interpreter.next()
		if interpreter.IsFinish {
			return
		}
	}

}
func (interpreter *Interpreter) Evaluate(statement lexer.Statement, scope *Scope) *EvalValue {
	switch statement.Kind {
	case lexer.K_OPEN_TAG:
		return interpreter.EvaluateOpenTag(statement.Body.(lexer.OpenTag), scope)
	case lexer.K_CLOSE_TAG:
		return interpreter.EvaluateCloseTag(statement.Body.(lexer.CloseTag), scope)
	case lexer.K_PARAMETER_VALUE:
		return interpreter.Evaluate(statement.Body.(lexer.Statement), scope)
	case lexer.K_IDENTIFIER:
		return interpreter.EvaluateIdentifier(statement.Body.(string), scope, false)
	case lexer.K_ARRAY:
		return interpreter.EvaluateArray(statement.Body.(lexer.Array), scope)
	case lexer.K_OBJECT:
		return interpreter.EvaluateObjectMemberCall(statement.Body.(lexer.Object), scope)
	case lexer.K_EXPRESSION:
		return interpreter.EvaluateExpression(statement.Body.(lexer.Expression), scope)
	case lexer.K_NUMBER:
		return &EvalValue{Type: VAR_TYPE_NUMBER, Value: statement.Body.(int)}
	case lexer.K_STRING:
		return &EvalValue{Type: VAR_TYPE_STRING, Value: statement.Body.(string)}
	case lexer.K_EOF:
		return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
	default:
		println(statement.Kind.String(), " unknown")
		return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
	}
}

func (interpreter *Interpreter) EvaluateOpenTag(openTag lexer.OpenTag, scope *Scope) *EvalValue {
	newScope := &Scope{Prev: scope}
	if openTag.Name == "Function" {
		return interpreter.EvaluateFunctionDeclaration(openTag, newScope)
	}
	if openTag.Name == "If" {
		return interpreter.EvaluateIfStatement(openTag, newScope)
	}
	if openTag.Name == "For" {
		return interpreter.EvaluateForLoop(openTag, newScope)
	}
	children := openTag.Children
	for _, child := range children {
		switch child.Kind {
		case lexer.K_OPEN_TAG:
			if child.Body.(lexer.OpenTag).Name == "For" {
				interpreter.EvaluateForLoop(child.Body.(lexer.OpenTag), newScope)
				continue
			}
			if child.Body.(lexer.OpenTag).Name == "If" {
				interpreter.EvaluateIfStatement(child.Body.(lexer.OpenTag), newScope)
				continue
			}
			interpreter.threwError(fmt.Sprintf("tag '%v' is undefined", openTag.Name))
		case lexer.K_CLOSE_TAG:
			interpreter.EvaluateCloseTag(child.Body.(lexer.CloseTag), scope)
		}
	}
	return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}

func (interpreter *Interpreter) EvaluateCloseTag(closeTag lexer.CloseTag, scope *Scope) *EvalValue {
	name := closeTag.Name
	if name == "Let" {
		return interpreter.EvaluateLetDeclaration(closeTag, scope)
	}
	if name == "Break" {
		return &EvalValue{Value: "break", Type: VAR_TYPE_NATIVE_FUNCTION}
	}

	if name == "Continue" {
		return &EvalValue{Value: "continue", Type: VAR_TYPE_NATIVE_FUNCTION}
	}
	if name == "Set" {
		return interpreter.EvaluateSet(closeTag, scope)
	}
	result := interpreter.EvaluateIdentifier(name, scope, true)
	if result.Type == VAR_TYPE_NATIVE_FUNCTION {
		result := interpreter.EvaluateFunctionCall(
			result.Value.(*RuntimeFunctionCall),
			interpreter.EvaluateParameters(closeTag.Params, scope))
		return result
	}
	if result.Type == VAR_TYPE_FUNCTION {
		result := interpreter.EvaluateFunctionCall(
			result.Value.(*RuntimeFunctionCall),
			interpreter.EvaluateParameters(closeTag.Params, scope))
		return result
	}
	interpreter.threwError(fmt.Sprintf("function '%v' is undefined", name))
	return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}

func (interpreter *Interpreter) EvaluateParameters(parameters []lexer.Parameter, scope *Scope) Parameters {
	params := make(Parameters)
	for _, param := range parameters {
		params[param.Key] = interpreter.Evaluate(param.Value, scope)
	}
	return params
}
func (interpreter *Interpreter) EvaluateIdentifier(name string, scope *Scope, isFunction bool) *EvalValue {
	if scope == nil {
		interpreter.threwError(fmt.Sprintf("'%v' is not defined", name))
		return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
	}

	if isFunction {
		isFoundInCurrentScope, currentScopeFunction := interpreter.Scope.GetFunction(name)
		if isFoundInCurrentScope {
			if currentScopeFunction.IsNative {
				return &EvalValue{Value: currentScopeFunction, Type: VAR_TYPE_NATIVE_FUNCTION}
			}
			return &EvalValue{Value: currentScopeFunction, Type: VAR_TYPE_FUNCTION}
		}
	}
	isFoundInCurrentScope, currentScopeVariable := scope.GetVariable(name)
	if isFoundInCurrentScope {
		return &EvalValue{Type: currentScopeVariable.ValueType, Value: currentScopeVariable.Value}
	} else {
		return interpreter.EvaluateIdentifier(name, scope.Prev, false)
	}
}
func (interpreter *Interpreter) EvaluateFunctionCall(function *RuntimeFunctionCall, params Parameters) *EvalValue {
	newScope := &Scope{
		Variables: function.Scope.Variables,
		Stack:     function.Scope.Stack,
		Prev:      function.Scope,
	}

	function.Scope = newScope
	interpreter.ApplyParamsToFunction(function, params)
	if function.IsNative {
		return function.Call(params)
	}
	interpreter.CallStack.Push(function)
	for _, child := range function.Nodes {
		interpreter.Evaluate(child, newScope)
	}
	interpreter.CallStack.Pop()
	return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}

func (Interpreter *Interpreter) ApplyParamsToFunction(function *RuntimeFunctionCall, params Parameters) {
	newVariables := []*Variable{}
	for _, variable := range function.Scope.Variables {
		matched := params[variable.Name]
		if matched == nil {
			Interpreter.threwError(fmt.Sprintf("Expecting param '%v' for calling '%v'", variable.Name, function.Name))
		}
		if matched.Type == VAR_TYPE_UNDEFINED {
			Interpreter.threwError(fmt.Sprintf("Expected to have '%v' param for calling function '%v'", variable.Name, function.Name))
		}
		temp := Variable{}
		temp.Name = variable.Name
		temp.Value = matched.Value
		temp.ValueType = matched.Type
		newVariables = append(newVariables, &temp)
	}
	function.Scope.Variables = newVariables
}

func (interpreter *Interpreter) EvaluateCondition(param lexer.Parameter, scope *Scope) *EvalValue {
	if param.Key != "condition" {
		interpreter.threwError(fmt.Sprintf("Expect 'condition' param for if statement found '%v'", param.Key))
	}
	return interpreter.EvaluateExpression(param.Value.Body.(lexer.Statement).Body.(lexer.Expression), scope)
}

func (interpreter *Interpreter) EvaluateExpression(expr lexer.Expression, scope *Scope) *EvalValue {
	for _, ex := range expr.Statements {
		if ex.Kind == lexer.K_OPERATOR {
			scope.Stack.Push(interpreter.EvaluateOperator(ex, scope))
			continue
		}
		scope.Stack.Push(interpreter.Evaluate(ex, scope))
	}
	return scope.Stack.Pop()
}

func (interpreter *Interpreter) EvaluateArray(expr lexer.Array, scope *Scope) *EvalValue {
	array := &ArrayRuntime{}
	for _, ex := range expr.Items {
		array.Push(interpreter.Evaluate(ex, scope))
	}
	return &EvalValue{Value: array, Type: VAR_TYPE_ARRAY}
}

func (interpreter *Interpreter) EvaluateObjectMemberCall(obj lexer.Object, scope *Scope) *EvalValue {
	isFound, object := interpreter.Scope.GetObject(obj.Name)
	if !isFound {
		interpreter.threwError(fmt.Sprintf("'%v' is not defined", obj.Name))
		return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
	}
	memberName := obj.Members[0]
	isMember, member := object.GetObjectMember(memberName)
	if !isMember {
		interpreter.threwError(fmt.Sprintf("'%v' is not a member of object '%v'", memberName, obj.Name))
		return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
	}
	return member.Call(&scope.Stack)
}
