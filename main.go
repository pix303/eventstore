package main

import (
	"log"

	"github.com/pix303/eventstore/store"
)

func main() {

	log.Println("init event store")

	_, err := store.NewStore(store.WithBBoltRepository())

	if err != nil {
		log.Fatal(err)
	}

}
