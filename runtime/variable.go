package runtime

type VarType int

const (
	VAR_TYPE_UNDEFINED = iota
	VAR_TYPE_STRING
	VAR_TYPE_IDENTIFIER
	VAR_TYPE_NUMBER
	VAR_TYPE_FUNCTION
	VAR_TYPE_NATIVE_FUNCTION
)

type Variable struct {
	Name      string
	Value     interface{}
	ValueType VarType
}
