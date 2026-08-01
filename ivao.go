package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

func handleIVAO(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	aircraft := getAllAircraft()

	// Output whazzup format
	var builder strings.Builder

	builder.WriteString("!GENERAL\r\n")
	builder.WriteString("VERSION = 1\r\n")
	builder.WriteString("RELOAD = 1\r\n")

	now := time.Now().UTC()
	nowFormatted := now.Format("20060102150405")
	fmt.Fprintf(&builder, "UPDATE = %s\r\n", nowFormatted)

	fmt.Fprintf(&builder, "CONNECTED CLIENTS = %d\r\n", len(aircraft))
	builder.WriteString("CONNECTED SERVERS = 0\r\n")
	builder.WriteString("!CLIENTS\r\n")

	for _, data := range aircraft {
		data := formatAircraft(data)
		fmt.Fprintf(&builder, "%s:%s:%s:PILOT::%s:%s:%s:%s:%s:::::%s:B:6:%s:0:50:0:I:::::::::VFR:::::::%s:%s:1 :1:1::S:0:%s:0:40:\r\n",
			data.Callsign,
			data.Callsign,
			data.PilotName,
			data.Latitude,
			data.Longitude,
			data.AltFormatted,
			data.GroundspeedFormatted,
			data.AircraftType,
			data.MsfsServer,
			data.TransponderCode,
			nowFormatted,
			data.GroupName,
			data.Heading,
		)
	}

	fmt.Fprint(w, builder.String())
}

func getAllAircraft() []AircraftData {
	acCache.mu.RLock()
	defer acCache.mu.RUnlock()

	var aircraft []AircraftData
	for _, item := range acCache.items {
		aircraftData := item.data
		age := time.Since(aircraftData.Modified).Minutes()

		// Only include aircraft updated in the last 1 minute
		if age <= 1 {
			aircraft = append(aircraft, aircraftData)
		}
	}

	return aircraft
}
