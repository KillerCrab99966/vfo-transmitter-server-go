package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/pelletier/go-toml/v2"
)

// Server pin for authentication (set this to secure your endpoint).
// Leave empty to disable pin authentication
var serverPin = ""

// Global cache with 30min ttl
var cache = newAircraftCache(30 * time.Minute)

// Embed static files
//
//go:embed all:static
var staticFiles embed.FS

// If `development`, then `config.toml` is looked for in cwd,
// otherwise `config.toml` is looked for in binary's path.
var Environment = "development"

func main() {
	// Read the config.toml
	cfg := readConfig()
	serverPin = cfg.Pin

	// Initialise the router
	mux := initRoutes(cfg.Debug)

	// Create the server
	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Println("Server listening on:", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
	}
}

type Config struct {
	Addr  string `toml:"address"`
	Pin   string `toml:"pin"`
	Debug bool   `toml:"debug"`
}

func readConfig() Config {
	var configPath string

	// Get the path of the binary
	if Environment == "development" {
		configPath = "config.toml"
		cwd, _ := os.Getwd()
		fmt.Println("Looking for config in:", cwd)
	} else {

		// Get path to the running executable file
		execPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}

		// Get the directory containing the executable
		execDir := filepath.Dir(execPath)

		// Build absolute path
		fmt.Println("Looking for config in", execDir)
		configPath = filepath.Join(execDir, "config.toml")
	}

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Unmarshal
	var cfg Config
	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Error unmarshaling TOML: %v", err)
	}

	fmt.Println("Config found and parsed!")
	return cfg
}
