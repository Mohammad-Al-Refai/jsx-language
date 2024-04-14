package runtime

import "fmt"

func (interpreter *Interpreter) Sum(x EvalValue, y EvalValue) EvalValue {
	if !(x.IsNumber() || x.IsString() || y.IsNumber() || y.IsString()) {
		interpreter.threwError(fmt.Sprintf("expect string or number found %v and %v", y.Type.String(), x.Type.String()))
	}
	if x.IsString() && y.IsString() {
		return EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(string) + x.Value.(string)}
	}
	if x.IsNumber() && y.IsNumber() {
		return EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(int) + x.Value.(int)}
	}
	if x.IsNumber() || x.IsString() && y.IsNumber() || y.IsString() {
		interpreter.threwError(fmt.Sprintf("expect left and right to be number found %v and %v", y.Type.String(), x.Type.String()))
	}
	return EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(int) + x.Value.(int)}
}
func (interpreter *Interpreter) Mul(x EvalValue, y EvalValue) EvalValue {
	if !(x.IsNumber() || y.IsNumber()) {
		interpreter.threwError(fmt.Sprintf("expect both number found %v and %v", y.Type.String(), x.Type.String()))
	}
	return EvalValue{Type: VAR_TYPE_NUMBER, Value: x.Value.(int) * y.Value.(int)}
}
func (interpreter *Interpreter) Div(x EvalValue, y EvalValue) EvalValue {
	if !(x.IsNumber() || y.IsNumber()) {
		interpreter.threwError(fmt.Sprintf("expect both  number found %v and %v", y.Type.String(), x.Type.String()))
	}
	return EvalValue{Type: VAR_TYPE_NUMBER, Value: x.Value.(int) / y.Value.(int)}
}
func (interpreter *Interpreter) Sub(x EvalValue, y EvalValue) EvalValue {
	if !(x.IsNumber() || y.IsNumber()) {
		interpreter.threwError(fmt.Sprintf("expect both number found %v and %v", y.Type.String(), x.Type.String()))
	}
	return EvalValue{Type: VAR_TYPE_NUMBER, Value: x.Value.(int) - y.Value.(int)}
}
func (interpreter *Interpreter) GreaterThan(x EvalValue, y EvalValue) EvalValue {
	if !(x.IsNumber() || y.IsNumber()) {
		interpreter.threwError(fmt.Sprintf("expect both number or boolean found %v and %v", y.Type.String(), x.Type.String()))
	}
	return EvalValue{Type: VAR_TYPE_BOOLEAN, Value: y.Value.(int) > x.Value.(int)}
}
func (interpreter *Interpreter) SmallerThan(x EvalValue, y EvalValue) EvalValue {
	if !(x.IsNumber() || y.IsNumber()) {
		interpreter.threwError(fmt.Sprintf("expect both number or boolean found %v and %v", y.Type.String(), x.Type.String()))
	}
	return EvalValue{Type: VAR_TYPE_BOOLEAN, Value: y.Value.(int) < x.Value.(int)}
}
func (interpreter *Interpreter) Equal(x EvalValue, y EvalValue) EvalValue {
	if !(x.IsNumber() || x.IsString() || y.IsNumber() || y.IsString()) {
		interpreter.threwError(fmt.Sprintf("expect string or number found %v and %v", y.Type.String(), x.Type.String()))
	}
	if x.IsString() && y.IsString() {
		return EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(string) == x.Value.(string)}
	}
	return EvalValue{Type: VAR_TYPE_BOOLEAN, Value: y.Value.(int) == x.Value.(int)}
}
func (interpreter *Interpreter) NotEqual(x EvalValue, y EvalValue) EvalValue {
	if !(x.IsNumber() || x.IsString() || y.IsNumber() || y.IsString()) {
		interpreter.threwError(fmt.Sprintf("expect string or number found %v and %v", y.Type.String(), x.Type.String()))
	}
	if x.IsString() && y.IsString() {
		return EvalValue{Type: VAR_TYPE_NUMBER, Value: y.Value.(string) != x.Value.(string)}
	}
	return EvalValue{Type: VAR_TYPE_BOOLEAN, Value: y.Value.(int) != x.Value.(int)}
}
