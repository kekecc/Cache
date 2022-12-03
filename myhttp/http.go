package myhttp

import (
	"MyCache/group"
	"log"
	"net/http"
	"strings"
)

const API = "/cache"

type HTTP struct {
	IP   string
	Port string
	Api  string
}

func NewHTTP(ip string, port string) *HTTP {
	return &HTTP{
		IP:   ip,
		Port: port,
		Api:  API,
	}
}

func (h *HTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) { //实现ServeHTTP方法
	log.Printf("%s, %s\n", r.Method, r.URL.Path)
	log.Println(r.URL.Path)
	info := strings.Split(r.URL.Path[len(h.Api):], "/")
	if len(info) != 3 {
		log.Println(len(info))
		log.Println(info[0])
		log.Println(info[1])
		log.Println(info[2])
		http.Error(w, "参数传输有误", http.StatusBadRequest)
		return
	}

	//获取name
	g := group.GetFormGroupName(info[1])
	if g == nil {
		http.Error(w, "group不存在", http.StatusBadRequest)
		return
	}

	//获取key
	value, err := g.Get(info[2])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(value.ToByte())
}
