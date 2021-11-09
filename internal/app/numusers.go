package numusers

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jhyoo1220/abtest-distributer/internal/pkg/cache"
	"github.com/jhyoo1220/abtest-distributer/internal/pkg/dbs"
	"strconv"
	"sync"
)

var (
	c cache.Cache
)

func Init() {
	c.Init()
}

func Read(testName string, groupName string) (int64, error) {
	var numUsers int64
	var err error

	key := dbs.GetTestNumUsersKey(testName, groupName)

	numUsersStr, err := c.Read(key, false)
	if err != nil {
		return numUsers, err
	}

	if numUsers, err = strconv.ParseInt(numUsersStr, 10, 64); err != nil {
		log.Println(err.Error())
		return numUsers, err
	}

	if numUsersCache, exists := c.Load(key); exists {
		numUsers += int64(numUsersCache.(int))
	}

	if numUsers < 0 {
		numUsers = 0
	}

	return numUsers, nil
}
