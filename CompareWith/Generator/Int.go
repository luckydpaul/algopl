package Generator

import (
	"math/rand"
	"reflect"
	"sync"
	"time"
)

// 这个东西是我们在构建 intgenerator 时必要传进来的东西，设置的是需要生成的字符大小，等等
type IntConfig struct {
	Min *int //最小值
	Max *int //最大值

	Parity *string

	DIYS []func(int) bool
}

type IntGenerator struct {
	mu    sync.Mutex //保证并发安全
	r     *rand.Rand
	cfg   IntConfig //配置信息
	descs []struct {
		Even   *bool //奇偶
		MinMax *[2]int
	}
}

func NewIntGenerator(cfg IntConfig) *IntGenerator {
	g := &IntGenerator{}
	g.r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func (i *IntGenerator) Exec() (interface{}, error) {
	return 8, nil
}

func (i *IntGenerator) GetType() (reflect.Type, error) {
	return reflect.TypeOf(0), nil
}
