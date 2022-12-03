package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type HashFunc func([]byte) uint32

type consistentHash struct {
	hashFunc   HashFunc
	virtualNum int
	keys       []int
	hashMap    map[int]string
}

func NewHash(num int, hashFunc HashFunc) *consistentHash {
	ch := &consistentHash{
		hashFunc:   hashFunc,
		virtualNum: num,
		hashMap:    make(map[int]string),
	}
	if hashFunc == nil {
		ch.hashFunc = crc32.ChecksumIEEE
	}
	return ch
}

func (c *consistentHash) IsEmpty() bool { //if any items available
	return len(c.keys) == 0
}

func (c *consistentHash) AddKey(keys ...string) {
	for _, key := range keys {
		for i := 0; i < c.virtualNum; i++ {
			num := int(c.hashFunc([]byte(strconv.Itoa(i) + key)))
			c.keys = append(c.keys, num)
			//维护map
			c.hashMap[num] = key
		}
	}
	sort.Ints(c.keys)
}

// 选择节点
func (c *consistentHash) GetNode(key string) string {
	if c.IsEmpty() {
		return ""
	}
	num := int(c.hashFunc([]byte(key)))
	index := sort.Search(len(c.keys), func(n int) bool {
		return c.keys[n] >= num
	})
	if index == len(c.keys) {
		index = 0
	}
	return c.hashMap[c.keys[index]]
}
