package lru

import "testing"

// 为string 实现Len方法
type String string

func (s String) Len() int {
	return len(s)
}

// get函数
func TestGetElement(t *testing.T) {
	lru_cache := NewCache(int64(20))
	lru_cache.AddElement("name", String("kk"))
	value, ok := lru_cache.GetElement("name")
	if !ok {
		t.Fatalf("没有正确获取到kv")
	}
	if string(value.(String)) != "kk" {
		t.Fatalf("获取到的value不正确")
	}

	if _, ok = lru_cache.GetElement("id"); ok {
		t.Fatalf("出现了不存在的key")
	}
}

func TestRemoveTheLast(t *testing.T) {
	cap := int64(len("name" + "kk" + "number" + "U202111248"))
	lru_cache := NewCache(cap)
	lru_cache.AddElement("name", String("kk")) //最近最少使用
	lru_cache.AddElement("number", String("U202111248"))
	lru_cache.AddElement("id", String("100"))
	if _, ok := lru_cache.GetElement("name"); ok {
		t.Fatalf("没有正确找到最近最少使用的结点")
	}
	if lru_cache.Length() != 2 { //还有几个结点
		t.Fatalf("误删了某些数据")
	}
	lru_cache.AddElement("dwad", String("1234"))
	if _, ok := lru_cache.GetElement("number"); ok {
		t.Fatalf("xx没有正确找到删除的结点")
	}
	if lru_cache.Length() != 2 { //还有几个结点
		t.Fatalf("误删了某些数据")
	}
}
