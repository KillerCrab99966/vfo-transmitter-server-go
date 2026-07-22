package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// Global cache with 30min ttl
var cache = newAircraftCache(30 * time.Minute)

func main() {
	// Initialise the router
	mux := initRoutes()

	// Create the server
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Server listening on", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
	}
}
