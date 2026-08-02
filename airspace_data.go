package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"strconv"
)

const rateLimitMax = 60 // 60 requests per minute

func handleAirspaceData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "public, max-age=86400")

	// Rate limiting
	if rateLimit {
		clientIP, _, _ := net.SplitHostPort(r.RemoteAddr)
		reqCount, ok := airspaceRateCache.get(clientIP)
		if ok {
			// Increment requests
			airspaceRateCache.set(clientIP, reqCount+1)

			if reqCount >= rateLimitMax {
				w.WriteHeader(http.StatusTooManyRequests)
				fmt.Fprint(w, `{"error":"Rate limit exceeded"}`)
				return
			}
		} else {
			// If the key does not exist, this is the first request
			airspaceRateCache.set(clientIP, 1)
		}
	}

	// Input validation
	source := getOrDefault(r.URL.Query()["source"], 0, "faa")
	dataset := getOrDefault(r.URL.Query()["dataset"], 0, "class")

	validSources := [2]string{"faa", "openaip"}
	validDatasets := [2]string{"class", "boundary"}

	if !slices.Contains(validSources[:], source) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":"Invalid source"}`)
		return
	}

	if !slices.Contains(validDatasets[:], dataset) {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":"Invalid dataset"}`)
		return
	}

	// Bounding box — all four values required
	keys := [4]string{"minLat", "minLon", "maxLat", "maxLon"}
	values := [4]float64{}
	for i, k := range keys {
		// Get v and parse to number
		v := getOrDefault(r.URL.Query()[k], 0, "")
		num, err := strconv.ParseFloat(v, 64)

		if v == "" || err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"error":"Missing or non-numeric parameter: %s"}`, k)
			return
		}

		values[i] = num
	}

	minLat := values[0]
	minLon := values[1]
	maxLat := values[2]
	maxLon := values[3]

	if minLat < -90 || maxLat > 90 || minLon < -180 || maxLon > 180 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":"Bounding box out of range"}`)
		return
	}

	if (maxLat-minLat) > 90 || (maxLon-minLon) > 90 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, `{"error":"Bounding box too large (max 90 degrees per side)"}`)
		return
	}

	// Cache key
	bboxString := fmt.Sprintf("%v,%v,%v,%v", minLat, minLon, maxLat, maxLon)
	hash := md5.Sum([]byte(bboxString))
	bboxHash := hex.EncodeToString(hash[:])

	var cacheKey string
	if source == "faa" {
		cacheKey = fmt.Sprintf("airspace_faa_%s_%s", dataset, bboxHash)
	} else {
		cacheKey = fmt.Sprintf("airspace_openaip_%s", bboxHash)
	}

	// Check for cached value
	cached, ok := airpaceCache.get(cacheKey)
	if ok {
		w.Header().Set("X-Cache", "HIT")
		fmt.Fprint(w, cached)
		return
	}

	// Fetch from upstream
	var result string
	var err error
	if source == "faa" {
		result, err = fetchFAA(dataset, minLat, minLon, maxLat, maxLon)
	} else {
		// Validate OpenAIP key
		apiKey := getOrDefault(r.URL.Query()["key"], 0, "")

		format := regexp.MustCompile(`^[a-zA-Z0-9]{24,64}$`)
		if !format.MatchString(apiKey) {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, `{"error":"Invalid or missing OpenAIP API key"}`)
			return
		}

		result, err = fetchOpenAIP(apiKey, minLat, minLon, maxLat, maxLon)
	}

	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		fmt.Fprint(w, `{"error":"Failed to fetch airspace data from upstream"}`)
		return
	}

	// Cache and return value
	airpaceCache.set(cacheKey, result)
	w.Header().Set("X-Cache", "MISS")
	fmt.Fprint(w, result)
}

func fetchFAA(dataset string, minLat, minLon, maxLat, maxLon float64) (string, error) {
	endpoints := map[string]string{
		"class":    "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/Class_Airspace/FeatureServer/0/query",
		"boundary": "https://services6.arcgis.com/ssFJjBXIUyZDrSYZ/arcgis/rest/services/Boundary_Airspace/FeatureServer/0/query",
	}

	// Add URL parameters
	reqURL, err := url.Parse(endpoints[dataset])
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("where", "1=1")
	params.Add("geometry", fmt.Sprintf("%v,%v,%v,%v", minLon, minLat, maxLon, maxLat))
	params.Add("geometryType", "esriGeometryEnvelope")
	params.Add("inSR", "4326")
	params.Add("spatialRel", "esriSpatialRelIntersects")
	params.Add("outFields", "*")
	params.Add("f", "geojson")
	params.Add("outSR", "4326")
	reqURL.RawQuery = params.Encode()

	// Make the request
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func fetchOpenAIP(apiKey string, minLat, minLon, maxLat, maxLon float64) (string, error) {
	centerLat := (minLat + maxLat) / 2
	centerLon := (minLon + maxLon) / 2

	// Approximate radius in NM — 1 degree ≈ 60 NM; add margin
	distNM := int(math.Max(maxLat-minLat, maxLon-minLon)*60/2 + 10)

	// Add URL parameters
	reqURL, err := url.Parse("https://api.core.openaip.net/api/airspaces")
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("pos", fmt.Sprintf("%v,%v", centerLat, centerLon))
	params.Add("dist", strconv.Itoa(distNM))
	reqURL.RawQuery = params.Encode()

	// Make the request
	req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
	if err != nil {
		return "", err
	}

	// API key header
	req.Header.Set("x-openaip-api-key", apiKey)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// Decode JSON
	var data map[string]any
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	rawItems, ok := data["items"]
	if !ok || rawItems == nil {
		return "", errors.New("missing items field in response")
	}

	itemsList, ok := rawItems.([]any)
	if !ok {
		return "", errors.New("invalid items format")
	}

	// Normalise to the same GeoJSON FeatureCollection schema as FAA data
	features := make([]map[string]any, 0, len(itemsList))

	for _, rawItem := range itemsList {
		item, ok := rawItem.(map[string]any)
		if !ok || item["geometry"] == nil {
			continue
		}

		// Key extraction
		name, _ := item["name"].(string)

		class := ""
		if icao, ok := item["icaoClass"].(string); ok {
			class = icao
		} else if t, ok := item["type"].(string); ok {
			class = t
		}

		typeCode, _ := item["type"].(string)

		var upperVal any = nil
		var upperUom string = ""
		if upper, ok := item["upperLimit"].(map[string]any); ok {
			upperVal = upper["value"]
			upperUom, _ = upper["unit"].(string)
		}

		var lowerVal any = nil
		var lowerUom string = ""
		if lower, ok := item["lowerLimit"].(map[string]any); ok {
			lowerVal = lower["value"]
			lowerUom, _ = lower["unit"].(string)
		}

		feature := map[string]any{
			"type":     "Feature",
			"geometry": item["geometry"],
			"properties": map[string]any{
				"NAME":      name,
				"CLASS":     class,
				"TYPE_CODE": typeCode,
				"UPPER_VAL": upperVal,
				"UPPER_UOM": upperUom,
				"LOWER_VAL": lowerVal,
				"LOWER_UOM": lowerUom,
				"_source":   "openaip",
			},
		}
		features = append(features, feature)
	}

	geoJSON, err := json.Marshal(map[string]any{
		"type":     "FeatureCollection",
		"features": features,
	})
	if err != nil {
		return "", err
	}

	return string(geoJSON), nil
}
