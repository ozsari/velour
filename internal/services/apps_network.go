package services

import "github.com/ozsari/velour/internal/models"

func init() {
	Registry = append(Registry,
		models.ServiceDefinition{
			ID: "wireguard", Name: "WireGuard", Description: "Fast, modern, secure VPN tunnel. Simple yet powerful.", Icon: "wireguard", Category: "network", Image: "lscr.io/linuxserver/wireguard:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 51820, Container: 51820, Protocol: "udp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/wireguard/config", Container: "/config"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
			Native: &models.NativeConfig{
				Method:      "apt", ServiceName: "wg-quick@wg0", Port: 0,
				AptPackages: []string{"wireguard", "wireguard-tools"},
				ConfigDir:   "/etc/wireguard",
			}},
		models.ServiceDefinition{
			ID: "thelounge", Name: "The Lounge", Description: "Modern, self-hosted IRC client with always-on connectivity and web interface.", Icon: "thelounge", Category: "network", Image: "lscr.io/linuxserver/thelounge:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 9000, Container: 9000, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/thelounge/config", Container: "/config"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
			Native: &models.NativeConfig{
				Method:       "script", ServiceName: "thelounge", Port: 9000,
				Dependencies: []string{"curl"},
				User:         "thelounge", ConfigDir: "${DATA_DIR}/thelounge",
				InstallScript: `#!/bin/bash
set -e
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
apt-get install -y nodejs
npm install -g thelounge`,
				ServiceUnit: `[Unit]
Description=The Lounge IRC
After=network.target

[Service]
Type=simple
User=thelounge
Environment=THELOUNGE_HOME=/opt/velour/thelounge
ExecStart=/usr/bin/thelounge start
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
			}},
		models.ServiceDefinition{
			ID: "znc", Name: "ZNC", Description: "Advanced IRC bouncer that stays connected and buffers messages while you're offline.", Icon: "znc", Category: "network", Image: "lscr.io/linuxserver/znc:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 6501, Container: 6501, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/znc/config", Container: "/config"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
			Native: &models.NativeConfig{
				Method:      "apt", ServiceName: "znc", Port: 6501,
				AptPackages: []string{"znc"},
				User:        "znc", ConfigDir: "${DATA_DIR}/znc",
			}},
	)
}
