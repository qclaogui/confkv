package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qclaogui/kv"
)

var prefix = "/app"

// 当然 backend（默认) 中需要有这些配置
var keysHasPrefix = []string{
	"/upstream/host1",
	"/upstream/host2"}

func main() {
	ckv := kv.Watch(kv.Prefix(prefix), kv.Keys(keysHasPrefix),
		kv.Zookeeper("127.0.0.1:2181", "127.0.0.1:2182"))

	// 等待从后端获取配置 然后第一次加载到内存 浪费点启动内存
	time.Sleep(time.Second)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	for range time.Tick(5 * time.Second) {
		vs, err := kv.DB.GetVs("/app/upstream/*")
		if err != nil {
			fmt.Printf("GetVs error %v \n\n", err)
		}
		fmt.Printf("%v \n\n", vs)
		select {
		case <-quit:
			return
		default:
		}
	}

	ckv.Stop()
}
