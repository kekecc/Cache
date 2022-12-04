package main

import (
	"MyCache/group"
	"errors"
	"flag"
	"log"
	"net/http"
)

var m = map[string]string{
	"kk": "ddd",
	"cc": "awd",
	"cs": "dwadawdaw",
}

func CreateGroup() *group.Group {
	return group.New("name", 2<<10, group.HelpFunc(func(key string) ([]byte, error) {
		log.Println("search from DB")
		if value, ok := m[key]; ok {
			return []byte(value), nil
		} else {
			return nil, errors.New("数据库无此信息")
		}
	}))
}

func CacheServer(addr string, addrs []string, g *group.Group) {
	peer := group.NewHTTP("", "", addr)
	peer.Settings(addrs...)
	g.RegisterRoute(peer)
	log.Println("addr:", addr, " up")
	log.Fatal(http.ListenAndServe(addr[7:], peer))
}

func LocalServer(addr string, g *group.Group) {
	http.Handle("/cache", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		value, err := g.Get(key)
		if err != nil {
			http.Error(w, "未查找到该数据", http.StatusBadRequest)
			return
		}
		w.Write(value.ToByte())
	}))
	log.Println("启动本地服务", addr)
	log.Fatal(http.ListenAndServe(addr[7:], nil))
}

func main() {
	var port int
	var local bool
	flag.IntVar(&port, "port", 8001, "输入port")
	flag.BoolVar(&local, "local", false, "本地服务")
	flag.Parse()

	local_addr := "http://127.0.0.1:8999"
	Map := map[int]string{
		8001: "http://127.0.0.1:8001",
		8002: "http://127.0.0.1:8002",
		8003: "http://127.0.0.1:8003",
	}
	var addrs []string
	for _, value := range Map {
		addrs = append(addrs, value)
	}

	g := CreateGroup()

	if local {
		go LocalServer(local_addr, g)
	}
	CacheServer(Map[port], addrs, g)
}
