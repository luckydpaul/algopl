package CompareWith

import "fmt"

type Operator struct {
	name     string
	function func()
}

// GAI
func (o *Operator) Do() {
	fmt.Printf("运行调度 [%s] 任务\n", o.name)
	o.function()
	fmt.Printf("任务调用结束")
}

type Dispatcher struct {
	funcs map[string]*Operator
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		funcs: make(map[string]*Operator),
	}
}

func (d *Dispatcher) Register(Name string, function func()) {
	d.funcs[Name] = &Operator{name: Name, function: function}
}

func (d *Dispatcher) Run(Name string) {
	o, ok := d.funcs[Name]
	if !ok {
		panic("没有找到这个名字的调度")
	}
	o.Do()
}
