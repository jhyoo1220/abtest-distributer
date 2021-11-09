package testlist

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jhyoo1220/abtest-distributer/internal/pkg/cache"
	"github.com/jhyoo1220/abtest-distributer/internal/pkg/dbs"
	"strconv"
	"sync"
)

var (
	hc cache.HashCache
	m  sync.Mutex
)

func Init() {
	hc.Init()
}

func Read(refresh bool) (map[string]string, error) {
	key := dbs.GetTestlistKey()
	return hc.Read(key, refresh)
}

func Add(newName string) error {
	testlist, err := Read(true)
	if err != nil {
		return err
	}

	for _, name := range testlist {
		if name == newName {
			return fmt.Errorf("%s already exists!", name)
		}
	}

	cntKey := dbs.GetTestCountKey()

	m.Lock()

	id, err := dbs.Get(cntKey)
	if err != nil {
		m.Unlock()
		return err
	}

	err = dbs.IncrBy(cntKey, 1)
	if err != nil {
		m.Unlock()
		return err
	}

	m.Unlock()

	listKey := dbs.GetTestlistKey()
	return dbs.HSet(listKey, id, name)
}
