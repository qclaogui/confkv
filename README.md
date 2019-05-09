
> This project is under development, do not use in production!

<div align="center">
  <h1>kv</h1>
</div>

</div>
<p align="center">
<a href="https://travis-ci.org/qclaogui/kv"><img src="https://travis-ci.org/qclaogui/kv.svg?branch=master"></a>
<a href="https://goreportcard.com/report/github.com/qclaogui/kv"><img src="https://goreportcard.com/badge/github.com/qclaogui/kv?v=1" /></a>
<a href="https://godoc.org/github.com/qclaogui/kv"><img src="https://godoc.org/github.com/qclaogui/kv?status.svg"></a>
<a href="https://github.com/qclaogui/kv/blob/master/LICENSE"><img src="https://img.shields.io/github/license/qclaogui/kv.svg" alt="License"></a>
</p>


## Usage

```go
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

```

深度参考 [confd](https://github.com/kelseyhightower/confd) 的实现。