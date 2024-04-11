package runtime

import (
	"fmt"
	"os"

	"m.shebli.refaai/ht/lexer"
)

type Parameters map[string]EvalValue

type Interpreter struct {
	Scope Scope
	AST   lexer.Program
}

func NewInterpreter(ast lexer.Program) *Interpreter {
	return &Interpreter{
		AST:   ast,
		Scope: Scope{},
	}
}
func (interpreter *Interpreter) threwError(message string) {
	fmt.Println(fmt.Errorf(fmt.Sprintf("[RuntimeError] %v", message)))
	os.Exit(1)
}
func (interpreter *Interpreter) Run() {
	for _, stm := range interpreter.AST.Statements {
		interpreter.Evaluate(stm)
	}
}
func (interpreter *Interpreter) Evaluate(statement lexer.Statement) EvalValue {
	println(statement.Kind.String())
	switch statement.Kind {
	case lexer.K_OPEN_TAG:
		return interpreter.EvaluateOpenTag(statement.Body.(lexer.OpenTag))
	case lexer.K_CLOSE_TAG:
		return interpreter.EvaluateCloseTag(statement.Body.(lexer.CloseTag))
	case lexer.K_PARAMETER_VALUE:
		return interpreter.Evaluate(statement.Body.(lexer.Statement))
	case lexer.K_IDENTIFIER:
		return EvalValue{Value: statement.Body.(string), Type: VAR_TYPE_STRING}
	case lexer.K_NUMBER:
		return EvalValue{Value: statement.Body.(int), Type: VAR_TYPE_NUMBER}
	case lexer.K_STRING:
		return EvalValue{Value: statement.Body.(string), Type: VAR_TYPE_STRING}
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
			return interpreter.EvaluateCloseTag(child.Body.(lexer.CloseTag))
		}
	}
	return EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}

func (interpreter *Interpreter) EvaluateCloseTag(closeTag lexer.CloseTag) EvalValue {
	name := closeTag.Name
	isKeyword, _ := lexer.IsKeyword(name)
	if isKeyword && name == "Let" {
		interpreter.EvaluateLetDeclaration(closeTag)
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
	fmt.Printf("Declare %+v\n", evaluatedParams)
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
