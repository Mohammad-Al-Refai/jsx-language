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
	Last         string
	Program      Program
}

func NewAST(tokens []Tokenized) *AST {
	return &AST{
		Tokens:       tokens,
		CurrentToken: tokens[0],
		IsEnd:        false,
		CurrentIndex: 0,
		Program:      Program{},
	}
}

type Parameter struct {
	Key   string    `json:"key"`
	Value Statement `json:"value"`
}
type OpenTag struct {
	Name     string      `json:"name"`
	Params   []Parameter `json:"params"`
	Children []Statement `json:"children"`
}
type Array struct {
	Items []Statement `json:"Items"`
}
type Object struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}
type CloseTag struct {
	Name   string      `json:"name"`
	Params []Parameter `json:"params"`
}
type Expression struct {
	Statements []Statement
}
type Statement struct {
	Kind StatementKind `json:"kind"`
	Body interface{}   `json:"body"`
}

type Program struct {
	Declarations []Statement `json:"declarations"`
	Statements   []Statement `json:"statements"`
}

func (ast *AST) expect(token Token, message string) {
	ast.next()
	if ast.CurrentToken.Token == token {
		return
	}
	ast.threwError(fmt.Sprintf("Expect %v in %v", token, message))
}
func (ast *AST) expectKeyWordOrAny(message string) {
	ast.next()
	ok, _ := IsKeyword(ast.CurrentToken.Literal)
	if ast.CurrentToken.Token == IDENT || ok {
		return
	}
	ast.threwError(fmt.Sprintf("Expect %v in %v", ast.CurrentToken.Literal, message))
}
func (ast *AST) threwError(message string) {
	fmt.Println(fmt.Errorf(fmt.Sprintf("[ParseError][%v] %v got '%v' at [%v:%v]",
		ast.Last,
		message,
		ast.CurrentToken.Literal,
		ast.CurrentToken.Pos.Line,
		ast.CurrentToken.Pos.Column)))
	os.Exit(1)
}
func (ast *AST) threwIllegalError(message string) {
	fmt.Printf("[ParserError Illegal]: %+v\n", message)
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
	for {
		if ast.IsEnd {
			return ast.Program
		}
		ast.Program.Statements = append(ast.Program.Statements, ast.Parse())
		ast.next()
	}
}

func (ast *AST) Parse() Statement {
	ast.Last = "Parse"
	switch ast.CurrentToken.Token {
	case OPEN_TAG:
		return ast.ParseOpenTag()
	case CLOSE_TAG:
		return ast.ParseCloseTag()
	case EOF:
		return Statement{Kind: EOF}
	default:
		ast.threwError(fmt.Sprintf("Invalid expression '%v' ", ast.CurrentToken))
		return Statement{}
	}

}

func (ast *AST) ParseOpenTag() Statement {
	ast.Last = "ParseOpenTag"
	statement := Statement{Kind: K_OPEN_TAG}
	children := []Statement{}
	isNotFoundClose := true
	ast.expectKeyWordOrAny("Tag name in OpenTag")
	openTag := OpenTag{}
	openTag.Name = ast.CurrentToken.Literal
	openTag.Params = ast.ParseParameter()
	ast.next()
	ast.next()
	for isNotFoundClose {
		newNode := ast.Parse()
		if ast.CurrentToken.Token == EOF {
			ast.threwError(fmt.Sprintf("Expect </ %v >", openTag.Name))
		}
		if newNode.Kind == K_OPEN_TAG && openTag.Name == "App" && newNode.Body.(OpenTag).Name == "Function" {
			ast.Program.Declarations = append(ast.Program.Declarations, newNode)
			ast.next()
			continue
		}
		if newNode.Kind == K_OPEN_TAG && openTag.Name != "App" && newNode.Body.(OpenTag).Name == "Function" {
			ast.threwIllegalError("Can't declare function outside App scope!")
			return statement
		}
		if newNode.Kind == K_OPEN_TAG {
			children = append(children, newNode)
			ast.next()
			continue
		}
		if newNode.Kind == K_CLOSE_TAG && newNode.Body.(CloseTag).Name != openTag.Name {
			children = append(children, newNode)
			ast.next()
		} else {
			isNotFoundClose = false
		}

	}
	openTag.Children = children
	statement.Body = openTag
	return statement

}
func (ast *AST) ParseCloseTag() Statement {
	ast.Last = "ParseCloseTag"
	statement := Statement{Kind: K_CLOSE_TAG}
	ast.expectKeyWordOrAny("Tag name in CloseTag")
	closeTag := CloseTag{}
	closeTag.Name = ast.CurrentToken.Literal
	closeTag.Params = ast.ParseParameter()
	statement.Body = closeTag
	ast.next()
	return statement
}

func (ast *AST) ParseParameter() []Parameter {
	ast.Last = "ParseParameter"
	params := []Parameter{}
	for ast.checkForward().Token != CLOSE_OPEN_TAG {
		ast.expect(IDENT, "for parameter")
		param := Parameter{}
		param.Key = ast.CurrentToken.Literal
		ast.expect(EQUAL, "after key")
		param.Value = ast.ParseParameterValue()
		params = append(params, param)
	}

	return params
}
func (ast *AST) ParseParameterValue() Statement {
	ast.Last = "ParseParameterValue"
	statement := Statement{Kind: K_PARAMETER_VALUE}
	ast.expect(LBRACE, "after '='")
	if ast.CurrentToken.Token == LBRACE {
		ast.next()
		statement.Body = ast.ParseParameterValueExpr()
	}
	if ast.CurrentToken.Token != RBRACE {
		ast.threwError("Expect } at the end of parameter")
	}
	return statement
}

func (ast *AST) ParseParameterValueExpr() Statement {
	ast.Last = "ParseParameterValueExpr"
	stmts := []Statement{}
	stmt := Statement{Kind: K_EXPRESSION}
	if ast.CurrentToken.Token == LBRACK {
		ast.next()
		stmt = ast.ParseArray()
		return stmt
	}
	for ast.CurrentToken.Token != RBRACE {
		stmts = append(stmts, ast.ParseExpr())
		ast.next()
		if ast.CurrentToken.Token == COMMA {
			ast.next()
		}
	}
	stmt.Body = Expression{Statements: stmts}
	return stmt
}
func (ast *AST) ParseArray() Statement {
	ast.Last = "ParseArray"
	statement := Statement{Kind: K_ARRAY}
	items := []Statement{}
	for ast.CurrentToken.Token != RBRACK {
		items = append(items, ast.ParseExpr())
		ast.next()
		if ast.CurrentToken.Token == COMMA {
			ast.next()
		}
	}
	ast.next()
	statement.Body = Array{
		Items: items,
	}
	return statement
}

// Add checks for ()
func (ast *AST) ParseObject() Statement {
	ast.Last = "ParseObject"
	statement := Statement{Kind: K_OBJECT}
	objName := ast.CurrentToken.Literal
	ast.next()
	members := []string{}

	for ast.CurrentToken.Token != RPAREN {
		if ast.CurrentToken.Token == RBRACE {
			ast.threwError("Expect ')'")
		}
		if ast.CurrentToken.Token == DOT {
			ast.next()
		}
		if ast.CurrentToken.Token == IDENT {
			members = append(members, ast.CurrentToken.Literal)
			ast.next()
		}
		if ast.CurrentToken.Token == LPAREN {
			ast.next()
		}
		if ast.CurrentToken.Token == RPAREN {
			break
		}
		ast.next()
	}
	statement.Body = Object{
		Name:    objName,
		Members: members,
	}
	return statement
}
func (ast *AST) ParseExpr() Statement {
	ast.Last = "ParseExpr"
	token := ast.CurrentToken.Token
	switch token {
	case IDENT:
		return ast.ParseIdentifier()
	case INT:
		return ast.ParseInt()
	case STRING:
		return ast.ParseString()
	default:
		if isOperator(token) {
			return ast.ParseOperator()
		}
		ast.threwError(fmt.Sprintf("Invalid expression '%v'", token))
	}
	return Statement{}
}
func (ast *AST) ParseInt() Statement {
	stmt := Statement{}
	stmt.Kind = K_NUMBER
	n, err := strconv.Atoi(ast.CurrentToken.Literal)
	if err != nil {
		panic(err)
	}
	stmt.Body = n
	return stmt
}
func (ast *AST) ParseString() Statement {
	stmt := Statement{}
	stmt.Kind = K_STRING
	stmt.Body = ast.CurrentToken.Literal
	return stmt
}
func (ast *AST) ParseIdentifier() Statement {
	if ast.checkForward().Token == DOT {
		return ast.ParseObject()
	}
	stmt := Statement{}
	stmt.Kind = K_IDENTIFIER
	stmt.Body = ast.CurrentToken.Literal
	return stmt
}
func (ast *AST) ParseOperator() Statement {
	stmt := Statement{}
	stmt.Kind = K_OPERATOR
	stmt.Body = ast.CurrentToken.Literal
	return stmt
}
