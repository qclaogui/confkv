package confkv

import (
	"log"
	"path"
	"time"

	"github.com/qclaogui/confkv/backends"
)

var defaultPrefix = "/"

var Store *StoreX

func init() { Store = NewStoreX() }

type conf struct {
	prefix string
	keys   []string
	bs     backends.BackendStore
	// closer holds a cleanup function that run after
	closer func()
}

// 设置keys的Prefix
func Prefix(prefix string) func(*conf) {
	return func(c *conf) { c.prefix = path.Join("/", prefix) }
}

// 设置Keys
func Keys(keys []string) func(*conf) {
	return func(c *conf) { c.keys = keys }
}

// 设置后端KV获取位置为Zookeeper
func Zookeeper(nodes ...string) func(*conf) {
	bs, err := backends.New(&backends.BackendConfig{
		Backend:      "zookeeper",
		BackendNodes: nodes})

	if err != nil {
		log.Fatalf("backends.New error %v\n", err)
	}
	return func(c *conf) { c.bs = bs }
}

func (c *conf) Stop() {
	c.closer()
	time.Sleep(500 * time.Millisecond)
}

func Watch(options ...func(*conf)) interface {
	Stop()
} {
	var c conf
	for _, opt := range options {
		opt(&c)
	}

	if c.prefix == "" {
		c.prefix = defaultPrefix
	}
	// 假如没有设置后端数据源，默认使用Zookeeper
	if c.bs == nil {
		Zookeeper()(&c)
	}

	stopChan := make(chan bool)
	// 设置closer 发送一个退出的信号
	c.closer = func() { close(stopChan) }

	// 运行
	errChan := make(chan error, 10)
	p := &watchProcessor{c, stopChan, errChan}
	go p.Process()

	select {
	case err := <-errChan:
		log.Printf("ERROR %v\n", err.Error())
	default:
	}
	return &c
}
