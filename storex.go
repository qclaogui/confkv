package kv

import (
	"path"
	"sort"
	"sync"
)

type StoreX struct{ m *sync.Map }

func NewStoreX() *StoreX { return &StoreX{new(sync.Map)} }

func (s *StoreX) Set(key, value string) { s.m.Store(key, KVPair{key, value}) }

func (s *StoreX) Exists(key string) bool {
	if _, err := s.Get(key); err != nil {
		return false
	}
	return true
}

func (s *StoreX) Get(key string) (KVPair, error) {
	if v, ok := s.m.Load(key); !ok {
		return KVPair{}, ErrNotExist
	} else {
		return v.(KVPair), nil
	}
}

func (s *StoreX) GetValue(key string, defaultValue ...string) (string, error) {
	kv, err := s.Get(key)
	if err != nil {
		// 如果有设置默认值,将返回defaultValue中的第一个作为默认值
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return "", err
	}
	return kv.Value, nil
}

func (s *StoreX) GetAll(pattern string) (KVPairs, error) {
	kvs := make(KVPairs, 0)
	s.m.Range(func(_, value interface{}) bool {
		kv := value.(KVPair)
		if matched, _ := path.Match(pattern, kv.Key); matched {
			kvs = append(kvs, kv)
		}
		return true
	})
	// 查看是否匹配到
	if len(kvs) == 0 {
		return nil, ErrNoMatch
	}
	sort.Sort(kvs)
	return kvs, nil
}

func (s *StoreX) GetAllValues(pattern string) ([]string, error) {
	vs := make([]string, 0)
	kvs, err := s.GetAll(pattern)
	if err != nil {
		return nil, err
	}
	for _, kv := range kvs {
		vs = append(vs, kv.Value)
	}
	sort.Strings(vs)
	return vs, nil
}

func (s *StoreX) Del(key string) { s.m.Delete(key) }

func (s *StoreX) Purge() {
	s.m.Range(func(key, _ interface{}) bool {
		s.m.Delete(key)
		return true
	})
}
