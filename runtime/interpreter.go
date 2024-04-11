package runtime

import (
	"fmt"
	"os"

	"m.shebli.refaai/ht/lexer"
)

type Parameters map[string]EvalValue

type Interpreter struct {
	Scope            Scope
	AST              lexer.Program
	IsFinish         bool
	CurrentStatement lexer.Statement
	CurrentIndex     int
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
	globalScope := Scope{}
	globalScope.DefineVariable(Variable{Name: "Print", ValueType: VAR_TYPE_NATIVE_FUNCTION, Value: RuntimeFunctionCall{
		IsNative: true,
		Name:     "Print",
		Call:     NativePrint,
	}})
	return &Interpreter{
		AST:              ast,
		Scope:            globalScope,
		CurrentStatement: ast.Statements[0],
		CurrentIndex:     0,
	}
}
func (interpreter *Interpreter) threwError(message string) {
	fmt.Println(fmt.Errorf(fmt.Sprintf("[RuntimeError] %v", message)))
	os.Exit(1)
}
func (interpreter *Interpreter) Run() {
	for {
		interpreter.Evaluate(interpreter.CurrentStatement)
		interpreter.next()
		if interpreter.IsFinish {
			return
		}
	}

}
func (interpreter *Interpreter) Evaluate(statement lexer.Statement) EvalValue {
	switch statement.Kind {
	case lexer.K_OPEN_TAG:
		return interpreter.EvaluateOpenTag(statement.Body.(lexer.OpenTag))
	case lexer.K_CLOSE_TAG:
		return interpreter.EvaluateCloseTag(statement.Body.(lexer.CloseTag))
	case lexer.K_PARAMETER_VALUE:
		return interpreter.Evaluate(statement.Body.(lexer.Statement))
	case lexer.K_IDENTIFIER:
		return interpreter.EvaluateIdentifier(statement.Body.(string))
	case lexer.K_NUMBER:
		return EvalValue{Value: statement.Body.(int), Type: VAR_TYPE_NUMBER}
	case lexer.K_STRING:
		return EvalValue{Value: statement.Body.(string), Type: VAR_TYPE_STRING}
	case lexer.K_EOF:
		return EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
	default:
		println(statement.Kind.String(), " unknown")
		return EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
	}
}

func (interpreter *Interpreter) EvaluateOpenTag(openTag lexer.OpenTag) EvalValue {
	children := openTag.Children
	for _, child := range children {
		switch child.Kind {
		case lexer.K_CLOSE_TAG:
			interpreter.EvaluateCloseTag(child.Body.(lexer.CloseTag))
		case lexer.K_OPEN_TAG:
			interpreter.Evaluate(child)
		}

	}
	return EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}

func (interpreter *Interpreter) EvaluateCloseTag(closeTag lexer.CloseTag) EvalValue {
	name := closeTag.Name
	isKeyword, _ := lexer.IsKeyword(name)
	isInScope, variable := interpreter.Scope.GetVariable(name)
	if isKeyword && name == "Let" {
		return interpreter.EvaluateLetDeclaration(closeTag)
	}
	if isInScope && variable.ValueType == VAR_TYPE_NATIVE_FUNCTION {
		return interpreter.EvaluateNativeFunction(
			variable.Value.(RuntimeFunctionCall),
			interpreter.EvaluateLetParameters(closeTag.Params))
	}
	return EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}

func (interpreter *Interpreter) EvaluateLetDeclaration(closeTag lexer.CloseTag) EvalValue {
	params := closeTag.Params
	evaluatedParams := interpreter.EvaluateLetParameters(params)
	id, isId := evaluatedParams["id"]
	value, isValue := evaluatedParams["value"]

	if !isValue {
		interpreter.threwError("Expect 'value' param")
		return EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
	}
	if !isId {
		interpreter.threwError("Expect 'id' param")
		return EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
	}

	isOk := interpreter.Scope.DefineVariable(Variable{Name: id.Value.(string), Value: value})
	if !isOk {
		interpreter.threwError(fmt.Sprintf("%v is already declared", id))
	}
	return EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}

func (interpreter *Interpreter) EvaluateLetParameters(parameters []lexer.Parameter) Parameters {
	params := make(Parameters)
	for _, param := range parameters {
		params[param.Key] = interpreter.Evaluate(param.Value)
	}
	return params
}
func (interpreter *Interpreter) EvaluateIdentifier(name string) EvalValue {
	isDefined, variable := interpreter.Scope.GetVariable(name)
	if isDefined {
		return EvalValue{Type: variable.ValueType, Value: variable.Value}
	}
	interpreter.threwError(fmt.Sprintf("'%v' is undefined", name))
	return EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}

func (interpreter *Interpreter) EvaluateNativeFunction(function RuntimeFunctionCall, params Parameters) EvalValue {
	return function.Call(params)
}
