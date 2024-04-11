package runtime

type EvalValue struct {
	Value interface{}
	Type  VarType
}

type RuntimeFunction struct {
	Name  string
	Scope Scope
	Nodes []EvalValue
}
type RuntimeFunctionCall struct {
	Name     string
	Params   Parameters
	IsNative bool
	Call     func(Parameters) EvalValue
}
type RuntimeIfStatement struct {
	Condition EvalValue
	Scope     Scope
	Nodes     []EvalValue
}
