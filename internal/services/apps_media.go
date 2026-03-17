package services

import "github.com/ozsari/velour/internal/models"

func init() {
	Registry = append(Registry,
		models.ServiceDefinition{
			ID: "sonarr", Name: "Sonarr", Description: "Smart PVR for newsgroup and bittorrent users. Monitors RSS feeds for new episodes.", Icon: "sonarr", Category: "media", Image: "lscr.io/linuxserver/sonarr:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8989, Container: 8989, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/sonarr/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
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
			},
		},
		models.ServiceDefinition{
			ID: "sonarr2", Name: "Sonarr (2nd)", Description: "Second instance of Sonarr for managing a separate TV library.", Icon: "sonarr", Category: "media", Image: "lscr.io/linuxserver/sonarr:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8990, Container: 8989, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/sonarr2/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
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
			},
		},
		models.ServiceDefinition{
			ID: "radarr", Name: "Radarr", Description: "Movie collection manager. Automatically searches, downloads and manages movies.", Icon: "radarr", Category: "media", Image: "lscr.io/linuxserver/radarr:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 7878, Container: 7878, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/radarr/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
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
			},
		},
		models.ServiceDefinition{
			ID: "radarr2", Name: "Radarr (2nd)", Description: "Second instance of Radarr for managing a separate movie library.", Icon: "radarr", Category: "media", Image: "lscr.io/linuxserver/radarr:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 7879, Container: 7878, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/radarr2/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
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
			},
		},
		models.ServiceDefinition{
			ID: "bazarr", Name: "Bazarr", Description: "Companion to Sonarr and Radarr for managing and downloading subtitles.", Icon: "bazarr", Category: "media", Image: "lscr.io/linuxserver/bazarr:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 6767, Container: 6767, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/bazarr/config", Container: "/config"}, {Host: "${DATA_DIR}/data/media", Container: "/data/media"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
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
			},
		},
		models.ServiceDefinition{
			ID: "lidarr", Name: "Lidarr", Description: "Music collection manager. Monitors RSS feeds for new albums and manages your library.", Icon: "lidarr", Category: "media", Image: "lscr.io/linuxserver/lidarr:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8686, Container: 8686, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/lidarr/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
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
			},
		},
		models.ServiceDefinition{
			ID: "medusa", Name: "Medusa", Description: "Automatic video library manager for TV shows with multi-source support.", Icon: "medusa", Category: "media", Image: "lscr.io/linuxserver/medusa:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8081, Container: 8081, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/medusa/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
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
			},
		},
		models.ServiceDefinition{
			ID: "sickchill", Name: "SickChill", Description: "Automatic video library manager for TV shows from various sources.", Icon: "sickchill", Category: "media", Image: "lscr.io/linuxserver/sickchill:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8083, Container: 8081, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/sickchill/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
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
			},
		},
		models.ServiceDefinition{
			ID: "sickgear", Name: "SickGear", Description: "SickBeard fork with improved stability, performance and features.", Icon: "sickgear", Category: "media", Image: "lscr.io/linuxserver/sickgear:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8082, Container: 8081, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/sickgear/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
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
			},
		},
		models.ServiceDefinition{
			ID: "mylar3", Name: "Mylar3", Description: "Automated comic book downloader for usenet and torrent users.", Icon: "mylar3", Category: "media", Image: "lscr.io/linuxserver/mylar3:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8090, Container: 8090, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/mylar3/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
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
			},
		},
		models.ServiceDefinition{
			ID: "filebot", Name: "FileBot", Description: "Ultimate tool for organizing and renaming movies, TV shows, anime and music.", Icon: "filebot", Category: "media", Image: "jlesage/filebot:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 5800, Container: 5800, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/filebot/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
			Native: &models.NativeConfig{
				Method: "script", ServiceName: "filebot", Port: 0,
				Dependencies: []string{"default-jre-headless", "curl", "libmediainfo0v5"},
				User: "filebot", ConfigDir: "${DATA_DIR}/filebot",
				InstallScript: `#!/bin/bash
set -e
curl -sL https://get.filebot.net/filebot/FileBot_5.1.6/FileBot_5.1.6_amd64.deb -o /tmp/filebot.deb
dpkg -i /tmp/filebot.deb || apt-get install -f -y
rm /tmp/filebot.deb`,
			},
		},
		models.ServiceDefinition{
			ID: "lazylibrarian", Name: "LazyLibrarian", Description: "Automated book downloader for eBooks and audiobooks from usenet and torrents.", Icon: "lazylibrarian", Category: "media", Image: "lscr.io/linuxserver/lazylibrarian:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 5299, Container: 5299, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/lazylibrarian/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
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
			},
		},
		models.ServiceDefinition{
			ID: "maintainerr", Name: "Maintainerr", Description: "Automated media maintenance for Plex. Clean up old or unwatched content.", Icon: "maintainerr", Category: "media", Image: "ghcr.io/jorenn92/maintainerr:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 6246, Container: 6246, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/maintainerr/config", Container: "/opt/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
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
			},
		},
		models.ServiceDefinition{
			ID: "unpackerr", Name: "Unpackerr", Description: "Extracts downloaded archives for Sonarr, Radarr, Lidarr and Readarr.", Icon: "unpackerr", Category: "media", Image: "golift/unpackerr:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/unpackerr/config", Container: "/config"}, {Host: "${DATA_DIR}/data/downloads", Container: "/data/downloads"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
			Native: &models.NativeConfig{
				Method: "apt", ServiceName: "unpackerr", Port: 0,
				AptPackages: []string{"unpackerr"},
				AptRepo: &models.AptRepo{
					KeyURL:   "https://packagecloud.io/golift/pkgs/gpgkey",
					RepoLine: "deb [signed-by=/usr/share/keyrings/velour-golift.gpg] https://packagecloud.io/golift/pkgs/debian/ any main",
				},
				User: "unpackerr", ConfigDir: "${DATA_DIR}/unpackerr",
			},
		},
		models.ServiceDefinition{
			ID: "readarr", Name: "Readarr", Description: "Book, audiobook and comic collection manager for usenet and torrent users. Part of the Arr stack.", Icon: "readarr", Category: "media", Image: "lscr.io/linuxserver/readarr:develop",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8787, Container: 8787, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/readarr/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
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
			},
		},
		models.ServiceDefinition{
			ID: "whisparr", Name: "Whisparr", Description: "Adult content collection manager. Automatically searches, downloads and manages content.", Icon: "whisparr", Category: "media", Image: "cr.hotio.dev/hotio/whisparr:nightly",
			InstallTypes: []models.InstallType{models.InstallDocker},
			Ports:        []models.PortMapping{{Host: 6969, Container: 6969, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/whisparr/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
		},
		models.ServiceDefinition{
			ID: "notifiarr", Name: "Notifiarr", Description: "Unified notification client for Sonarr, Radarr, Lidarr, Readarr, Prowlarr, Plex and Tautulli.", Icon: "notifiarr", Category: "media", Image: "golift/notifiarr:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 5454, Container: 5454, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/notifiarr/config", Container: "/config"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
			Native: &models.NativeConfig{
				Method: "apt", ServiceName: "notifiarr", Port: 5454,
				AptPackages: []string{"notifiarr"},
				AptRepo: &models.AptRepo{
					KeyURL:   "https://packagecloud.io/golift/pkgs/gpgkey",
					RepoLine: "deb [signed-by=/usr/share/keyrings/velour-golift.gpg] https://packagecloud.io/golift/pkgs/debian/ any main",
				},
				User: "notifiarr", ConfigDir: "/etc/notifiarr",
			},
		},
		models.ServiceDefinition{
			ID: "tdarr", Name: "Tdarr", Description: "Distributed transcoding system for automating media library transcode and remux management.", Icon: "tdarr", Category: "media", Image: "ghcr.io/haveagitgat/tdarr:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8265, Container: 8265, Protocol: "tcp"}, {Host: 8266, Container: 8266, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/tdarr/server", Container: "/app/server"}, {Host: "${DATA_DIR}/tdarr/configs", Container: "/app/configs"}, {Host: "${DATA_DIR}/tdarr/logs", Container: "/app/logs"}, {Host: "${DATA_DIR}/data", Container: "/media"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000", "serverIP": "0.0.0.0", "serverPort": "8266", "webUIPort": "8265"},
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
			},
		},
		models.ServiceDefinition{
			ID: "handbrake", Name: "HandBrake", Description: "Open-source video transcoder with web GUI. Convert videos to modern formats with hardware acceleration.", Icon: "handbrake", Category: "media", Image: "jlesage/handbrake:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 5800, Container: 5800, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/handbrake/config", Container: "/config"}, {Host: "${DATA_DIR}/data", Container: "/storage"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
			Native: &models.NativeConfig{
				Method: "apt", ServiceName: "handbrake", Port: 0,
				AptPackages: []string{"handbrake-cli"},
				User: "root",
			},
		},
	)
}
