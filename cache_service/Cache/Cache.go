package Cache

import (
	"context"
	"encoding/json"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"io/ioutil"
	"net/http"
	"strconv"
	"sub_cache/Models"
)

type IpgDataBase interface {
	GetLasts(context.Context, int) ([]Models.OrderInfo, error)
}

type Postgres struct {
	DB    IpgDataBase
}

func CreateCache(cache *lru.Cache, cacheSize int) error {
	resp, errGetLasts := http.Get("http://127.0.0.1:3000/sub_db/get/lasts/" + strconv.Itoa(cacheSize))
	if errGetLasts != nil {
		return fmt.Errorf("get lasts: %v", errGetLasts)
	}

	var orders []Models.OrderInfo
	body, errGetBody := ioutil.ReadAll(resp.Body)
	if errGetBody != nil {
		fmt.Errorf("get body: %v", errGetBody.Error())
	}

	if errUnmarshalBody := json.Unmarshal(body, &orders); errUnmarshalBody != nil || errGetBody != nil {
		return fmt.Errorf("unmarshal body: %v", errUnmarshalBody)
	}

	for _, v := range orders {
		cache.Add(v.ID, v)
	}

	return nil
}
