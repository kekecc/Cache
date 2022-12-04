package group

import (
	singleflight "MyCache/singleFlight"
	"errors"
	"log"
	"sync"
	"time"
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
	route     *HTTP
	flight    *singleflight.Group
}

var global_map = make(map[string]Byte)

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
		flight: &singleflight.Group{},
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

func (g *Group) RegisterRoute(route *HTTP) {
	if g.route != nil {
		panic("group已经注册路由")
	}
	g.route = route
}

func (g *Group) GetFromAnother(key string) (Byte, error) {
	data, err := g.flight.Do(key, func() (interface{}, error) {
		if g.route != nil {
			if peer, err := g.route.Pick(key); err != nil { //获取对应的客户端
				if bytes, err := peer.Get(g.name, key); err == nil {
					return Byte{data: bytes}, nil
				}
				log.Println("远程节点不存在该数据")
			}
		}
		return g.FromLocal(key)
	})
	if err == nil {
		return data.(Byte), err
	}

	return Byte{}, errors.New("找不到该数据")
}

func (g *Group) FromLocal(key string) (Byte, error) {
	bytes, err := g.call_back.Get(key)
	if err != nil {
		return Byte{}, err
	}
	value := Byte{data: Copy(bytes)}
	g.mainCache.Add(key, value) //添加
	return value, nil
}

func GetFromGroupName(name string) *Group {
	//map不是并发安全的
	mutex.RLock()
	g := groups[name]
	mutex.RUnlock()
	return g
}

func (g *Group) Update(key string, value Byte) error {
	mutex.Lock()
	global_map[key] = value
	mutex.Unlock()
	g.mainCache.Add(key, value)
	return nil
}

func (g *Group) Golang() {
	for {
		time.Sleep(1000 * time.Hour)
		for key, value := range global_map {
			//写回逻辑
			log.Println("更新：", key, value)
			delete(global_map, key)
		}
	}
}
