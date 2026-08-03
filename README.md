# Virtual Flight Online Transmitter Server - Golang version

This repository is a complete backend port of Jonathan Beckett's [Virtual Flight Online Transmitter Server](https://github.com/jonbeckett/vfo-transmitter-server) from PHP to Go.

It is designed to be fully compatible with the MSFS and X-Plane clients as it maintains the same API.

About the server:\
A real-time aircraft tracking server for Microsoft Flight Simulator and X-Plane. Receives position data from transmitter clients and serves it via an interactive web-based radar display.

## Setup/Deployment

1. **Download the Release:** Download the latest binary for your operating system from [Releases](https://github.com/KillerCrab99966/vfo-transmitter-server-go/releases) or build for production from source (with `-ldflags="-X 'main.Environment=production'"`).

2. **Set Up the Configuration File:** Place the downloaded binary and a `config.toml` in the same directory:
	```text
   ├── vfo-transmitter-server
   └── config.toml
	```

3. **Execute the binary:** Run the server via the command line or double-clicking the binary.

## API Endpoints

See [API.md](docs/API.md)

## Configuration

See [CONFIG.md](docs/CONFIG.md)

## Related Projects

- **[Original Project](https://github.com/jonbeckett/vfo-transmitter-server)** — Jonathan Beckett's PHP-Based server
- **[VFO Transmitter MSFS Client](https://github.com/jonbeckett/vfo-transmitter-client-msfs)** — Windows app that reads MSFS data and posts to this server
- **[VFO Transmitter X-Plane Client](https://github.com/jonbeckett/vfo-transmitter-client-xplane)** — FlyWithLua plugin that reads X-Plane data and posts to this server
- **[Virtual Flight Online](https://virtualflight.online)** — Community homepage

## License

Same as original:

> Open source — for educational and simulation use.
