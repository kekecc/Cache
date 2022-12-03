package myhttp_test

import (
	"MyCache/group"
	"MyCache/myhttp"
	"errors"
	"log"
	"net/http"
	"testing"
)

var m = map[string]string{
	"kk": "dawdaw",
	"cc": "dwadawawd",
	"we": "dawdaw",
}

func Test(t *testing.T) {
	//new一个group
	group.New("student", 100000000, group.HelpFunc(func(key string) ([]byte, error) {
		log.Println("从map里面获取")
		if value, ok := m[key]; ok {
			return []byte(value), nil
		}
		return nil, errors.New("不存在这个数据")
	}))

	IP := "127.0.0.1"
	Port := "3344"
	hp := myhttp.NewHTTP(IP, Port)
	http.ListenAndServe(IP+":"+Port, hp)
}
