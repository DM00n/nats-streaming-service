package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"wbService/cache"
	"wbService/pg"
	"wbService/stan"
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
	db, err := pg.NewDB(psqlconn)
	CheckError(err)

	rows, _ := db.SelectAll()
	var list []cache.Order
	var asd cache.Order
	for rows.Next() {
		rows.Scan(&asd.Id, &asd.Data)
		list = append(list, asd)
	}
	rows.Close()

	cch := cache.NewCache(list)

	natsCon, err := stan.NewStanConnect()
	CheckError(err)

	//http.HandleFunc("/", web.Handler)
	//err = http.ListenAndServe(":8080", nil)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}

	sub, err := natsCon.Sub(
		func(m []byte) {
			if !json.Valid(m) {
				fmt.Println("invalid JSON string:", string(m))
				return
			}
			var order cache.Order
			err = json.Unmarshal(m, &order)
			if err != nil {
				println(err.Error())
				return
			}
			updated := cch.AddOrder(order.Id, m)
			if updated {
				err = db.UpdateValue(order.Id, m)
				if err != nil {
					println(err.Error())
					return
				}
			} else {
				err = db.InsertValue(order.Id, m)
				if err != nil {
					println(err.Error())
					return
				}
			}
		})
	CheckError(err)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	signal.Notify(signalChan, os.Interrupt)
	go func() {
		for range signalChan {
			fmt.Printf("\nReceived an interrupt, unsubscribing and closing connection...\n\n")
			sub.Unsubscribe()
			sub.Close()
			natsCon.Close()
			cleanupDone <- true
		}
	}()
	<-cleanupDone

}

//func MessageHandler(m []byte) {
//	println(string(m))
//	if !json.Valid(m) {
//		fmt.Println("invalid JSON string:", string(m))
//		return
//	}
//	var str cache.Order
//	err := json.Unmarshal(m, &str)
//
//	if err != nil {
//		println(err.Error())
//		return
//	}
//}
