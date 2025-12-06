# horizonx-server

A server for managing and monitoring your infrastructure, with features for serving server metrics, SSH in the browser, and deployment management.

## Features

- **Server Metrics:** Real-time monitoring of CPU, GPU, memory, disk, network, and uptime. Metrics are available via a REST API and a websocket stream.
- **SSH in Browser (WIP):** Access your server's terminal directly from your web browser.
- **Deployment Management (WIP):** Manage your application deployments.

## Requirements

- Go 1.25.4+.
- Linux host with `/proc` and `/sys` available.

## Configuration

The agent can be configured via environment variables, which can also be specified in a `.env` file at the root of the project.

| Variable          | Description                                                                                            | Default |
| ----------------- | ------------------------------------------------------------------------------------------------------ | ------- |
| `HTTP_ADDR`       | Sets the address and port for the HTTP server (e.g., `:8080`).                                         | `:3000` |
| `SCRAPE_INTERVAL` | Defines the interval for metric collection. This should be a duration string, such as `1s` or `500ms`. | `1s`    |
| `LOG_LEVEL`       | Sets the logging level. Supported values are `debug`, `info`, `warn`, and `error`.                     | `info`  |
| `LOG_FORMAT`      | Determines the log output format. Supported values are `text` and `json`.                              | `text`  |

## Building and Running

### Commands

The project includes a `Makefile` that simplifies common tasks, but standard Go commands can also be used.

- **To run the server locally**:
  - `go run ./cmd/horizonx-server`

- **To build the server**:
  - `make build`
  - This command compiles the server and places the binary at `bin/horizonx-server`.

- **To run tests**:
  - `go test ./...`
  - (Note: There are currently no tests.)

- **To clean the build artifacts**:
  - `make clean`
  - This removes the `bin` directory.

## API Endpoints

### Metrics

#### `GET /metrics`

Returns the latest snapshot of all metrics.

#### `GET /ws`

Upgrades the connection to a websocket to stream metrics. The websocket uses a channel-based system. To receive metrics, the client must send a JSON message to subscribe to the `metrics` channel.

**Subscribe Message:**

```json
{
  "type": "subscribe",
  "channel": "metrics"
}
```

**Broadcast Message:**
The server will then stream messages with the following format:

```json
{
  "channel": "metrics",
  "payload": { ... }
}
```
