package services

import "github.com/ozsari/velour/internal/models"

func init() {
	Registry = append(Registry,
		models.ServiceDefinition{
			ID: "resilio", Name: "Resilio Sync", Description: "Fast, reliable file sync using peer-to-peer technology.", Icon: "resilio", Category: "sync", Image: "lscr.io/linuxserver/resilio-sync:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8888, Container: 8888, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/resilio/config", Container: "/config"}, {Host: "${DATA_DIR}/resilio/sync", Container: "/sync"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
			Native: &models.NativeConfig{
				Method: "apt", ServiceName: "resilio-sync", Port: 8888,
				AptPackages: []string{"resilio-sync"},
				AptRepo: &models.AptRepo{
					KeyURL:   "https://linux-packages.resilio.com/resilio-sync/key.asc",
					RepoLine: "deb [signed-by=/usr/share/keyrings/velour-resilio.gpg] https://linux-packages.resilio.com/resilio-sync/deb resilio-sync non-free",
				},
			}},
		models.ServiceDefinition{
			ID: "nextcloud", Name: "Nextcloud", Description: "Self-hosted cloud storage, contacts, calendar and collaboration platform.", Icon: "nextcloud", Category: "sync", Image: "lscr.io/linuxserver/nextcloud:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8443, Container: 443, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/nextcloud/config", Container: "/config"}, {Host: "${DATA_DIR}/nextcloud/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
			Native: &models.NativeConfig{
				Method: "script", ServiceName: "apache2", Port: 8443,
				Dependencies: []string{"apache2", "libapache2-mod-php", "php-gd", "php-json", "php-mysql", "php-curl", "php-mbstring", "php-intl", "php-imagick", "php-xml", "php-zip", "mariadb-server"},
				User: "www-data", ConfigDir: "${DATA_DIR}/nextcloud",
				InstallScript: `#!/bin/bash
curl -sL https://download.nextcloud.com/server/releases/latest.tar.bz2 | tar xj -C /var/www/
chown -R www-data:www-data /var/www/nextcloud
systemctl restart apache2`,
			}},
		models.ServiceDefinition{
			ID: "syncthing", Name: "Syncthing", Description: "Continuous, decentralized file synchronization between devices.", Icon: "syncthing", Category: "sync", Image: "lscr.io/linuxserver/syncthing:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8384, Container: 8384, Protocol: "tcp"}, {Host: 22000, Container: 22000, Protocol: "tcp"}, {Host: 21027, Container: 21027, Protocol: "udp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/syncthing/config", Container: "/config"}, {Host: "${DATA_DIR}/syncthing/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul", "PUID": "1000", "PGID": "1000"},
			Native: &models.NativeConfig{
				Method: "apt", ServiceName: "syncthing@velour", Port: 8384,
				AptPackages: []string{"syncthing"},
				AptRepo: &models.AptRepo{
					KeyURL:   "https://syncthing.net/release-key.gpg",
					RepoLine: "deb [signed-by=/usr/share/keyrings/velour-release-key.gpg] https://apt.syncthing.net/ syncthing stable",
				},
				PostInstallCmds: []string{`sed -i 's|<address>127.0.0.1:8384</address>|<address>0.0.0.0:8384</address>|' /home/velour/.config/syncthing/config.xml 2>/dev/null; true`},
			}},
		models.ServiceDefinition{
			ID: "rclone", Name: "Rclone", Description: "Command-line tool to manage files on cloud storage with web GUI.", Icon: "rclone", Category: "sync", Image: "rclone/rclone:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 5572, Container: 5572, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/rclone/config", Container: "/config/rclone"}, {Host: "${DATA_DIR}/data", Container: "/data"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
			Native: &models.NativeConfig{
				Method: "script", ServiceName: "rclone-webgui", Port: 5572,
				User: "rclone", ConfigDir: "${DATA_DIR}/rclone",
				InstallScript: `#!/bin/bash
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
		models.ServiceDefinition{
			ID: "filebrowser", Name: "File Browser", Description: "Web-based file manager with a clean interface for uploading, deleting, previewing and sharing files.", Icon: "filebrowser", Category: "sync", Image: "filebrowser/filebrowser:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 8092, Container: 8092, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/filebrowser/config", Container: "/config"}, {Host: "${DATA_DIR}", Container: "/srv"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
			Native: &models.NativeConfig{
				Method: "script", ServiceName: "filebrowser", Port: 8092,
				User: "filebrowser", ConfigDir: "${DATA_DIR}/filebrowser",
				InstallScript: `#!/bin/bash
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
					`python3 -c "
import subprocess, os
user = os.environ.get('VELOUR_USER', '')
passwd = os.environ.get('VELOUR_PASS', '')
db = '/opt/velour/filebrowser/filebrowser.db'
if user and passwd:
    subprocess.run(['filebrowser', 'users', 'update', 'admin', '-p', passwd, '-d', db])
    subprocess.run(['filebrowser', 'users', 'update', 'admin', '--username', user, '-d', db])
    print('OK: File Browser credentials set')
else:
    print('SKIP: no credentials')
"`,
				},
			}},
	)
}
