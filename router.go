package main

import "net/http"

func initRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// Endpoints
	mux.HandleFunc("GET /transmit", handleTransmit)
	mux.HandleFunc("GET /transmit.php", handleTransmit)

	mux.HandleFunc("POST /transmit", handleTransmit)
	mux.HandleFunc("POST /transmit.php", handleTransmit)

	// Both JSON endpoints give same data - both kept for backwards-compatibility
	mux.HandleFunc("GET /radar_data", handleJSON)
	mux.HandleFunc("GET /radar_data.php", handleJSON)
	mux.HandleFunc("GET /status_json", handleJSON)
	mux.HandleFunc("GET /status_json.php", handleJSON)

	return mux
}
