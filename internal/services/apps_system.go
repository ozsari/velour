package services

import "github.com/ozsari/velour/internal/models"

func init() {
	Registry = append(Registry,
		// ── System / Infrastructure ──
		models.ServiceDefinition{ID: "uptimekuma", Name: "Uptime Kuma", Description: "Self-hosted monitoring tool with beautiful status pages and notifications.", Icon: "uptimekuma", Category: "system", Image: "louislam/uptime-kuma:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 3010, Container: 3001, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/uptimekuma/data", Container: "/app/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
			Native: &models.NativeConfig{
				Method: "script", ServiceName: "uptime-kuma", Port: 3001,
				Dependencies: []string{"git", "curl"},
				User: "uptimekuma", ConfigDir: "${DATA_DIR}/uptimekuma",
				InstallScript: `#!/bin/bash
set -e
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
apt-get install -y nodejs
git clone https://github.com/louislam/uptime-kuma.git /opt/uptime-kuma
cd /opt/uptime-kuma && npm run setup
chown -R uptimekuma:uptimekuma /opt/uptime-kuma`,
				ServiceUnit: `[Unit]
Description=Uptime Kuma
After=network.target

[Service]
Type=simple
User=uptimekuma
WorkingDirectory=/opt/uptime-kuma
Environment=DATA_DIR=/opt/velour/uptimekuma
Environment=UPTIME_KUMA_HOST=0.0.0.0
ExecStart=/usr/bin/node /opt/uptime-kuma/server/server.js
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
			}},
		models.ServiceDefinition{ID: "mariadb", Name: "MariaDB", Description: "Community-developed MySQL fork. Reliable, high performance database server.", Icon: "mariadb", Category: "system", Image: "lscr.io/linuxserver/mariadb:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 3306, Container: 3306, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/mariadb/config", Container: "/config"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000", "MYSQL_ROOT_PASSWORD": "changeme"},
			Native: &models.NativeConfig{
				Method: "apt", ServiceName: "mariadb", Port: 3306,
				AptPackages: []string{"mariadb-server"},
				User:        "mysql",
			}},
	)
}
