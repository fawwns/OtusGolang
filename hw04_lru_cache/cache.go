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

// Set - добавиление значения в кэш по ключу.
func (lru *lruCache) Set(key Key, value interface{}) bool {
	item, ok := lru.items[key]
	if ok {
		// Ключ уже есть — обновим значение и передвинем вверх
		item.Value = &cacheItem{
			key:   key,
			value: value,
		}
		lru.queue.MoveToFront(item)
		return true
	}

	// Ключа нет — создаём элемент
	newItem := lru.queue.PushFront(&cacheItem{
		key:   key,
		value: value,
	})
	lru.items[key] = newItem
	// Если переполнен — удаляем последний
	if lru.capacity < lru.queue.Len() {
		last := lru.queue.Back()
		if last != nil {
			lru.queue.Remove(last)
			keyDel := last.Value.(*cacheItem)
			delete(lru.items, keyDel.key)
		}
	}
	return false
}

// Get - получение значения из кэша по ключу.
func (lru *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := lru.items[key]
	if ok {
		lru.queue.MoveToFront(item)
		return item.Value.(*cacheItem).value, true
	}

	return nil, false
}

// Clear - очищает кэш.
func (lru *lruCache) Clear() {
	lru.items = map[Key]*ListItem{}
	lru.queue = NewList()

}
