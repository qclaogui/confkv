package kv

import (
	"path"
	"time"
)

type watchProcessor struct {
	config   conf
	stopChan chan bool
	errChan  chan error
}

func appendPrefix(prefix string, keys []string) []string {
	s := make([]string, len(keys))
	for i, k := range keys {
		s[i] = path.Join(prefix, k)
	}
	return s
}

func (wp *watchProcessor) Process() {
	var waitIndex = uint64(0)
	keys := appendPrefix(wp.config.prefix, wp.config.keys)
	for {
		index, err := wp.config.bs.WatchPrefix(wp.config.prefix, keys, waitIndex, wp.stopChan)
		if err != nil {
			wp.errChan <- err
			// Prevent backend errors from consuming all resources.
			time.Sleep(time.Second * 2)
			continue
		}
		// stopChan 发送信号 index 返回500 表示退出
		if index == 500 {
			return
		}

		waitIndex = index
		result, err := wp.config.bs.GetValues(keys)
		if err != nil {
			wp.errChan <- err
		}
		wp.config.log.Infof("Got the following map from backend: %v\n\n", result)

		// 重新赋值
		mem.Flush()
		mem.PutMany(result)
	}
}
