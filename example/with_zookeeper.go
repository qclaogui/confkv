package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qclaogui/kv/log"

	"github.com/qclaogui/kv"
)

var prefix = "/app"

// 当然 backend（默认) 中需要有这些配置
var keys = []string{
	"/upstream/host1",
	"/upstream/host2"}

func main() {
	defer kv.Watch(prefix, keys,
		kv.Options.WithLogger(log.StdLogger),
		kv.Options.Zookeeper(),
	).Stop()

	// 等待从后端获取配置 然后第一次加载到内存 浪费点启动内存
	time.Sleep(time.Second)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	for range time.Tick(3 * time.Second) {
		v, err := kv.GetV("/app/upstream/host1")
		if err != nil {
			fmt.Printf("Get error %v \n\n", err)
		}
		fmt.Printf("Get %v \n\n", v)

		vs, err := kv.GetVs("/app/upstream/*")
		if err != nil {
			fmt.Printf("GetMany error %v \n\n", err)
		}

		fmt.Printf("GetMany %v \n\n", vs)
		select {
		case <-quit:
			return
		default:
		}
	}
}
