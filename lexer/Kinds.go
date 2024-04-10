package lexer

type StatementKind int

const (
	//KINDS
	PROGRAM = iota
	K_LET_DECLARATION
	K_FUN_DECLARATION
	K_END_OF_FILE
	K_BINARY_EXPR
	K_IDENTIFIER
	K_NUMERIC_LITERAL
	K_FUN_PARAM
	K_FUN_ARGS_LIST
	K_FUN_CALL
	K_BLOCK
	K_STRING
)

var kinds = []string{
	PROGRAM:           "Program",
	K_LET_DECLARATION: "LetDeclaration",
	K_END_OF_FILE:     "EOF",
	K_BINARY_EXPR:     "BinaryExpr",
	K_IDENTIFIER:      "Identifier",
	K_NUMERIC_LITERAL: "NumericLiteral",
	K_STRING:          "String",
	K_FUN_DECLARATION: "FunDeclaration",
	K_FUN_PARAM:       "FunParam",
	K_BLOCK:           "Block",
	K_FUN_CALL:        "FunCall",
	K_FUN_ARGS_LIST:   "FunArgsList",
}

func (k StatementKind) String() string {
	return kinds[k]
}

// convert iota to string
func (l StatementKind) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}
