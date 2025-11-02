package CompareWith

import (
	"algopl/CompareWith/Generator"
	"algopl/CompareWith/Paramstor"
	"errors"
	"fmt"
	"reflect"
)

func toReflectValue(data []interface{}) []reflect.Value {
	paramsValue := make([]reflect.Value, len(data))
	for i, param := range data {
		paramsValue[i] = reflect.ValueOf(param)
	}
	return paramsValue
}

// CloneValues 深度克隆一个 []reflect.Value，返回全新的、与原数据完全隔离的副本。
func CloneValues(src []reflect.Value) []reflect.Value {
	dst := make([]reflect.Value, len(src))
	for i, v := range src {
		dst[i] = deepCopyValue(v)
	}
	return dst
}

// deepCopyValue 递归深拷贝单个 reflect.Value
func deepCopyValue(v reflect.Value) reflect.Value {
	if !v.IsValid() {
		return v // 零值直接返回
	}

	// 处理接口和指针：递归解引用
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return reflect.New(v.Type()).Elem()
		}
		v = v.Elem()
	}

	switch v.Kind() {

	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.String, reflect.Complex64, reflect.Complex128:
		// 基础类型 —— 直接返回副本
		return v

	case reflect.Array:
		t := v.Type()
		l := v.Len()
		newArr := reflect.New(t).Elem()
		for i := 0; i < l; i++ {
			newArr.Index(i).Set(deepCopyValue(v.Index(i)))
		}
		return newArr

	case reflect.Slice:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		l := v.Len()
		newSlice := reflect.MakeSlice(v.Type(), l, l)
		for i := 0; i < l; i++ {
			newSlice.Index(i).Set(deepCopyValue(v.Index(i)))
		}
		return newSlice

	case reflect.Map:
		if v.IsNil() {
			return reflect.Zero(v.Type())
		}
		t := v.Type()
		newMap := reflect.MakeMapWithSize(t, v.Len())
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			newKey := deepCopyValue(key)
			newVal := deepCopyValue(val)
			newMap.SetMapIndex(newKey, newVal)
		}
		return newMap

	case reflect.Struct:
		t := v.Type()
		newStruct := reflect.New(t).Elem()
		for i := 0; i < t.NumField(); i++ {
			if t.Field(i).PkgPath != "" {
				// 未导出字段跳过
				continue
			}
			newStruct.Field(i).Set(deepCopyValue(v.Field(i)))
		}
		return newStruct

	case reflect.Chan, reflect.Func, reflect.UnsafePointer:
		// 这些类型无法深拷贝，返回零值
		return reflect.Zero(v.Type())

	default:
		// 未知类型，返回零值
		return reflect.Zero(v.Type())
	}
}

type Comparator struct {
	funcA  interface{}
	funcB  interface{}
	params Paramstor.Paramstor //调度器，用于param的生成

	count    int  //生成器循环次数
	nowCount int  //当前循环的次数
	log      bool //是否打印日志
}

func NewComparator(a, b interface{}, params Paramstor.Paramstor, count int) *Comparator {
	return &Comparator{funcA: a, funcB: b, params: params, count: count}
}

// 设置日志
func (c *Comparator) SetLog(log bool) *Comparator {
	c.log = log
	return c
}
func (c *Comparator) validate() error {
	//判断 函数 a 是不是函数
	if reflect.TypeOf(c.funcA).Kind() != reflect.Func {
		return errors.New("传入的函数 A 不是函数")
	}
	//判断 函数 b 是不是函数
	if reflect.TypeOf(c.funcB).Kind() != reflect.Func {
		return errors.New("传入的函数 B 不是函数")
	}

	funcAType := reflect.TypeOf(c.funcA)
	funcBType := reflect.TypeOf(c.funcB)
	//比较 params 的数量是否 a,b 函数所需的传参是一致的数量
	if funcAType.NumIn() != funcBType.NumIn() {
		return errors.New("传入的函数 A,B 所需的参数数量不一致")
	}
	for i := 0; i < funcAType.NumIn(); i++ {
		if funcAType.In(i) != funcBType.In(i) {
			return errors.New(fmt.Sprintf("传入的参数 第%d项 所需要的类型不一致", i+1))
		}
	}
	//比较 params 的数量是否 a,b 函数所输出的参数是一致的数量
	if funcAType.NumOut() != funcBType.NumOut() {
		return errors.New("传入的函数 A,B 输出的参数数量不一致")
	}
	for i := 0; i < funcAType.NumOut(); i++ {
		if funcAType.Out(i) != funcBType.Out(i) {
			return errors.New(fmt.Sprintf("传入的参数 第%d项 输出的类型不一致", i+1))
		}
	}
	// 检查传入的参数数量是否与函数参数数量一致
	_paramsDemo, err := c.params.Value()
	if err != nil {
		return errors.New(fmt.Sprintf("获取的 params value 错误 error:%v", err))
	}
	//验一下总数对不对的上
	if len(_paramsDemo) != funcAType.NumIn() {
		return errors.New("生成的TempParams 与 函数参数数量不一致")
	}
	//挨个判断一下 挨个的
	for i, param := range _paramsDemo {
		expectedType := funcAType.In(i)
		if generate, ok := param.(Generator.Generator); ok {
			paramType, err := generate.GetType()
			if err != nil {
				return errors.New(fmt.Sprintf("获取 generate的GetType函数错误 error:%v", err))
			}
			if paramType != expectedType {
				return errors.New(fmt.Sprintf("1第 %d 个参数类型不匹配，期望:%v ,实际:%v", i+1, expectedType, paramType))
			}
		} else {
			paramType := reflect.TypeOf(param)
			if paramType != expectedType {
				return errors.New(fmt.Sprintf("2第 %d 个参数类型不匹配，期望:%v ,实际:%v", i+1, expectedType, paramType))
			}
		}
	}
	return nil
}

// 调用运行函数
func (c *Comparator) Run() (bool, error) {
	var err error
	err = c.validate()
	if err != nil {
		return false, err
	}
	for i := 0; i < c.count; i++ {
		c.nowCount++
		//通过参数调用器生成 []interface 参数
		newParam, err := c.params.Value()
		if err != nil {
			return false, err
		}
		ResultA, ResultB := c.doFunction(toReflectValue(newParam))
		err = c.finallyCheck(newParam, ResultA, ResultB)
		if err != nil {
			return false, err
		}
	}
	fmt.Println("完成任务,校验通过!!!")
	return true, nil
}
func (c *Comparator) doFunction(data []reflect.Value) ([]reflect.Value, []reflect.Value) {
	//深度拷贝一个隔离数据的放进去
	ResultA := reflect.ValueOf(c.funcA).Call(CloneValues(data))
	ResultB := reflect.ValueOf(c.funcB).Call(CloneValues(data))
	return ResultA, ResultB
}
func (c *Comparator) finallyCheck(params []interface{}, ResultA, ResultB []reflect.Value) error {
	var retInterface = func(data []reflect.Value) []interface{} {
		newdata := make([]interface{}, len(data))
		for i, v := range data {
			newdata[i] = v.Interface()
		}
		return newdata
	}
	var err error
	defer func() {
		//
		if (c.log && err == nil) || err != nil {
			fmt.Printf("︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻\n")
			fmt.Printf("| 【总次数:%d次】【当前是第%d次】\n", c.count, c.nowCount)
			fmt.Printf("| 传入的参数是 %v\n", params)
			fmt.Printf("| a的输出结果是 %v\n", retInterface(ResultA))
			fmt.Printf("| b的输出结果是 %v\n", retInterface(ResultB))
			fmt.Printf("︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻︻\n")
		}
	}()
	if len(ResultA) != len(ResultB) {
		return errors.New("函数输出数量不一致")
	}
	for i := 0; i < len(ResultA); i++ {
		if !reflect.DeepEqual(ResultA[i].Interface(), ResultA[i].Interface()) {
			err = errors.New("函数输出结果不一致")
			return err
		}
	}
	return nil
}
