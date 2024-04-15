package runtime

import (
	"fmt"

	"m.shebli.refaai/ht/lexer"
)

func (interpreter *Interpreter) EvaluateOperator(expr lexer.Statement, scope *Scope) *EvalValue {
	if expr.Body == "+" {
		return interpreter.Sum(scope.Pop(), scope.Pop())
	}
	if expr.Body == "*" {
		return interpreter.Mul(scope.Pop(), scope.Pop())
	}
	if expr.Body == "/" {
		return interpreter.Div(scope.Pop(), scope.Pop())
	}
	if expr.Body == "-" {
		return interpreter.Sub(scope.Pop(), scope.Pop())
	}
	if expr.Body == "greater" {
		return interpreter.GreaterThan(scope.Pop(), scope.Pop())
	}
	if expr.Body == "smaller" {
		return interpreter.SmallerThan(scope.Pop(), scope.Pop())
	}
	if expr.Body == "==" {
		return interpreter.Equal(scope.Pop(), scope.Pop())
	}
	if expr.Body == "!=" {
		return interpreter.NotEqual(scope.Pop(), scope.Pop())
	}
	return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}

func (interpreter *Interpreter) Sum(x *EvalValue, y *EvalValue) *EvalValue {
	if !(x.IsNumber() || x.IsString() || y.IsNumber() || y.IsString()) {
		interpreter.threwError(fmt.Sprintf("expect string or number found %v and %v", y.Type.String(), x.Type.String()))
	}
	if x.IsString() && y.IsString() {
		return &EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(string) + x.Value.(string)}
	}
	if x.IsNumber() && y.IsNumber() {
		return &EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(int) + x.Value.(int)}
	}
	if x.IsNumber() || x.IsString() && y.IsNumber() || y.IsString() {
		interpreter.threwError(fmt.Sprintf("expect left and right to be number or string found %v and %v", y.Type.String(), x.Type.String()))
	}
	return &EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(int) + x.Value.(int)}
}
func (interpreter *Interpreter) Mul(x *EvalValue, y *EvalValue) *EvalValue {
	if !(x.IsNumber() || y.IsNumber()) {
		interpreter.threwError(fmt.Sprintf("expect both number found %v and %v", y.Type.String(), x.Type.String()))
	}
	return &EvalValue{Type: VAR_TYPE_NUMBER, Value: x.Value.(int) * y.Value.(int)}
}
func (interpreter *Interpreter) Div(x *EvalValue, y *EvalValue) *EvalValue {
	if !(x.IsNumber() || y.IsNumber()) {
		interpreter.threwError(fmt.Sprintf("expect both  number found %v and %v", y.Type.String(), x.Type.String()))
	}
	return &EvalValue{Type: VAR_TYPE_NUMBER, Value: x.Value.(int) / y.Value.(int)}
}
func (interpreter *Interpreter) Sub(x *EvalValue, y *EvalValue) *EvalValue {
	if !(x.IsNumber() || y.IsNumber()) {
		interpreter.threwError(fmt.Sprintf("expect both number found %v and %v", y.Type.String(), x.Type.String()))
	}
	return &EvalValue{Type: VAR_TYPE_NUMBER, Value: x.Value.(int) - y.Value.(int)}
}
func (interpreter *Interpreter) GreaterThan(x *EvalValue, y *EvalValue) *EvalValue {
	if !(x.IsNumber() || y.IsNumber()) {
		interpreter.threwError(fmt.Sprintf("expect both number or boolean found %v and %v", y.Type.String(), x.Type.String()))
	}
	return &EvalValue{Type: VAR_TYPE_BOOLEAN, Value: y.Value.(int) > x.Value.(int)}
}
func (interpreter *Interpreter) SmallerThan(x *EvalValue, y *EvalValue) *EvalValue {
	if !(x.IsNumber() || y.IsNumber()) {
		interpreter.threwError(fmt.Sprintf("expect both number or boolean found %v and %v", y.Type.String(), x.Type.String()))
	}
	return &EvalValue{Type: VAR_TYPE_BOOLEAN, Value: y.Value.(int) < x.Value.(int)}
}
func (interpreter *Interpreter) Equal(x *EvalValue, y *EvalValue) *EvalValue {
	if !(x.IsNumber() || x.IsString() || y.IsNumber() || y.IsString() || y.IsBoolean() || y.IsBoolean()) {
		interpreter.threwError(fmt.Sprintf("expect string or number found %v and %v", y.Type.String(), x.Type.String()))
	}
	if x.IsString() && y.IsString() {
		return &EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(string) == x.Value.(string)}
	}
	if x.IsBoolean() && y.IsBoolean() {
		return &EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(bool) == x.Value.(bool)}
	}
	return &EvalValue{Type: VAR_TYPE_BOOLEAN, Value: y.Value.(int) == x.Value.(int)}
}
func (interpreter *Interpreter) NotEqual(x *EvalValue, y *EvalValue) *EvalValue {
	if !(x.IsNumber() || x.IsString() || y.IsNumber() || y.IsString() || y.IsBoolean() || y.IsBoolean()) {
		interpreter.threwError(fmt.Sprintf("expect string or number found %v and %v", y.Type.String(), x.Type.String()))
	}
	if x.IsString() && y.IsString() {
		return &EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(string) != x.Value.(string)}
	}
	if x.IsBoolean() && y.IsBoolean() {
		return &EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(bool) != x.Value.(bool)}
	}
	return &EvalValue{Type: VAR_TYPE_BOOLEAN, Value: y.Value.(int) != x.Value.(int)}
}
