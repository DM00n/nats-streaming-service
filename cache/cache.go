package cache

type Order struct {
	Id   string `json:"order_uid"`
	Data []byte
}

type Cache struct {
	//sync.RWMutex
	items map[string][]byte
}

func NewCache(orders []Order) *Cache {
	items := make(map[string][]byte)
	for _, item := range orders {
		items[item.Id] = item.Data
	}
	cache := Cache{
		items: items,
	}
	return &cache
}

func (c *Cache) AddOrder(id string, order []byte) bool {
	_, ok := c.items[id]
	//c.RWMutex.Lock()
	c.items[id] = order
	//c.RWMutex.Unlock()
	return ok
}

func (c *Cache) Get(id string) []byte {
	order, ok := c.items[id]
	if ok {
		return order
	} else {
		return nil
	}
}
