package lexer

type StatementKind int

const (
	//KINDS
	K_EOF = iota
	K_OPEN_TAG
	K_CLOSE_TAG
	K_PARAMETERS
	K_PARAMETER_VALUE
	K_IDENTIFIER
	K_OPERATOR
	K_EXPRESSION
	K_NUMBER
	K_STRING
	K_OBJECT
	K_BINARY_EXPR
	K_LOGICAL_EXPR
)

var kinds = []string{
	K_EOF:             "EOF",
	K_OPEN_TAG:        "OpenTag",
	K_CLOSE_TAG:       "CloseTag",
	K_PARAMETERS:      "Parameters",
	K_PARAMETER_VALUE: "ParameterValue",
	K_IDENTIFIER:      "Identifier",
	K_NUMBER:          "Number",
	K_STRING:          "String",
	K_OBJECT:          "Object",
	K_BINARY_EXPR:     "BinaryExpr",
	K_LOGICAL_EXPR:    "LogicalExpr",
	K_OPERATOR:        "Operator",
	K_EXPRESSION:      "Expression",
}

func (k StatementKind) String() string {
	return kinds[k]
}

// convert iota to string
func (l StatementKind) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}
