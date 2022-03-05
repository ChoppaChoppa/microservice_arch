package BrokerConnection

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"sub_db/Models"
	"time"

	"github.com/nats-io/stan.go"
)

const secretKey string = "sub_db"

type IPgDataBase interface {
	Add(ctx context.Context, order Models.OrderInfo) (Models.OrderInfo, error)
}

type DataBase struct {
	DB IPgDataBase
}

func KeepAliveSub(pg DataBase, url, clusterID, clientID, subject string) error {
	for {
		fmt.Println("connection to stan...")

		stanConn, errConn := stan.Connect(clusterID, clientID, stan.NatsURL(url))
		if errConn != nil {
			fmt.Println("conn: ", errConn)
			return errConn
		}
		fmt.Println("connected")



		var msg Models.Message
		_, errSub := stanConn.Subscribe(subject, func(m *stan.Msg) {
			if errUnmarshal := json.Unmarshal(m.Data, &msg); errUnmarshal != nil {
				fmt.Println("unmarshal: ", errUnmarshal)
				return
			}

			if msg.SecretKey != secretKey {
				fmt.Println("wrong key")
				return
			}

			_, errAdd := pg.DB.Add(context.Background(), msg.Order)
			if errAdd != nil {
				fmt.Println("failed to add in db: ", errAdd)
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
