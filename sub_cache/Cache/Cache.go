package Cache

import (
	"context"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"sub_cache/Models"
)

type IpgDataBase interface {
	GetLasts(context.Context, int) ([]Models.OrderInfo, error)
}

type Postgres struct {
	DB    IpgDataBase
}

func CreateCache(cache *lru.Cache, pg Postgres, cacheSize int) error {
	fmt.Println("cache len ", cache.Len())
	orders, errGetLasts := pg.DB.GetLasts(context.Background(), cacheSize)
	if errGetLasts != nil {
		return fmt.Errorf("get lasts: %v", errGetLasts)
	}

	fmt.Println(orders)
	for _, v := range orders {
		cache.Add(v.ID, v)
		fmt.Println(cache.Get(v.ID))
	}

	return nil
}
