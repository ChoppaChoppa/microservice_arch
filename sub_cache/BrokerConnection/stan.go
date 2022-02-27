package BrokerConnection

import (
	"encoding/json"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/nats-io/nats.go"
	"log"
	"time"

	"github.com/nats-io/stan.go"
)

func KeepAliveSub(cache *lru.Cache, url, clusterID, clientID, subject string) error {
	for {
		fmt.Println("connection to stan...")

		stanConn, errConn := stan.Connect(clusterID, clientID, stan.NatsURL(url))
		if errConn != nil {
			fmt.Println("conn: ", errConn)
			return errConn
		}
		fmt.Println("connected")



		var order interface{}
		_, errSub := stanConn.Subscribe(subject, func(m *stan.Msg) {
			if errUnmarshal := json.Unmarshal(m.Data, &order); errUnmarshal != nil {
				log.Printf("unmarshal: %v", errUnmarshal)
				return
			}


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
		 return fmt.Errorf("ERROR: nats connection is closed")
	}

}
