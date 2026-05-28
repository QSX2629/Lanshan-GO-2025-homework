package context

import (
	"sync"
	"time"
)

type ChatContext struct {
	Messages []string
	Updated  time.Time
}

var store = make(map[string]*ChatContext)
var mu sync.RWMutex

func GetContext(key string) *ChatContext {
	mu.RLock()
	defer mu.RUnlock()
	return store[key]
}

func AppendContext(key, msg string) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := store[key]; !ok {
		store[key] = &ChatContext{}
	}
	store[key].Messages = append(store[key].Messages, msg)
	store[key].Updated = time.Now()
}

func FlushContext(key string) {
	mu.Lock()
	defer mu.Unlock()
	delete(store, key)
}
