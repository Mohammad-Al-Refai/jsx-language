package runtime

import (
	"fmt"
)

func GlobalScope() *Scope {
	globalScope := Scope{}
	globalScope.DefineFunction(Print(Parameters{"value": &EvalValue{}}))
	globalScope.DefineVariable(Variable{Name: "true", ValueType: VAR_TYPE_BOOLEAN, Value: true})
	globalScope.DefineVariable(Variable{Name: "false", ValueType: VAR_TYPE_BOOLEAN, Value: false})
	builtinArray := array{}
	globalScope.DefineObject(builtinArray.Init())
	return &globalScope
}
func Print(param Parameters) *RuntimeFunctionCall {
	return &RuntimeFunctionCall{
		IsNative: true,
		Name:     "Print",
		Scope:    &Scope{},
		Call: func(p Parameters) *EvalValue {
			value := p["value"]
			if value.Type == VAR_TYPE_ARRAY {
				arr := value.Value.(*ArrayRuntime)
				items := make([]interface{}, 0)
				for _, item := range arr.Items {
					items = append(items, item.Value)
				}
				fmt.Printf("%+v\n", items)
				return &EvalValue{}
			}
			fmt.Println(value.Value)
			return &EvalValue{}

		},
	}
}

type array struct{}

func (a *array) Init() *RuntimeObject {
	return &RuntimeObject{
		Name: "array",
		Members: []*RuntimeObjectMember{
			a.length(),
			a.at(),
			a.push(),
			a.pop(),
		},
	}
}
func (a *array) length() *RuntimeObjectMember {
	return &RuntimeObjectMember{
		Name: "length",
		Call: func(stack *ScopeStack) *EvalValue {
			targetArray := stack.Pop()
			if targetArray.Type == VAR_TYPE_ARRAY {
				return &EvalValue{Value: targetArray.Value.(*ArrayRuntime).Size, Type: VAR_TYPE_NUMBER}
			}
			return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
		},
	}
}
func (a *array) push() *RuntimeObjectMember {
	return &RuntimeObjectMember{
		Name: "push",
		Call: func(stack *ScopeStack) *EvalValue {
			targetArray := stack.Pop()
			if targetArray.Type == VAR_TYPE_ARRAY {
				value := stack.Pop()
				targetArray.Value.(*ArrayRuntime).Push(value)
				return &EvalValue{Value: targetArray.Value.(*ArrayRuntime), Type: VAR_TYPE_ARRAY}
			}
			return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}
		},
	}
}
func (a *array) at() *RuntimeObjectMember {
	return &RuntimeObjectMember{
		Name: "at",
		Call: func(stack *ScopeStack) *EvalValue {
			targetArray := stack.Pop()
			if targetArray.Type != VAR_TYPE_ARRAY {
				return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}

			}
			index := stack.Pop()
			if index.Type != VAR_TYPE_NUMBER {
				return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}

			}
			return targetArray.Value.(*ArrayRuntime).At(index.Value.(int))
		},
	}
}
func (a *array) pop() *RuntimeObjectMember {
	return &RuntimeObjectMember{
		Name: "pop",
		Call: func(stack *ScopeStack) *EvalValue {
			targetArray := stack.Pop()
			if targetArray.Type != VAR_TYPE_ARRAY {
				return &EvalValue{Type: VAR_TYPE_UNDEFINED, Value: "undefined"}

			}

			return targetArray.Value.(*ArrayRuntime).Pop()
		},
	}
}
