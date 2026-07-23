package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type AircraftJSON struct {
	AircraftData

	LatFormatted         string `json:"latitude_formatted"`
	LongFormatted        string `json:"longitude_formatted"`
	AltFormatted         string `json:"altitude_formatted"`
	HeadingFormatted     string `json:"heading_formatted"`
	AirspeedFormatted    string `json:"airspeed_formatted"`
	GroundspeedFormatted string `json:"groundspeed_formatted"`
	TDVelocityFormatted  string `json:"touchdown_velocity_formatted"`

	TimeOnline             string `json:"time_online"`
	SecondsSinceLastUpdate int    `json:"seconds_since_last_update"`
	Modified               string `json:"modified"`
}

func handleJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if len(cache.items) == 0 {
		fmt.Fprint(w, "[]")
		return
	}

	aircraft := []AircraftJSON{}

	// Collect the connected aircraft into the slice
	cache.mu.RLock()
	for _, data := range cache.items {
		// We only need the data, not the age
		data := data.data

		// Convert the touchdown velocity from ft/sec to fpm
		tv, _ := strconv.ParseFloat(data.TouchdownVelocity, 64)
		data.TouchdownVelocity = strconv.FormatFloat(tv*60, 'f', -1, 64)

		lat, _ := strconv.ParseFloat(data.Latitude, 64)
		dms := decToDMS(lat)
		latFormatted := fmt.Sprintf("%d&deg; %d' %.2f\" %s", dms.d, dms.m, dms.s, cardinalLat(lat))

		long, _ := strconv.ParseFloat(data.Longitude, 64)
		dms = decToDMS(long)
		longFormatted := fmt.Sprintf("%d&deg; %d' %.2f\" %s", dms.d, dms.m, dms.s, cardinalLong(long))

		timeOnline := formatTimeOnline(time.Since(data.Created))
		sinceLast := time.Since(data.Modified).Seconds()
		modified := formatModified(data.Modified)

		aircraft = append(aircraft, AircraftJSON{
			AircraftData: data,

			LatFormatted:         latFormatted,
			LongFormatted:        longFormatted,
			AltFormatted:         formatAndTrunc(data.Altitude),
			HeadingFormatted:     formatAndTrunc(data.Heading),
			AirspeedFormatted:    formatAndTrunc(data.Airspeed),
			GroundspeedFormatted: formatAndTrunc(data.Groundspeed),
			TDVelocityFormatted:  formatAndTrunc(data.TouchdownVelocity),

			TimeOnline:             timeOnline,
			SecondsSinceLastUpdate: int(sinceLast),
			Modified:               modified,
		})
	}
	cache.mu.RUnlock()

	jsonData, err := json.Marshal(aircraft)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{error:"%v"}`, err)
		return
	}

	fmt.Fprint(w, string(jsonData))
}

func formatTimeOnline(d time.Duration) string {
	// Round down to whole seconds to drop sub-second precision
	d = d.Round(time.Second)

	hours := d / time.Hour
	d -= hours * time.Hour

	minutes := d / time.Minute
	d -= minutes * time.Minute

	seconds := d / time.Second

	// Format with leading zeros
	return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
}

func formatModified(t time.Time) string {
	y, m, d := t.Date()

	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", y, m, d, t.Hour(), t.Minute(), t.Second())
}

type dms struct {
	d int
	m int
	s float64
}

// Converts decimal longitude / latitude to DMS (Degrees / minutes / seconds).
// The returned values will always be positive.
func decToDMS(dec float64) dms {
	if dec == 0 {
		return dms{}
	}

	// Work with positive fraction so minutes/seconds are positive
	dFrac, tempma := math.Modf(math.Abs(dec))

	tempma *= 3600
	m := math.Floor(tempma / 60)
	s := tempma - m*60

	return dms{int(dFrac), int(m), s}
}

func cardinalLat(lat float64) string {
	if lat > 0 {
		return "N"
	} else {
		return "S"
	}
}

func cardinalLong(long float64) string {
	if long > 0 {
		return "E"
	} else {
		return "W"
	}
}

// formatAndTrunc discards the fractional part and add thousands delimiters
func formatAndTrunc(num string) string {
	if len(num) == 0 {
		return num
	}

	// Truncate
	num = strings.Split(num, ".")[0]

	// Handle optional sign (+ or -)
	sign := ""
	if num[0] == '+' || num[0] == '-' {
		sign = string(num[0])
		num = num[1:]
	}

	if len(num) <= 3 {
		return sign + num
	}

	// Insert commas from right to left
	var result strings.Builder
	n := len(num)

	for i, char := range num {
		// Calculate remaining digits to the right
		if i > 0 && (n-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(char)
	}

	return sign + result.String()
}
