package services

import "github.com/ozsari/velour/internal/models"

var Registry = []models.ServiceDefinition{
	// ── Download Automation ──
	{ID: "autobrr", Name: "Autobrr", Description: "Modern automation tool for torrents and usenet. Monitors IRC, RSS feeds and more.", Icon: "autobrr", Category: "download", Image: "ghcr.io/autobrr/autobrr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 7474, Container: 7474, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/autobrr/config", Container: "/config"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "binary", ServiceName: "autobrr", Port: 7474,
			BinaryURL:  "https://github.com/autobrr/autobrr/releases/latest/download/autobrr_linux_${ARCH}.tar.gz",
			BinaryPath: "/usr/local/bin/autobrr",
			ConfigDir:  "${DATA_DIR}/autobrr", User: "autobrr",
			PostInstallCmds: []string{`sed -i 's/host = "127.0.0.1"/host = "0.0.0.0"/' /opt/velour/autobrr/config.toml`},
			ServiceUnit: `[Unit]
Description=autobrr
After=network.target

[Service]
Type=simple
User=autobrr
ExecStart=/usr/local/bin/autobrr --config /opt/velour/autobrr
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},

	// ── Download Clients ──
	{ID: "deluge", Name: "Deluge", Description: "Lightweight, free, cross-platform BitTorrent client with a web interface.", Icon: "deluge", Category: "client", Image: "lscr.io/linuxserver/deluge:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8112, Container: 8112, Protocol: "tcp"}, {Host: 6881, Container: 6881, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/deluge/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "deluged", Port: 8112,
			AptPackages: []string{"deluged", "deluge-web"},
			ConfigDir: "${DATA_DIR}/deluge", User: "deluge",
			PostInstallCmds: []string{
				// Auth file for deluged daemon RPC
				`echo "${VELOUR_USER}:${VELOUR_PASS}:10" > /opt/velour/deluge/auth && chown deluge:deluge /opt/velour/deluge/auth`,
				// Create deluge-web systemd service
				`cat > /etc/systemd/system/deluge-web.service << 'UNIT'
[Unit]
Description=Deluge Web Interface
After=network-online.target deluged.service

[Service]
Type=simple
User=deluge
ExecStart=/usr/bin/deluge-web -d -c /opt/velour/deluge
Restart=on-failure

[Install]
WantedBy=multi-user.target
UNIT
systemctl daemon-reload && systemctl enable deluge-web`,
				// Start deluged+deluge-web briefly so deluge-web generates default web.conf
				`systemctl start deluged && sleep 2 && systemctl start deluge-web && sleep 5 && systemctl stop deluge-web && systemctl stop deluged`,
				// Set web UI password using Deluge's exact hashing method (all in python3, no bash variable issues)
				`python3 << 'PYEOF'
import hashlib, os, subprocess
password = """${VELOUR_PASS}"""
salt = hashlib.sha1(os.urandom(32)).hexdigest()
s = hashlib.sha1(salt.encode("utf-8"))
s.update(password.encode("utf-8"))
pwd_hash = s.hexdigest()
conf = "/opt/velour/deluge/web.conf"
subprocess.run(["sed", "-i", f's|"pwd_salt": "[^"]*"|"pwd_salt": "{salt}"|', conf])
subprocess.run(["sed", "-i", f's|"pwd_sha1": "[^"]*"|"pwd_sha1": "{pwd_hash}"|', conf])
subprocess.run(["sed", "-i", 's|"first_login": true|"first_login": false|', conf])
PYEOF`,
				// Restart deluge-web with patched config (deluged will be started by main flow)
				`(sleep 2 && systemctl start deluge-web) &`,
			},
			ServiceUnit: `[Unit]
Description=Deluge Bittorrent Client Daemon
After=network-online.target

[Service]
Type=simple
User=deluge
ExecStart=/usr/bin/deluged -d -c /opt/velour/deluge
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "flood", Name: "Flood", Description: "Modern web UI for rTorrent, qBittorrent, and Transmission with a clean interface.", Icon: "flood", Category: "client", Image: "jesec/flood:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 3001, Container: 3000, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/flood/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "flood", Port: 3001,
			Dependencies: []string{"curl"},
			User: "flood", ConfigDir: "${DATA_DIR}/flood",
			InstallScript: `#!/bin/bash
set -e
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
apt-get install -y nodejs
npm install -g flood`,
			ServiceUnit: `[Unit]
Description=Flood
After=network.target

[Service]
Type=simple
User=flood
ExecStart=/usr/bin/flood --rundir /opt/velour/flood --port 3001 --host 0.0.0.0
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "qbittorrent", Name: "qBittorrent", Description: "Free, open-source BitTorrent client with a feature-rich web interface.", Icon: "qbittorrent", Category: "client", Image: "lscr.io/linuxserver/qbittorrent:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8085, Container: 8080, Protocol: "tcp"}, {Host: 6882, Container: 6881, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/qbittorrent/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000", "WEBUI_PORT": "8080"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "qbittorrent-nox", Port: 8085,
			AptPackages: []string{"qbittorrent-nox"},
			User: "qbittorrent",
			PostInstallCmds: []string{
				// Start qBittorrent, login with default creds, set new creds via API
				`systemctl start qbittorrent-nox && sleep 3 && SID=$(curl -s -c - 'http://localhost:8085/api/v2/auth/login' -d 'username=admin&password=adminadmin' | grep -oP 'SID\s+\K\S+') && curl -s -b "SID=$SID" 'http://localhost:8085/api/v2/app/setPreferences' -d 'json={"web_ui_username":"${VELOUR_USER}","web_ui_password":"${VELOUR_PASS}"}'; true`,
			},
			ServiceUnit: `[Unit]
Description=qBittorrent-nox
After=network.target

[Service]
Type=simple
User=qbittorrent
ExecStart=/usr/bin/qbittorrent-nox --webui-port=8085
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "qui", Name: "qui", Description: "Fast, modern web interface for qBittorrent with multi-instance support and automations.", Icon: "qui", Category: "client", Image: "ghcr.io/autobrr/qui:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 7476, Container: 7476, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/qui/config", Container: "/config"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "binary", ServiceName: "qui", Port: 7476,
			BinaryURL:  "https://github.com/autobrr/qui/releases/latest/download/qui_linux_${ARCH}.tar.gz",
			BinaryPath: "/usr/local/bin/qui",
			User: "qui", ConfigDir: "${DATA_DIR}/qui",
			PostInstallCmds: []string{`sed -i 's/host = "127.0.0.1"/host = "0.0.0.0"/' /opt/velour/qui/config.toml 2>/dev/null; true`},
			ServiceUnit: `[Unit]
Description=qui
After=network.target

[Service]
Type=simple
User=qui
ExecStart=/usr/local/bin/qui --config /opt/velour/qui
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "rutorrent", Name: "ruTorrent", Description: "rTorrent + ruTorrent web UI + autodl-irssi bundle. Complete torrent solution with IRC automation.", Icon: "rutorrent", Category: "client", Image: "crazymax/rtorrent-rutorrent:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8087, Container: 8080, Protocol: "tcp"}, {Host: 50000, Container: 50000, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/rutorrent/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "rtorrent", Port: 8087,
			Dependencies: []string{"rtorrent", "nginx", "php-fpm", "php-cli", "mediainfo", "unrar-free", "curl"},
			User: "rtorrent", ConfigDir: "${DATA_DIR}/rutorrent",
			InstallScript: `#!/bin/bash
set -e
# Install rTorrent
apt-get install -y rtorrent
# Install ruTorrent
mkdir -p /var/www/rutorrent
git clone https://github.com/Novik/ruTorrent.git /var/www/rutorrent
chown -R www-data:www-data /var/www/rutorrent
# Install autodl-irssi
apt-get install -y irssi
mkdir -p /home/rtorrent/.irssi/scripts/autorun
cd /home/rtorrent/.irssi/scripts
curl -sL https://git.io/vlcND | grep -Po '(?<="browser_download_url": ")(.*-v[\d.]+.zip)' | xargs curl -sL -o autodl-irssi.zip
unzip -o autodl-irssi.zip
cp autodl-irssi.pl autorun/
chown -R rtorrent:rtorrent /home/rtorrent/.irssi`,
			ServiceUnit: `[Unit]
Description=rTorrent
After=network.target

[Service]
Type=simple
User=rtorrent
ExecStart=/usr/bin/rtorrent
WorkingDirectory=/home/rtorrent
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "transmission", Name: "Transmission", Description: "Fast, easy and free BitTorrent client with minimal resource usage.", Icon: "transmission", Category: "client", Image: "lscr.io/linuxserver/transmission:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 9091, Container: 9091, Protocol: "tcp"}, {Host: 51413, Container: 51413, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/transmission/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "transmission-daemon", Port: 9091,
			AptPackages: []string{"transmission-daemon"},
			User: "debian-transmission",
			PostInstallCmds: []string{
				`sed -i 's/"rpc-whitelist-enabled":.*/"rpc-whitelist-enabled": false,/' /etc/transmission-daemon/settings.json 2>/dev/null; true`,
				`sed -i 's/"rpc-authentication-required":.*/"rpc-authentication-required": true,/' /etc/transmission-daemon/settings.json 2>/dev/null; true`,
				`sed -i 's/"rpc-username":.*/"rpc-username": "${VELOUR_USER}",/' /etc/transmission-daemon/settings.json 2>/dev/null; true`,
				`sed -i 's/"rpc-password":.*/"rpc-password": "${VELOUR_PASS}",/' /etc/transmission-daemon/settings.json 2>/dev/null; true`,
			},
		}},
	{ID: "nzbget", Name: "NZBGet", Description: "Efficient usenet downloader written in C++ for maximum performance.", Icon: "nzbget", Category: "client", Image: "lscr.io/linuxserver/nzbget:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 6789, Container: 6789, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/nzbget/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "binary", ServiceName: "nzbget", Port: 6789,
			BinaryURL:  "https://github.com/nzbgetcom/nzbget/releases/latest/download/nzbget-linux.run",
			BinaryPath: "/opt/nzbget/nzbget",
			ConfigDir: "${DATA_DIR}/nzbget", User: "nzbget",
			InstallScript: `#!/bin/bash
set -e
curl -sL https://github.com/nzbgetcom/nzbget/releases/latest/download/nzbget-linux.run -o /tmp/nzbget.run
sh /tmp/nzbget.run --destdir /opt/nzbget
rm /tmp/nzbget.run
chown -R nzbget:nzbget /opt/nzbget`,
			ServiceUnit: `[Unit]
Description=NZBGet
After=network.target

[Service]
Type=forking
User=nzbget
ExecStart=/opt/nzbget/nzbget -D -c /opt/velour/nzbget/nzbget.conf
ExecStop=/opt/nzbget/nzbget -Q
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
			PostInstallCmds: []string{
				`sed -i 's/^ControlUsername=.*/ControlUsername=${VELOUR_USER}/' /opt/velour/nzbget/nzbget.conf 2>/dev/null; true`,
				`sed -i 's/^ControlPassword=.*/ControlPassword=${VELOUR_PASS}/' /opt/velour/nzbget/nzbget.conf 2>/dev/null; true`,
			},
		}},
	{ID: "sabnzbd", Name: "SABnzbd", Description: "Free, open-source usenet downloader with web-based interface.", Icon: "sabnzbd", Category: "client", Image: "lscr.io/linuxserver/sabnzbd:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8088, Container: 8080, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/sabnzbd/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "sabnzbdplus", Port: 8080,
			AptPackages: []string{"sabnzbdplus"},
			AptRepo: &models.AptRepo{
				KeyURL:   "https://ppa.launchpadcontent.net/jcfp/nobetas/ubuntu/dists/jammy/Release.gpg",
				RepoLine: "deb [signed-by=/usr/share/keyrings/velour-jcfp.gpg] https://ppa.launchpadcontent.net/jcfp/nobetas/ubuntu jammy main",
			},
			User: "sabnzbd", ConfigDir: "${DATA_DIR}/sabnzbd",
		}},
	{ID: "pyload", Name: "pyLoad", Description: "Free, open-source download manager for HTTP, FTP, and other protocols.", Icon: "pyload", Category: "client", Image: "lscr.io/linuxserver/pyload-ng:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8089, Container: 8000, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/pyload/config", Container: "/config"}, {Host: "${DATA_DIR}/data/downloads", Container: "/downloads"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "pyload", Port: 8089,
			Dependencies: []string{"python3", "python3-pip", "python3-venv"},
			User: "pyload", ConfigDir: "${DATA_DIR}/pyload",
			InstallScript: `#!/bin/bash
set -e
python3 -m venv /opt/pyload
/opt/pyload/bin/pip install pyload-ng[all]
chown -R pyload:pyload /opt/pyload`,
			ServiceUnit: `[Unit]
Description=pyLoad
After=network.target

[Service]
Type=simple
User=pyload
ExecStart=/opt/pyload/bin/pyload --storagedir /opt/velour/downloads --userdir /opt/velour/pyload
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "jdownloader2", Name: "JDownloader 2", Description: "Open-source download manager with web interface for direct downloads.", Icon: "jdownloader2", Category: "client", Image: "jlesage/jdownloader-2:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 5801, Container: 5800, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/jdownloader2/config", Container: "/config"}, {Host: "${DATA_DIR}/data/downloads", Container: "/output"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "jdownloader2", Port: 0,
			Dependencies: []string{"default-jre-headless", "curl"},
			User: "jdownloader", ConfigDir: "${DATA_DIR}/jdownloader2",
			InstallScript: `#!/bin/bash
set -e
mkdir -p /opt/jdownloader
curl -sL http://installer.jdownloader.org/JDownloader.jar -o /opt/jdownloader/JDownloader.jar
chown -R jdownloader:jdownloader /opt/jdownloader`,
			ServiceUnit: `[Unit]
Description=JDownloader 2
After=network.target

[Service]
Type=simple
User=jdownloader
ExecStart=/usr/bin/java -Djava.awt.headless=true -jar /opt/jdownloader/JDownloader.jar -norestart
WorkingDirectory=/opt/velour/jdownloader2
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},

	// ── Media Management ──
	{ID: "sonarr", Name: "Sonarr", Description: "Smart PVR for newsgroup and bittorrent users. Monitors RSS feeds for new episodes.", Icon: "sonarr", Category: "media", Image: "lscr.io/linuxserver/sonarr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8989, Container: 8989, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/sonarr/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "sonarr", Port: 8989,
			AptPackages: []string{"sonarr"},
			AptRepo: &models.AptRepo{
				KeyURL:   "https://apt.sonarr.tv/sonarr.asc",
				RepoLine: "deb [signed-by=/usr/share/keyrings/velour-sonarr.gpg] https://apt.sonarr.tv/debian buster main",
			},
			User: "sonarr", ConfigDir: "/var/lib/sonarr",
			PostInstallCmds: []string{
				`sed -i 's|<AuthenticationMethod>.*</AuthenticationMethod>|<AuthenticationMethod>Forms</AuthenticationMethod>|' /var/lib/sonarr/config.xml 2>/dev/null; true`,
				`sed -i 's|<AuthenticationRequired>.*</AuthenticationRequired>|<AuthenticationRequired>Enabled</AuthenticationRequired>|' /var/lib/sonarr/config.xml 2>/dev/null; true`,
			},
		}},
	{ID: "sonarr2", Name: "Sonarr (2nd)", Description: "Second instance of Sonarr for managing a separate TV library.", Icon: "sonarr", Category: "media", Image: "lscr.io/linuxserver/sonarr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8990, Container: 8989, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/sonarr2/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "sonarr2", Port: 8990,
			AptPackages: []string{"sonarr"},
			AptRepo: &models.AptRepo{
				KeyURL:   "https://apt.sonarr.tv/sonarr.asc",
				RepoLine: "deb [signed-by=/usr/share/keyrings/velour-sonarr.gpg] https://apt.sonarr.tv/debian buster main",
			},
			User: "sonarr", ConfigDir: "/var/lib/sonarr2",
			PostInstallCmds: []string{
				`sed -i 's|<AuthenticationMethod>.*</AuthenticationMethod>|<AuthenticationMethod>Forms</AuthenticationMethod>|' /var/lib/sonarr2/config.xml 2>/dev/null; true`,
				`sed -i 's|<AuthenticationRequired>.*</AuthenticationRequired>|<AuthenticationRequired>Enabled</AuthenticationRequired>|' /var/lib/sonarr2/config.xml 2>/dev/null; true`,
			},
			ServiceUnit: `[Unit]
Description=Sonarr (2nd Instance)
After=network.target

[Service]
Type=simple
User=sonarr
ExecStart=/opt/Sonarr/Sonarr -data=/var/lib/sonarr2 -port=8990 -nobrowser
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "radarr", Name: "Radarr", Description: "Movie collection manager. Automatically searches, downloads and manages movies.", Icon: "radarr", Category: "media", Image: "lscr.io/linuxserver/radarr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 7878, Container: 7878, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/radarr/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "radarr", Port: 7878,
			AptPackages: []string{"radarr"},
			AptRepo: &models.AptRepo{
				KeyURL:   "https://apt.sonarr.tv/sonarr.asc",
				RepoLine: "deb [signed-by=/usr/share/keyrings/velour-sonarr.gpg] https://apt.sonarr.tv/debian buster main",
			},
			User: "radarr", ConfigDir: "/var/lib/radarr",
			PostInstallCmds: []string{
				`sed -i 's|<AuthenticationMethod>.*</AuthenticationMethod>|<AuthenticationMethod>Forms</AuthenticationMethod>|' /var/lib/radarr/config.xml 2>/dev/null; true`,
				`sed -i 's|<AuthenticationRequired>.*</AuthenticationRequired>|<AuthenticationRequired>Enabled</AuthenticationRequired>|' /var/lib/radarr/config.xml 2>/dev/null; true`,
			},
		}},
	{ID: "radarr2", Name: "Radarr (2nd)", Description: "Second instance of Radarr for managing a separate movie library.", Icon: "radarr", Category: "media", Image: "lscr.io/linuxserver/radarr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 7879, Container: 7878, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/radarr2/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "radarr2", Port: 7879,
			AptPackages: []string{"radarr"},
			AptRepo: &models.AptRepo{
				KeyURL:   "https://apt.sonarr.tv/sonarr.asc",
				RepoLine: "deb [signed-by=/usr/share/keyrings/velour-sonarr.gpg] https://apt.sonarr.tv/debian buster main",
			},
			User: "radarr", ConfigDir: "/var/lib/radarr2",
			PostInstallCmds: []string{
				`sed -i 's|<AuthenticationMethod>.*</AuthenticationMethod>|<AuthenticationMethod>Forms</AuthenticationMethod>|' /var/lib/radarr2/config.xml 2>/dev/null; true`,
				`sed -i 's|<AuthenticationRequired>.*</AuthenticationRequired>|<AuthenticationRequired>Enabled</AuthenticationRequired>|' /var/lib/radarr2/config.xml 2>/dev/null; true`,
			},
			ServiceUnit: `[Unit]
Description=Radarr (2nd Instance)
After=network.target

[Service]
Type=simple
User=radarr
ExecStart=/opt/Radarr/Radarr -data=/var/lib/radarr2 -port=7879 -nobrowser
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "bazarr", Name: "Bazarr", Description: "Companion to Sonarr and Radarr for managing and downloading subtitles.", Icon: "bazarr", Category: "media", Image: "lscr.io/linuxserver/bazarr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 6767, Container: 6767, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/bazarr/config", Container: "/config"}, {Host: "${DATA_DIR}/data/media", Container: "/data/media"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "bazarr", Port: 6767,
			Dependencies: []string{"python3", "python3-pip", "python3-venv"},
			User: "bazarr", ConfigDir: "${DATA_DIR}/bazarr",
			InstallScript: `#!/bin/bash
set -e
mkdir -p /opt/bazarr
curl -sL https://github.com/morpheus65535/bazarr/releases/latest/download/bazarr.zip -o /tmp/bazarr.zip
unzip -o /tmp/bazarr.zip -d /opt/bazarr
rm /tmp/bazarr.zip
cd /opt/bazarr && python3 -m venv venv && venv/bin/pip install -r requirements.txt
chown -R bazarr:bazarr /opt/bazarr`,
			ServiceUnit: `[Unit]
Description=Bazarr
After=network.target

[Service]
Type=simple
User=bazarr
ExecStart=/opt/bazarr/venv/bin/python3 /opt/bazarr/bazarr.py --config /opt/velour/bazarr --address 0.0.0.0
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "lidarr", Name: "Lidarr", Description: "Music collection manager. Monitors RSS feeds for new albums and manages your library.", Icon: "lidarr", Category: "media", Image: "lscr.io/linuxserver/lidarr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8686, Container: 8686, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/lidarr/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "lidarr", Port: 8686,
			AptPackages: []string{"lidarr"},
			AptRepo: &models.AptRepo{
				KeyURL:   "https://apt.sonarr.tv/sonarr.asc",
				RepoLine: "deb [signed-by=/usr/share/keyrings/velour-sonarr.gpg] https://apt.sonarr.tv/debian buster main",
			},
			User: "lidarr", ConfigDir: "/var/lib/lidarr",
			PostInstallCmds: []string{
				`sed -i 's|<AuthenticationMethod>.*</AuthenticationMethod>|<AuthenticationMethod>Forms</AuthenticationMethod>|' /var/lib/lidarr/config.xml 2>/dev/null; true`,
				`sed -i 's|<AuthenticationRequired>.*</AuthenticationRequired>|<AuthenticationRequired>Enabled</AuthenticationRequired>|' /var/lib/lidarr/config.xml 2>/dev/null; true`,
			},
		}},
	{ID: "medusa", Name: "Medusa", Description: "Automatic video library manager for TV shows with multi-source support.", Icon: "medusa", Category: "media", Image: "lscr.io/linuxserver/medusa:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8081, Container: 8081, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/medusa/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "medusa", Port: 8081,
			Dependencies: []string{"python3", "python3-pip", "git"},
			User: "medusa", ConfigDir: "${DATA_DIR}/medusa",
			InstallScript: `#!/bin/bash
set -e
git clone https://github.com/pymedusa/Medusa.git /opt/medusa
chown -R medusa:medusa /opt/medusa`,
			ServiceUnit: `[Unit]
Description=Medusa
After=network.target

[Service]
Type=simple
User=medusa
ExecStart=/usr/bin/python3 /opt/medusa/SickBeard.py --datadir /opt/velour/medusa --host 0.0.0.0
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "sickchill", Name: "SickChill", Description: "Automatic video library manager for TV shows from various sources.", Icon: "sickchill", Category: "media", Image: "lscr.io/linuxserver/sickchill:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8083, Container: 8081, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/sickchill/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "sickchill", Port: 8083,
			Dependencies: []string{"python3", "python3-pip", "git"},
			User: "sickchill", ConfigDir: "${DATA_DIR}/sickchill",
			InstallScript: `#!/bin/bash
set -e
pip3 install sickchill
mkdir -p /opt/velour/sickchill
chown -R sickchill:sickchill /opt/velour/sickchill`,
			ServiceUnit: `[Unit]
Description=SickChill
After=network.target

[Service]
Type=simple
User=sickchill
ExecStart=/usr/local/bin/SickChill --datadir /opt/velour/sickchill --port 8083
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "sickgear", Name: "SickGear", Description: "SickBeard fork with improved stability, performance and features.", Icon: "sickgear", Category: "media", Image: "lscr.io/linuxserver/sickgear:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8082, Container: 8081, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/sickgear/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "sickgear", Port: 8082,
			Dependencies: []string{"python3", "python3-pip", "git"},
			User: "sickgear", ConfigDir: "${DATA_DIR}/sickgear",
			InstallScript: `#!/bin/bash
set -e
git clone https://github.com/SickGear/SickGear.git /opt/sickgear
chown -R sickgear:sickgear /opt/sickgear`,
			ServiceUnit: `[Unit]
Description=SickGear
After=network.target

[Service]
Type=simple
User=sickgear
ExecStart=/usr/bin/python3 /opt/sickgear/sickgear.py --datadir /opt/velour/sickgear --port 8082 --host 0.0.0.0
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "mylar3", Name: "Mylar3", Description: "Automated comic book downloader for usenet and torrent users.", Icon: "mylar3", Category: "media", Image: "lscr.io/linuxserver/mylar3:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8090, Container: 8090, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/mylar3/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "mylar3", Port: 8090,
			Dependencies: []string{"python3", "python3-pip", "git"},
			User: "mylar3", ConfigDir: "${DATA_DIR}/mylar3",
			InstallScript: `#!/bin/bash
set -e
git clone https://github.com/mylar3/mylar3.git /opt/mylar3
cd /opt/mylar3 && pip3 install -r requirements.txt
chown -R mylar3:mylar3 /opt/mylar3`,
			ServiceUnit: `[Unit]
Description=Mylar3
After=network.target

[Service]
Type=simple
User=mylar3
ExecStart=/usr/bin/python3 /opt/mylar3/Mylar.py --datadir /opt/velour/mylar3 --host 0.0.0.0
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "filebot", Name: "FileBot", Description: "Ultimate tool for organizing and renaming movies, TV shows, anime and music.", Icon: "filebot", Category: "media", Image: "jlesage/filebot:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 5800, Container: 5800, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/filebot/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "filebot", Port: 0,
			Dependencies: []string{"default-jre-headless", "curl", "libmediainfo0v5"},
			User: "filebot", ConfigDir: "${DATA_DIR}/filebot",
			InstallScript: `#!/bin/bash
set -e
curl -sL https://get.filebot.net/filebot/FileBot_5.1.6/FileBot_5.1.6_amd64.deb -o /tmp/filebot.deb
dpkg -i /tmp/filebot.deb || apt-get install -f -y
rm /tmp/filebot.deb`,
		}},
	{ID: "lazylibrarian", Name: "LazyLibrarian", Description: "Automated book downloader for eBooks and audiobooks from usenet and torrents.", Icon: "lazylibrarian", Category: "media", Image: "lscr.io/linuxserver/lazylibrarian:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 5299, Container: 5299, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/lazylibrarian/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "lazylibrarian", Port: 5299,
			Dependencies: []string{"python3", "python3-pip", "git"},
			User: "lazylibrarian", ConfigDir: "${DATA_DIR}/lazylibrarian",
			InstallScript: `#!/bin/bash
set -e
git clone https://gitlab.com/LazyLibrarian/LazyLibrarian.git /opt/lazylibrarian
chown -R lazylibrarian:lazylibrarian /opt/lazylibrarian`,
			ServiceUnit: `[Unit]
Description=LazyLibrarian
After=network.target

[Service]
Type=simple
User=lazylibrarian
ExecStart=/usr/bin/python3 /opt/lazylibrarian/LazyLibrarian.py --datadir /opt/velour/lazylibrarian --host 0.0.0.0
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "maintainerr", Name: "Maintainerr", Description: "Automated media maintenance for Plex. Clean up old or unwatched content.", Icon: "maintainerr", Category: "media", Image: "ghcr.io/jorenn92/maintainerr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 6246, Container: 6246, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/maintainerr/config", Container: "/opt/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "maintainerr", Port: 6246,
			Dependencies: []string{"git", "curl"},
			User: "maintainerr", ConfigDir: "${DATA_DIR}/maintainerr",
			InstallScript: `#!/bin/bash
set -e
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
apt-get install -y nodejs
git clone https://github.com/jorenn92/Maintainerr.git /opt/maintainerr
cd /opt/maintainerr && npm ci && npm run build
chown -R maintainerr:maintainerr /opt/maintainerr`,
			ServiceUnit: `[Unit]
Description=Maintainerr
After=network.target

[Service]
Type=simple
User=maintainerr
WorkingDirectory=/opt/maintainerr
Environment=DATA_DIR=/opt/velour/maintainerr
ExecStart=/usr/bin/node /opt/maintainerr/dist/main.js
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "unpackerr", Name: "Unpackerr", Description: "Extracts downloaded archives for Sonarr, Radarr, Lidarr and Readarr.", Icon: "unpackerr", Category: "media", Image: "golift/unpackerr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/unpackerr/config", Container: "/config"}, {Host: "${DATA_DIR}/data/downloads", Container: "/data/downloads"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "unpackerr", Port: 0,
			AptPackages: []string{"unpackerr"},
			AptRepo: &models.AptRepo{
				KeyURL:   "https://packagecloud.io/golift/pkgs/gpgkey",
				RepoLine: "deb [signed-by=/usr/share/keyrings/velour-golift.gpg] https://packagecloud.io/golift/pkgs/debian/ any main",
			},
			User: "unpackerr", ConfigDir: "${DATA_DIR}/unpackerr",
		}},

	// ── Media Servers ──
	{ID: "plex", Name: "Plex", Description: "Stream movies, TV, music and more to any device from your own media library.", Icon: "plex", Category: "server", Image: "lscr.io/linuxserver/plex:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 32400, Container: 32400, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/plex/config", Container: "/config"}, {Host: "${DATA_DIR}/data/media", Container: "/data/media"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "plexmediaserver", Port: 32400,
			AptPackages: []string{"plexmediaserver"},
			AptRepo: &models.AptRepo{
				KeyURL:   "https://downloads.plex.tv/plex-keys/PlexSign.key",
				RepoLine: "deb https://downloads.plex.tv/repo/deb public main",
			},
			User: "plex",
		}},
	{ID: "jellyfin", Name: "Jellyfin", Description: "Free, open-source media server. Stream to any device from your own server.", Icon: "jellyfin", Category: "server", Image: "lscr.io/linuxserver/jellyfin:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8097, Container: 8096, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/jellyfin/config", Container: "/config"}, {Host: "${DATA_DIR}/data/media", Container: "/data/media"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "jellyfin", Port: 8096,
			AptPackages: []string{"jellyfin"},
			AptRepo: &models.AptRepo{
				KeyURL:   "https://repo.jellyfin.org/jellyfin_team.gpg.key",
				RepoLine: "deb [signed-by=/usr/share/keyrings/velour-jellyfin_team.gpg] https://repo.jellyfin.org/debian bookworm main",
			},
			User: "jellyfin",
		}},
	{ID: "emby", Name: "Emby", Description: "Media server that organizes and streams your personal media collection.", Icon: "emby", Category: "server", Image: "lscr.io/linuxserver/emby:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8096, Container: 8096, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/emby/config", Container: "/config"}, {Host: "${DATA_DIR}/data/media", Container: "/data/media"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "emby-server", Port: 8096,
			AptPackages: []string{"emby-server"},
			AptRepo: &models.AptRepo{
				KeyURL:   "https://mb3admin.com/startupapikey/keyfile.asc",
				RepoLine: "deb [signed-by=/usr/share/keyrings/velour-emby.gpg] https://packages.emby.media/deb stable main",
			},
			User: "emby",
		}},
	{ID: "airsonic", Name: "Airsonic", Description: "Free, web-based media streamer for your music collection.", Icon: "airsonic", Category: "server", Image: "lscr.io/linuxserver/airsonic-advanced:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 4040, Container: 4040, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/airsonic/config", Container: "/config"}, {Host: "${DATA_DIR}/data/media/music", Container: "/music"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "airsonic", Port: 4040,
			Dependencies: []string{"default-jre-headless", "curl"},
			User: "airsonic", ConfigDir: "${DATA_DIR}/airsonic",
			InstallScript: `#!/bin/bash
set -e
mkdir -p /opt/airsonic
curl -sL https://github.com/airsonic-advanced/airsonic-advanced/releases/latest/download/airsonic.war -o /opt/airsonic/airsonic.war
chown -R airsonic:airsonic /opt/airsonic`,
			ServiceUnit: `[Unit]
Description=Airsonic Advanced
After=network.target

[Service]
Type=simple
User=airsonic
ExecStart=/usr/bin/java -jar /opt/airsonic/airsonic.war --port=4040 --airsonic.home=/opt/velour/airsonic
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "navidrome", Name: "Navidrome", Description: "Modern, open-source music server and streamer compatible with Subsonic/Airsonic.", Icon: "navidrome", Category: "server", Image: "deluan/navidrome:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 4533, Container: 4533, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/navidrome/config", Container: "/navidrome"}, {Host: "${DATA_DIR}/data/media/music", Container: "/music"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "binary", ServiceName: "navidrome", Port: 4533,
			BinaryURL:  "https://github.com/navidrome/navidrome/releases/latest/download/navidrome_0.52.5_linux_${ARCH}.tar.gz",
			BinaryPath: "/usr/local/bin/navidrome",
			User: "navidrome", ConfigDir: "${DATA_DIR}/navidrome",
			ServiceUnit: `[Unit]
Description=Navidrome
After=network.target

[Service]
Type=simple
User=navidrome
ExecStart=/usr/local/bin/navidrome --datafolder /opt/velour/navidrome --musicfolder /opt/velour/music --address 0.0.0.0
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "calibreweb", Name: "Calibre-Web", Description: "Web app for browsing, reading and downloading eBooks from a Calibre database.", Icon: "calibreweb", Category: "server", Image: "lscr.io/linuxserver/calibre-web:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8084, Container: 8083, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/calibreweb/config", Container: "/config"}, {Host: "${DATA_DIR}/data/media/books", Container: "/books"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "calibre-web", Port: 8083,
			Dependencies: []string{"python3", "python3-pip", "python3-venv"},
			User: "calibreweb", ConfigDir: "${DATA_DIR}/calibreweb",
			InstallScript: `#!/bin/bash
set -e
mkdir -p /opt/calibre-web
python3 -m venv /opt/calibre-web/venv
/opt/calibre-web/venv/bin/pip install calibreweb
chown -R calibreweb:calibreweb /opt/calibre-web`,
			ServiceUnit: `[Unit]
Description=Calibre-Web
After=network.target

[Service]
Type=simple
User=calibreweb
ExecStart=/opt/calibre-web/venv/bin/cps -p 8083 -s /opt/velour/calibreweb
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "tautulli", Name: "Tautulli", Description: "Monitoring and tracking tool for Plex Media Server usage and statistics.", Icon: "tautulli", Category: "server", Image: "lscr.io/linuxserver/tautulli:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8181, Container: 8181, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/tautulli/config", Container: "/config"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "tautulli", Port: 8181,
			Dependencies: []string{"python3", "python3-pip", "git"},
			User: "tautulli", ConfigDir: "${DATA_DIR}/tautulli",
			InstallScript: `#!/bin/bash
set -e
git clone https://github.com/Tautulli/Tautulli.git /opt/tautulli
chown -R tautulli:tautulli /opt/tautulli`,
			ServiceUnit: `[Unit]
Description=Tautulli
After=network.target

[Service]
Type=simple
User=tautulli
ExecStart=/usr/bin/python3 /opt/tautulli/Tautulli.py --datadir /opt/velour/tautulli --host 0.0.0.0
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "audiobookshelf", Name: "Audiobookshelf", Description: "Self-hosted audiobook and podcast server with web-based player.", Icon: "audiobookshelf", Category: "server", Image: "ghcr.io/advplyr/audiobookshelf:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 13378, Container: 80, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/audiobookshelf/config", Container: "/config"}, {Host: "${DATA_DIR}/audiobookshelf/metadata", Container: "/metadata"}, {Host: "${DATA_DIR}/data/media/audiobooks", Container: "/audiobooks"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "audiobookshelf", Port: 13378,
			AptPackages: []string{"audiobookshelf"},
			AptRepo: &models.AptRepo{
				KeyURL:   "https://advplyr.github.io/audiobookshelf-ppa/KEY.gpg",
				RepoLine: "deb [signed-by=/usr/share/keyrings/velour-audiobookshelf.gpg] https://advplyr.github.io/audiobookshelf-ppa ./",
			},
		}},
	{ID: "kavita", Name: "Kavita", Description: "Fast, feature-rich manga, comic, and book reader server.", Icon: "kavita", Category: "server", Image: "jvmilazz0/kavita:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 5000, Container: 5000, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/kavita/config", Container: "/kavita/config"}, {Host: "${DATA_DIR}/data/media/books", Container: "/books"}, {Host: "${DATA_DIR}/data/media/comics", Container: "/comics"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "kavita", Port: 5000,
			User: "kavita", ConfigDir: "${DATA_DIR}/kavita",
			InstallScript: `#!/bin/bash
set -e
mkdir -p /opt/kavita
curl -sL https://github.com/Kareadita/Kavita/releases/latest/download/kavita-linux-x64.tar.gz | tar xz -C /opt/kavita
chmod +x /opt/kavita/Kavita
chown -R kavita:kavita /opt/kavita`,
			ServiceUnit: `[Unit]
Description=Kavita
After=network.target

[Service]
Type=simple
User=kavita
WorkingDirectory=/opt/kavita
ExecStart=/opt/kavita/Kavita
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "ubooquity", Name: "Ubooquity", Description: "Home server for comics and eBooks with OPDS feed and web reader.", Icon: "ubooquity", Category: "server", Image: "lscr.io/linuxserver/ubooquity:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 2202, Container: 2202, Protocol: "tcp"}, {Host: 2203, Container: 2203, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/ubooquity/config", Container: "/config"}, {Host: "${DATA_DIR}/data/media/books", Container: "/books"}, {Host: "${DATA_DIR}/data/media/comics", Container: "/comics"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "ubooquity", Port: 2202,
			Dependencies: []string{"default-jre-headless", "curl"},
			User: "ubooquity", ConfigDir: "${DATA_DIR}/ubooquity",
			InstallScript: `#!/bin/bash
set -e
mkdir -p /opt/ubooquity
curl -sL https://vaemendis.net/ubooquity/downloads/Ubooquity-2.1.2.zip -o /tmp/ubooquity.zip
unzip -o /tmp/ubooquity.zip -d /opt/ubooquity
rm /tmp/ubooquity.zip
chown -R ubooquity:ubooquity /opt/ubooquity`,
			ServiceUnit: `[Unit]
Description=Ubooquity
After=network.target

[Service]
Type=simple
User=ubooquity
WorkingDirectory=/opt/velour/ubooquity
ExecStart=/usr/bin/java -jar /opt/ubooquity/Ubooquity.jar --headless --port 2202 --adminport 2203
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},

	// ── Indexers ──
	{ID: "prowlarr", Name: "Prowlarr", Description: "Indexer manager that integrates with Sonarr, Radarr, Lidarr and more.", Icon: "prowlarr", Category: "indexer", Image: "lscr.io/linuxserver/prowlarr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 9696, Container: 9696, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/prowlarr/config", Container: "/config"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "prowlarr", Port: 9696,
			AptPackages: []string{"prowlarr"},
			AptRepo: &models.AptRepo{
				KeyURL:   "https://apt.sonarr.tv/sonarr.asc",
				RepoLine: "deb [signed-by=/usr/share/keyrings/velour-sonarr.gpg] https://apt.sonarr.tv/debian buster main",
			},
			User: "prowlarr", ConfigDir: "/var/lib/prowlarr",
			PostInstallCmds: []string{
				`sed -i 's|<AuthenticationMethod>.*</AuthenticationMethod>|<AuthenticationMethod>Forms</AuthenticationMethod>|' /var/lib/prowlarr/config.xml 2>/dev/null; true`,
				`sed -i 's|<AuthenticationRequired>.*</AuthenticationRequired>|<AuthenticationRequired>Enabled</AuthenticationRequired>|' /var/lib/prowlarr/config.xml 2>/dev/null; true`,
			},
		}},
	{ID: "jackett", Name: "Jackett", Description: "Proxy server that translates queries from Sonarr/Radarr into tracker-site-specific queries.", Icon: "jackett", Category: "indexer", Image: "lscr.io/linuxserver/jackett:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 9117, Container: 9117, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/jackett/config", Container: "/config"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "binary", ServiceName: "jackett", Port: 9117,
			BinaryURL:  "https://github.com/Jackett/Jackett/releases/latest/download/Jackett.Binaries.LinuxAMDx64.tar.gz",
			BinaryPath: "/opt/jackett/jackett",
			User: "jackett", ConfigDir: "${DATA_DIR}/jackett",
			ServiceUnit: `[Unit]
Description=Jackett
After=network.target

[Service]
Type=simple
User=jackett
ExecStart=/opt/jackett/jackett --DataFolder=/opt/velour/jackett
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "nzbhydra2", Name: "NZBHydra2", Description: "Meta search for usenet indexers. Provides unified access to multiple newznab indexers.", Icon: "nzbhydra2", Category: "indexer", Image: "lscr.io/linuxserver/nzbhydra2:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 5076, Container: 5076, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/nzbhydra2/config", Container: "/config"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "nzbhydra2", Port: 5076,
			Dependencies: []string{"default-jre-headless", "unzip", "curl"},
			User: "nzbhydra2", ConfigDir: "${DATA_DIR}/nzbhydra2",
			InstallScript: `#!/bin/bash
set -e
mkdir -p /opt/nzbhydra2
curl -sL https://github.com/theotherp/nzbhydra2/releases/latest/download/nzbhydra2-linux-amd64-release.zip -o /tmp/nzbhydra2.zip
unzip -o /tmp/nzbhydra2.zip -d /opt/nzbhydra2
rm /tmp/nzbhydra2.zip
chmod +x /opt/nzbhydra2/nzbhydra2
chown -R nzbhydra2:nzbhydra2 /opt/nzbhydra2`,
			ServiceUnit: `[Unit]
Description=NZBHydra2
After=network.target

[Service]
Type=simple
User=nzbhydra2
ExecStart=/opt/nzbhydra2/nzbhydra2 --datafolder /opt/velour/nzbhydra2
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "flaresolverr", Name: "FlareSolverr", Description: "Proxy server to bypass Cloudflare and DDoS-GUARD protection for scrapers.", Icon: "flaresolverr", Category: "indexer", Image: "ghcr.io/flaresolverr/flaresolverr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8191, Container: 8191, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "flaresolverr", Port: 8191,
			Dependencies: []string{"git", "curl", "chromium-browser"},
			User: "flaresolverr",
			InstallScript: `#!/bin/bash
set -e
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
apt-get install -y nodejs chromium-browser
git clone https://github.com/FlareSolverr/FlareSolverr.git /opt/flaresolverr
cd /opt/flaresolverr && npm ci
chown -R flaresolverr:flaresolverr /opt/flaresolverr`,
			ServiceUnit: `[Unit]
Description=FlareSolverr
After=network.target

[Service]
Type=simple
User=flaresolverr
WorkingDirectory=/opt/flaresolverr
Environment=LOG_LEVEL=info
ExecStart=/usr/bin/node /opt/flaresolverr/index.js
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},

	// ── Request Management ──
	{ID: "ombi", Name: "Ombi", Description: "Self-hosted web app for Plex/Emby users to request content. Integrates with Sonarr/Radarr.", Icon: "ombi", Category: "request", Image: "lscr.io/linuxserver/ombi:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 3579, Container: 3579, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/ombi/config", Container: "/config"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "ombi", Port: 3579,
			User: "ombi", ConfigDir: "${DATA_DIR}/ombi",
			InstallScript: `#!/bin/bash
set -e
mkdir -p /opt/ombi
curl -sL https://github.com/Ombi-app/Ombi/releases/latest/download/linux-x64.tar.gz | tar xz -C /opt/ombi
chown -R ombi:ombi /opt/ombi`,
			ServiceUnit: `[Unit]
Description=Ombi
After=network.target

[Service]
Type=simple
User=ombi
ExecStart=/opt/ombi/Ombi --storage /opt/velour/ombi
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "overseerr", Name: "Overseerr", Description: "Modern media request and discovery tool for Plex, Sonarr and Radarr.", Icon: "overseerr", Category: "request", Image: "lscr.io/linuxserver/overseerr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 5055, Container: 5055, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/overseerr/config", Container: "/config"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "overseerr", Port: 5055,
			Dependencies: []string{"git", "curl"},
			User: "overseerr", ConfigDir: "${DATA_DIR}/overseerr",
			InstallScript: `#!/bin/bash
set -e
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
apt-get install -y nodejs
git clone https://github.com/sct/overseerr.git /opt/overseerr
cd /opt/overseerr && npm ci && npx tsc && npm run build
chown -R overseerr:overseerr /opt/overseerr`,
			ServiceUnit: `[Unit]
Description=Overseerr
After=network.target

[Service]
Type=simple
User=overseerr
WorkingDirectory=/opt/overseerr
Environment=CONFIG_DIRECTORY=/opt/velour/overseerr
ExecStart=/usr/bin/node /opt/overseerr/dist/index.js
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
	{ID: "doplarr", Name: "Doplarr", Description: "Discord bot for requesting media through Overseerr or Ombi via Discord.", Icon: "doplarr", Category: "request", Image: "lscr.io/linuxserver/doplarr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/doplarr/config", Container: "/config"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "doplarr", Port: 0,
			Dependencies: []string{"git", "curl"},
			User: "doplarr", ConfigDir: "${DATA_DIR}/doplarr",
			InstallScript: `#!/bin/bash
set -e
curl -fsSL https://deb.nodesource.com/setup_20.x | bash -
apt-get install -y nodejs
git clone https://github.com/kiranshila/Doplarr.git /opt/doplarr
cd /opt/doplarr && npm ci && npm run build
chown -R doplarr:doplarr /opt/doplarr`,
			ServiceUnit: `[Unit]
Description=Doplarr
After=network.target

[Service]
Type=simple
User=doplarr
WorkingDirectory=/opt/doplarr
EnvironmentFile=/opt/velour/doplarr/doplarr.env
ExecStart=/usr/bin/node /opt/doplarr/dist/index.js
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},

	// ── Sync & Cloud ──
	{ID: "resilio", Name: "Resilio Sync", Description: "Fast, reliable file sync using peer-to-peer technology.", Icon: "resilio", Category: "sync", Image: "lscr.io/linuxserver/resilio-sync:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8888, Container: 8888, Protocol: "tcp"}, {Host: 55555, Container: 55555, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/resilio/config", Container: "/config"}, {Host: "${DATA_DIR}/resilio/sync", Container: "/sync"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "resilio-sync", Port: 8888,
			AptPackages: []string{"resilio-sync"},
			AptRepo: &models.AptRepo{
				KeyURL:   "https://linux-packages.resilio.com/resilio-sync/key.asc",
				RepoLine: "deb [signed-by=/usr/share/keyrings/velour-resilio.gpg] https://linux-packages.resilio.com/resilio-sync/deb resilio-sync non-free",
			},
			User: "rslsync",
		}},
	{ID: "nextcloud", Name: "Nextcloud", Description: "Self-hosted cloud storage, contacts, calendar and collaboration platform.", Icon: "nextcloud", Category: "sync", Image: "lscr.io/linuxserver/nextcloud:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8443, Container: 443, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/nextcloud/config", Container: "/config"}, {Host: "${DATA_DIR}/nextcloud/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "apache2", Port: 8443,
			Dependencies: []string{"apache2", "libapache2-mod-php", "php-gd", "php-json", "php-mysql", "php-curl", "php-mbstring", "php-intl", "php-imagick", "php-xml", "php-zip", "php-bcmath", "php-gmp", "unzip", "curl"},
			User: "www-data", ConfigDir: "${DATA_DIR}/nextcloud",
			InstallScript: `#!/bin/bash
set -e
curl -sL https://download.nextcloud.com/server/releases/latest.tar.bz2 | tar xj -C /var/www/
chown -R www-data:www-data /var/www/nextcloud
a2enmod rewrite headers env dir mime
systemctl restart apache2`,
		}},
	{ID: "syncthing", Name: "Syncthing", Description: "Continuous, decentralized file synchronization between devices.", Icon: "syncthing", Category: "sync", Image: "lscr.io/linuxserver/syncthing:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8384, Container: 8384, Protocol: "tcp"}, {Host: 22000, Container: 22000, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/syncthing/config", Container: "/config"}, {Host: "${DATA_DIR}/syncthing/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "syncthing@velour", Port: 8384,
			AptPackages: []string{"syncthing"},
			AptRepo: &models.AptRepo{
				KeyURL:   "https://syncthing.net/release-key.gpg",
				RepoLine: "deb [signed-by=/usr/share/keyrings/velour-release-key.gpg] https://apt.syncthing.net/ syncthing stable",
			},
			PostInstallCmds: []string{`sed -i 's|<address>127.0.0.1:8384</address>|<address>0.0.0.0:8384</address>|' /home/velour/.config/syncthing/config.xml 2>/dev/null; true`},
		}},
	{ID: "rclone", Name: "Rclone", Description: "Command-line tool to manage files on cloud storage with web GUI.", Icon: "rclone", Category: "sync", Image: "rclone/rclone:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 5572, Container: 5572, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/rclone/config", Container: "/config/rclone"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "rclone-webgui", Port: 5572,
			User: "rclone", ConfigDir: "${DATA_DIR}/rclone",
			InstallScript: `#!/bin/bash
set -e
curl -fsSL https://rclone.org/install.sh | bash`,
			ServiceUnit: `[Unit]
Description=Rclone Web GUI
After=network.target

[Service]
Type=simple
User=rclone
ExecStart=/usr/bin/rclone rcd --rc-web-gui --rc-addr=0.0.0.0:5572 --config /opt/velour/rclone/rclone.conf
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},

	{ID: "filebrowser", Name: "File Browser", Description: "Web-based file manager with a clean interface for uploading, deleting, previewing and sharing files.", Icon: "filebrowser", Category: "sync", Image: "filebrowser/filebrowser:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8092, Container: 80, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/filebrowser/config", Container: "/config"}, {Host: "${DATA_DIR}", Container: "/srv"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "filebrowser", Port: 8092,
			User: "filebrowser", ConfigDir: "${DATA_DIR}/filebrowser",
			InstallScript: `#!/bin/bash
set -e
curl -fsSL https://raw.githubusercontent.com/filebrowser/get/master/get.sh | bash`,
			ServiceUnit: `[Unit]
Description=File Browser
After=network.target

[Service]
Type=simple
User=filebrowser
ExecStart=/usr/local/bin/filebrowser -a 0.0.0.0 -p 8092 -r /opt/velour -d /opt/velour/filebrowser/filebrowser.db
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
			PostInstallCmds: []string{
				`filebrowser users update admin -p "${VELOUR_PASS}" -d /opt/velour/filebrowser/filebrowser.db 2>/dev/null; true`,
				`filebrowser users update admin --username "${VELOUR_USER}" -d /opt/velour/filebrowser/filebrowser.db 2>/dev/null; true`,
			},
		}},

	// ── Network ──
	{ID: "wireguard", Name: "WireGuard", Description: "Fast, modern, secure VPN tunnel. Simple yet powerful.", Icon: "wireguard", Category: "network", Image: "lscr.io/linuxserver/wireguard:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 51820, Container: 51820, Protocol: "udp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/wireguard/config", Container: "/config"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "wg-quick@wg0", Port: 0,
			AptPackages: []string{"wireguard", "wireguard-tools"},
			ConfigDir: "/etc/wireguard",
		}},
	{ID: "thelounge", Name: "The Lounge", Description: "Modern, self-hosted IRC client with always-on connectivity and web interface.", Icon: "thelounge", Category: "network", Image: "lscr.io/linuxserver/thelounge:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 9000, Container: 9000, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/thelounge/config", Container: "/config"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "script", ServiceName: "thelounge", Port: 9000,
			Dependencies: []string{"curl"},
			User: "thelounge", ConfigDir: "${DATA_DIR}/thelounge",
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
	{ID: "znc", Name: "ZNC", Description: "Advanced IRC bouncer that stays connected and buffers messages while you're offline.", Icon: "znc", Category: "network", Image: "lscr.io/linuxserver/znc:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 6501, Container: 6501, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/znc/config", Container: "/config"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "znc", Port: 6501,
			AptPackages: []string{"znc"},
			User: "znc", ConfigDir: "${DATA_DIR}/znc",
		}},

	// ── System ──
	{ID: "uptimekuma", Name: "Uptime Kuma", Description: "Self-hosted monitoring tool with beautiful status pages and notifications.", Icon: "uptimekuma", Category: "system", Image: "louislam/uptime-kuma:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 3010, Container: 3001, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/uptimekuma/data", Container: "/app/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
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
	{ID: "mariadb", Name: "MariaDB", Description: "Community-developed MySQL fork. Reliable, high performance database server.", Icon: "mariadb", Category: "system", Image: "lscr.io/linuxserver/mariadb:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 3306, Container: 3306, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/mariadb/config", Container: "/config"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000", "MYSQL_ROOT_PASSWORD": "changeme"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "mariadb", Port: 3306,
			AptPackages: []string{"mariadb-server"},
			User: "mysql",
		}},

	// ── Books & Reading ──
	{ID: "readarr", Name: "Readarr", Description: "Book, audiobook and comic collection manager for usenet and torrent users. Part of the Arr stack.", Icon: "readarr", Category: "media", Image: "lscr.io/linuxserver/readarr:develop",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8787, Container: 8787, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/readarr/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "readarr", Port: 8787,
			AptPackages: []string{"readarr"},
			AptRepo: &models.AptRepo{
				KeyURL:   "https://apt.sonarr.tv/sonarr.asc",
				RepoLine: "deb [signed-by=/usr/share/keyrings/velour-sonarr.gpg] https://apt.sonarr.tv/debian buster main",
			},
			User: "readarr", ConfigDir: "/var/lib/readarr",
			PostInstallCmds: []string{
				`sed -i 's|<AuthenticationMethod>.*</AuthenticationMethod>|<AuthenticationMethod>Forms</AuthenticationMethod>|' /var/lib/readarr/config.xml 2>/dev/null; true`,
				`sed -i 's|<AuthenticationRequired>.*</AuthenticationRequired>|<AuthenticationRequired>Enabled</AuthenticationRequired>|' /var/lib/readarr/config.xml 2>/dev/null; true`,
			},
		}},

	{ID: "komga", Name: "Komga", Description: "Free and open-source manga, comic and book media server with OPDS support and web reader.", Icon: "komga", Category: "server", Image: "gotson/komga:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 25600, Container: 25600, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/komga/config", Container: "/config"}, {Host: "${DATA_DIR}/data/comics", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "binary", ServiceName: "komga", Port: 25600,
			BinaryURL:  "https://github.com/gotson/komga/releases/latest/download/komga.jar",
			BinaryPath: "/opt/komga/komga.jar",
			ConfigDir:  "${DATA_DIR}/komga", User: "komga",
			InstallScript: `apt-get install -y -qq openjdk-17-jre-headless
mkdir -p /opt/komga ${DATA_DIR}/komga
useradd -r -s /usr/sbin/nologin komga 2>/dev/null || true`,
			ServiceUnit: `[Unit]
Description=Komga
After=network.target

[Service]
Type=simple
User=komga
ExecStart=/usr/bin/java -jar /opt/komga/komga.jar --komga.config-dir=/opt/velour/komga --server.port=25600
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},

	// ── Adult Content ──
	{ID: "whisparr", Name: "Whisparr", Description: "Adult content collection manager. Automatically searches, downloads and manages content.", Icon: "whisparr", Category: "media", Image: "cr.hotio.dev/hotio/whisparr:nightly",
		InstallTypes: []models.InstallType{models.InstallDocker},
		Ports: []models.PortMapping{{Host: 6969, Container: 6969, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/whisparr/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"}},

	// ── Automation & Sync Tools ──
	{ID: "recyclarr", Name: "Recyclarr", Description: "Automatically sync TRaSH Guides recommended settings to Sonarr and Radarr instances.", Icon: "recyclarr", Category: "download", Image: "ghcr.io/recyclarr/recyclarr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/recyclarr/config", Container: "/config"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "binary", ServiceName: "recyclarr", Port: 0,
			BinaryURL:  "https://github.com/recyclarr/recyclarr/releases/latest/download/recyclarr-linux-x64.tar.xz",
			BinaryPath: "/usr/local/bin/recyclarr",
			ConfigDir:  "${DATA_DIR}/recyclarr", User: "recyclarr",
			ServiceUnit: `[Unit]
Description=Recyclarr
After=network.target

[Service]
Type=oneshot
User=recyclarr
ExecStart=/usr/local/bin/recyclarr sync --config /opt/velour/recyclarr/recyclarr.yml
WorkingDirectory=/opt/velour/recyclarr

[Install]
WantedBy=multi-user.target`,
		}},

	// ── Notifications ──
	{ID: "notifiarr", Name: "Notifiarr", Description: "Unified notification client for Sonarr, Radarr, Lidarr, Readarr, Prowlarr, Plex and Tautulli.", Icon: "notifiarr", Category: "media", Image: "golift/notifiarr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 5454, Container: 5454, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/notifiarr/config", Container: "/config"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "notifiarr", Port: 5454,
			AptPackages: []string{"notifiarr"},
			AptRepo: &models.AptRepo{
				KeyURL:   "https://packagecloud.io/golift/pkgs/gpgkey",
				RepoLine: "deb [signed-by=/usr/share/keyrings/velour-golift.gpg] https://packagecloud.io/golift/pkgs/debian/ any main",
			},
			User: "notifiarr", ConfigDir: "/etc/notifiarr",
		}},

	{ID: "requestrr", Name: "Requestrr", Description: "Discord chatbot for requesting movies and TV shows via Sonarr, Radarr and Overseerr.", Icon: "requestrr", Category: "request", Image: "thomst08/requestrr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 4545, Container: 4545, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/requestrr/config", Container: "/root/config"}},
		Env: map[string]string{"TZ": "Europe/Istanbul"},
		Native: &models.NativeConfig{
			Method: "binary", ServiceName: "requestrr", Port: 4545,
			BinaryURL:  "https://github.com/thomst08/requestrr/releases/latest/download/requestrr-linux-x64.tar.gz",
			BinaryPath: "/opt/requestrr/Requestrr.WebApi",
			ConfigDir:  "${DATA_DIR}/requestrr", User: "requestrr",
			InstallScript: `apt-get install -y -qq libicu-dev
mkdir -p /opt/requestrr ${DATA_DIR}/requestrr
useradd -r -s /usr/sbin/nologin requestrr 2>/dev/null || true`,
			ServiceUnit: `[Unit]
Description=Requestrr
After=network.target

[Service]
Type=simple
User=requestrr
WorkingDirectory=/opt/requestrr
ExecStart=/opt/requestrr/Requestrr.WebApi
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},

	// ── Media Processing ──
	{ID: "tdarr", Name: "Tdarr", Description: "Distributed transcoding system for automating media library transcode and remux management.", Icon: "tdarr", Category: "media", Image: "ghcr.io/haveagitgat/tdarr:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 8265, Container: 8265, Protocol: "tcp"}, {Host: 8266, Container: 8266, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/tdarr/server", Container: "/app/server"}, {Host: "${DATA_DIR}/tdarr/configs", Container: "/app/configs"}, {Host: "${DATA_DIR}/tdarr/logs", Container: "/app/logs"}, {Host: "${DATA_DIR}/data", Container: "/media"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000", "serverIP": "0.0.0.0", "serverPort": "8266", "webUIPort": "8265"},
		Native: &models.NativeConfig{
			Method: "binary", ServiceName: "tdarr", Port: 8265,
			BinaryURL:  "https://f000.backblazeb2.com/file/tdarrs/versions/latest/linux_x64/Tdarr_Server.zip",
			BinaryPath: "/opt/tdarr/Tdarr_Server/Tdarr_Server",
			ConfigDir:  "${DATA_DIR}/tdarr", User: "tdarr",
			InstallScript: `apt-get install -y -qq unzip handbrake-cli ffmpeg
mkdir -p /opt/tdarr ${DATA_DIR}/tdarr
useradd -r -s /usr/sbin/nologin tdarr 2>/dev/null || true`,
			ServiceUnit: `[Unit]
Description=Tdarr Server
After=network.target

[Service]
Type=simple
User=tdarr
ExecStart=/opt/tdarr/Tdarr_Server/Tdarr_Server
Environment=NODE_ENV=production
WorkingDirectory=/opt/tdarr/Tdarr_Server
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},

	{ID: "handbrake", Name: "HandBrake", Description: "Open-source video transcoder with web GUI. Convert videos to modern formats with hardware acceleration.", Icon: "handbrake", Category: "media", Image: "jlesage/handbrake:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 5800, Container: 5800, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/handbrake/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/storage"}},
		Env: map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		Native: &models.NativeConfig{
			Method: "apt", ServiceName: "handbrake", Port: 0,
			AptPackages: []string{"handbrake-cli"},
			User: "root",
		}},

	// ── Media Organization ──
	{ID: "stash", Name: "Stash", Description: "Organizer and metadata scraper for your personal media collection with tagging and filtering.", Icon: "stash", Category: "server", Image: "stashapp/stash:latest",
		InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
		Ports: []models.PortMapping{{Host: 9999, Container: 9999, Protocol: "tcp"}},
		Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/stash/config", Container: "/root/.stash"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
		Env: map[string]string{"STASH_STASH": "/data", "STASH_GENERATED": "/root/.stash/generated", "STASH_METADATA": "/root/.stash/metadata", "STASH_CACHE": "/root/.stash/cache"},
		Native: &models.NativeConfig{
			Method: "binary", ServiceName: "stash", Port: 9999,
			BinaryURL:  "https://github.com/stashapp/stash/releases/latest/download/stash-linux",
			BinaryPath: "/usr/local/bin/stash",
			ConfigDir:  "${DATA_DIR}/stash", User: "stash",
			InstallScript: `apt-get install -y -qq ffmpeg
mkdir -p ${DATA_DIR}/stash
useradd -r -s /usr/sbin/nologin stash 2>/dev/null || true`,
			ServiceUnit: `[Unit]
Description=Stash
After=network.target

[Service]
Type=simple
User=stash
ExecStart=/usr/local/bin/stash --port 9999 --host 0.0.0.0
Environment=STASH_CONFIG_FILE=/opt/velour/stash/config.yml
Restart=on-failure

[Install]
WantedBy=multi-user.target`,
		}},
}

func GetRegistry() []models.ServiceDefinition {
	// Ensure all entries have at least Docker as install type
	result := make([]models.ServiceDefinition, len(Registry))
	copy(result, Registry)
	for i := range result {
		if len(result[i].InstallTypes) == 0 {
			result[i].InstallTypes = []models.InstallType{models.InstallDocker}
		}
	}
	return result
}

func FindByID(id string) *models.ServiceDefinition {
	for _, s := range Registry {
		if s.ID == id {
			return &s
		}
	}
	return nil
}
