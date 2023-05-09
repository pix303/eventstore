package main

import (
	"fmt"
	"log"

	"github.com/pix303/eventstore/store"
)

func main() {

	log.Println("init event store")

	store, err := store.NewStore(store.WithPostgresRepository())
	if err != nil {
		log.Fatal(err)
	}

	//data := "{\"saluto\": \"ciao\", \"it\":true, \"num\":100}"
	//err = store.Add("ADD_MESSAGE", "MSG", "fb8e9dcf-32f5-4c92-8f8f-bb65f5b6a8b1", data)
	if err != nil {
		log.Fatal(err)
	}

	//time.Sleep(time.Millisecond * 100)
	//err = store.Add("ADD_MESSAGE", "MSG", "f0ad2add-be95-4635-93a1-6482860cd1cd", data)
	// time.Sleep(time.Second * 1)
	// store.Add("ADD_MESSAGE", "MSG", "123", "ciao")
	// store.Add("ADD_MESSAGE", "MSG", "443", "ciao ciao 443")

	events, err := store.GetAggregate("MSG", "fb8e9dcf-32f5-4c92-8f8f-bb65f5b6a8b1")
	for k, e := range events {
		fmt.Printf("%d: %s - %s \n", k, e.CreatedAt, e.AggregateID)
	}

	if err != nil {
		log.Println("error!")
		log.Fatal(err)
	}

}
