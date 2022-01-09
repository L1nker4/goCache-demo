package lru

import "container/list"

type Cache struct {
	maxBytes  int64                         //允许使用的最大内存
	nBytes    int64                         //当前已使用内存
	ll        *list.List                    //存放所有值
	cache     map[string]*list.Element      //键值对
	OnEvicted func(key string, value Value) //某条记录被移除时的回调函数
}

type entry struct {
	key   string //键
	value Value  //值
}

type Value interface {
	Len() int //返回值所占用的内存大小
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, status bool) {
	if ele, status := c.cache[key]; status {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	//如果元素存在
	if ele != nil {
		//list中删除
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		//map中删除该节点的映射关系
		delete(c.cache, kv.key)
		//修改nBytes
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		//检查是否需要callback
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add 新增/修改
func (c *Cache) Add(key string, value Value) {
	//如果存在，直接更新节点值，并且移动到队尾
	if ele, status := c.cache[key]; status {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		//不存在则新增节点，队尾添加新节点，并在map中添加key:value
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	//如果超出最大值，移除最少访问的节点
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
