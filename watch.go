package confkv

import (
	"log"
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
			log.Printf("WatchPrefix error: %v\n", err)
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
		log.Printf("Key prefix set to %v\n", wp.config.prefix)
		result, err := wp.config.bs.GetValues(keys)
		if err != nil {
			wp.errChan <- err
		}
		log.Printf("Got the following map from backend: %v\n\n", result)

		// 重新赋值
		Store.Purge()
		for k, v := range result {
			Store.Set(k, v)
		}
	}
}
