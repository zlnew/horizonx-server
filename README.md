# go-monitor-agent

Lightweight Linux metrics agent that scrapes `/proc` and `/sys` and can either emit a single snapshot, stream snapshots, or serve the latest JSON at `/metrics`. Collectors are defensive: unreadable files just zero the value so the agent keeps running.

## Requirements

- Go 1.25.4+.
- Linux host with `/proc` and `/sys` available.
- Optional kernel interfaces/sensors:
  - RAPL or hwmon power/temperature files for CPU and GPU.
  - `/sys/class/drm/*/gpu_busy_percent` and `mem_info_vram_*` for GPU utilization and VRAM.
  - NVMe hwmon temperature for disk temperature.
- Missing files are tolerated; collectors will emit zero values instead of failing.

## Configuration

The agent can be configured via environment variables, which can also be specified in a `.env` file at the root of the project.

| Variable          | Description                                                                                                                          | Default |
| ----------------- | ------------------------------------------------------------------------------------------------------------------------------------ | ------- |
| `HTTP_ADDR`       | Sets the address and port for the HTTP server in `serve` mode (e.g., `:8080`).                                                       | `:3000` |
| `SCRAPE_INTERVAL` | Defines the interval for metric collection in `stream` and `serve` modes. This should be a duration string, such as `1s` or `500ms`. | `1s`    |
| `LOG_LEVEL`       | Sets the logging level. Supported values are `debug`, `info`, `warn`, and `error`.                                                   | `info`  |
| `LOG_FORMAT`      | Determines the log output format. Supported values are `text` and `json`.                                                            | `text`  |

## Building and Running

### Commands

The project includes a `Makefile` that simplifies common tasks, but standard Go commands can also be used.

- **To run the agent locally**:
  - `go run ./cmd/agent [mode]`
  - Where `[mode]` can be `snapshot`, `stream`, or `serve`.
  - If no mode is specified, it defaults to `serve`.

- **To build the agent**:
  - `make build`
  - This command compiles the agent and places the binary at `bin/monitor-agent`.

- **To run tests**:
  - `go test ./...`
  - (Note: The `README.md` indicates that there are currently no tests.)

- **To clean the build artifacts**:
  - `make clean`
  - This removes the `bin` directory.

## Metrics

Collectors refresh every second and degrade gracefully when inputs are unreadable.

- CPU: usage and per-core usage percent (averaged since boot from `/proc/stat`), temperature via hwmon CPU sensors, current frequency (MHz) from cpufreq, and `power_watt` from RAPL or hwmon when available.
- GPU: discovers DRM cards under `/sys/class/drm/card*` (skips render nodes) and reports vendor/model (using PCI IDs where possible), per-card name, `core_usage_percent` from `gpu_busy_percent`, VRAM total/used/percent from `mem_info_vram_*`, temperature/power from hwmon, and `fan_speed_percent` when PWM or fan RPM are readable. Still returned as a slice even for a single GPU.
- Memory: total, available, used, and swap total/free/used in GiB from `/proc/meminfo`.
- Disk: enumerates block devices under `/sys/class/block` (skips loop/ram/dm). For each disk it reports `raw_size_gb` and temperature via `/sys/block/<disk>/device` hwmon, plus per-filesystem usage for each mounted partition discovered in `/proc/self/mountinfo`. Filesystem entries include device name, mountpoint, total/used/free, and percent used.
- Network: aggregate of non-loopback interfaces from `/proc/net/dev` with cumulative `rx_bytes`/`tx_bytes` and Mbps rates computed between samples (first sample reports `0` speeds).
- Uptime: `uptime_seconds` from `/proc/uptime`.

### Endpoint

`GET /metrics` returns the latest snapshot:

`GET /ws` upgrades the connection to a websocket and streams the latest snapshot at the `SCRAPE_INTERVAL`.

```json
{
  "cpu": {
    "usage": 12.5,
    "per_core": [8.2, 14.1, 13.3, 9.9],
    "temperature": 55.0,
    "frequency": 2200.0,
    "power_watt": 7.8
  },
  "gpu": [
    {
      "id": 0,
      "card": "card0",
      "vendor": "AMD",
      "model": "Radeon 6800 XT",
      "temperature": 48.3,
      "core_usage_percent": 15.7,
      "vram_total_gb": 8.0,
      "vram_used_gb": 2.1,
      "vram_percent": 26.2,
      "power_watt": 65.2,
      "fan_speed_percent": 38.0
    }
  ],
  "memory": {
    "total_gb": 15.8,
    "used_gb": 6.4,
    "available_gb": 9.4,
    "swap_total_gb": 2.0,
    "swap_free_gb": 2.0,
    "swap_used_gb": 0.0
  },
  "disk": [
    {
      "name": "nvme0n1",
      "raw_size_gb": 953.8,
      "temperature": 41.5,
      "filesystems": [
        {
          "device": "nvme0n1p2",
          "mountpoint": "/",
          "total_gb": 476.3,
          "used_gb": 174.2,
          "free_gb": 302.1,
          "percent": 56.0
        },
        {
          "device": "nvme0n1p3",
          "mountpoint": "/home",
          "total_gb": 447.5,
          "used_gb": 120.2,
          "free_gb": 327.3,
          "percent": 26.8
        }
      ]
    }
  ],
  "network": {
    "rx_bytes": 123456789,
    "tx_bytes": 987654321,
    "rx_speed": 1.2,
    "tx_speed": 8.5
  },
  "uptime_seconds": 12345
}
```

## Notes

- The HTTP server runs until process exit; the scheduler stops when context is canceled.
- Build artifacts land in `bin/`; avoid committing binaries unless intentional.
