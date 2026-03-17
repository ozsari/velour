package services

import "github.com/ozsari/velour/internal/models"

func init() {
	Registry = append(Registry,
		models.ServiceDefinition{
			ID: "prowlarr", Name: "Prowlarr", Description: "Indexer manager that integrates with Sonarr, Radarr, Lidarr.", Icon: "prowlarr", Category: "indexer", Image: "lscr.io/linuxserver/prowlarr:latest",
			InstallTypes: []models.InstallType{models.InstallDocker},
			Ports:        []models.PortMapping{{Host: 9696, Container: 9696, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/prowlarr/config", Container: "/config"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
		},
		models.ServiceDefinition{
			ID: "jackett", Name: "Jackett", Description: "Proxy server translating queries for tracker-site-specific searches.", Icon: "jackett", Category: "indexer", Image: "lscr.io/linuxserver/jackett:latest",
			InstallTypes: []models.InstallType{models.InstallDocker},
			Ports:        []models.PortMapping{{Host: 9117, Container: 9117, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/jackett/config", Container: "/config"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
		},
		models.ServiceDefinition{
			ID: "nzbhydra2", Name: "NZBHydra2", Description: "Meta search for usenet indexers. Unified access to multiple newznab indexers.", Icon: "nzbhydra2", Category: "indexer", Image: "lscr.io/linuxserver/nzbhydra2:latest",
			InstallTypes: []models.InstallType{models.InstallDocker},
			Ports:        []models.PortMapping{{Host: 5076, Container: 5076, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{{Host: "${DATA_DIR}/nzbhydra2/config", Container: "/config"}},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
		},
		models.ServiceDefinition{
			ID: "flaresolverr", Name: "FlareSolverr", Description: "Proxy server to bypass Cloudflare protection for scrapers.", Icon: "flaresolverr", Category: "indexer", Image: "ghcr.io/flaresolverr/flaresolverr:latest",
			InstallTypes: []models.InstallType{models.InstallDocker},
			Ports:        []models.PortMapping{{Host: 8191, Container: 8191, Protocol: "tcp"}},
			Volumes:      []models.VolumeMapping{},
			Env:          map[string]string{"TZ": "Europe/Istanbul"},
		},
	)
}
