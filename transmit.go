package main

import "net/http"

// Server pin for authentication (set this to secure your endpoint).
// Leave empty to disable pin authentication
var serverPin = ""

type aircraftData struct {
	UserPin           string  `json:"pin"`
	Callsign          string  `json:"callsign"`
	AircraftType      string  `json:"aircraft_type"`
	PilotName         string  `json:"pilot_name"`
	GroupName         string  `json:"group_name"`
	MsfsServer        string  `json:"msfs_server"`
	TransponderCode   string  `json:"transponder_code"`
	Latitude          float64 `json:"latitude"`
	Longitude         float64 `json:"longitude"`
	Altitude          float64 `json:"altitude"`
	Heading           float64 `json:"heading"`
	Airspeed          float64 `json:"airspeed"`
	Groundspeed       float64 `json:"groundspeed"`
	TouchdownVelocity float64 `json:"touchdown_velocity"`
	Notes             string  `json:"notes"`
	Version           string  `json:"version"`
}

func handleGetTransmit(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func handlePostTransmit(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}
