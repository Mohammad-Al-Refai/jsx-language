package runtime

type ArrayRuntime struct {
	Size  int
	Items []*EvalValue
}

func (a *ArrayRuntime) Push(value *EvalValue) {
	a.Items = append(a.Items, value)
}

func (a *ArrayRuntime) Pop() *EvalValue {
	last := a.Items[len(a.Items)-1]
	a.Items = a.Items[:len(a.Items)-1]
	return last
}
