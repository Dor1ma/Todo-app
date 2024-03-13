package main

import (
	"Todo-app/internal/server"
	"log"
)

func main() {
	log.Println("Server is starting on port 8080")

	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
