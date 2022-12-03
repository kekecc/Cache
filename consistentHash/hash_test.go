package consistenthash_test

import (
	consistenthash "MyCache/consistentHash"
	"log"
	"strconv"
	"testing"
)

func TestHash(t *testing.T) {
	hash_circle := consistenthash.NewHash(4, func(b []byte) uint32 {
		num, _ := strconv.Atoi(string(b))
		return uint32(num)
	})

	hash_circle.AddKey("1", "3", "5")
	// 1 3 5 11 13 15 21 23 25

	m := map[string]string{
		"14": "5",
		"2":  "3",
		"16": "1",
	}

	for key, value := range m {
		if s := hash_circle.GetNode(key); s != value {
			log.Println(key, s)
			t.Fatalf("键值对不匹配")
		}
	}
}
