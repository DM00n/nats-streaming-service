package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"os/signal"
	"wbService/nats"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "pguser"
	password = "pgpass"
	dbname   = "db"
)

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)
	defer db.Close()

	rows, err := db.Query(`SELECT * FROM "orders"`)
	CheckError(err)

	defer rows.Close()
	for rows.Next() {
		var id int
		var data []byte
		err = rows.Scan(&id, &data)
		CheckError(err)

		//err = json.Unmarshal(data, &res)
		//CheckError(err)

		fmt.Println(id, string(data))
	}

	connect, err := nats.NewStanConnect()
	CheckError(err)
	defer connect.Close()

	sub, err := nats.Sub(connect)
	CheckError(err)

	defer sub.Close()
	defer sub.Unsubscribe()

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			// Do not unsubscribe a durable on exit, except if asked to.
			connect.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone

}
