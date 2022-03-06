package main

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"net/http"
	"sub_cache/Cache"
	"sub_cache/Router"
)

func main() {
	var cacheSize int = 200
	cache, errCreateCache := lru.New(cacheSize)
	if errCreateCache != nil {
		fmt.Println("failed to create cache: ", errCreateCache)
		return
	}

	if errGetItems := Cache.CreateCache(cache, cacheSize); errGetItems != nil {
		fmt.Println(errGetItems)
		return
	}
	fmt.Println("cache created")
	fmt.Println("start server")
	router := Router.Route(cache)
	http.ListenAndServe(":3001", router)
}
