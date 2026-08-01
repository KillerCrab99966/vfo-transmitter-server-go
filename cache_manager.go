package main

import (
	"fmt"
	"net/http"
	"time"
)

func handleCacheManager(w http.ResponseWriter, r *http.Request) {
	message := ""

	// Carry out action
	switch getOrDefault(r.URL.Query()["action"], 0, "") {
	case "clear_all":
		acCache.clear()
		airpaceCache.clear()
		message = "<p style='color: green;'>✅ Cache cleared successfully</p>"

	case "clear_aircraft":
		aircraft := acCache.len()
		acCache.clear()
		message = fmt.Sprintf("<p style='color: green;'>✅ Cleared %v aircraft position entries</p>", aircraft)
	}

	// Return HTML
	fmt.Fprintf(w, `<!DOCTYPE html>
	<html>
	<head>
		<title>VFO APCu Cache Management</title>
		<style>
			body { font-family: Arial, sans-serif; margin: 20px; }
			table { border-collapse: collapse; width: 100%%; margin: 10px 0; }
			th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
			th { background-color: #f2f2f2; }
			.stats { background-color: #f9f9f9; padding: 15px; border-radius: 5px; margin: 10px 0; }
			.button { display: inline-block; padding: 10px 20px; background-color: #007cba; color: white; text-decoration: none; border-radius: 5px; margin: 5px; }
			.button:hover { background-color: #005a87; }
			.danger { background-color: #dc3545; }
			.danger:hover { background-color: #c82333; }
			.nav { margin-bottom: 20px; }
			.nav a { margin-right: 15px; color: #007cba; text-decoration: none; }
		</style>
	</head>
	<body>
		<div class="nav">
			<a href="test_aircraft">← Back to Aircraft Test</a>
			<a href="/">Home</a>
			<a href="radar">Radar Display</a>
		</div>

		<h1>🚀 VFO Cache Management (Database-Free)</h1>

		%s

		<h2>🛠️ Cache Operations</h2>
		<a href="?action=clear_aircraft" class="button" onclick="return confirm('Clear all aircraft position data?')">Clear Aircraft Data</a>
		<a href="?action=clear_all" class="button danger" onclick="return confirm('Clear ALL cache entries?')">Clear All Cache</a>
		<a href="?" class="button">Refresh Stats</a>

		<div class="stats">
			<p><strong>VFO Cache Breakdown:</strong></p>
			<ul>
				<li>Aircraft Positions: %v</li>
				<li>Airspace Data: %v</li>
			</ul>
		</div>

		<hr>
		<p><small>Last updated: %v</small></p>
	</body>
	</html>`, message, acCache.len(), airpaceCache.len(), time.Now().UTC().String())
}
