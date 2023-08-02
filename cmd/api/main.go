package main

import (
	"log"
	"net/http"

	"github.com/mertess/cryptocurrency-statistics-api-go/internal/api/handlers"
)

func main() {
	http.HandleFunc("/lastPrice", handlers.GetLastPrice)

	log.Println("Started listening...")

	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
