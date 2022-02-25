package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"net/http"
	"pub/Route"
)

const pgDB string = "postgresql://maui:maui@172.20.10.3:5432/postgres"

func main() {
	stanConn, errConn := stan.Connect("test-cluster", "publisher", stan.NatsURL("172.20.10.3:4222"))
	if errConn != nil {
		fmt.Println("conn: ", errConn)
		return
	}

	router := Route.Router(Route.NatsConn{Conn: stanConn})
	http.ListenAndServe(":3000", router)
}
