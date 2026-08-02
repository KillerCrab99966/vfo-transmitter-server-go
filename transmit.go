package main

import (
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func handleTransmit(w http.ResponseWriter, r *http.Request) {
	userPin := getOrDefault(r.URL.Query()["Pin"], 0, "")

	// If the server pin is used, the user pin must match the server pin
	if len(serverPin) != 0 && userPin != serverPin {
		fmt.Fprint(w, "invalid pin")
		return
	}

	// Get the data from the request
	now := time.Now().UTC()
	aircraft := AircraftData{
		Callsign:          parameterOrDefault(r, "Callsign", ""),
		AircraftType:      parameterOrDefault(r, "AircraftType", ""),
		PilotName:         parameterOrDefault(r, "PilotName", ""),
		GroupName:         parameterOrDefault(r, "GroupName", ""),
		MsfsServer:        parameterOrDefault(r, "MSFSServer", ""),
		TransponderCode:   parameterOrDefault(r, "TransponderCode", ""),
		Latitude:          parameterOrDefault(r, "Latitude", "0"),
		Longitude:         parameterOrDefault(r, "Longitude", "0"),
		Altitude:          parameterOrDefault(r, "Altitude", "0"),
		Heading:           parameterOrDefault(r, "Heading", "0"),
		Airspeed:          parameterOrDefault(r, "Airspeed", "0"),
		Groundspeed:       parameterOrDefault(r, "Groundspeed", "0"),
		TouchdownVelocity: parameterOrDefault(r, "TouchdownVelocity", "0"),
		Notes:             parameterOrDefault(r, "Notes", ""),
		Version:           parameterOrDefault(r, "Version", "1.0.0.n"),
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

	// Check rate limiting
	if rateLimit {
		// NOTE: The RemoteAddr may not be the client's IP if this server
		// is behind another proxy/s such as NGINX or Cloudflare.
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		if isRateLimited(aircraft.Callsign, ip) {
			fmt.Fprint(w, "rate limited")
			return
		}
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
	if data, ok := acCache.get(aircraft.Callsign); ok {
		aircraft.Created = data.Created
	} else {
		aircraft.Created = now
	}

	acCache.set(aircraft.Callsign, aircraft)

	fmt.Fprint(w, "updated")
}

func getRateLimitKey(callsign, ip string) string {
	ipFormatted := strings.ReplaceAll(ip, ".", "_")
	return callsign + "_" + ipFormatted
}

func isRateLimited(callsign, ip string) bool {
	rateKey := getRateLimitKey(callsign, ip)
	lastUpdate, ok := rateCache.get(rateKey)

	// 1 second minimum between updates
	if ok && time.Since(lastUpdate) < time.Second {
		return true
	}

	// Update rate limit timestamp
	rateCache.set(rateKey, time.Now())
	return false
}

// parameterOrDefault returns the GET or POST parameter or the provided default value.
func parameterOrDefault(r *http.Request, parameter string, defaultValue string) string {
	queryParams := r.URL.Query()
	if values, exists := queryParams[parameter]; exists {
		return getOrDefault(values, 0, defaultValue)
	}

	r.ParseForm()
	if values, exists := r.PostForm[parameter]; exists {
		return getOrDefault(values, 0, defaultValue)
	}

	return defaultValue
}

// getOrDefault returns the slice element at index, or the provided default value if out of bounds.
func getOrDefault[T any](slice []T, index int, defaultValue T) T {
	if index < 0 || index >= len(slice) {
		return defaultValue
	}
	return slice[index]
}
