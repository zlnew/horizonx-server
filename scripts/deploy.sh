#!/bin/bash
set -e

# --- Configuration ---
APP_NAME="horizonx-server"
AGENT_NAME="horizonx-agent"
MIGRATE_TOOL="horizonx-migrate"
SEED_TOOL="horizonx-seed"

SERVER_SERVICE="horizonx-server"
AGENT_SERVICE="horizonx-agent"

INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/horizonx"
LOG_DIR="/var/log/horizonx"

# User System
SYS_USER="horizonx"
SYS_GROUP="horizonx"

# Source Paths
BIN_SRC="./bin"
ENV_SERVER_SRC="./.env.server.prod"
ENV_AGENT_SRC="./.env.agent.prod"

echo "=== HorizonX Secure Deployment (Split Config) ==="

# 1. Build
echo "🛠️  Building project..."
make build

# 2. Stop Services
echo "🛑 Stopping existing services..."
sudo systemctl stop $SERVER_SERVICE || true
sudo systemctl stop $AGENT_SERVICE || true

# 3. Create System User
if ! id -u $SYS_USER >/dev/null 2>&1; then
    echo "👤 Creating system user $SYS_USER..."
    sudo useradd -r -s /bin/false $SYS_USER
fi

# 4. Directories
echo "📂 Setting up directories..."
sudo mkdir -p $CONFIG_DIR $LOG_DIR
# Log dir owned by sys user
sudo chown -R $SYS_USER:$SYS_GROUP $LOG_DIR
sudo chmod 755 $LOG_DIR
# Config dir owned by root initially
sudo mkdir -p $CONFIG_DIR
sudo chmod 755 $CONFIG_DIR

# 5. Deploy Binaries
echo "🚀 Copying binaries..."
sudo cp $BIN_SRC/server $INSTALL_DIR/$APP_NAME
sudo cp $BIN_SRC/agent $INSTALL_DIR/$AGENT_NAME
sudo cp $BIN_SRC/migrate $INSTALL_DIR/$MIGRATE_TOOL
sudo cp $BIN_SRC/seed $INSTALL_DIR/$SEED_TOOL
sudo chmod +x $INSTALL_DIR/*

# 6. Deploy Configs
echo "📄 Deploying Separate Configurations..."

# A. Server Config
if [ -f "$ENV_SERVER_SRC" ]; then
    sudo cp $ENV_SERVER_SRC $CONFIG_DIR/server.env
    sudo chown root:$SYS_GROUP $CONFIG_DIR/server.env
    sudo chmod 640 $CONFIG_DIR/server.env
    echo "   -> server.env deployed (Secure)"
else
    echo "⚠️  FATAL: $ENV_SERVER_SRC not found!"
    exit 1
fi

# B. Agent Config
if [ -f "$ENV_AGENT_SRC" ]; then
    sudo cp $ENV_AGENT_SRC $CONFIG_DIR/agent.env
    sudo chown root:root $CONFIG_DIR/agent.env
    sudo chmod 600 $CONFIG_DIR/agent.env
    echo "   -> agent.env deployed (Secure)"
else
    echo "⚠️  WARNING: $ENV_AGENT_SRC not found! Agent service might fail."
fi

# 7. Run Migrations
echo "📦 Running Database Migrations..."
sudo sh -c "set -a; source $CONFIG_DIR/server.env; set +a; $INSTALL_DIR/$MIGRATE_TOOL -op=up"

# 8. Setup Systemd: SERVER
echo "⚙️  Configuring Server Service..."
SERVER_UNIT="/etc/systemd/system/${SERVER_SERVICE}.service"
sudo tee $SERVER_UNIT >/dev/null <<EOF
[Unit]
Description=HorizonX Core Server
After=network.target postgresql.service

[Service]
Type=simple
# Load SERVER specific env
EnvironmentFile=$CONFIG_DIR/server.env
ExecStart=$INSTALL_DIR/$APP_NAME
Restart=always
RestartSec=5
User=$SYS_USER
Group=$SYS_GROUP
StandardOutput=append:${LOG_DIR}/server.log
StandardError=append:${LOG_DIR}/server.error.log

[Install]
WantedBy=multi-user.target
EOF

# 9. Setup Systemd: AGENT
echo "⚙️  Configuring Agent Service..."
AGENT_UNIT="/etc/systemd/system/${AGENT_SERVICE}.service"
sudo tee $AGENT_UNIT >/dev/null <<EOF
[Unit]
Description=HorizonX Metrics Agent
After=network.target ${SERVER_SERVICE}.service

[Service]
Type=simple
# Load AGENT specific env
EnvironmentFile=$CONFIG_DIR/agent.env
ExecStart=$INSTALL_DIR/$AGENT_NAME
Restart=always
RestartSec=5
# Run as root for hardware access
User=root
Group=root
StandardOutput=append:${LOG_DIR}/agent.log
StandardError=append:${LOG_DIR}/agent.error.log

[Install]
WantedBy=multi-user.target
EOF

# 10. Start
echo "🔥 Reloading and Starting..."
sudo systemctl daemon-reload
sudo systemctl enable $SERVER_SERVICE $AGENT_SERVICE
sudo systemctl start $SERVER_SERVICE
sudo systemctl start $AGENT_SERVICE

echo "✅ Deployment Success! Configurations are isolated."
echo ""
echo "---------------------------------------------------------"
echo "🌱  HINT: FIRST TIME DEPLOY?"
echo "    If you need to seed the database, run this command manually:"
echo ""
echo "    sudo sh -c \"set -a; source $CONFIG_DIR/server.env; set +a; $INSTALL_DIR/$SEED_TOOL\""
echo "---------------------------------------------------------"
