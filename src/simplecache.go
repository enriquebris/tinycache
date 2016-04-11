package tinycache

import (
	"fmt"
	"strings"
	"sync"
)

type SimpleCache struct {
	cache map[string]string
	l     sync.RWMutex
}

func NewSimpleCache() (sc SimpleCache, err error) {
	sc = SimpleCache{
		cache: make(map[string]string),
	}
	return sc, nil
}

func (sc *SimpleCache) Exists(key string) (ret bool, err error) {

	if strings.Trim(key, " ") == "" {
		return false, &TinyCacheError{
			msg:      errorEmptyKey,
			emptykey: true,
		}
	}

	_, ret = sc.cache[key]
	return ret, nil
}

func (sc *SimpleCache) Set(key string, value string) (err error) {

	if strings.Trim(key, " ") == "" {
		return fmt.Errorf(errorEmptyKey)
	}

	sc.l.Lock()
	defer sc.l.Unlock()
	sc.cache[key] = value

	return nil
}

func (sc *SimpleCache) Get(key string) (value string, err error) {

	if strings.Trim(key, " ") == "" {
		return "", fmt.Errorf(errorEmptyKey)
	}

	sc.l.RLock()
	defer sc.l.RUnlock()
	v, exists := sc.cache[key]

	if !exists {
		return "", nil
	}

	return v, nil
}

func (sc *SimpleCache) Del(key string) (err error) {

	sc.l.Lock()
	defer sc.l.Unlock()

	ex, err := sc.Exists(key)
	if err != nil {
		return err
	}

	if !ex {
		return &TinyCacheError{
			missingkey: true,
			msg:        missingKeyError(key),
		}
	}

	delete(sc.cache, key)

	return nil
}

func (sc *SimpleCache) Total() (total int) {
	return len(sc.cache)
}
