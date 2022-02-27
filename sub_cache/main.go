package main

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"sub_cache/BrokerConnection"
)

func main() {
	cache, errNewCache := lru.New(300)
	if errNewCache != nil {
		fmt.Errorf("new chache: %v", errNewCache)
		return
	}

	BrokerConnection.KeepAliveSub(cache, "172.20.10.3:4222", "test-cluster",
		"subscriber_cache", "cache_service")
}