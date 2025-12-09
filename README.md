# HorizonX Server

**HorizonX Server** is a robust, high-performance system monitoring backend built with Go. It provides real-time visibility into your infrastructure's health, streaming comprehensive metrics (CPU, GPU, Memory, Disk, Network) to clients via WebSockets and providing snapshots via a REST API.

It is designed for both **development** ease and **production** reliability, featuring built-in authentication, automatic database migrations, and a systemd deployment script.

## 🚀 Key Features

- **Real-Time Streaming:** High-frequency metric updates pushed to clients via WebSockets.
- **Comprehensive Metrics:**
  - **CPU:** Per-core usage, frequency, temperature, power consumption.
  - **GPU:** Full NVIDIA support (Usage, VRAM, Fan Speed, Power, Temp) via `nvidia-smi`.
  - **Memory:** RAM usage, Swap usage, and availability.
  - **Disk:** Usage stats, partition mapping, and temperature monitoring.
  - **Network:** Real-time throughput (RX/TX) and link speeds.
  - **System:** Uptime, kernel version, OS details.
- **Secure:**
  - JWT-based authentication (HttpOnly/Secure cookies).
  - CSRF protection.
  - SQLite-backed user management with `bcrypt` password hashing.
- **Production Ready:**
  - Includes a `deploy.sh` script for auto-configuring systemd services.
  - Runs as `root` for full sensor access while dropping privileges for logging/data where possible.
  - Configurable logging (JSON/Text) and CORS.

## 🛠️ Technology Stack

- **Language:** Go (1.25+)
- **Database:** SQLite (embedded, auto-migrating)
- **Communication:** REST (Standard Lib) & WebSockets (`gorilla/websocket`)
- **Auth:** JWT (`golang-jwt/jwt`) & Bcrypt
- **Logging:** `log/slog` (Structured Logging)

## 📋 Prerequisites

- **Go 1.25+**
- **C Compiler (GCC):** Required for SQLite (`go-sqlite3` uses CGO).
- **NVIDIA Drivers (Optional):** Required for GPU metrics.

## ⚡ Getting Started (Development)

### 1. Installation

```bash
git clone https://github.com/yourusername/horizonx-server.git
cd horizonx-server
```

### 2. Configuration

Create a `.env` file based on the example. **You must set `JWT_SECRET`.**

```bash
cp .env.example .env
```

**Recommended `.env` for Dev:**

```ini
HTTP_ADDR=:3000
JWT_SECRET=dev-secret-change-me
SCRAPE_INTERVAL=3s
LOG_LEVEL=debug
LOG_FORMAT=text
DB_PATH=./dev.db
ALLOWED_ORIGINS=http://localhost:5173
```

### 3. Running

You can use the provided `Makefile` or standard Go commands.

```bash
# Build binary
make build

# Run server
make run

# Clean artifacts
make clean
```

The server will start on `http://localhost:3000` (or your configured port).
The SQLite database will be automatically created at `./dev.db` (or your configured path).

## 🌍 Getting Started (Production)

For deploying to a Linux server, use the included deployment script. It handles:

1.  Building the binary.
2.  Creating a `horizonx` system user.
3.  Setting up directories (`/etc/horizonx`, `/var/lib/horizonx`, `/var/log/horizonx`).
4.  Installing a systemd service running as `root` (required for hardware sensors) but logging safely.

### Deployment Steps

1.  **Prepare Production Config:**
    Create a `.env.production` file with your production secrets.

    ```bash
    cp .env.example .env.production
    # Edit .env.production with real JWT_SECRET and paths
    ```

2.  **Run Deploy Script:**

    ```bash
    chmod +x scripts/deploy.sh
    ./scripts/deploy.sh
    ```

3.  **Manage Service:**
    ```bash
    sudo systemctl status horizonx-server
    sudo journalctl -u horizonx-server -f
    ```

## ⚙️ Configuration Reference

| Environment Variable | Description                                          | Default    |
| :------------------- | :--------------------------------------------------- | :--------- |
| `HTTP_ADDR`          | Host and port to listen on.                          | `:3000`    |
| `JWT_SECRET`         | **REQUIRED.** Secret key for signing tokens.         | (None)     |
| `JWT_EXPIRY`         | Token validity duration.                             | `24h`      |
| `SCRAPE_INTERVAL`    | Frequency of metric collection.                      | `3s`       |
| `LOG_LEVEL`          | Log verbosity (`debug`, `info`, `warn`, `error`).    | `info`     |
| `LOG_FORMAT`         | Log output format (`text` for dev, `json` for prod). | `text`     |
| `ALLOWED_ORIGINS`    | CSV list of allowed CORS origins.                    | (Empty)    |
| `DB_PATH`            | Path to the SQLite database file.                    | `./app.db` |

## 🔌 API Reference

### Authentication

The API uses **HttpOnly Cookies** for session management (`access_token`).

#### Register

- **POST** `/auth/register`
- **Body:** `{"email": "user@example.com", "password": "securepassword"}`

#### Login

- **POST** `/auth/login`
- **Body:** `{"email": "user@example.com", "password": "securepassword"}`
- **Response:** Sets `access_token` cookie. Returns user info.

#### Logout

- **POST** `/auth/logout`
- **Action:** Clears auth cookies.

### Metrics (REST)

- **GET** `/metrics`
- **Auth Required:** Yes (Cookie)
- **Response:** JSON snapshot of the latest system metrics.

### Real-Time Stream (WebSocket)

- **URL:** `ws://your-server:3000/ws`
- **Auth:** Handled via Cookie (browser) or headers.

**Protocol:**

1.  **Connect** to the endpoint.
2.  **Subscribe** to the metrics channel.
    ```json
    {
      "type": "subscribe",
      "channel": "metrics"
    }
    ```
3.  **Receive** events.
    ```json
    {
      "channel": "metrics",
      "payload": {
        "cpu": { ... },
        "memory": { ... },
        "gpu": [ ... ],
        "network": { ... },
        "disk": [ ... ]
      }
    }
    ```

## ❓ Troubleshooting

- **`nvidia-smi` not found:**
  GPU metrics will be skipped if NVIDIA drivers are not installed. This is normal on non-NVIDIA systems.

- **Permission Denied (Sensors):**
  Reading detailed hardware info (like power/temps) often requires root. In production, the systemd service runs as root to ensure full visibility.

- **SQLite CGO Errors:**
  Ensure you have `gcc` installed. Windows users might need MinGW or TDM-GCC.

## License

[MIT](LICENSE)

