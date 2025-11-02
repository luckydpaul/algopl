package CompareWith

import (
	"algopl/CompareWith/Generator"
	"algopl/CompareWith/Paramstor"
	"fmt"
	"log"
	"testing"
)

func TestDemo(t *testing.T) {
	const demoName = "测试任务"
	dispatcher := NewDispatcher()
	dispatcher.Register(demoName, Demo)
	dispatcher.Run(demoName)
}

func A(a, b int) int {
	return a + b
}
func B(a, b int) int {
	return a + b
}
func Demo() {
	//采用 自定义的 params格式
	manParams := Paramstor.NewManParamstor(func() ([]interface{}, error) {
		data := make([]interface{}, 2)
		data[0] = 1
		data[1] = 2
		return data, nil
	})
	fmt.Println("ManParams:", manParams)
	//auto
	autoParams := Paramstor.NewAutoParamstor([]interface{}{
		Generator.NewIntGenerator(),
		2,
	})
	fmt.Println("ManParams:", manParams)

	comparator := NewComparator(A, B, autoParams, 100).SetLog(true)
	_, err := comparator.Run()
	if err != nil {
		log.Fatal(err)
	}
}
