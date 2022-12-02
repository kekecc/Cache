package mycache

import (
	lru "MyCache/LRU"
	"errors"
	"sync"
)

type Cache struct {
	mutex    sync.Mutex //互斥锁
	lru      *lru.Cache
	capacity int64
}

func (c *Cache) Add(key string, value Byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		c.lru = lru.NewCache(c.capacity)
	}
	c.lru.AddElement(key, value) //差点多加了&
}

func (c *Cache) Get(key string) (Byte, error) {
	c.mutex.Lock()
	//var return_byte Byte
	defer c.mutex.Unlock()
	if c.lru == nil {
		c.lru = lru.NewCache(c.capacity)
	}
	if value, ok := c.lru.GetElement(key); ok {
		return value.(Byte), nil
	}
	return Byte{}, errors.New("出现某些问题")
}

func (c *Cache) Del(key string) (Byte, error) { //删除数据
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if c.lru == nil {
		c.lru = lru.NewCache(c.capacity)
	}
	if value, ok := c.lru.RmElememt(key); ok {
		return value.(Byte), nil
	}
	return Byte{}, errors.New("不存在这个元素")
}
