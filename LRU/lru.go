package lru

import (
	"container/list"
	"log"
)

/*
type List struct {
	Prior *List
	Next  *List
	val   interface{}
}
*/

type Cache struct { //隐藏
	maxSize int64
	nowSize int64
	lists   *list.List //双向链表
	reflect map[string]*list.Element
}

type Entry struct { //键值对
	key   string
	value Value
}

type Value interface {
	Len() int
}

func NewCache(maxSize int64) *Cache {
	return &Cache{
		maxSize: maxSize,
		nowSize: 0,
		lists:   list.New(),
		reflect: make(map[string]*list.Element),
	}
}

func (c *Cache) GetElement(key string) (Value, bool) {
	if element, ok := c.reflect[key]; ok {
		//访问过 移动到前面
		c.lists.MoveToFront(element)         //放到前列
		kvalue, ok := element.Value.(*Entry) //类型断言
		if !ok {
			log.Println("类型断言失败!")
			return nil, false
		}
		return kvalue.value, true
	}
	return nil, false
}

func (c *Cache) RemoveTheLast() {
	element := c.lists.Back()
	if element != nil {
		c.lists.Remove(element)
		kvalue := element.Value.(*Entry)
		delete(c.reflect, kvalue.key)
		c.nowSize -= int64(len(kvalue.key)) + int64(kvalue.value.Len())
	}
}

func (c *Cache) AddElement(key string, value Value) {
	if element, ok := c.reflect[key]; ok { //已经有这个key的话
		c.lists.MoveToFront(element)
		kvalue := element.Value.(*Entry)
		c.nowSize += int64(value.Len()) - int64(kvalue.value.Len()) //修改
		kvalue.value = value
	} else {
		ele := c.lists.PushFront(&Entry{key: key, value: value})
		c.reflect[key] = ele
		c.nowSize += int64(len(key)) + int64(value.Len())
	}
	for {
		if c.maxSize != 0 && c.nowSize > c.maxSize { //淘汰
			c.RemoveTheLast()
		} else {
			break
		}
	}
}

func (c *Cache) Length() int {
	return c.lists.Len()
}
