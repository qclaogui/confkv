package kv

import (
	"strings"

	"github.com/qclaogui/kv/log"

	"github.com/qclaogui/kv/backends"
)

var Options options

type options struct{}

func (options) WithLog(l log.Logger) func(*conf) {
	return func(c *conf) { c.log = l }
}

// Zookeeper 设置后端KV获取位置为Zookeeper
func (options) Zookeeper(nodes ...string) func(*conf) {
	cnf := backends.BackendConfig{"zookeeper", parseNodes(nodes...)}
	bs, err := backends.New(&cnf)
	if err != nil {
		panic(err.Error())
	}
	return func(c *conf) { c.bs = bs }
}

func parseNodes(nodes ...string) (nodesAddrs []string) {
	for _, node := range nodes {
		for _, addr := range strings.Split(node, ",") {
			if strings.Index(addr, ":") == -1 {
				addr += ":2181"
			}
			nodesAddrs = append(nodesAddrs, addr)
		}
	}
	return nodesAddrs
}
