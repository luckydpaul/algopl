package Paramstor

import "algopl/CompareWith/Generator"

type AutoParamstor struct {
	params []interface{}
}

func NewAutoParamstor(data []interface{}) *AutoParamstor {
	return &AutoParamstor{params: data}
}

func (param *AutoParamstor) Value() ([]interface{}, error) {
	newParams := make([]interface{}, len(param.params))
	for i, params := range param.params {
		if generator, ok := params.(Generator.Generator); ok {
			result, err := generator.Exec()
			if err != nil {
				return nil, err
			}
			newParams[i] = result
		} else {
			newParams[i] = params
		}
	}
	return newParams, nil
}
