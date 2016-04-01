package main

import (
	"fmt"
	"strings"
	"sync"
)

type tinyCache struct {
	cache map[string]string
	l     sync.RWMutex
}

func newTinyCache() (sc tinyCache, err error) {
	sc = tinyCache{
		cache: make(map[string]string),
	}
	return sc, nil
}

func (sc *tinyCache) set(key string, value string) error {

	if strings.Trim(key, " ") == "" {
		return fmt.Errorf("simpleCache.set(key, value) needs a key different than an empty string")
	}

	sc.l.Lock()
	defer sc.l.Unlock()
	sc.cache[key] = value

	return nil
}

func (sc *tinyCache) get(key string) (value string, err error) {

	if strings.Trim(key, " ") == "" {
		return "", fmt.Errorf("simpleCache.get(key string) needs a key different than an empty string")
	}

	sc.l.RLock()
	defer sc.l.RUnlock()
	v, exists := sc.cache[key]

	if !exists {
		return "", nil
	}

	return v, nil
}

func main() {
	fmt.Println("Hello, playground")
}
