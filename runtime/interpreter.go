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
		return interpreter.EvaluateIdentifier(statement.Body.(string), scope)
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
	if openTag.Name == "Function" {
		return interpreter.EvaluateFunctionDeclaration(openTag, &Scope{})
	}
	if openTag.Name == "If" {
		return interpreter.EvaluateIfStatement(openTag, scope)
	}
	if openTag.Name == "For" {
		return interpreter.EvaluateForLoop(openTag, scope)
	}
	children := openTag.Children
	for _, child := range children {
		switch child.Kind {
		case lexer.K_OPEN_TAG:
			if child.Body.(lexer.OpenTag).Name == "For" {
				interpreter.EvaluateForLoop(child.Body.(lexer.OpenTag), scope)
				continue
			}
			if child.Body.(lexer.OpenTag).Name == "If" {
				interpreter.EvaluateIfStatement(child.Body.(lexer.OpenTag), scope)
				continue
			}
			interpreter.EvaluateOpenTag(child.Body.(lexer.OpenTag), &Scope{})
		case lexer.K_CLOSE_TAG:
			interpreter.EvaluateCloseTag(child.Body.(lexer.CloseTag), scope)
		}
	}
	return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}

func (interpreter *Interpreter) EvaluateCloseTag(closeTag lexer.CloseTag, scope *Scope) *EvalValue {
	name := closeTag.Name
	isKeyword, _ := lexer.IsKeyword(name)
	isInScope, variable := interpreter.Scope.GetVariable(name)
	if isKeyword && name == "Let" {
		return interpreter.EvaluateLetDeclaration(closeTag, scope)
	}
	if isKeyword && name == "Break" {
		return &EvalValue{Value: "break", Type: VAR_TYPE_NATIVE_FUNCTION}
	}
	if isKeyword && name == "Set" {
		return interpreter.EvaluateSet(closeTag, scope)
	}
	if isInScope && variable.ValueType == VAR_TYPE_NATIVE_FUNCTION {
		result := interpreter.EvaluateNativeFunction(
			variable.Value.(RuntimeNativeFunctionCall),
			interpreter.EvaluateParameters(closeTag.Params, scope))
		//TODO CLEAR
		return result
	}
	if isInScope && variable.ValueType == VAR_TYPE_FUNCTION {
		result := interpreter.EvaluateFunctionCall(
			variable.Value.(*RuntimeFunctionCall),
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
func (interpreter *Interpreter) EvaluateIdentifier(name string, scope *Scope) *EvalValue {
	localIsDefined, local_variable := scope.GetVariable(name)
	globalScopeIsDefined, global_variable := interpreter.Scope.GetVariable(name)
	if localIsDefined {
		return &EvalValue{Type: local_variable.ValueType, Value: local_variable.Value}
	}
	if globalScopeIsDefined {
		return &EvalValue{Type: global_variable.ValueType, Value: global_variable.Value}
	}
	interpreter.threwError(fmt.Sprintf("'%v' is undefined", name))
	return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}

func (interpreter *Interpreter) EvaluateNativeFunction(function RuntimeNativeFunctionCall, params Parameters) *EvalValue {
	return function.Call(params)
}
func (interpreter *Interpreter) EvaluateFunctionCall(function *RuntimeFunctionCall, params Parameters) *EvalValue {
	newScope := &Scope{}
	newScope.Variables = function.Scope.Variables
	newScope.Stack = function.Scope.Stack
	function.Scope = newScope
	interpreter.ApplyParamsToFunction(function, params)
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
			scope.Push(interpreter.EvaluateOperator(ex, scope))
			continue
		}
		scope.Push(interpreter.Evaluate(ex, scope))
	}
	return scope.Pop()
}
