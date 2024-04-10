package lexer

import (
	"fmt"
	"os"
	"strconv"
)

type Statement struct {
	Kind StatementKind `json:"kind"`
	Body interface{}   `json:"body"`
}

type LetDeclaration struct {
	Id   string    `json:"id"`
	Init Statement `json:"init"`
}
type FunCall struct {
	Caller string      `json:"caller"`
	Args   []Statement `json:"args"`
}
type BinaryExpr struct {
	Left     Statement `json:"left"`
	Right    Statement `json:"right"`
	Operator Token     `json:"operator"`
}
type FunctionDeclaration struct {
	Id     string
	Params []Statement
	Body   []Statement
}
type FunParam struct {
	Name string
	Type Tokenized
}
type Program struct {
	Type       StatementKind `json:"type,string"`
	Statements []Statement   `json:"statements"`
}
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
		CurrentIndex: 0,
		IsEnd:        false,
	}
}
func (ast *AST) threwError(message string) {
	fmt.Println(fmt.Errorf(fmt.Sprintf("[ParserError] %v", message)))
	os.Exit(1)
}
func (ast *AST) expect(kind StatementKind, token Token) {
	if token != ast.CurrentToken.Token {
		ast.threwError(fmt.Sprintf("[%v] expect '%v' got '%v' at %v:%v", kind, token, ast.CurrentToken.Token, ast.CurrentToken.Pos.Line, ast.CurrentToken.Pos.Column))
	}
}
func (ast *AST) expectOneOf(kind StatementKind, token ...Token) {
	for _, t := range token {
		if ast.CurrentToken.Token == t {
			return
		}
	}
	ast.threwError(fmt.Sprintf("[%v] expect '%v' got '%v' at %v:%v", kind, token, ast.CurrentToken.Token, ast.CurrentToken.Pos.Line, ast.CurrentToken.Pos.Column))
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
	program := Program{Type: PROGRAM}
	for !ast.IsEnd {
		stmt := ast.Parse()
		program.Statements = append(program.Statements, stmt)
		ast.next()
	}
	return program
}
func (ast *AST) Parse() Statement {
	switch ast.CurrentToken.Token {
	case LET:
		return ast.ParseLet()
	case INT:
		return ast.ParseExpr()
	case IDENT:
		return ast.ParsePrimaryExpr()
	case FUN:
		return ast.ParseFun()
	case EOF:
		ast.IsEnd = true
	default:
		ast.threwError(fmt.Sprintf("Unsupported expiration '%+v'", ast.CurrentToken))
	}
	return Statement{Kind: K_END_OF_FILE}
}
func (ast *AST) ParseLet() Statement {
	stmt := Statement{Kind: StatementKind(K_LET_DECLARATION)}
	decl := LetDeclaration{}
	ast.next()
	ast.expect(K_LET_DECLARATION, IDENT)
	decl.Id = ast.CurrentToken.Literal
	ast.next()
	ast.expect(K_LET_DECLARATION, ASSIGN)
	ast.next()
	decl.Init = ast.ParseExpr()
	stmt.Body = decl
	return stmt
}
func (ast *AST) ParseExpr() Statement {
	stmt := ast.ParseBinaryExpr()
	return stmt
}

func (ast *AST) ParseBinaryExpr() Statement {
	left := ast.ParsePrimaryExpr()
	if !isOperator(ast.checkForward().Token) {
		return left
	}
	for {
		if !isOperator(ast.checkForward().Token) {
			return left
		} else {
			ast.next()
		}
		operator := ast.CurrentToken
		ast.next()
		right := ast.ParsePrimaryExpr()
		left = Statement{
			Kind: K_BINARY_EXPR,
			Body: BinaryExpr{
				Left:     left,
				Operator: operator.Token,
				Right:    right,
			},
		}

	}
}

func (ast *AST) ParsePrimaryExpr() Statement {
	stmt := Statement{}
	token := ast.CurrentToken.Token
	switch token {
	case IDENT:
		if ast.checkForward().Token == LPAREN {
			stmt.Kind = K_FUN_CALL
			stmt.Body = ast.ParseFunCall(ast.CurrentToken)
			return stmt
		}
		stmt.Kind = K_IDENTIFIER
		stmt.Body = ast.CurrentToken.Literal
		return stmt
	case INT:
		stmt.Kind = K_NUMERIC_LITERAL
		n, err := strconv.Atoi(ast.CurrentToken.Literal)
		if err != nil {
			panic(err)
		}
		stmt.Body = n
		return stmt
	case STRING:
		stmt.Kind = K_STRING
		stmt.Body = ast.CurrentToken.Literal
	default:
		ast.threwError(fmt.Sprintf("Invalid expression '%v' expect 'identifier' or 'number' at %v:%v", token, ast.CurrentToken.Pos.Line, ast.CurrentToken.Pos.Column))
	}

	return stmt
}
func (ast *AST) ParseFun() Statement {
	stmt := Statement{Kind: K_FUN_DECLARATION}
	fun := FunctionDeclaration{}
	ast.next()
	ast.expect(K_FUN_DECLARATION, IDENT)
	fun.Id = ast.CurrentToken.Literal
	ast.next()
	params := ast.ParseFunParams()
	fun.Params = params
	ast.next()
	fun.Body = ast.ParseBlock()
	stmt.Body = fun
	return stmt
}
func (ast *AST) ParseFunParams() []Statement {
	stmts := []Statement{}
	ast.expect(K_FUN_PARAM, LPAREN)
	ast.next()
	if ast.checkForward().Token == RPAREN {
		ast.next()
		return stmts
	}
	for ast.CurrentToken.Token != RPAREN {
		currentParam := FunParam{}
		ast.expect(K_FUN_PARAM, IDENT)
		currentParam.Name = ast.CurrentToken.Literal
		ast.next()
		ast.expect(K_FUN_PARAM, COLON)
		ast.next()
		ast.expect(K_FUN_PARAM, IDENT)
		currentParam.Type = ast.CurrentToken
		ast.next()
		stmts = append(stmts, Statement{Kind: K_FUN_PARAM, Body: currentParam})
		// Skip comma
		if ast.CurrentToken.Token == COMMA {
			ast.next()
		}
	}
	return stmts
}
func (ast *AST) ParseBlock() []Statement {
	stmts := []Statement{}
	ast.expect(K_BLOCK, LBRACE)
	ast.next()
	for ast.CurrentToken.Token != RBRACE {
		stmts = append(stmts, ast.Parse())
		if ast.CurrentToken.Token == EOF {
			ast.expect(K_BLOCK, RBRACE)
		}
		ast.next()
	}
	return stmts
}
func (ast *AST) ParseFunCall(caller Tokenized) FunCall {
	funCall := FunCall{}
	funCall.Caller = caller.Literal
	ast.expect(K_FUN_CALL, IDENT)
	ast.next()
	args := ast.ParseFunArgs()
	funCall.Args = args
	return funCall
}
func (ast *AST) ParseFunArgs() []Statement {
	args := []Statement{}
	ast.expect(K_FUN_ARGS_LIST, LPAREN)
	ast.next()
	for ast.CurrentToken.Token != RPAREN {
		ast.expectOneOf(K_FUN_ARGS_LIST, IDENT, INT, STRING)
		currentParam := ast.ParseExpr()
		args = append(args, currentParam)
		ast.next()
		if ast.CurrentToken.Token == RPAREN {
			return args
		}
		ast.expect(K_FUN_ARGS_LIST, COMMA)
		// Skip comma
		if ast.CurrentToken.Token == COMMA {
			ast.next()
		}
	}
	return args
}
