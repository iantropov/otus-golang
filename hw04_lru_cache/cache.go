package hw04lrucache

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
	i, exists := lc.items[key]
	if exists {
		lc.queue.MoveToFront(i)
		cacheItem := i.Value.(cacheItem)
		cacheItem.value = value
		return true
	}

	newCacheItem := &cacheItem{key, value}
	i = lc.queue.PushFront(newCacheItem)
	if lc.queue.Len() > lc.capacity {
		back := lc.queue.Back()
		backCacheItem := back.Value.(cacheItem)

		lc.queue.Remove(back)
		delete(lc.items, backCacheItem.key)
	}

	lc.items[key] = i

	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	i, exists := lc.items[key]
	if !exists {
		return nil, false
	}

	lc.queue.MoveToFront(i)
	cacheItem := i.Value.(cacheItem)
	return cacheItem.value, true
}

func (lc *lruCache) Clear() {
	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}
