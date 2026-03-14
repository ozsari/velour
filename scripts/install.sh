#!/usr/bin/env bash
set -euo pipefail

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

VELOUR_DIR="/opt/velour"
VELOUR_PORT="${VELOUR_PORT:-8585}"

banner() {
  echo -e "${CYAN}"
  echo ' __     __   _'
  echo ' \ \   / /__| | ___  _   _ _ __'
  echo '  \ \ / / _ \ |/ _ \| | | | '\''__|'
  echo '   \ V /  __/ | (_) | |_| | |'
  echo '    \_/ \___|_|\___/ \__,_|_|'
  echo ''
  echo -e "${NC}"
  echo -e "${BLUE}  Server Management Panel - Installer${NC}"
  echo ''
}

log() { echo -e "${GREEN}[+]${NC} $1"; }
warn() { echo -e "${YELLOW}[!]${NC} $1"; }
err() { echo -e "${RED}[x]${NC} $1"; exit 1; }

check_root() {
  if [ "$EUID" -ne 0 ]; then
    err "Please run as root: sudo bash install.sh"
  fi
}

check_os() {
  if [ ! -f /etc/os-release ]; then
    err "Unsupported OS. Velour requires a modern Linux distribution."
  fi
  . /etc/os-release
  log "Detected OS: $PRETTY_NAME"
}

install_docker() {
  if command -v docker &> /dev/null; then
    log "Docker already installed: $(docker --version)"
    return
  fi

  log "Installing Docker..."
  curl -fsSL https://get.docker.com | bash
  systemctl enable docker
  systemctl start docker
  log "Docker installed successfully"
}

install_docker_compose() {
  if docker compose version &> /dev/null; then
    log "Docker Compose already available"
    return
  fi

  log "Installing Docker Compose plugin..."
  apt-get update -qq && apt-get install -y -qq docker-compose-plugin 2>/dev/null || {
    # Fallback: install compose standalone
    COMPOSE_VERSION=$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep tag_name | cut -d '"' -f 4)
    curl -fsSL "https://github.com/docker/compose/releases/download/${COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
  }
  log "Docker Compose installed"
}

setup_velour() {
  log "Setting up Velour..."
  mkdir -p "$VELOUR_DIR"

  # Create docker-compose.yml
  cat > "$VELOUR_DIR/docker-compose.yml" << 'COMPOSE'
services:
  velour:
    image: ghcr.io/ozsari/velour:latest
    container_name: velour
    restart: unless-stopped
    ports:
      - "${VELOUR_PORT:-8585}:8585"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - velour_data:/opt/velour
    environment:
      - VELOUR_HOST=0.0.0.0
      - VELOUR_PORT=8585
      - VELOUR_DATA_DIR=/opt/velour

volumes:
  velour_data:
COMPOSE

  # Create .env
  cat > "$VELOUR_DIR/.env" << EOF
VELOUR_PORT=${VELOUR_PORT}
EOF

  log "Configuration created at $VELOUR_DIR"
}

start_velour() {
  log "Starting Velour..."
  cd "$VELOUR_DIR"
  docker compose pull
  docker compose up -d

  echo ''
  echo -e "${GREEN}========================================${NC}"
  echo -e "${GREEN}  Velour installed successfully!${NC}"
  echo -e "${GREEN}========================================${NC}"
  echo ''
  echo -e "  Open: ${CYAN}http://$(hostname -I | awk '{print $1}'):${VELOUR_PORT}${NC}"
  echo ''
  echo -e "  Commands:"
  echo -e "    Stop:    ${YELLOW}cd $VELOUR_DIR && docker compose down${NC}"
  echo -e "    Start:   ${YELLOW}cd $VELOUR_DIR && docker compose up -d${NC}"
  echo -e "    Logs:    ${YELLOW}cd $VELOUR_DIR && docker compose logs -f${NC}"
  echo -e "    Update:  ${YELLOW}cd $VELOUR_DIR && docker compose pull && docker compose up -d${NC}"
  echo ''
}

# Main
banner
check_root
check_os
install_docker
install_docker_compose
setup_velour
start_velour
