package group

import (
	consistenthash "MyCache/consistentHash"
	"errors"
	"log"
	"net/http"
	"strings"
	"sync"
)

const API string = "/cache/"

type HTTP struct {
	url             string
	IP              string
	Port            string
	Api             string
	mu              sync.Mutex //保护
	consistent_hash *consistenthash.ConsistentHash
	Clients         map[string]*HTTPClient
}

func NewHTTP(ip string, port string, url string) *HTTP {
	return &HTTP{
		url:     url,
		IP:      ip,
		Port:    port,
		Api:     API,
		Clients: make(map[string]*HTTPClient),
	}
}

func (h *HTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) { //实现ServeHTTP方法
	log.Printf("%s, %s\n", r.Method, r.URL.Path)
	log.Println(r.URL.Path)
	info := strings.Split(r.URL.Path[len(h.Api):], "/")
	if len(info) != 2 {
		log.Println(len(info))
		log.Println(info[0])
		log.Println(info[1])
		//log.Println(info[2])
		http.Error(w, "参数传输有误", http.StatusBadRequest)
		return
	}

	//获取name
	g := GetFromGroupName(info[0])
	if g == nil {
		http.Error(w, "group不存在", http.StatusBadRequest)
		return
	}

	//获取key
	value, err := g.Get(info[1])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(value.ToByte())
}

func (h *HTTP) Settings(urls ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.consistent_hash = consistenthash.NewHash(10, nil)
	h.consistent_hash.AddKey(urls...)
	h.Clients = make(map[string]*HTTPClient, len(urls))
	for _, url := range urls {
		h.Clients[url] = &HTTPClient{url: url + h.Api}
	}
}

func (h *HTTP) Pick(key string) (*HTTPClient, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if peer := h.consistent_hash.GetNode(key); peer != "" && peer != h.url {
		log.Println("选择了节点", peer)
		return h.Clients[peer], nil
	}
	return nil, errors.New("找不到节点")
}

//var _ Picker = (*HTTP)(nil)
