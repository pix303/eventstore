package main

import (
	"log"
	"time"

	"github.com/pix303/eventstore/store"
)

func main() {

	log.Println("init event store")

	store, err := store.NewStore(store.WithBBoltRepository())
	if err != nil {
		log.Fatal(err)
	}
	store.Add("ADD_MESSAGE", "MSG", "123", "ciao ciao ciao")
	time.Sleep(time.Millisecond * 1)
	store.Add("ADD_MESSAGE", "MSG", "123", "ciao ciao")
	time.Sleep(time.Second * 1)
	store.Add("ADD_MESSAGE", "MSG", "123", "ciao")
	store.Add("ADD_MESSAGE", "MSG", "443", "ciao ciao 443")

	r, err := store.GetAggregate("MSG", "123")

	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(r))
	log.Println(r)
}
