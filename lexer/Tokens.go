package lexer

import (
	"fmt"
)

type Token int

func (l Token) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

const (
	EOF = iota
	ILLEGAL
	IDENT
	INT
	COMMA //,
	GREATER_THAN
	SMALLER_THAN
	CLOSE_TAG
	OPEN_TAG
	CLOSE_OPEN_TAG

	// Infix ops
	ADD // +
	SUB // -
	MUL // *
	DIV // /

	LPAREN // (
	LBRACK // [
	LBRACE // {

	RPAREN  // )
	RBRACE  // }
	COMMENT // #
	OPERATOR
	STRING // 'string'

	//keywords
	IF
	LET
	FUNCTION
	PRINT
	RETURN
	FOR
	// operators
	EQUAL_EQUAL // ==
	EQUAL       // =
	NOT_EQUAL   // !=

)

var keywords = map[string]Token{"Let": LET, "If": IF, "Function": FUNCTION, "Print": PRINT, "Return": RETURN, "For": FOR, "greater": GREATER_THAN, "smaller": SMALLER_THAN}
var operators = map[Token]string{ADD: "+", SUB: "-", MUL: "*", DIV: "/", EQUAL_EQUAL: "==", NOT_EQUAL: "!="}

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	IDENT:   "IDENT",
	INT:     "INT",
	COMMA:   ",",
	LBRACE:  "{",
	RBRACE:  "}",
	LPAREN:  "(",
	RPAREN:  ")",
	STRING:  "STRING",
	// Infix ops
	ADD:            "+",
	SUB:            "-",
	MUL:            "*",
	DIV:            "/",
	EQUAL_EQUAL:    "==",
	NOT_EQUAL:      "!=",
	EQUAL:          "=",
	IF:             "If",
	LET:            "Let",
	FUNCTION:       "Function",
	FOR:            "For",
	RETURN:         "Return",
	COMMENT:        "#",
	GREATER_THAN:   "greater",
	SMALLER_THAN:   "smaller",
	CLOSE_OPEN_TAG: ">",
	CLOSE_TAG:      "</",
}

func (t Token) String() string {
	return fmt.Sprintf("%+v\n", tokens[t])
}

func IsKeyword(value string) (bool, Token) {
	if token, ok := keywords[value]; ok {
		return true, token
	}
	return false, ILLEGAL
}
func isOperator(value Token) bool {
	_, ok := operators[value]
	return ok
}
