package main

import (
	"fmt"
	"net/http"

	"sub_db/BrokerConnection"
	"sub_db/PgDataBase"
	"sub_db/Route"
)

func main() {
	connection, errConnection := PgDataBase.Connection("postgresql://maui:maui@192.168.0.139:5432/postgres")
	if errConnection != nil {
		//log.Panic("failed to connect: ", errConnection)
		//return
	}
	go BrokerConnection.KeepAliveSub(BrokerConnection.DataBase{DB: connection}, "172.20.10.3:4222", "test-cluster",
		"subscriber_cache", "db_service")

	fmt.Println("start server")
	router := Route.Router(Route.DataBase{DB: connection})
	http.ListenAndServe(":3000", router)


}

//TODO обработка ошибок, чтобы приложение не падало
//