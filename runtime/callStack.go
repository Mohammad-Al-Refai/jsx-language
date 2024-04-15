package runtime

import (
	"fmt"
	"os"
)

const MAX_CALL_STACK = 1024

type Call *RuntimeFunctionCall

type CallStack struct {
	Calls []Call
	Size  int
}

func (callStack *CallStack) threwError(message string) {
	fmt.Println(fmt.Errorf(fmt.Sprintf("[RuntimeError] %v", message)))
	os.Exit(1)
}
func NewCallStack() *CallStack {
	return &CallStack{
		Calls: make([]Call, 0, MAX_CALL_STACK),
	}
}
func (stack *CallStack) Push(function Call) {
	if stack.IsFull() {
		stack.threwError(fmt.Sprintf("Maximum call stack size [%+v] exceeded for function '%v'", MAX_CALL_STACK, function.Name))
	}
	stack.Calls = append(stack.Calls, function)

}

func (stack *CallStack) Pop() {
	stack.Calls = stack.Calls[:len(stack.Calls)-1]
}
func (stack *CallStack) IsFull() bool {
	return len(stack.Calls) == MAX_CALL_STACK
}
