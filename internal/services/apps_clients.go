package services

import "github.com/ozsari/velour/internal/models"

func init() {
	Registry = append(Registry,
		// ── Download Clients ──
		models.ServiceDefinition{
			ID: "deluge", Name: "Deluge", Description: "Lightweight, free, cross-platform BitTorrent client with a web interface.", Icon: "deluge", Category: "client", Image: "lscr.io/linuxserver/deluge:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8112, Container: 8112, Protocol: "tcp"}, {Host: 6881, Container: 6881, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/deluge/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
			Native: &models.NativeConfig{
				Method: "apt", ServiceName: "deluged", Port: 8112,
				AptPackages: []string{"deluged", "deluge-web"},
				ConfigDir:   "${DATA_DIR}/deluge", User: "deluge",
				PostInstallCmds: []string{
					// Create deluge-web service
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
					// Start both briefly so deluge-web creates default web.conf
					`systemctl start deluged && sleep 2 && systemctl start deluge-web && sleep 5 && systemctl stop deluge-web && systemctl stop deluged`,
					// Set web UI password from env var
					`python3 -c "
import hashlib, os, subprocess
password = os.environ.get('VELOUR_PASS', '')
if password:
    salt = hashlib.sha1(os.urandom(32)).hexdigest()
    s = hashlib.sha1(salt.encode('utf-8'))
    s.update(password.encode('utf-8'))
    pwd_hash = s.hexdigest()
    conf = '/opt/velour/deluge/web.conf'
    subprocess.run(['sed', '-i', 's|\"pwd_salt\": \"[^\"]*\"|\"pwd_salt\": \"' + salt + '\"|', conf])
    subprocess.run(['sed', '-i', 's|\"pwd_sha1\": \"[^\"]*\"|\"pwd_sha1\": \"' + pwd_hash + '\"|', conf])
    subprocess.run(['sed', '-i', 's|\"first_login\": true|\"first_login\": false|', conf])
    print('OK: password set')
else:
    print('SKIP: no VELOUR_PASS')
"`,
					// Start deluge-web with patched config
					`systemctl start deluge-web`,
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
			},
		},
		models.ServiceDefinition{
			ID: "flood", Name: "Flood", Description: "Modern web UI for rTorrent, qBittorrent, and Transmission with a clean interface.", Icon: "flood", Category: "client", Image: "jesec/flood:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 3001, Container: 3000, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/flood/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
			Native: &models.NativeConfig{
				Method:       "script", ServiceName: "flood", Port: 3001,
				Dependencies: []string{"curl"},
				User:         "flood", ConfigDir: "${DATA_DIR}/flood",
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
			},
		},
		models.ServiceDefinition{
			ID: "qbittorrent", Name: "qBittorrent", Description: "Free, open-source BitTorrent client with a feature-rich web interface.", Icon: "qbittorrent", Category: "client", Image: "lscr.io/linuxserver/qbittorrent:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8085, Container: 8080, Protocol: "tcp"}, {Host: 6882, Container: 6881, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/qbittorrent/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000", "WEBUI_PORT": "8080"},
			Native: &models.NativeConfig{
				Method: "apt", ServiceName: "qbittorrent-nox", Port: 8085,
				AptPackages: []string{"qbittorrent-nox"},
				User:        "qbittorrent",
				PostInstallCmds: []string{
					// Start qBittorrent, login with default creds, set new creds via API
					`systemctl start qbittorrent-nox && sleep 3 && python3 -c "
import urllib.request, urllib.parse, os, http.cookiejar
user = os.environ.get('VELOUR_USER', '')
passwd = os.environ.get('VELOUR_PASS', '')
if not user or not passwd:
    print('SKIP: no credentials')
    exit(0)
cj = http.cookiejar.CookieJar()
opener = urllib.request.build_opener(urllib.request.HTTPCookieProcessor(cj))
# Login with default creds
opener.open(urllib.request.Request('http://localhost:8085/api/v2/auth/login', data=b'username=admin&password=adminadmin'))
# Set new credentials
import json
prefs = json.dumps({'web_ui_username': user, 'web_ui_password': passwd})
opener.open(urllib.request.Request('http://localhost:8085/api/v2/app/setPreferences', data=('json=' + prefs).encode()))
print('OK: qBittorrent credentials set')
"`,
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
			},
		},
		models.ServiceDefinition{
			ID: "qui", Name: "qui", Description: "Fast, modern web interface for qBittorrent with multi-instance support and automations.", Icon: "qui", Category: "client", Image: "ghcr.io/autobrr/qui:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 7476, Container: 7476, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/qui/config", Container: "/config"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
			Native: &models.NativeConfig{
				Method:     "binary", ServiceName: "qui", Port: 7476,
				BinaryURL:  "https://github.com/autobrr/qui/releases/latest/download/qui_linux_${ARCH}.tar.gz",
				BinaryPath: "/usr/local/bin/qui",
				User:       "qui", ConfigDir: "${DATA_DIR}/qui",
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
			},
		},
		models.ServiceDefinition{
			ID: "rutorrent", Name: "ruTorrent", Description: "rTorrent + ruTorrent web UI + autodl-irssi bundle. Complete torrent solution with IRC automation.", Icon: "rutorrent", Category: "client", Image: "crazymax/rtorrent-rutorrent:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8087, Container: 8080, Protocol: "tcp"}, {Host: 50000, Container: 50000, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/rutorrent/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
			Native: &models.NativeConfig{
				Method:       "script", ServiceName: "rtorrent", Port: 8087,
				Dependencies: []string{"rtorrent", "nginx", "php-fpm", "php-cli", "mediainfo", "unrar-free", "curl"},
				User:         "rtorrent", ConfigDir: "${DATA_DIR}/rutorrent",
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
			},
		},
		models.ServiceDefinition{
			ID: "transmission", Name: "Transmission", Description: "Fast, easy and free BitTorrent client with minimal resource usage.", Icon: "transmission", Category: "client", Image: "lscr.io/linuxserver/transmission:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 9091, Container: 9091, Protocol: "tcp"}, {Host: 51413, Container: 51413, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/transmission/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
			Native: &models.NativeConfig{
				Method:      "apt", ServiceName: "transmission-daemon", Port: 9091,
				AptPackages: []string{"transmission-daemon"},
				User:        "debian-transmission",
				PostInstallCmds: []string{
					`sed -i 's/"rpc-whitelist-enabled":.*/"rpc-whitelist-enabled": false,/' /etc/transmission-daemon/settings.json 2>/dev/null; true`,
					`sed -i 's/"rpc-authentication-required":.*/"rpc-authentication-required": true,/' /etc/transmission-daemon/settings.json 2>/dev/null; true`,
					`python3 -c "
import json, os
user = os.environ.get('VELOUR_USER', '')
passwd = os.environ.get('VELOUR_PASS', '')
conf = '/etc/transmission-daemon/settings.json'
if user and passwd:
    try:
        with open(conf) as f:
            s = json.load(f)
        s['rpc-authentication-required'] = True
        s['rpc-whitelist-enabled'] = False
        s['rpc-username'] = user
        s['rpc-password'] = passwd
        with open(conf, 'w') as f:
            json.dump(s, f, indent=4)
        print('OK: Transmission credentials set')
    except Exception as e:
        print(f'ERR: {e}')
else:
    print('SKIP: no credentials')
"`,
				},
			},
		},
		models.ServiceDefinition{
			ID: "nzbget", Name: "NZBGet", Description: "Efficient usenet downloader written in C++ for maximum performance.", Icon: "nzbget", Category: "client", Image: "lscr.io/linuxserver/nzbget:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 6789, Container: 6789, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/nzbget/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
			Native: &models.NativeConfig{
				Method:     "binary", ServiceName: "nzbget", Port: 6789,
				BinaryURL:  "https://github.com/nzbgetcom/nzbget/releases/latest/download/nzbget-linux.run",
				BinaryPath: "/opt/nzbget/nzbget",
				ConfigDir:  "${DATA_DIR}/nzbget", User: "nzbget",
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
					`python3 -c "
import os, re
user = os.environ.get('VELOUR_USER', '')
passwd = os.environ.get('VELOUR_PASS', '')
conf = '/opt/velour/nzbget/nzbget.conf'
if user and passwd:
    try:
        with open(conf) as f:
            txt = f.read()
        txt = re.sub(r'^ControlUsername=.*', 'ControlUsername=' + user, txt, flags=re.MULTILINE)
        txt = re.sub(r'^ControlPassword=.*', 'ControlPassword=' + passwd, txt, flags=re.MULTILINE)
        with open(conf, 'w') as f:
            f.write(txt)
        print('OK: NZBGet credentials set')
    except Exception as e:
        print(f'ERR: {e}')
else:
    print('SKIP: no credentials')
"`,
				},
			},
		},
		models.ServiceDefinition{
			ID: "sabnzbd", Name: "SABnzbd", Description: "Free, open-source usenet downloader with web-based interface.", Icon: "sabnzbd", Category: "client", Image: "lscr.io/linuxserver/sabnzbd:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8088, Container: 8080, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/sabnzbd/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
			Native: &models.NativeConfig{
				Method:      "apt", ServiceName: "sabnzbdplus", Port: 8080,
				AptPackages: []string{"sabnzbdplus"},
				AptRepo: &models.AptRepo{
					KeyURL:   "https://ppa.launchpadcontent.net/jcfp/nobetas/ubuntu/dists/jammy/Release.gpg",
					RepoLine: "deb [signed-by=/usr/share/keyrings/velour-jcfp.gpg] https://ppa.launchpadcontent.net/jcfp/nobetas/ubuntu jammy main",
				},
				User: "sabnzbd", ConfigDir: "${DATA_DIR}/sabnzbd",
			},
		},
		models.ServiceDefinition{
			ID: "pyload", Name: "pyLoad", Description: "Free, open-source download manager for HTTP, FTP, and other protocols.", Icon: "pyload", Category: "client", Image: "lscr.io/linuxserver/pyload-ng:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8089, Container: 8000, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/pyload/config", Container: "/config"}, {Host: "${DATA_DIR}/data/downloads", Container: "/downloads"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
			Native: &models.NativeConfig{
				Method:       "script", ServiceName: "pyload", Port: 8089,
				Dependencies: []string{"python3", "python3-pip", "python3-venv"},
				User:         "pyload", ConfigDir: "${DATA_DIR}/pyload",
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
			},
		},
		models.ServiceDefinition{
			ID: "jdownloader2", Name: "JDownloader 2", Description: "Open-source download manager with web interface for direct downloads.", Icon: "jdownloader2", Category: "client", Image: "jlesage/jdownloader-2:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 5801, Container: 5800, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/jdownloader2/config", Container: "/config"}, {Host: "${DATA_DIR}/data/downloads", Container: "/output"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
			Native: &models.NativeConfig{
				Method:       "script", ServiceName: "jdownloader2", Port: 0,
				Dependencies: []string{"default-jre-headless", "curl"},
				User:         "jdownloader", ConfigDir: "${DATA_DIR}/jdownloader2",
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
			},
		},
	)
}
