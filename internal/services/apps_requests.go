package services

import "github.com/ozsari/velour/internal/models"

func init() {
	Registry = append(Registry,
		models.ServiceDefinition{
			ID: "ombi", Name: "Ombi", Description: "Self-hosted app for Plex/Emby users to request content.", Icon: "ombi", Category: "request", Image: "lscr.io/linuxserver/ombi:latest",
			Ports:   []models.PortMapping{{Host: 5000, Container: 3579, Protocol: "tcp"}},
			Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/ombi/config", Container: "/config"}},
			Env:     map[string]string{"TZ": "Europe/Istanbul"},
		},
		models.ServiceDefinition{
			ID: "overseerr", Name: "Overseerr", Description: "Modern media request and discovery tool for Plex, Sonarr and Radarr.", Icon: "overseerr", Category: "request", Image: "lscr.io/linuxserver/overseerr:latest",
			Ports:   []models.PortMapping{{Host: 5055, Container: 5055, Protocol: "tcp"}},
			Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/overseerr/config", Container: "/config"}},
			Env:     map[string]string{"TZ": "Europe/Istanbul"},
		},
		models.ServiceDefinition{
			ID: "doplarr", Name: "Doplarr", Description: "Discord bot for requesting media through Overseerr or Ombi.", Icon: "doplarr", Category: "request", Image: "lscr.io/linuxserver/doplarr:latest",
			Ports:   []models.PortMapping{},
			Volumes: []models.VolumeMapping{{Host: "${DATA_DIR}/doplarr/config", Container: "/config"}},
			Env:     map[string]string{"TZ": "Europe/Istanbul"},
		},
		models.ServiceDefinition{
			ID: "requestrr", Name: "Requestrr", Description: "Discord chatbot for requesting movies and TV shows via Sonarr, Radarr and Overseerr.", Icon: "requestrr", Category: "request", Image: "thomst08/requestrr:latest",
			InstallTypes: []models.InstallType{models.InstallDocker, models.InstallNative},
			Ports:        []models.PortMapping{{Host: 4545, Container: 4545, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/requestrr/config", Container: "/root/config"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
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
			},
		},
	)
}
