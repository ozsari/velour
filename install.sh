#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
YELLOW='\033[1;33m'
BOLD='\033[1m'
NC='\033[0m'

VELOUR_VERSION="0.1.0"
VELOUR_DIR="/opt/velour"
VELOUR_BIN="/usr/local/bin/velour"
VELOUR_USER="velour"

echo -e "${CYAN}"
echo " __     __   _                  "
echo " \ \   / /__| | ___  _   _ _ __ "
echo "  \ \ / / _ \ |/ _ \| | | | '__|"
echo "   \ V /  __/ | (_) | |_| | |   "
echo "    \_/ \___|_|\___/ \__,_|_|   "
echo -e "${NC}"
echo -e "${BOLD}Velour Installer v${VELOUR_VERSION}${NC}"
echo ""

# Check root
if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}Error: Please run as root (sudo)${NC}"
  exit 1
fi

# Check OS
if [ ! -f /etc/os-release ]; then
  echo -e "${RED}Error: Unsupported operating system${NC}"
  exit 1
fi

. /etc/os-release
echo -e "${GREEN}OS:${NC} $PRETTY_NAME"
echo -e "${GREEN}Arch:${NC} $(uname -m)"
echo ""

# Select install mode
echo -e "${BOLD}Select installation mode:${NC}"
echo ""
echo -e "  ${BLUE}1)${NC} ${BOLD}Docker${NC}  - Apps run in containers"
echo -e "     ${CYAN}Easier setup, isolated apps, works on any Linux${NC}"
echo ""
echo -e "  ${GREEN}2)${NC} ${BOLD}Native${NC}  - Apps installed directly on system"
echo -e "     ${CYAN}Better performance, lower overhead, like Swizzin${NC}"
echo ""

while true; do
  read -p "$(echo -e ${BOLD}Choice [1/2]: ${NC})" choice </dev/tty
  case $choice in
    1) INSTALL_MODE="docker"; break;;
    2) INSTALL_MODE="native"; break;;
    *) echo -e "${RED}Please enter 1 or 2${NC}";;
  esac
done

echo ""
echo -e "${GREEN}Selected:${NC} ${BOLD}${INSTALL_MODE}${NC} mode"
echo ""

# ── Common setup ──

echo -e "${BLUE}[1/5]${NC} Creating velour user and directories..."
id -u $VELOUR_USER &>/dev/null || useradd -r -s /usr/sbin/nologin -d $VELOUR_DIR -m $VELOUR_USER
mkdir -p $VELOUR_DIR/{downloads,tv,movies,music,books,comics,audiobooks}
chown -R $VELOUR_USER:$VELOUR_USER $VELOUR_DIR

echo -e "${BLUE}[2/5]${NC} Installing dependencies..."
apt-get update -qq
apt-get install -y -qq curl wget sqlite3 > /dev/null

# ── Mode-specific setup ──

if [ "$INSTALL_MODE" = "docker" ]; then
  echo -e "${BLUE}[3/5]${NC} Setting up Docker..."

  if command -v docker &>/dev/null; then
    echo -e "  ${GREEN}Docker already installed:${NC} $(docker --version)"
  else
    echo -e "  Installing Docker..."
    curl -fsSL https://get.docker.com | sh
    systemctl enable docker
    systemctl start docker
    usermod -aG docker $VELOUR_USER
    echo -e "  ${GREEN}Docker installed successfully${NC}"
  fi

else
  echo -e "${BLUE}[3/5]${NC} Preparing native environment..."
  # Install common tools needed for native installs
  apt-get install -y -qq git unzip gnupg apt-transport-https > /dev/null
  echo -e "  ${GREEN}Native dependencies ready${NC}"
fi

# ── Download and install Velour binary ──

echo -e "${BLUE}[4/5]${NC} Downloading Velour..."

ARCH=$(uname -m)
case $ARCH in
  x86_64) ARCH="amd64";;
  aarch64|arm64) ARCH="arm64";;
  *) echo -e "${RED}Unsupported architecture: $ARCH${NC}"; exit 1;;
esac

curl -fsSL "https://github.com/ozsari/velour/releases/latest/download/velour_linux_${ARCH}" -o $VELOUR_BIN
chmod +x $VELOUR_BIN

# Write config
cat > $VELOUR_DIR/config.json <<EOF
{
  "version": "$VELOUR_VERSION",
  "host": "0.0.0.0",
  "port": 8585,
  "data_dir": "$VELOUR_DIR",
  "db_path": "$VELOUR_DIR/velour.db",
  "install_mode": "$INSTALL_MODE"
}
EOF

chown $VELOUR_USER:$VELOUR_USER $VELOUR_DIR/config.json

echo -e "${BLUE}[5/5]${NC} Creating systemd service..."

cat > /etc/systemd/system/velour.service <<EOF
[Unit]
Description=Velour - Server Management Panel
After=network.target$([ "$INSTALL_MODE" = "docker" ] && echo " docker.service")
$([ "$INSTALL_MODE" = "docker" ] && echo "Requires=docker.service")

[Service]
Type=simple
User=root
Environment=VELOUR_DATA_DIR=$VELOUR_DIR
Environment=VELOUR_INSTALL_MODE=$INSTALL_MODE
ExecStart=$VELOUR_BIN
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable velour
systemctl start velour

echo ""
echo -e "${GREEN}${BOLD}Velour installed successfully!${NC}"
echo ""
echo -e "  ${BOLD}Mode:${NC}      $INSTALL_MODE"
echo -e "  ${BOLD}Data:${NC}      $VELOUR_DIR"
echo -e "  ${BOLD}Config:${NC}    $VELOUR_DIR/config.json"
echo -e "  ${BOLD}Panel:${NC}     http://$(hostname -I | awk '{print $1}'):8585"
echo ""
echo -e "  ${CYAN}Start:${NC}     systemctl start velour"
echo -e "  ${CYAN}Status:${NC}    systemctl status velour"
echo -e "  ${CYAN}Logs:${NC}      journalctl -u velour -f"
echo ""
echo -e "  Open the panel in your browser to create your admin account."
echo ""
