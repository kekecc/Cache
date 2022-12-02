package mycache

import (
	"errors"
	"reflect"
	"testing"
)

func TestCallBack(t *testing.T) {
	var f CallBack = HelpFunc(func(key string) ([]byte, error) { return []byte("kk"), nil })

	var temp = []byte("kk")
	if data, _ := f.Get("kk"); !reflect.DeepEqual(temp, data) {
		t.Fatalf("测试回调函数出错")
	}
}

func TestGroup(t *testing.T) {
	m := map[string]string{
		"kk": "udwadaw",
		"cs": "123",
	}
	count := make(map[string]int, 2)
	g := New("mygroup", 100000000, HelpFunc(func(key string) ([]byte, error) {
		if value, ok := m[key]; ok {
			if _, OK := count[key]; !OK {
				count[key] = 0
			}
			count[key]++
			return []byte(value), nil
		}
		return nil, errors.New("找不到该数据")
	}))

	for k, v := range m {
		if value, err := g.Get(k); err != nil || value.ToString() != v {
			t.Fatalf("从外地获取数据错误")
		}
		if _, err := g.Get(k); err != nil || count[k] > 1 {
			t.Fatalf("获取缓存错误")
		}
	}
	if value, err := g.Del("kk"); err != nil {
		t.Fatalf("删除结点数据%s失败", value)
	}
	if value, err := g.Get("kk"); err != nil || count["kk"] != 2 {
		t.Fatalf("删除结点数据%s失败", value)
	}
}
