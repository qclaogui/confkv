package main

import (
	"fmt"
	"time"

	"github.com/qclaogui/kv"
)

func main() {
	var prefix = "/app"
	var keysHasPrefix = []string{
		"/upstream/host1",
		"/upstream/host2"}

	defer kv.Watch(kv.Prefix(prefix), kv.Keys(keysHasPrefix)).Stop()
	time.Sleep(time.Second)

	vs, err := kv.DB.GetAllValues("/app/upstream/*")
	if err != nil {
		fmt.Printf("GetAllValues error %v \n\n", err)
	}

	fmt.Printf("%v \n\n", vs)
}
