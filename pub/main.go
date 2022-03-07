package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	"net/http"
	"pub/Route"
)

func main() {
	stanConn, errConn := stan.Connect("test-cluster", "publisher",
		stan.NatsURL("127.0.0.1:4222"))
	if errConn != nil {
		fmt.Println("conn: ", errConn)
		return
	}

	router := Route.Router(Route.NatsConn{Conn: stanConn})
	http.ListenAndServe(":3000", router)
}
