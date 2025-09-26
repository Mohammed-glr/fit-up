package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Message Service starting...")
	
	// TODO: Implement message service
	port := ":8082"
	
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Message Service is running")
	})
	
	log.Printf("Message Service listening on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}