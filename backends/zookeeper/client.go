package zookeeper

import (
	"log"
	"path/filepath"
	"strings"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

// Client provides a wrapper around the zookeeper client
type Client struct{ client *zk.Conn }

// NewZookeeperClient new connection to zookeeper
func NewZookeeperClient(nodes []string) (*Client, error) {
	c, _, err := zk.Connect(nodes, time.Second)
	if err != nil {
		panic(err)
	}
	return &Client{c}, nil
}
func nodeWalk(prefix string, c *Client, vars map[string]string) error {
	var s string
	l, stat, err := c.client.Children(prefix)
	if err != nil {
		return err
	}
	if stat.NumChildren == 0 {
		b, _, err := c.client.Get(prefix)
		if err != nil {
			return err
		}
		vars[prefix] = string(b)
	} else {
		for _, key := range l {
			if prefix == "/" {
				s = "/" + key
			} else {
				s = prefix + "/" + key
			}
			_, stat, err := c.client.Exists(s)
			if err != nil {
				return err
			}
			if stat.NumChildren == 0 {
				b, _, err := c.client.Get(s)
				if err != nil {
					return err
				}
				vars[s] = string(b)
			} else {
				nodeWalk(s, c, vars)
			}
		}
	}
	return nil
}
func (c *Client) GetValues(keys []string) (map[string]string, error) {
	vars := make(map[string]string)
	for _, v := range keys {
		v = strings.Replace(v, "/*", "", -1)
		_, _, err := c.client.Exists(v)
		if err != nil {
			return vars, err
		}

		err = nodeWalk(v, c, vars)
		if err != nil {
			return vars, err
		}
	}

	return vars, nil
}

type watchResponse struct {
	waitIndex uint64
	err       error
}

func (c *Client) WatchPrefix(prefix string, keys []string, waitIndex uint64, stopChan chan bool) (uint64, error) {
	log.Printf("waitIndex: %v\n", waitIndex)
	// return something > 0 to trigger a key retrieval from the store
	if waitIndex == 0 {
		return 1, nil
	}
	// List the children first
	entries, err := c.GetValues([]string{prefix})
	//log.Printf("entries: %v\n", entries)
	if err != nil {
		log.Printf("c.GetValues error: %v\n", err)
		return 0, err
	}

	respChan := make(chan watchResponse)
	cancel := make(chan bool)
	defer close(cancel)

	//watch all sub folders for changes
	watchMap := make(map[string]string)
	for k, _ := range entries {
		for _, v := range keys {
			if strings.HasPrefix(k, v) {
				for dir := filepath.Dir(k); dir != "/"; dir = filepath.Dir(dir) {
					if _, ok := watchMap[dir]; !ok {
						watchMap[dir] = ""
						log.Printf("Watching: %v\n", dir)
						go c.watch(dir, respChan, cancel)
					}
				}
				break
			}
		}
	}

	//watch all keys in prefix for changes
	for k, _ := range entries {
		for _, v := range keys {
			if strings.HasPrefix(k, v) {
				log.Printf("Watching: %v\n", k)
				go c.watch(k, respChan, cancel)
				break
			}
		}
	}

	for {
		select {
		case <-stopChan:
			log.Println("stopChan got sig")
			return 500, nil
		case r := <-respChan:
			return r.waitIndex, r.err
		}
	}
}

func (c *Client) watch(key string, respChan chan watchResponse, cancel chan bool) {
	_, _, keyEventCh, err := c.client.GetW(key)
	if err != nil {
		respChan <- watchResponse{0, err}
	}
	_, _, childEventCh, err := c.client.ChildrenW(key)
	if err != nil {
		respChan <- watchResponse{0, err}
	}

	for {
		select {
		case e := <-keyEventCh:
			if e.Type == zk.EventNodeDataChanged {
				respChan <- watchResponse{3, e.Err}
			}
		case e := <-childEventCh:
			if e.Type == zk.EventNodeChildrenChanged {
				respChan <- watchResponse{4, e.Err}
			}
		case <-cancel:
			log.Printf("Stop Watching: %v\n", key)
			// There is no way to stop GetW/ChildrenW so just quit
			return
		}
	}
}
