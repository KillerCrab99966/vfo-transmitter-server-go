package main

import "net/http"

func initRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("GET /transmit", handleGetTransmit)
	mux.HandleFunc("GET /transmit.php", handleGetTransmit)

	// TODO: POST support
	// mux.HandleFunc("POST /transmit", handlePostTransmit)
	// mux.HandleFunc("POST /transmit.php", handlePostTransmit)

	return mux
}
