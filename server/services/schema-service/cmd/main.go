package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Println("Schema Service starting...")

	// TODO: Implement schema service
	port := ":8083"

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Schema Service is running")
	})

	log.Printf("Schema Service listening on %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
