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

type aircraftData struct {
	userPin           string
	callsign          string
	aircraftType      string
	pilotName         string
	groupName         string
	msfsServer        string
	transponderCode   string
	latitude          string
	longitude         string
	altitude          string
	heading           string
	airspeed          string
	groundspeed       string
	touchdownVelocity string
	notes             string
	version           string
	timestamp         time.Time
	created           time.Time
	modified          time.Time
}

func handleGetTransmit(w http.ResponseWriter, r *http.Request) {
	// Get the data from the request
	aircraft := aircraftData{
		userPin:           GetOrElse(r.URL.Query()["Pin"], 0, ""),
		callsign:          GetOrElse(r.URL.Query()["Callsign"], 0, ""),
		aircraftType:      GetOrElse(r.URL.Query()["AircraftType"], 0, ""),
		pilotName:         GetOrElse(r.URL.Query()["PilotName"], 0, ""),
		groupName:         GetOrElse(r.URL.Query()["GroupName"], 0, ""),
		msfsServer:        GetOrElse(r.URL.Query()["MSFSServer"], 0, ""),
		transponderCode:   GetOrElse(r.URL.Query()["TransponderCode"], 0, ""),
		latitude:          GetOrElse(r.URL.Query()["Latitude"], 0, "0"),
		longitude:         GetOrElse(r.URL.Query()["Longitude"], 0, "0"),
		altitude:          GetOrElse(r.URL.Query()["Altitude"], 0, "0"),
		heading:           GetOrElse(r.URL.Query()["Heading"], 0, "0"),
		airspeed:          GetOrElse(r.URL.Query()["Airspeed"], 0, "0"),
		groundspeed:       GetOrElse(r.URL.Query()["Groundspeed"], 0, "0"),
		touchdownVelocity: GetOrElse(r.URL.Query()["TouchdownVelocity"], 0, "0"),
		notes:             GetOrElse(r.URL.Query()["Notes"], 0, ""),
		version:           GetOrElse(r.URL.Query()["Version"], 0, "1.0.0.n"),
		timestamp:         time.Now(),
	}

	// If the server pin is used, the user pin must match the server pin
	if len(serverPin) != 0 && aircraft.userPin != serverPin {
		fmt.Fprint(w, "invalid pin")
		return
	}

	// Default groundspeed to airspeed if it is not supplied
	if aircraft.groundspeed == "0" {
		aircraft.groundspeed = aircraft.airspeed
	}

	// Check we have everything we need to store the data
	if aircraft.callsign == "" || aircraft.aircraftType == "" || aircraft.pilotName == "" || aircraft.groupName == "" {
		fmt.Fprint(w, "Insufficient data received")
		return
	}

	// Clean up the data (only allow alphanumeric characters, spaces and hyphens)
	regex := regexp.MustCompile(`[^A-Za-z0-9. -]`)
	aircraft.callsign = regex.ReplaceAllString(aircraft.callsign, "")
	aircraft.pilotName = regex.ReplaceAllString(aircraft.pilotName, "")
	aircraft.groupName = regex.ReplaceAllString(aircraft.groupName, "")
	aircraft.aircraftType = regex.ReplaceAllString(aircraft.aircraftType, "")

	// Allow newlines and commas in notes
	regex = regexp.MustCompile(`[^A-Za-z0-9.,\n -]`)
	aircraft.notes = regex.ReplaceAllString(aircraft.notes, "")

	// Check if aircraft already exists and preserve created time
	if data, ok := cache.get(aircraft.callsign); ok {
		aircraft.created = data.created
	} else {
		aircraft.created = time.Now()
	}
	aircraft.modified = time.Now()

	cache.set(aircraft.callsign, aircraft)

	fmt.Fprint(w, "updated")
}

// GetOrElse returns the slice element at index, or the provided default value if out of bounds.
func GetOrElse[T any](slice []T, index int, defaultValue T) T {
	if index < 0 || index >= len(slice) {
		return defaultValue
	}
	return slice[index]
}
