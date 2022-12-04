package singleflight

import "sync"

type call struct {
	wg    sync.WaitGroup
	value interface{}
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func (g *Group) Do(key string, f func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		c.wg.Wait() //阻塞在这 等待某一个请求完成
		return c.value, nil
	}
	call := new(call)
	call.wg.Add(1)
	g.m[key] = call
	g.mu.Unlock()
	value, err := f()
	if err != nil {
		return nil, err
	}
	call.value = value
	call.wg.Done()

	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()
	return call.value, nil
}
