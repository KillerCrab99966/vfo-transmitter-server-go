# Virtual Flight Online Transmitter Server - Golang version

This repository is a complete backend port of the Jonathan Beckett's [Virtual Flight Online Transmitter Server](https://github.com/jonbeckett/vfo-transmitter-server) from PHP to Go.

It is designed to be fully compatible with the MSFS and X-Plane clients as it maintains the same API.

> A real-time aircraft tracking server for Microsoft Flight Simulator and X-Plane. Receives position data from transmitter clients and serves it via an interactive web-based radar display.


## TODOs

- Rate limiting
- Flexability with a JSON config file or env variables

## API Endpoints

Same functions and parameters as the original: [https://github.com/jonbeckett/vfo-transmitter-server/wiki/api-endpoints](https://github.com/jonbeckett/vfo-transmitter-server/wiki/api-endpoints)

Each endpoint has a version with no suffix and a version with the original suffix for backwards-compatibility.

### Data Ingestion

```
GET /transmit
GET /transmit.php

POST /transmit
POST /transmit.php
```
Receives aircraft position data from the VFO Transmitter client.

### Data APIs

```
GET /radar_data
GET /radar_data.php

GET /status_json
GET /status_json.php
```
Both return the aircraft array as JSON, with the same fields as the original `status_json.php` ([https://github.com/jonbeckett/vfo-transmitter-server/wiki/api-endpoints#get-status_jsonphp](https://github.com/jonbeckett/vfo-transmitter-server/wiki/api-endpoints#get-status_jsonphp)).

## Related Projects

- **[Original Project](https://github.com/jonbeckett/vfo-transmitter-server)** — Jonathan Beckett's PHP-Based server
- **[VFO Transmitter MSFS Client](https://github.com/jonbeckett/vfo-transmitter-client-msfs)** — Windows app that reads MSFS data and posts to this server
- **[VFO Transmitter X-Plane Client](https://github.com/jonbeckett/vfo-transmitter-client-xplane)** — FlyWithLua plugin that reads X-Plane data and posts to this server
- **[Virtual Flight Online](https://virtualflight.online)** — Community homepage

## License

Same as original:
> Open source — for educational and simulation use.