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
	CLOSE_TAG      //</
	OPEN_TAG       //<
	CLOSE_OPEN_TAG // >

	// Infix ops
	ADD    // +
	SUB    // -
	MUL    // *
	DIV    // /
	MOD    // %
	DOT    // .
	LPAREN // (
	LBRACK // [
	RBRACK // ]
	LBRACE // {

	RPAREN  // )
	RBRACE  // }
	COMMENT // #
	OPERATOR
	STRING // 'string'

	//keywords
	IF
	LET
	SET
	FUNCTION
	PRINT
	RETURN
	BREAK
	FOR
	// operators
	EQUAL_EQUAL // ==
	EQUAL       // =
	NOT_EQUAL   // !=
	OR          //or
	AND         //and

)

var keywords = map[string]Token{"Let": LET, "Set": SET, "If": IF, "Function": FUNCTION, "Print": PRINT, "Break": BREAK, "Return": RETURN, "For": FOR, "greater": GREATER_THAN, "smaller": SMALLER_THAN, "or": OR, "and": AND}
var operators = map[Token]string{ADD: "+", SUB: "-", MUL: "*", DIV: "/", EQUAL_EQUAL: "==", MOD: "%", NOT_EQUAL: "!=", OR: "or", AND: "and", GREATER_THAN: "greater", SMALLER_THAN: "smaller"}

var tokens = []string{
	EOF:            "EOF",
	ILLEGAL:        "ILLEGAL",
	IDENT:          "IDENT",
	INT:            "INT",
	CLOSE_OPEN_TAG: ">",
	CLOSE_TAG:      "</",
	COMMA:          ",",
	LBRACE:         "{",
	RBRACE:         "}",
	LPAREN:         "(",
	RPAREN:         ")",
	LBRACK:         "[",
	RBRACK:         "]",
	DOT:            ".",
	STRING:         "STRING",
	// Infix ops
	ADD:          "+",
	SUB:          "-",
	MUL:          "*",
	DIV:          "/",
	MOD:          "%",
	EQUAL_EQUAL:  "==",
	NOT_EQUAL:    "!=",
	EQUAL:        "=",
	IF:           "If",
	LET:          "Let",
	FUNCTION:     "Function",
	FOR:          "For",
	RETURN:       "Return",
	BREAK:        "Break",
	COMMENT:      "#",
	GREATER_THAN: "greater",
	SET:          "Set",
	SMALLER_THAN: "smaller",

	OR:  "or",
	AND: "and",
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
