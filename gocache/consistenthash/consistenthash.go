package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	hash     Hash           //hash函数
	replicas int            //虚拟节点的倍数
	keys     []int          //哈希环
	hashMap  map[int]string //虚拟节点与真实节点的映射表,键为虚拟节点哈希值，值为真实节点的名称
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add 对每个真实节点key，对应创建m.replicas个虚拟节点，虚拟节点的名称是:strconv.Itoa(i) + key
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			//使用m.hash()计算哈希值，使用append添加到环上
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			//增加虚拟节点和真实节点的映射关系
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}
