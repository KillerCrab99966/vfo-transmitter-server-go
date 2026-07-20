package main

import "net/http"

func initRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("GET /api/transmit", handleGetTransmit)
	mux.HandleFunc("POST /api/transmit", handlePostTransmit)

	return mux
}
