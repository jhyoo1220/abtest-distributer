package abtest

import (
	"github.com/jhyoo1220/abtest-distributer/internal/app/numusers"
	"github.com/jhyoo1220/abtest-distributer/internal/pkg/cache"
	"github.com/jhyoo1220/abtest-distributer/internal/pkg/dbs"
)

type ABTest struct {
	ID            *int64            `json:"id, number"`
	Name          *string           `json:"name, string"`
	Description   *string           `json:"description, string"`
	Type          *string           `json:"description, string"`
	Status        *string           `json:"status, string"`
	WinnerGroup   *string           `json:"winner_group, string"`
	SelectedGroup *string           `json:"selected_group, string"`
	Groups        []testgroup.Group `json:"test_groups, list"`
	Created       *int64            `json:"created, number"`
}

type Group struct {
	Name        *string  `json:"name, string"`
	TargetRatio *float64 `json:"target_ratio, number"`
	NumUsers    *int64   `json:"num_users, number"`
}

var (
	c cache.Cache
)

func Init() {
	c.Init()
}

func Read(name string, refresh bool) (ABTest, error) {
	key := dbs.GetTestKey(name)
	abTest := ABTest{}

	abTestStr, err := c.Read(key, refresh)
	if err != nil {
		return abTest, err
	}

	if err := json.Unmarshal([]byte(abTestStr), &abTest); err != nil {
		log.Println(err.Error())
		return abTest, err
	}

	if err := updateNumUsers(&abTest); err != nil {
		log.Println(err.Error())
		return abTest, err
	}

	return abTest, nil
}

func updateNumUsers(abTest *ABTest) error {
	for i, group := range abTest.Groups {
		numUsers, err := numusers.Read(abTest.Name, group.Name)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		abTest.Groups[i].NumUsers = numUsers
	}

	return nil
}
