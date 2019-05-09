package main

import (
	"fmt"
	"time"

	"github.com/qclaogui/kv"
)

var keys = []string{
	"/app/upstream/host1",
	"/app/upstream/host2",
}

func main() {
	defer kv.Watch(kv.Keys(keys)).Stop()
	time.Sleep(time.Second)

	vs, err := kv.Store().GetMany("/app/upstream/*")
	if err != nil {
		fmt.Printf("GetMany error %v \n\n", err)
	}

	fmt.Printf("%v \n\n", vs)
}
