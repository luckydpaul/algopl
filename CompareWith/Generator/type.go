package Generator

import "reflect"

type Generator interface {
	Exec() (interface{}, error)
	GetType() (reflect.Type, error)
}

const (
	Generator_sturct_odd  = "odd"  //奇数
	Generator_sturct_even = "even" //偶数
)
