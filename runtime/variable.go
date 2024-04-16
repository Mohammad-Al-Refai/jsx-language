package runtime

type VarType int

const (
	VAR_TYPE_UNDEFINED = iota
	VAR_TYPE_STRING
	VAR_TYPE_IDENTIFIER
	VAR_TYPE_NUMBER
	VAR_TYPE_BOOLEAN
	VAR_TYPE_FUNCTION
	VAR_TYPE_NATIVE_FUNCTION
	VAR_TYPE_BREAK
)

var types = []string{
	VAR_TYPE_UNDEFINED:       "undefined",
	VAR_TYPE_STRING:          "string",
	VAR_TYPE_IDENTIFIER:      "id",
	VAR_TYPE_NUMBER:          "number",
	VAR_TYPE_BOOLEAN:         "boolean",
	VAR_TYPE_FUNCTION:        "function",
	VAR_TYPE_NATIVE_FUNCTION: "n-function",
	VAR_TYPE_BREAK:           "break",
}

func (t VarType) String() string {
	return types[t]
}

type Variable struct {
	Name      string
	Value     interface{}
	ValueType VarType
}

func (v *Variable) IsUndefined() bool {
	return v.ValueType == VAR_TYPE_UNDEFINED
}
