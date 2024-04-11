package runtime

import "m.shebli.refaai/ht/lexer"

type Interpreter struct {
	Scope Scope
	AST   lexer.Program
}

func (interpreter *Interpreter) Evaluate(statement lexer.Statement) {
	switch statement.Kind{
	case lexer.K_OPEN_TAG:

	}
}

func (interpreter *Interpreter) EvaluateOpenTag(openTag lexer.OpenTag) {
	
}

