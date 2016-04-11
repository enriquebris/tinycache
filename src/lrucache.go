package tinycache

import (
	"container/list"
	"fmt"
	"strings"
	"sync"
)

type LRUCache struct {
	cache   map[string]cacheValue
	clist   *list.List
	l       sync.RWMutex
	limit   Limit
	current Limit
}

type cacheValue struct {
	value   string
	element *list.Element
}

// Returns a new LRUCache instance
func NewLRUCache(limit Limit) (lru LRUCache) {
	lru = LRUCache{
		cache:   make(map[string]cacheValue),
		clist:   list.New(),
		limit:   limit,
		current: Limit{},
	}

	return lru
}

func (lc *LRUCache) Exists(key string) (ret bool, err error) {

	if strings.Trim(key, " ") == "" {
		return false, &TinyCacheError{
			msg:      errorEmptyKey,
			emptykey: true,
		}
	}

	_, ret = lc.cache[key]
	return ret, nil
}

func (lc *LRUCache) Set(key string, value string) (err error) {

	// possible cases

	ex, err := lc.Exists(key)
	if err != nil {
		return err
	}

	lc.l.Lock()
	defer lc.l.Unlock()

	// length to add
	valueSize := len(value)

	// only adds the new key / value if the value's size is under the limits
	if valueSize <= lc.limit.SizeInBytes {

		// check the limits && clean whether it is needed
		lc.cleanBySize(valueSize)
		lc.cleanByTotalElements(1)

		// the key already exists, so let's replace the value and move the element to the last position into the list
		if ex {
			cv := lc.cache[key]
			// remove the size of the old value
			lc.addSizeInBytes(-len(cv.value))
			cv.value = value
			// add the size of the new value
			lc.addSizeInBytes(valueSize)
			lc.cache[key] = cv

			lc.clist.MoveToBack(cv.element)

		} else {
			// it is a new key ==> add the value and push the element into the list (as the last element)

			// add size & element
			lc.addSizeInBytes(valueSize)
			lc.addElements(1)

			cv := cacheValue{
				value: value,
			}
			// put the element in the last position
			cv.element = lc.clist.PushBack(key)

			lc.cache[key] = cv

		}

		return nil
	}

	return &TinyCacheError{
		msg:         fmt.Sprintf(errorExceedLimit, key),
		exceedlimit: true,
	}

}

func (lc *LRUCache) Get(key string) (value string, err error) {

	var ex bool
	ex, err = lc.Exists(key)
	if err != nil {
		// return error coming from lc.Exists
		return "", err
	}

	if !ex {
		// return No Exists error
		return "", &TinyCacheError{
			msg:        missingKeyError(key),
			missingkey: true,
		}
	}

	// the requested element exists
	lc.l.Lock()
	defer lc.l.Unlock()

	cv := lc.cache[key]
	// flag the element as the 'Most Recently Used' ==> move it to the first position onto the list
	lc.clist.MoveToFront(cv.element)

	return cv.value, nil
}

func (lc *LRUCache) Del(key string) (err error) {
	lc.l.Lock()
	defer lc.l.Unlock()

	ex, err := lc.Exists(key)
	if err != nil {
		return err
	}

	if !ex {
		return &TinyCacheError{
			missingkey: true,
			msg:        missingKeyError(key),
		}
	}

	cv := lc.cache[key]
	lc.clist.Remove(cv.element)
	delete(lc.cache, key)

	return nil
}

func (lc *LRUCache) Total() (total int) {
	return len(lc.cache)
}

// Adds size to *LRUCache.current.SizeInBytes
func (lc *LRUCache) addSizeInBytes(size int) {
	lc.current.SizeInBytes = lc.current.SizeInBytes + size
}

// Adds elements to *LRUCache.current.TotalElements
func (lc *LRUCache) addElements(n int) {
	lc.current.TotalElements = lc.current.TotalElements + n
}

// Remove the Least Recently Used element (the last onto lc.clist)
func (lc *LRUCache) removeLast() bool {
	// get the Least Recently Used element
	last := lc.clist.Back()

	if last != nil {
		// get the Least Recently Used key
		key := last.Value.(string)
		// get the associated values
		cv := lc.cache[key]

		// decrement element && size
		lc.addElements(-1)
		lc.addSizeInBytes(-len(cv.value))

		// delete entry in cache map && element in linked list
		delete(lc.cache, key)
		lc.clist.Remove(last)

		return true
	}

	return false
}

// Removes the needed Least Recently Used elements in order to have the used memory under the limits (lc.limit.SizeInBytes)
func (lc *LRUCache) cleanBySize(plusSize int) {
	for lc.current.SizeInBytes+plusSize > lc.limit.SizeInBytes {
		if !lc.removeLast() {
			// break the loop if there is not any last element to remove
			break
		}
	}
}

// Removes the needed Least Recently Used elements in order to not have more elements than the limit (lc.limit.TotalElements)
func (lc *LRUCache) cleanByTotalElements(plusElements int) {
	for lc.current.TotalElements+plusElements > lc.limit.TotalElements {
		if !lc.removeLast() {
			// break the loop if there is not any last element to remove
			break
		}
	}
}
