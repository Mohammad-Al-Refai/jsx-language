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
func NewInterpreter(ast lexer.Program){
	
}
func (interpreter *Interpreter) threwError(message string) {
	fmt.Println(fmt.Errorf(fmt.Sprintf("[RuntimeError] %v", message)))
	os.Exit(1)
}
func (interpreter *Interpreter) Evaluate(statement lexer.Statement) EvalValue {
	switch statement.Kind {
	case lexer.K_OPEN_TAG:
		return interpreter.EvaluateOpenTag(statement.Body.(lexer.OpenTag))
	case lexer.K_IDENTIFIER:
		return EvalValue{Value: statement.Body.(string), Type: VAR_TYPE_STRING}
	default:
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
	isExist := interpreter.Scope.DefineVariable(Variable{Name: id.Value.(string), Value: value})
	if isExist {
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
