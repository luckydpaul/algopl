package Generator

import (
	"reflect"
)

type IntGenerator struct {
}

func NewIntGenerator() *IntGenerator {
	return &IntGenerator{}
}

func (i *IntGenerator) Exec() (interface{}, error) {
	return 8, nil
}

func (i *IntGenerator) GetType() (reflect.Type, error) {
	return reflect.TypeOf(0), nil
}
