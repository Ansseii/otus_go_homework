package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mutex    sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	node, ok := l.items[key]

	if ok {
		node.Value = cacheItem{key: key, value: value}
		l.queue.MoveToFront(node)
	} else {
		if l.queue.Len() == l.capacity {
			last := l.queue.Back()
			l.queue.Remove(last)
			delete(l.items, last.Value.(cacheItem).key)
		}
		node = l.queue.PushFront(cacheItem{key: key, value: value})
		l.items[key] = node
	}

	return ok
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	node, ok := l.items[key]

	if ok {
		l.queue.MoveToFront(node)
		return node.Value.(cacheItem).value, true
	}

	return nil, false
}

func (l *lruCache) Clear() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.queue = NewList()
	l.items = make(map[Key]*ListItem, l.capacity)
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
