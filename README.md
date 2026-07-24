# Virtual Flight Online Transmitter Server - Golang version

This repository is a complete backend port of the Jonathan Beckett's [Virtual Flight Online Transmitter Server](https://github.com/jonbeckett/vfo-transmitter-server) from PHP to Go.

It is designed to be fully compatible with the MSFS and X-Plane clients as it maintains the same API.

About the server:\
A real-time aircraft tracking server for Microsoft Flight Simulator and X-Plane. Receives position data from transmitter clients and serves it via an interactive web-based radar display.


## TODOs

- Rate limiting
- Flexability with a JSON config file or env variables
- Airspace data
	- Enable button once enabled ([static/js/radar.js:2151](static/js/radar.js#L2151))

## API Endpoints

See [API.md](API.md)

## Related Projects

- **[Original Project](https://github.com/jonbeckett/vfo-transmitter-server)** — Jonathan Beckett's PHP-Based server
- **[VFO Transmitter MSFS Client](https://github.com/jonbeckett/vfo-transmitter-client-msfs)** — Windows app that reads MSFS data and posts to this server
- **[VFO Transmitter X-Plane Client](https://github.com/jonbeckett/vfo-transmitter-client-xplane)** — FlyWithLua plugin that reads X-Plane data and posts to this server
- **[Virtual Flight Online](https://virtualflight.online)** — Community homepage

## License

Same as original:
> Open source — for educational and simulation use.