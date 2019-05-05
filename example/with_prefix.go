package main

import (
	"fmt"
	"time"

	"github.com/qclaogui/confkv"
)

func main() {
	var prefix = "/app"
	var keysHasPrefix = []string{
		"/upstream/host1",
		"/upstream/host2"}

	defer confkv.Watch(confkv.Prefix(prefix), confkv.Keys(keysHasPrefix)).Stop()
	time.Sleep(time.Second)

	vs, err := confkv.Store.GetAllValues("/app/upstream/*")
	if err != nil {
		fmt.Printf("GetAllValues error %v \n\n", err)
	}

	fmt.Printf("%v \n\n", vs)
}
