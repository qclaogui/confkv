package main

import (
	"fmt"

	"github.com/qclaogui/confkv"
)

var prefix = "/app"

// 当然 backend（默认) 中需要有这些配置
var keysHasPrefix = []string{
	"/upstream/host1",
	"/upstream/host2"}

func main() {
	ckv := confkv.Watch(confkv.Prefix(prefix), confkv.Keys(keysHasPrefix),
		confkv.Zookeeper("127.0.0.1:2181", "127.0.0.1:2182"))

	vs, err := confkv.Store.GetAllValues("/app/upstream/*")
	if err != nil {
		fmt.Printf("GetAllValues error %v \n\n", err)
	}
	fmt.Printf("%v \n\n", vs)

	ckv.Stop()
}
