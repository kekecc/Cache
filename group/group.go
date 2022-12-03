package group

import (
	"errors"
	"log"
	"sync"
)

type CallBack interface {
	Get(key string) ([]byte, error)
}

type HelpFunc func(key string) ([]byte, error) //这是传入的函数 实现了CallBack接口

func (h HelpFunc) Get(key string) ([]byte, error) {
	return h(key) //调用自己
}

type Group struct {
	name      string
	call_back CallBack //注册回调函数
	mainCache Cache
}

var (
	mutex  sync.RWMutex
	groups = make(map[string]*Group) //全局的group
)

func New(name string, cap int64, call CallBack) *Group {
	if call == nil {
		panic("必须有一个回调函数")
	}
	mutex.Lock()
	defer mutex.Unlock()
	g := &Group{
		name:      name,
		call_back: call,
		mainCache: Cache{
			capacity: cap,
		},
	}
	groups[name] = g
	return g
}

func (g *Group) Del(key string) (Byte, error) {
	value, err := g.mainCache.Del(key)
	if err != nil {
		log.Println("删除不存在的缓存")
		return Byte{}, err
	}
	return value, nil
}

func (g *Group) Get(key string) (Byte, error) {
	if value, err := g.mainCache.Get(key); err == nil {
		log.Println("cache hit")
		return value, err
	}
	if key == "" {
		return Byte{}, errors.New("需要一个完整的key")
	}
	//未命中则需要调用回调函数
	return g.GetFromAnother(key)
}

func Copy(b []byte) []byte {
	temp := make([]byte, len(b))
	copy(temp, b)
	return temp
}

func (g *Group) GetFromAnother(key string) (Byte, error) {
	return g.FromLocal(key)
}

func (g *Group) FromLocal(key string) (Byte, error) {
	bytes, err := g.call_back.Get(key)
	if err != nil {
		return Byte{}, err
	}
	value := Byte{data: Copy(bytes)}
	g.mainCache.Add(key, value)
	return value, nil
}

func GetFormGroupName(name string) *Group {
	//map不是并发安全的
	mutex.RLock()
	g := groups[name]
	mutex.RUnlock()
	return g
}
