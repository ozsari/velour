package config

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Version     string `json:"version"`
	Host        string `json:"host"`
	Port        int    `json:"port"`
	DataDir     string `json:"data_dir"`
	DBPath      string `json:"db_path"`
	JWTSecret   string `json:"jwt_secret"`
	InstallMode string `json:"install_mode"` // "docker" or "native"
}

func defaultConfig() *Config {
	return &Config{
		Version:   "0.1.0",
		Host:      "0.0.0.0",
		Port:      8585,
		DataDir:   "/opt/velour",
		DBPath:    "/opt/velour/velour.db",
		JWTSecret: "",
	}
}

func Load() (*Config, error) {
	cfg := defaultConfig()

	if host := os.Getenv("VELOUR_HOST"); host != "" {
		cfg.Host = host
	}
	if port := os.Getenv("VELOUR_PORT"); port != "" {
		var p int
		if _, err := fmt.Sscanf(port, "%d", &p); err == nil {
			cfg.Port = p
		}
	}
	if dataDir := os.Getenv("VELOUR_DATA_DIR"); dataDir != "" {
		cfg.DataDir = dataDir
		cfg.DBPath = filepath.Join(dataDir, "velour.db")
	}
	if secret := os.Getenv("VELOUR_JWT_SECRET"); secret != "" {
		cfg.JWTSecret = secret
	}
	if mode := os.Getenv("VELOUR_INSTALL_MODE"); mode != "" {
		cfg.InstallMode = mode
	}

	configPath := filepath.Join(cfg.DataDir, "config.json")
	if data, err := os.ReadFile(configPath); err == nil {
		json.Unmarshal(data, cfg)
	}

	if cfg.JWTSecret == "" {
		cfg.JWTSecret = generateSecret()
	}

	return cfg, nil
}

func generateSecret() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "velour-default-secret-change-me"
	}
	return hex.EncodeToString(b)
}
