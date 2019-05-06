package backends

import (
	"errors"

	"github.com/qclaogui/kv/backends/zookeeper"
)

type BackendStore interface {
	GetValues(keys []string) (map[string]string, error)
	WatchPrefix(prefix string, keys []string, waitIndex uint64, stopChan chan bool) (uint64, error)
}

// BackendConfig
type BackendConfig struct {
	Backend      string
	BackendNodes []string
}

var defaultConfig = &BackendConfig{
	Backend:      "zookeeper",
	BackendNodes: []string{"127.0.0.1:2181"},
}

func New(config *BackendConfig) (BackendStore, error) {
	if config == nil {
		config = defaultConfig
	} else {
		if len(config.BackendNodes) < 1 {
			config = defaultConfig
		}
	}

	switch config.Backend {
	case "zookeeper":
		return zookeeper.NewZookeeperClient(config.BackendNodes)
	}
	return nil, errors.New("Invalid backend")
}
