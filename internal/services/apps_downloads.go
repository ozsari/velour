package services

import "github.com/ozsari/velour/internal/models"

func init() {
	Registry = append(Registry,
		models.ServiceDefinition{
			ID: "autobrr", Name: "Autobrr", Description: "Modern automation tool for torrents and usenet. Monitors IRC, RSS feeds and more.", Icon: "autobrr", Category: "download", Image: "ghcr.io/autobrr/autobrr:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 7474, Container: 7474, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/autobrr/config", Container: "/config"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
			Native: &models.NativeConfig{
				Method: "binary", ServiceName: "autobrr", Port: 7474,
				BinaryURL:       "https://github.com/autobrr/autobrr/releases/latest/download/autobrr_linux_${ARCH}.tar.gz",
				BinaryPath:      "/usr/local/bin/autobrr",
				ConfigDir:       "${DATA_DIR}/autobrr", User: "autobrr",
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
			},
		},
		models.ServiceDefinition{
			ID: "recyclarr", Name: "Recyclarr", Description: "Automatically sync TRaSH Guides recommended settings to Sonarr and Radarr instances.", Icon: "recyclarr", Category: "download", Image: "ghcr.io/recyclarr/recyclarr:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/recyclarr/config", Container: "/config"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
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
			},
		},
	)
}
