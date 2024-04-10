package lexer

type Token int

func (l Token) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

const (
	EOF = iota
	ILLEGAL
	IDENT
	INT
	SEMI  // ;
	COLON //;
	COMMA //,

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
	ASSIGN  // =
	COMMENT // #
	OPERATOR
	STRING // 'string'

	//keywords
	IF   // if
	ELSE // else
	LET  // let
	FUN  // fun

	// operators
	EQUAL_EQUAL // ==
)

var keywords = map[string]Token{"let": LET, "if": IF, "else": ELSE, "fun": FUN}
var operators = map[Token]string{ADD: "+", SUB: "-", MUL: "*", DIV: "/"}

var tokens = []string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	IDENT:   "IDENT",
	INT:     "INT",
	SEMI:    ";",
	COLON:   ":",
	COMMA:   ",",
	LBRACE:  "{",
	RBRACE:  "}",
	LPAREN:  "(",
	RPAREN:  ")",
	STRING:  "STRING",
	// Infix ops
	ADD:    "+",
	SUB:    "-",
	MUL:    "*",
	DIV:    "/",
	ASSIGN: "=",

	EQUAL_EQUAL: "==",

	IF:   "if",
	ELSE: "else",
	LET:  "let",
	FUN:  "fun",

	COMMENT: "#",
}

func (t Token) String() string {
	return tokens[t]
}

func isKeyword(value string) (bool, Token) {
	if token, ok := keywords[value]; ok {
		return true, token
	}
	return false, ILLEGAL
}
func isOperator(value Token) bool {
	_, ok := operators[value]
	return ok
}
