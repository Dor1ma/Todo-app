package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("Server is running on port 8080")
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		log.Fatal(err)
	}

}
