package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	mutex    sync.RWMutex
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	lc.mutex.RLock()
	i, exists := lc.items[key]
	lc.mutex.RUnlock()

	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	if exists {
		lc.queue.MoveToFront(i)
		cacheItem := getCacheItem(i)
		cacheItem.value = value
		return true
	}

	newCacheItem := &cacheItem{key, value}
	i = lc.queue.PushFront(newCacheItem)
	if lc.queue.Len() > lc.capacity {
		back := lc.queue.Back()
		backCacheItem := getCacheItem(back)

		lc.queue.Remove(back)
		delete(lc.items, backCacheItem.key)
	}

	lc.items[key] = i

	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	lc.mutex.RLock()
	i, exists := lc.items[key]
	lc.mutex.RUnlock()
	if !exists {
		return nil, false
	}

	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.queue.MoveToFront(i)
	cacheItem := getCacheItem(i)
	return cacheItem.value, true
}

func (lc *lruCache) Clear() {
	lc.mutex.Lock()
	defer lc.mutex.Unlock()

	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}

func getCacheItem(i *ListItem) *cacheItem {
	return i.Value.(*cacheItem)
}
