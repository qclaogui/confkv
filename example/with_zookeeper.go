package main

import (
	"fmt"
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
	time.Sleep(time.Second)

	vs, err := kv.DB.GetAllValues("/app/upstream/*")
	if err != nil {
		fmt.Printf("GetAllValues error %v \n\n", err)
	}
	fmt.Printf("%v \n\n", vs)

	ckv.Stop()
}
