package services

import "github.com/ozsari/velour/internal/models"

func init() {
	Registry = append(Registry,
		models.ServiceDefinition{
			ID: "plex", Name: "Plex", Description: "Stream movies, TV, music and more to any device from your own media library.", Icon: "plex", Category: "server", Image: "lscr.io/linuxserver/plex:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 32400, Container: 32400, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/plex/config", Container: "/config"}, {Host: "${DATA_DIR}/data/media", Container: "/data/media"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
			Native: &models.NativeConfig{
				Method: "apt", ServiceName: "plexmediaserver", Port: 32400,
				AptPackages: []string{"plexmediaserver"},
				AptRepo: &models.AptRepo{
					KeyURL:   "https://downloads.plex.tv/plex-keys/PlexSign.key",
					RepoLine: "deb https://downloads.plex.tv/repo/deb public main",
				},
				User: "plex",
			}},
		models.ServiceDefinition{
			ID: "jellyfin", Name: "Jellyfin", Description: "Free, open-source media server. Stream to any device from your own server.", Icon: "jellyfin", Category: "server", Image: "lscr.io/linuxserver/jellyfin:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8097, Container: 8096, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/jellyfin/config", Container: "/config"}, {Host: "${DATA_DIR}/data/media", Container: "/data/media"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
			Native: &models.NativeConfig{
				Method: "apt", ServiceName: "jellyfin", Port: 8096,
				AptPackages: []string{"jellyfin"},
				AptRepo: &models.AptRepo{
					KeyURL:   "https://repo.jellyfin.org/jellyfin_team.gpg.key",
					RepoLine: "deb [signed-by=/usr/share/keyrings/velour-jellyfin_team.gpg] https://repo.jellyfin.org/debian bookworm main",
				},
				User: "jellyfin",
			}},
		models.ServiceDefinition{
			ID: "emby", Name: "Emby", Description: "Media server that organizes and streams your personal media collection.", Icon: "emby", Category: "server", Image: "lscr.io/linuxserver/emby:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8096, Container: 8096, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/emby/config", Container: "/config"}, {Host: "${DATA_DIR}/data/media", Container: "/data/media"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
			Native: &models.NativeConfig{
				Method: "apt", ServiceName: "emby-server", Port: 8096,
				AptPackages: []string{"emby-server"},
				AptRepo: &models.AptRepo{
					KeyURL:   "https://mb3admin.com/startupapikey/keyfile.asc",
					RepoLine: "deb [signed-by=/usr/share/keyrings/velour-emby.gpg] https://packages.emby.media/deb stable main",
				},
				User: "emby",
			}},
		models.ServiceDefinition{
			ID: "airsonic", Name: "Airsonic", Description: "Free, web-based media streamer for your music collection.", Icon: "airsonic", Category: "server", Image: "lscr.io/linuxserver/airsonic-advanced:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 4040, Container: 4040, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/airsonic/config", Container: "/config"}, {Host: "${DATA_DIR}/data/media/music", Container: "/music"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
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
		models.ServiceDefinition{
			ID: "navidrome", Name: "Navidrome", Description: "Modern, open-source music server and streamer compatible with Subsonic/Airsonic.", Icon: "navidrome", Category: "server", Image: "deluan/navidrome:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 4533, Container: 4533, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/navidrome/config", Container: "/navidrome"}, {Host: "${DATA_DIR}/data/media/music", Container: "/music"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
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
		models.ServiceDefinition{
			ID: "calibreweb", Name: "Calibre-Web", Description: "Web app for browsing, reading and downloading eBooks from a Calibre database.", Icon: "calibreweb", Category: "server", Image: "lscr.io/linuxserver/calibre-web:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8084, Container: 8083, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/calibreweb/config", Container: "/config"}, {Host: "${DATA_DIR}/data/media/books", Container: "/books"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
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
		models.ServiceDefinition{
			ID: "tautulli", Name: "Tautulli", Description: "Monitoring and tracking tool for Plex Media Server usage and statistics.", Icon: "tautulli", Category: "server", Image: "lscr.io/linuxserver/tautulli:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8181, Container: 8181, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/tautulli/config", Container: "/config"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
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
		models.ServiceDefinition{
			ID: "audiobookshelf", Name: "Audiobookshelf", Description: "Self-hosted audiobook and podcast server with web-based player.", Icon: "audiobookshelf", Category: "server", Image: "ghcr.io/advplyr/audiobookshelf:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 13378, Container: 80, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/audiobookshelf/config", Container: "/config"}, {Host: "${DATA_DIR}/audiobookshelf/metadata", Container: "/metadata"}, {Host: "${DATA_DIR}/data/media/audiobooks", Container: "/audiobooks"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
			Native: &models.NativeConfig{
				Method: "apt", ServiceName: "audiobookshelf", Port: 13378,
				AptPackages: []string{"audiobookshelf"},
				AptRepo: &models.AptRepo{
					KeyURL:   "https://advplyr.github.io/audiobookshelf-ppa/KEY.gpg",
					RepoLine: "deb [signed-by=/usr/share/keyrings/velour-audiobookshelf.gpg] https://advplyr.github.io/audiobookshelf-ppa ./",
				},
			}},
		models.ServiceDefinition{
			ID: "kavita", Name: "Kavita", Description: "Fast, feature-rich manga, comic, and book reader server.", Icon: "kavita", Category: "server", Image: "jvmilazz0/kavita:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 5000, Container: 5000, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/kavita/config", Container: "/kavita/config"}, {Host: "${DATA_DIR}/data/media/books", Container: "/books"}, {Host: "${DATA_DIR}/data/media/comics", Container: "/comics"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
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
		models.ServiceDefinition{
			ID: "ubooquity", Name: "Ubooquity", Description: "Home server for comics and eBooks with OPDS feed and web reader.", Icon: "ubooquity", Category: "server", Image: "lscr.io/linuxserver/ubooquity:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 2202, Container: 2202, Protocol: "tcp"}, {Host: 2203, Container: 2203, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/ubooquity/config", Container: "/config"}, {Host: "${DATA_DIR}/data/media/books", Container: "/books"}, {Host: "${DATA_DIR}/data/media/comics", Container: "/comics"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
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
		models.ServiceDefinition{
			ID: "komga", Name: "Komga", Description: "Free and open-source manga, comic and book media server with OPDS support and web reader.", Icon: "komga", Category: "server", Image: "gotson/komga:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 25600, Container: 25600, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/komga/config", Container: "/config"}, {Host: "${DATA_DIR}/data/comics", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
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
		models.ServiceDefinition{
			ID: "stash", Name: "Stash", Description: "Organizer and metadata scraper for your personal media collection with tagging and filtering.", Icon: "stash", Category: "server", Image: "stashapp/stash:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 9999, Container: 9999, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/stash/config", Container: "/root/.stash"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"STASH_STASH": "/data", "STASH_GENERATED": "/root/.stash/generated", "STASH_METADATA": "/root/.stash/metadata", "STASH_CACHE": "/root/.stash/cache"},
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
	)
}
