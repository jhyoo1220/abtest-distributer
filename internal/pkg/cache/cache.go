package cache

import (
	"github.com/jhyoo1220/abtest-distributer/internal/pkg/dbs"
	"sync"
)

type Cache struct {
	Cache sync.Map
	Mutex sync.Mutex
}

func (c *Cache) Init() {
	c.Cache = sync.Map{}
	c.Mutex = sync.Mutex{}
}

func (c *HashCache) Read(key string, refresh bool) (string, error) {
	var val string
	var err error

	if refresh {
		if val, err = dbs.Get(key); err != nil {
			return val, err
		}

		c.Update(key, val)
	} else {
		if valCache, exists := c.Cache.Load(key); exists {
			val = valCache.(string)
		}
	}

	return val, nil
}

func (c *HashCache) Update(key string, val string) {
	c.Mutex.Lock()
	c.Cache.Store(key, val)
	c.Mutex.Unlock()
}

func (c *HashCache) Delete(key string) {
	c.Mutex.Lock()
	c.Cache.Delete(key)
	c.Mutex.Unlock()
}
