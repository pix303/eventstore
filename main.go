package main

import (
	"log"

	"github.com/pix303/eventstore/store"
)

func main() {

	log.Println("init event store")
	
	es, err := store.NewStore("dbtest.db")
	if err != nil {
		log.Fatal(err)
	}


	es.Close()
}
