# Configuring the Server

Server configuration is managed via a TOML config file (`config.toml`). An example configuration has been included [here](config.toml).

## Config File Location

The directory where the server looks for `config.toml` depends on the application's [environment](../main.go#L38) setting:

- **Development Mode (`Environment == "development"`) — Default**: The server searches for `config.toml` in the current working directory where the executable is launched.
- **Production Mode (`Environment == "production"`)**: The server searches for `config.toml` in the same directory as the compiled binary.

If `Environment` is not set to `production` or `development`, the program will panic.

## Settings

```toml
address = "0.0.0.0:8080"
```

The address for the server to listen on.
<br>

```toml
pin = ""
```

An optional pin to restrict who can submit data. Only affects the `/transmit` endpoint. Authentication is disabled if empty.
<br>

```toml
rate_limiting = true
```

Enable rate-limiting on the `/transmit` and `/airspace_data` endpoints.
- Minimum 1s per callsign and IP pair between transmissions
- Maximum 60 requests per minute per IP to `/airspace_data`
<br>

```toml
debug = false
```
If the debug/development pages (`debug_aircraft` & `test_aircraft`) should be served.
