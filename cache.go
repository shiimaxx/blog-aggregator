package main

import (
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/shiimaxx/blog-aggregator/structs"
)

type storage interface {
	Get(key string) []structs.Entry
	Set(key string, content []structs.Entry, duration time.Duration)
}

type item struct {
	content    []structs.Entry
	expiration int64
}

type memStorage struct {
	items map[string]item
	mu    *sync.RWMutex
}

func GenerateCacheKey(url, service string) string {
	return fmt.Sprintf("ba:%s:", service) + strings.TrimRight(base64.URLEncoding.EncodeToString([]byte(url)), "=")
}

func (m *memStorage) Get(key string) []structs.Entry {
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

func (m *memStorage) Set(key string, content []structs.Entry, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.items[key] = item{
		content:    content,
		expiration: time.Now().Add(duration).UnixNano(),
	}
}
