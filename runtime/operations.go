package runtime

import (
	"fmt"
	"strconv"

	"m.shebli.refaai/ht/lexer"
)

func (interpreter *Interpreter) EvaluateOperator(expr lexer.Statement, scope *Scope) *EvalValue {
	if expr.Body == "+" {
		return interpreter.Sum(scope.Stack.Pop(), scope.Stack.Pop())
	}
	if expr.Body == "*" {
		return interpreter.Mul(scope.Stack.Pop(), scope.Stack.Pop())
	}
	if expr.Body == "/" {
		return interpreter.Div(scope.Stack.Pop(), scope.Stack.Pop())
	}
	if expr.Body == "-" {
		return interpreter.Sub(scope.Stack.Pop(), scope.Stack.Pop())
	}
	if expr.Body == "%" {
		return interpreter.Mod(scope.Stack.Pop(), scope.Stack.Pop())
	}
	if expr.Body == "greater" {
		return interpreter.GreaterThan(scope.Stack.Pop(), scope.Stack.Pop())
	}
	if expr.Body == "smaller" {
		return interpreter.SmallerThan(scope.Stack.Pop(), scope.Stack.Pop())
	}
	if expr.Body == "==" {
		return interpreter.Equal(scope.Stack.Pop(), scope.Stack.Pop())
	}
	if expr.Body == "!=" {
		return interpreter.NotEqual(scope.Stack.Pop(), scope.Stack.Pop())
	}
	return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
}

func (interpreter *Interpreter) Sum(x *EvalValue, y *EvalValue) *EvalValue {
	if x.IsString() && y.IsString() {
		return &EvalValue{Type: VAR_TYPE_STRING, Value: y.Value.(string) + x.Value.(string)}
	}
	if x.IsNumber() && y.IsNumber() {
		return &EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(int) + x.Value.(int)}
	}
	if x.IsString() && y.IsNumber() {
		return &EvalValue{Type: VAR_TYPE_STRING, Value: strconv.Itoa(y.Value.(int)) + x.Value.(string)}
	}
	if x.IsNumber() && y.IsString() {
		return &EvalValue{Type: VAR_TYPE_STRING, Value: y.Value.(string) + strconv.Itoa(x.Value.(int))}
	}
	interpreter.threwError(fmt.Sprintf("expect left and right to be number or string found '%v' and '%v'", y.Type.String(), x.Type.String()))
	return &EvalValue{Type: VAR_TYPE_UNDEFINED}
}
func (interpreter *Interpreter) Mul(x *EvalValue, y *EvalValue) *EvalValue {
	if !x.IsNumber() || !y.IsNumber() {
		interpreter.threwError(fmt.Sprintf("expect both number found '%v' and '%v'", y.Type.String(), x.Type.String()))
	}
	return &EvalValue{Type: VAR_TYPE_NUMBER, Value: x.Value.(int) * y.Value.(int)}
}
func (interpreter *Interpreter) Div(x *EvalValue, y *EvalValue) *EvalValue {
	if !x.IsNumber() || !y.IsNumber() {
		interpreter.threwError(fmt.Sprintf("expect both number found '%v' and '%v'", y.Type.String(), x.Type.String()))
	}
	return &EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(int) / x.Value.(int)}
}
func (interpreter *Interpreter) Sub(x *EvalValue, y *EvalValue) *EvalValue {
	if !x.IsNumber() || !y.IsNumber() {
		interpreter.threwError(fmt.Sprintf("expect both number found '%v' and '%v'", y.Type.String(), x.Type.String()))
	}
	return &EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(int) - x.Value.(int)}
}
func (interpreter *Interpreter) Mod(x *EvalValue, y *EvalValue) *EvalValue {
	if !x.IsNumber() || !y.IsNumber() {
		interpreter.threwError(fmt.Sprintf("expect both number found '%v' and '%v'", y.Type.String(), x.Type.String()))
	}
	return &EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(int) % x.Value.(int)}
}
func (interpreter *Interpreter) GreaterThan(x *EvalValue, y *EvalValue) *EvalValue {
	if !x.IsNumber() || !y.IsNumber() {
		interpreter.threwError(fmt.Sprintf("expect both number or boolean found '%v' and '%v'", y.Type.String(), x.Type.String()))
	}
	return &EvalValue{Type: VAR_TYPE_BOOLEAN, Value: y.Value.(int) > x.Value.(int)}
}
func (interpreter *Interpreter) SmallerThan(x *EvalValue, y *EvalValue) *EvalValue {
	if !x.IsNumber() || !y.IsNumber() {
		interpreter.threwError(fmt.Sprintf("expect both number or boolean found '%v' and '%v'", y.Type.String(), x.Type.String()))
	}
	return &EvalValue{Type: VAR_TYPE_BOOLEAN, Value: y.Value.(int) < x.Value.(int)}
}
func (interpreter *Interpreter) Equal(x *EvalValue, y *EvalValue) *EvalValue {
	if x.IsString() && y.IsString() {
		return &EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(string) == x.Value.(string)}
	} else if x.IsBoolean() && y.IsBoolean() {
		return &EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(bool) == x.Value.(bool)}
	} else if x.IsNumber() && y.IsNumber() {
		return &EvalValue{Type: VAR_TYPE_BOOLEAN, Value: y.Value.(int) == x.Value.(int)}
	}
	interpreter.threwError(fmt.Sprintf("expect string or number found '%v' and '%v'", y.Type.String(), x.Type.String()))
	return &EvalValue{Type: VAR_TYPE_UNDEFINED}
}
func (interpreter *Interpreter) NotEqual(x *EvalValue, y *EvalValue) *EvalValue {
	if x.IsString() && y.IsString() {
		return &EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(string) != x.Value.(string)}
	} else if x.IsBoolean() && y.IsBoolean() {
		return &EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(bool) != x.Value.(bool)}
	} else if x.IsNumber() && y.IsNumber() {
		return &EvalValue{Type: VAR_TYPE_BOOLEAN, Value: y.Value.(int) != x.Value.(int)}
	}
	interpreter.threwError(fmt.Sprintf("expect string or number found '%v' and '%v'", y.Type.String(), x.Type.String()))
	return &EvalValue{Type: VAR_TYPE_UNDEFINED}
}
