package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
)

func main() {

	stanConn, errConn := stan.Connect("test-cluster", "publisher", stan.NatsURL("172.20.10.3:4222"))
	if errConn != nil {
		fmt.Println("conn: ", errConn)
		return
	}
	fmt.Println(stanConn)
	sub, errSub := stanConn.Subscribe("db_service", func(m *stan.Msg) {
		fmt.Println("msg from broker: ", string(m.Data))
	})
	if errSub != nil {
		fmt.Println(errSub)
		return
	}

	sub.Unsubscribe()
	sub.Close()
	fmt.Println("asd")

}