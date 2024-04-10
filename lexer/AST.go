package lexer

import (
	"fmt"
	"os"
	"strconv"
)

type AST struct {
	Tokens       []Tokenized
	CurrentToken Tokenized
	CurrentIndex int
	IsEnd        bool
}

func NewAST(tokens []Tokenized) *AST {
	return &AST{
		Tokens:       tokens,
		CurrentToken: tokens[0],
		IsEnd:        false,
		CurrentIndex: 0,
	}
}

type Parameter struct {
	Key   string    `json:"key"`
	Value Statement `json:"value"`
}
type OpenTag struct {
	Name   string      `json:"name"`
	Params []Parameter `json:"params"`
}
type CloseTag struct {
	Name string `json:"name"`
}
type Statement struct {
	Kind StatementKind `json:"kind"`
	Body interface{}   `json:"body"`
}

type Program struct {
	Statements []Statement `json:"statements"`
}

func (ast *AST) expect(token Token, message string) {
	ast.next()
	if ast.CurrentToken.Token == token {
		return
	}
	ast.threwError(fmt.Sprintf("Expect %v in %v", token, message))
}
func (ast *AST) threwError(message string) {
	fmt.Println(fmt.Errorf(fmt.Sprintf("[ParserError] %v", message)))
	os.Exit(1)
}
func (ast *AST) next() {
	ast.CurrentIndex++
	if ast.CurrentIndex < len(ast.Tokens) {
		ast.CurrentToken = ast.Tokens[ast.CurrentIndex]
	} else {
		ast.IsEnd = true
	}
}
func (ast *AST) checkForward() Tokenized {
	return ast.Tokens[ast.CurrentIndex+1]
}
func (ast *AST) ProduceAST() Program {
	program := Program{}
	for {
		println(ast.CurrentToken.Token)
		if ast.IsEnd {
			return program
		}
		program.Statements = append(program.Statements, ast.ParseExpression())
		ast.next()
	}
}

func (ast *AST) ParseExpression() Statement {
	switch ast.CurrentToken.Token {
	case OPEN_TAG:
		return ast.ParseOpenTag()
	case CLOSE_OPEN_TAG:
		return ast.ParseCloseTag()
	}
	return Statement{Kind: EOF}
}

func (ast *AST) ParseOpenTag() Statement {
	statement := Statement{Kind: K_OPEN_TAG}
	ast.expect(IDENT, "Tag name in OpenTag")
	openTag := OpenTag{}
	openTag.Name = ast.CurrentToken.Literal
	openTag.Params = ast.ParseParameter()
	statement.Body = openTag
	return statement
}
func (ast *AST) ParseCloseTag() Statement {
	statement := Statement{Kind: K_CLOSE_TAG}
	ast.expect(IDENT, "Tag name in CloseTag")
	closeTag := CloseTag{}
	closeTag.Name = ast.CurrentToken.Literal
	ast.expect(CLOSE_TAG, "CloseTag")
	statement.Body = closeTag
	return statement
}

func (ast *AST) ParseParameter() []Parameter {
	params := []Parameter{}
	for {
		println("ParseParameter ", ast.CurrentToken.Literal)

		if ast.CurrentToken.Token == CLOSE_TAG {
			return params
		}
		if ast.CurrentToken.Token == CLOSE_OPEN_TAG {
			return []Parameter{}
		}
		if ast.checkForward().Token == CLOSE_TAG {
			ast.next()
			return params
		}
		ast.expect(IDENT, "for parameter")
		param := Parameter{}
		param.Key = ast.CurrentToken.Literal
		ast.expect(EQUAL, "after key")
		param.Value = ast.ParseParameterValue()
		if ast.CurrentToken.Token == RBRACE {
			params = append(params, param)
		}
	}
}
func (ast *AST) ParseParameterValue() Statement {
	statement := Statement{Kind: K_PARAMETER_VALUE}
	ast.expect(LBRACE, "after '='")
	if ast.CurrentToken.Token == LBRACE {
		ast.next()
		statement.Body = ast.ParsePrimaryExpression()
	}
	return statement

}

func (ast *AST) ParsePrimaryExpression() Statement {
	stmt := Statement{}
	token := ast.CurrentToken.Token
	switch token {
	case IDENT:
		stmt.Kind = K_IDENTIFIER
		stmt.Body = ast.CurrentToken.Literal
		ast.next()
		return stmt
	case INT:
		stmt.Kind = K_NUMBER
		n, err := strconv.Atoi(ast.CurrentToken.Literal)
		if err != nil {
			panic(err)
		}
		stmt.Body = n
		return stmt
	case STRING:
		stmt.Kind = K_STRING
		stmt.Body = ast.CurrentToken.Literal
		ast.next()
	default:
		ast.threwError(fmt.Sprintf("Invalid expression '%v' expect 'identifier' or 'number' at %v:%v", token, ast.CurrentToken.Pos.Line, ast.CurrentToken.Pos.Column))
	}

	return stmt
}
