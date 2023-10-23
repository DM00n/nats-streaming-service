package cache

type Order struct {
	Id   string `json:"order_uid"`
	Data string
}

type Cache struct {
	items map[string]string
}

func NewCache(orders []Order) *Cache {
	items := make(map[string]string)
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
	c.items[id] = string(order)
	return ok
}

func (c *Cache) Get(id string) Order {
	order, ok := c.items[id]
	if ok {
		return Order{id, string(order)}
	} else {
		return Order{}
	}
}
