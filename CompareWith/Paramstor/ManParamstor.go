package Paramstor

import "errors"

// @这个方法是针对自己diy的func param生成器
type ManParamstor struct {
	function func() ([]interface{}, error)
}

func NewManParamstor(fn func() ([]interface{}, error)) *ManParamstor {
	return &ManParamstor{function: fn}
}

func (params *ManParamstor) Value() ([]interface{}, error) {
	if params.function == nil {
		return nil, errors.New("调用函数是空的")
	}
	data, err := params.function()
	if err != nil {
		return nil, err
	}
	return data, nil
}
