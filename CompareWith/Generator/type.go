package Generator

import "reflect"

type Generator interface {
	Exec() (interface{}, error)
	GetType() (reflect.Type, error)
}
