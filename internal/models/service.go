package models

import "time"

type ServiceStatus string

const (
	StatusRunning  ServiceStatus = "running"
	StatusStopped  ServiceStatus = "stopped"
	StatusInstalling ServiceStatus = "installing"
	StatusError    ServiceStatus = "error"
	StatusUnknown  ServiceStatus = "unknown"
)

type Service struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Icon        string        `json:"icon"`
	Category    string        `json:"category"`
	Port        int           `json:"port"`
	WebURL      string        `json:"web_url"`
	Status      ServiceStatus `json:"status"`
	Type        string        `json:"type"` // "docker" or "systemd"
	Image       string        `json:"image,omitempty"`
	Installed   bool          `json:"installed"`
	InstalledAt *time.Time    `json:"installed_at,omitempty"`
}

type InstallType string

const (
	InstallDocker InstallType = "docker"
	InstallNative InstallType = "native"
)

type ServiceDefinition struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Icon        string            `json:"icon"`
	Category    string            `json:"category"`
	// Docker install
	Image       string            `json:"image"`
	Ports       []PortMapping     `json:"ports"`
	Volumes     []VolumeMapping   `json:"volumes"`
	Env         map[string]string `json:"env"`
	DependsOn   []string          `json:"depends_on,omitempty"`
	// Native install
	Native      *NativeConfig     `json:"native,omitempty"`
	// Supported install types
	InstallTypes []InstallType    `json:"install_types"`
}

type NativeConfig struct {
	// How to install: "apt", "binary", "script"
	Method       string   `json:"method"`
	// APT packages to install
	AptPackages  []string `json:"apt_packages,omitempty"`
	// APT repo to add before install (e.g. PPA or custom repo)
	AptRepo      *AptRepo `json:"apt_repo,omitempty"`
	// Binary download URL (supports ${ARCH} placeholder)
	BinaryURL    string   `json:"binary_url,omitempty"`
	// Where to install the binary
	BinaryPath   string   `json:"binary_path,omitempty"`
	// Custom install script (embedded shell commands)
	InstallScript string  `json:"install_script,omitempty"`
	// Systemd service name
	ServiceName  string   `json:"service_name"`
	// Systemd unit file content
	ServiceUnit  string   `json:"service_unit,omitempty"`
	// Config directory
	ConfigDir    string   `json:"config_dir,omitempty"`
	// Data directory
	DataDir      string   `json:"data_dir,omitempty"`
	// User to run as
	User         string   `json:"user,omitempty"`
	// Port the native service listens on
	Port         int      `json:"port,omitempty"`
	// Dependencies (other native packages)
	Dependencies []string `json:"dependencies,omitempty"`
}

type AptRepo struct {
	// GPG key URL
	KeyURL  string `json:"key_url"`
	// Repo line (e.g. "deb https://apt.sonarr.tv/debian buster main")
	RepoLine string `json:"repo_line"`
}

type PortMapping struct {
	Host      int    `json:"host"`
	Container int    `json:"container"`
	Protocol  string `json:"protocol"`
}

type VolumeMapping struct {
	Host      string `json:"host"`
	Container string `json:"container"`
}
