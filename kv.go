package kv

import (
	"path"
	"time"

	"github.com/qclaogui/kv/log"

	"github.com/qclaogui/kv/backends"
	"github.com/qclaogui/kvdb"
)

var mem kvdb.Storage

func init() { mem = kvdb.NewMem() }

// 检查某个key是否存在
func Exists(key string) bool { return mem.Exists(key) }

// GetVs 获取相应key的value
func GetV(key string, defaultValue ...string) (string, error) {
	return mem.Get(key, defaultValue...)
}

// GetVs 获取匹配到pattern的所有keys的value
func GetVs(pattern string) ([]string, error) { return mem.GetMany(pattern) }

type conf struct {
	prefix string
	keys   []string
	log    log.Logger
	bs     backends.BackendStore
	// closer holds a cleanup function that run after
	closer func()
}

func (c *conf) Stop() {
	c.closer()
	time.Sleep(500 * time.Millisecond)
}

func Watch(prefix string, keys []string, options ...func(*conf)) interface {
	Stop()
} {
	var c conf
	if prefix == "" {
		c.prefix = "/"
	} else {
		c.prefix = path.Join("/", prefix)
	}
	c.keys = keys

	for _, opt := range options {
		opt(&c)
	}

	if c.bs == nil {
		Options.Zookeeper()(&c)
	}

	if c.log == nil {
		c.log = log.NullLogger
	}

	stopChan := make(chan bool)
	// 设置closer 发送一个退出的信号
	c.closer = func() { close(stopChan) }

	// 运行
	errChan := make(chan error, 10)
	p := &watchProcessor{c, stopChan, errChan}
	go p.Process()

	// Track errors
	go func() {
		for err := range errChan {
			c.log.Error(err.Error())
		}
	}()

	return &c
}
