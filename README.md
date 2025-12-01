# go-monitor-agent

Lightweight Go daemon that scrapes local CPU, memory, and disk stats from `/proc` and `/sys` and exposes a JSON snapshot at `/metrics`.

## Requirements

- Go 1.25.4+.
- Linux host with `/proc` and `/sys` available.
- Optional: `dmidecode` (for DIMM specs) and readable hwmon/RAPL entries for temps/power; some of these may need root/capabilities.

## Quick start

- Run locally: `go run ./cmd/agent` (binds to `:3000`).
- Configure address: set `HTTP_ADDR` in the env or a `.env` file, e.g. `HTTP_ADDR=:8080`.
- Build: `go build -o bin/agent ./cmd/agent`.
- Clean build artifact: `make clean`.
- Tests: `go test ./...` (none yet).

## Metrics endpoint

- `GET /metrics` returns the latest snapshot (refreshed every second by the scheduler).
- Example response:
  ```json
  {
    "cpu": {
      "spec": {
        "vendor": "GenuineIntel",
        "model": "11th Gen Intel(R) Core(TM) i7",
        "cores": 8,
        "threads": 16,
        "arch": "amd64",
        "base_freq": 1200,
        "max_freq": 4800
      },
      "usage": 12.5,
      "per_core": [8.2, 14.1, 13.3, 9.9],
      "watt": 7.8,
      "temp": 55.0,
      "frequency": 2200.0
    },
    "memory": {
      "specs": [
        {
          "size": "16 GB",
          "type": "DDR4",
          "speed": "3200 MT/s",
          "manufacturer": "Micron",
          "part_number": "ABC123",
          "form_factor": "SODIMM"
        }
      ],
      "mem_total": 15.8,
      "mem_available": 9.4,
      "mem_used": 6.4,
      "swap_total": 2.0,
      "swap_free": 2.0,
      "swap_used": 0.0
    },
    "disk": {
      "total": 476.3,
      "free": 302.1,
      "used": 174.2,
      "temp": 41.5
    }
  }
  ```

## Project layout

- `cmd/agent`: Entrypoint wiring config/logger and starting the agent.
- `internal/agent`: Builds the registry, registers CPU/memory/disk collectors, and starts HTTP.
- `internal/core`: Metric types, registry, and 1s scheduler that refreshes snapshots.
- `internal/collector/cpu`: Reads `/proc/stat`, `/proc/cpuinfo`, cpufreq, hwmon, and RAPL for usage, specs, temps, power, and frequency.
- `internal/collector/memory`: Pulls `/proc/meminfo` for usage and `dmidecode` for DIMM specs.
- `internal/collector/disk`: Uses `statfs` for capacity/usage and hwmon (e.g., NVMe) for temperature.
- `internal/transport/http`: Exposes `/metrics`.
- `internal/infra`: Config loader (`HTTP_ADDR`, `.env`) and logger interface/implementation.
- `pkg`: Small shared helpers.
- `bin/`: Build output; avoid committing binaries unless intentional.

## Notes

- Collectors fail softly: if a data source is missing/unreadable, values default to zero and logging records the error.
- The HTTP server runs until context cancellation; the scheduler stops when the agent shuts down.
