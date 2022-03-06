package main

import (
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"net/http"
	"sub_cache/Cache"
	"sub_cache/DataBase"
	"sub_cache/Router"
)

func main() {
	conn, errConn := DataBase.Connection("postgresql://maui:maui@192.168.0.139:5432/postgres")
	if errConn != nil {
		fmt.Println("failed to connect db: ", errConn)
		return
	}
	fmt.Println("connect to db")

	var cacheSize int = 200
	cache, errCreateCache := lru.New(cacheSize)
	if errCreateCache != nil {
		fmt.Println("failed to create cache: ", errCreateCache)
		return
	}

	if errGetItems := Cache.CreateCache(cache, Cache.Postgres{DB: conn}, cacheSize); errGetItems != nil {
		fmt.Println(errGetItems)
		return
	}
	fmt.Println("cache created")
	fmt.Println("start server")
	router := Router.Route(Router.Postgres{DB: conn}, cache)
	http.ListenAndServe(":3001", router)
}
