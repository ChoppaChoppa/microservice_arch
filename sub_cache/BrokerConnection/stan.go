package BrokerConnection

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"

	"github.com/nats-io/stan.go"
)

func KeepAliveSub(url, clusterID, clientID, subject string) {
	for {
		fmt.Println("connection to stan...")

		stanConn, errConn := stan.Connect(clusterID, clientID, stan.NatsURL(url))
		if errConn != nil {
			fmt.Println("conn: ", errConn)
			return
		}
		fmt.Println("connected")

		_, errSub := stanConn.Subscribe(subject, func(m *stan.Msg) {
			fmt.Println("msg from broker: ", string(m.Data))
			time.Sleep(time.Second * 5)
		})
		if errSub != nil {
			fmt.Println("failed subscribe")
			continue
		}
		fmt.Println("subscribed")

		chSubIsClosed := make(chan bool)

		go func(chSubIsClosed chan bool, conn *nats.Conn) {
			for conn.IsConnected() {
				time.Sleep(time.Millisecond * 100)
			}
			chSubIsClosed <- true
		}(chSubIsClosed, stanConn.NatsConn())

		<-chSubIsClosed
		fmt.Println("ERROR: nats connection is closed")
	}

}
