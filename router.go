package main

import (
	"io"
	"io/fs"
	"log"
	"net/http"
	"time"
)

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

	// Static files

	// Strip the "static" prefix so files are served relative to root
	subFS, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}

	// Helper function to serve specific HTML files
	serveHTML := func(filePath string) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			file, err := subFS.Open(filePath)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			defer file.Close()

			// Reads and serves the file with proper content headers
			http.ServeContent(w, r, filePath, time.Time{}, file.(io.ReadSeeker))
		}
	}

	fileServer := http.FileServer(http.FS(subFS))
	mux.Handle("GET /", fileServer)

	// Suffixless endpoints
	mux.Handle("GET /radar", serveHTML("radar.html"))
	mux.Handle("GET /status", serveHTML("status.html"))
	mux.Handle("GET /embed", serveHTML("embed.html"))

	return mux
}
