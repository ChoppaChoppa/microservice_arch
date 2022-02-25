package main

import (
	"fmt"
	"sub_cache/BrokerConnection"

	"github.com/nats-io/stan.go"
)

func main() {
	stanConn, errConn := stan.Connect("test-cluster", "sub", stan.NatsURL("172.20.10.3:4222"))
	if errConn != nil {
		fmt.Println("conn: ", errConn)
		return
	}

	sub, errSub := stanConn.QueueSubscribe("db_service", "queue", func(m *stan.Msg) {
		fmt.Println("msg from broker: ", string(m.Data))
	}, stan.DurableName("durable-name"))
	if errSub != nil {
		fmt.Println(errSub)
		return
	}

	BrokerConnection.KeepAliveSub("172.20.10.3:4222", "test-cluster",
		"subscriber_cache", "cache_service")

	sub.Unsubscribe()
	sub.Close()
}