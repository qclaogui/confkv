package main

import (
	"fmt"
	"time"

	"github.com/qclaogui/confkv"
)

var keys = []string{
	"/app/upstream/host1",
	"/app/upstream/host2",
}

func main() {
	defer confkv.Watch(confkv.Keys(keys)).Stop()
	time.Sleep(time.Second)

	vs, err := confkv.Store.GetAllValues("/app/upstream/*")
	if err != nil {
		fmt.Printf("GetAllValues error %v \n\n", err)
	}
	fmt.Printf("%v \n\n", vs)

}
