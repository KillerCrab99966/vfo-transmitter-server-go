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

	// Both JSON endpoints give same data - both kept for backwards-compatibility
	mux.HandleFunc("GET /radar_data", handleJSON)
	mux.HandleFunc("GET /radar_data.php", handleJSON)
	mux.HandleFunc("GET /status_json", handleJSON)
	mux.HandleFunc("GET /status_json.php", handleJSON)

	return mux
}
