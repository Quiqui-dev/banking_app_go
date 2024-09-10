package main

import (
	"log"

	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")

	store, err := NewPostgresStore()

	if err != nil {
		log.Fatalf("error: %s", err.Error())
	}

	if err := store.Init(); err != nil {
		log.Fatalf("couldn't initialise the store: %+v\n", err)
	}

	srv := NewAPIServer(":3000", store)

	srv.Run()
}
