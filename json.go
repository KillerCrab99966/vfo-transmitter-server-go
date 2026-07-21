package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type AircraftJSON struct {
	AircraftData
	TimeOnline             string `json:"time_online"`
	SecondsSinceLastUpdate int    `json:"seconds_since_last_update"`
	Modified               string `json:"modified"`
}

func handleJSON(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if len(cache.data) == 0 {
		fmt.Fprint(w, "[]")
		return
	}

	aircraft := []AircraftJSON{}

	// Collect the connected aircraft into the slice
	cache.mu.Lock()
	for _, data := range cache.data {
		timeOnline := formatTimeOnline(time.Since(data.Created))
		sinceLast := time.Since(data.Modified).Seconds()
		modified := formatModified(data.Modified)

		aircraft = append(aircraft, AircraftJSON{
			AircraftData:           data,
			TimeOnline:             timeOnline,
			SecondsSinceLastUpdate: int(sinceLast),
			Modified:               modified,
		})
	}
	cache.mu.Unlock()

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
