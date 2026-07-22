package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"
)

// Server pin for authentication (set this to secure your endpoint).
// Leave empty to disable pin authentication
var serverPin = ""

func handleGetTransmit(w http.ResponseWriter, r *http.Request) {
	userPin := GetOrElse(r.URL.Query()["Pin"], 0, "")

	// If the server pin is used, the user pin must match the server pin
	if len(serverPin) != 0 && userPin != serverPin {
		fmt.Fprint(w, "invalid pin")
		return
	}

	// Get the data from the request
	now := time.Now()
	aircraft := AircraftData{
		Callsign:          GetOrElse(r.URL.Query()["Callsign"], 0, ""),
		AircraftType:      GetOrElse(r.URL.Query()["AircraftType"], 0, ""),
		PilotName:         GetOrElse(r.URL.Query()["PilotName"], 0, ""),
		GroupName:         GetOrElse(r.URL.Query()["GroupName"], 0, ""),
		MsfsServer:        GetOrElse(r.URL.Query()["MSFSServer"], 0, ""),
		TransponderCode:   GetOrElse(r.URL.Query()["TransponderCode"], 0, ""),
		Latitude:          GetOrElse(r.URL.Query()["Latitude"], 0, "0"),
		Longitude:         GetOrElse(r.URL.Query()["Longitude"], 0, "0"),
		Altitude:          GetOrElse(r.URL.Query()["Altitude"], 0, "0"),
		Heading:           GetOrElse(r.URL.Query()["Heading"], 0, "0"),
		Airspeed:          GetOrElse(r.URL.Query()["Airspeed"], 0, "0"),
		Groundspeed:       GetOrElse(r.URL.Query()["Groundspeed"], 0, "0"),
		TouchdownVelocity: GetOrElse(r.URL.Query()["TouchdownVelocity"], 0, "0"),
		Notes:             GetOrElse(r.URL.Query()["Notes"], 0, ""),
		Version:           GetOrElse(r.URL.Query()["Version"], 0, "1.0.0.n"),
		Modified:          now,
	}

	// Default groundspeed to airspeed if it is not supplied
	if aircraft.Groundspeed == "0" {
		aircraft.Groundspeed = aircraft.Airspeed
	}

	// Check we have everything we need to store the data
	if aircraft.Callsign == "" || aircraft.AircraftType == "" || aircraft.PilotName == "" || aircraft.GroupName == "" {
		fmt.Fprint(w, "Insufficient data received")
		return
	}

	// Clean up the data (only allow alphanumeric characters, spaces and hyphens)
	regex := regexp.MustCompile(`[^A-Za-z0-9. -]`)
	aircraft.Callsign = regex.ReplaceAllString(aircraft.Callsign, "")
	aircraft.PilotName = regex.ReplaceAllString(aircraft.PilotName, "")
	aircraft.GroupName = regex.ReplaceAllString(aircraft.GroupName, "")
	aircraft.AircraftType = regex.ReplaceAllString(aircraft.AircraftType, "")

	// Allow newlines and commas in notes
	regex = regexp.MustCompile(`[^A-Za-z0-9.,\n -]`)
	aircraft.Notes = regex.ReplaceAllString(aircraft.Notes, "")

	// Check if aircraft already exists and preserve created time
	if data, ok := cache.get(aircraft.Callsign); ok {
		aircraft.Created = data.Created
	} else {
		aircraft.Created = now
	}

	cache.set(aircraft.Callsign, aircraft)

	fmt.Fprint(w, "updated")
}

// GetOrElse returns the slice element at index, or the provided default value if out of bounds.
func GetOrElse[T any](slice []T, index int, defaultValue T) T {
	if index < 0 || index >= len(slice) {
		return defaultValue
	}
	return slice[index]
}
