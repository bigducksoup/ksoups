package common

import (
	"container/list"
	"sync"
)

type LRUCache[T any] struct {
	capacity int
	cache    map[string]*list.Element
	lruList  *list.List
	mutex    sync.RWMutex
}

func NewLRUCache[T any](capacity int) *LRUCache[T] {
	return &LRUCache[T]{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		lruList:  list.New(),
	}
}

func (l *LRUCache[T]) Put(key string, value T) {

	l.mutex.Lock()
	defer l.mutex.Unlock()

	//key存在
	if element, ok := l.cache[key]; ok {
		l.lruList.MoveToFront(element)
		element.Value = value
		return
	}

	//key不存在
	//容量未满
	if l.lruList.Len() < l.capacity {
		element := l.lruList.PushFront(value)
		l.cache[key] = element
		return
	}

	//容量已满
	if l.lruList.Len() == l.capacity {
		backElement := l.lruList.Back()
		delete(l.cache, key)
		l.lruList.Remove(backElement)
		element := l.lruList.PushFront(value)
		l.cache[key] = element
	}

}

func (l *LRUCache[T]) Get(key string) (T, bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	if element, ok := l.cache[key]; ok {
		l.lruList.MoveToFront(element)
		return element.Value.(T), true
	}

	return *new(T), false

}

func (l *LRUCache[T]) Clear() {

	l.mutex.Lock()
	defer l.mutex.Unlock()
	clear(l.cache)
	l.lruList.Init()
	l.capacity = 0
}

func (l *LRUCache[T]) Range(f func(key string, value T) bool) {

	l.mutex.RLock()
	defer l.mutex.RUnlock()

	for k, e := range l.cache {
		if !f(k, e.Value.(T)) {
			break
		}
	}

}

func (l *LRUCache[T]) Size() int {
	return l.lruList.Len()
}
