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

// Server pin for authentication. Disabled when empty.
var serverPin = ""

// Enable rate-limiting
var rateLimit = false

// Global aircraft acCache with 30min ttl
var acCache = newCache[AircraftData](30 * time.Minute)

// Airspace data and rate limit cache
var airpaceCache = newCache[string](24 * time.Hour)
var airspaceRateCache = newCache[int](time.Minute)

// Rate limits cache
var rateCache = newCache[time.Time](10 * time.Second)

// Embed static files
//
//go:embed all:static
var staticFiles embed.FS

// If `development`, then `config.toml` is looked for in cwd,
// if `production` then config is looked for in binary's path. If neither,
// the program will panic.
var Environment = "development"

// Global HTTP client for airspace data requests
var client = &http.Client{
	Timeout: 60 * time.Second,
}

func main() {
	// Read the config.toml
	cfg := readConfig()
	serverPin = cfg.Pin
	rateLimit = cfg.RateLimit

	// Initialise the router
	mux := initRoutes(cfg.Debug)

	// Create the server
	server := &http.Server{
		Addr:         cfg.Addr,
		Handler:      mux,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	fmt.Println("Server listening on:", server.Addr)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
	}
}

type Config struct {
	Addr      string `toml:"address"`
	Pin       string `toml:"pin"`
	RateLimit bool   `toml:"rate_limiting"`
	Debug     bool   `toml:"debug"`
}

func readConfig() Config {
	// Check build environment
	if Environment != "development" && Environment != "production" {
		log.Fatalf("Invalid build: Environment must be 'development' or 'production', got %q", Environment)
	}

	var configPath string

	// Get the path of the binary
	if Environment == "development" {
		configPath = "config.toml"
		cwd, _ := os.Getwd()
		fullPath := filepath.Join(cwd, "config.toml")
		fmt.Println("Looking for config at:", fullPath)
	} else {
		// Production build

		// Get path to the running executable file
		execPath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}

		// Get the directory containing the executable
		execDir := filepath.Dir(execPath)

		// Build absolute path
		configPath = filepath.Join(execDir, "config.toml")
		fmt.Println("Looking for config at:", configPath)
	}

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Unmarshal and parse
	var cfg Config
	err = toml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Error parsing TOML: %v", err)
	}

	fmt.Println("Config found and parsed!")
	return cfg
}
