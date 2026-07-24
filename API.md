# API Endpoints

This port maintains the same base API as the original to ensure full compatibility. For more comprehensive documentation, visit [https://github.com/jonbeckett/vfo-transmitter-server/wiki/api-endpoints](https://github.com/jonbeckett/vfo-transmitter-server/wiki/api-endpoints).

Each endpoint has a version with no suffix and a version with the original suffix for backwards-compatibility.

## Data Ingestion

```
GET /transmit
GET /transmit.php

POST /transmit
POST /transmit.php
```
Receives aircraft position data from the VFO Transmitter client.

## Data APIs

```
GET /radar_data
GET /radar_data.php

GET /status_json
GET /status_json.php
```
Both return the aircraft array as JSON, with the same fields as the original `status_json.php` ([https://github.com/jonbeckett/vfo-transmitter-server/wiki/api-endpoints#get-status_jsonphp](https://github.com/jonbeckett/vfo-transmitter-server/wiki/api-endpoints#get-status_jsonphp)).

## Pages (HTML responses)

```
GET /radar
GET /radar.html
```
Full-screen interactive radar.
<br>
<br>
```
GET /status
GET /status.html
```
Aircraft status dashboard.
<br>
<br>
```
GET /embed
GET /embed.html
```
Minimal radar widget. Suitable for `<iframe>` embedding.