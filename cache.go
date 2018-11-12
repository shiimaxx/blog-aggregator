package main

import (
	"sync"
	"time"
)

type storage interface {
	Get(key string) []byte
	Set(key string, content []byte, duration time.Duration)
}

type item struct {
	content    []byte
	expiration int64
}

type memStorage struct {
	items map[string]item
	mu    *sync.RWMutex
}

func (m *memStorage) Get(key string) []byte {
	m.mu.RLock()
	defer m.mu.RUnlock()

	i, ok := m.items[key]
	if !ok {
		return nil
	}
	if time.Now().UnixNano() > i.expiration {
		return nil
	}

	return i.content
}

func (m *memStorage) Set(key string, content []byte, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[key] = item{
		content:    content,
		expiration: time.Now().Add(duration).UnixNano(),
	}
}
