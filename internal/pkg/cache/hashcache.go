package cache

import (
	"github.com/jhyoo1220/abtest-distributer/internal/pkg/dbs"
	"sync"
)

type HashCache struct {
	Cache sync.Map
	Mutex sync.Mutex
}

func (c *HashCache) Init() {
	c.Cache = sync.Map{}
	c.Mutex = sync.Mutex{}
}

func (c *HashCache) Read(hash string, refresh bool) (map[string]string, error) {
	var pairs map[string]string
	var err error

	if refresh {
		if pairs, err = dbs.HGetAll(hash); err != nil {
			return pairs, err
		}

		c.UpdateAll(pairs)
	} else {
		pairs = make(map[string]string)

		c.Cache.Range(func(key, value interface{}) bool {
			pairs[key.(string)] = value.(string)
			return true
		})
	}

	return pairs, nil
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

func (c *HashCache) UpdateAll(pairs map[string]string) {
	c.Mutex.Lock()

	for key, val := range pairs {
		c.Cache.Store(key, val)
	}

	c.Cache.Range(func(key, val interface{}) bool {
		if _, f := pairs[key.(string)]; !f {
			c.Cache.Delete(key.(string))
		}

		return true
	})

	c.Mutex.Unlock()
}
