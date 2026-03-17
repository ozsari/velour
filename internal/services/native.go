package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/ozsari/velour/internal/models"
)

type NativeManager struct {
	dataDir     string
	appUsername string
	appPassword string
}

func NewNativeManager(dataDir string) *NativeManager {
	return &NativeManager{dataDir: dataDir}
}

func (nm *NativeManager) SetCredentials(username, password string) {
	nm.appUsername = username
	nm.appPassword = password
}

func (nm *NativeManager) Install(ctx context.Context, def *models.ServiceDefinition) error {
	if def.Native == nil {
		return fmt.Errorf("no native config for service %s", def.ID)
	}

	native := def.Native

	// Create user if specified
	if native.User != "" {
		if err := nm.ensureUser(native.User); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	// Create directories
	if native.ConfigDir != "" {
		dir := nm.expandPath(native.ConfigDir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create config dir: %w", err)
		}
		if native.User != "" {
			exec.CommandContext(ctx, "chown", "-R", native.User+":"+native.User, dir).Run()
		}
	}
	if native.DataDir != "" {
		dir := nm.expandPath(native.DataDir)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create data dir: %w", err)
		}
		if native.User != "" {
			exec.CommandContext(ctx, "chown", "-R", native.User+":"+native.User, dir).Run()
		}
	}

	// Install dependencies
	for _, dep := range native.Dependencies {
		if err := nm.aptInstall(ctx, dep); err != nil {
			return fmt.Errorf("failed to install dependency %s: %w", dep, err)
		}
	}

	// Install based on method
	switch native.Method {
	case "apt":
		if err := nm.installApt(ctx, native); err != nil {
			return err
		}
	case "binary":
		if err := nm.installBinary(ctx, native); err != nil {
			return err
		}
	case "script":
		if err := nm.installScript(ctx, native); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown install method: %s", native.Method)
	}

	// Run post-install commands (e.g. patch config to bind 0.0.0.0)
	if len(native.PostInstallCmds) > 0 {
		if err := nm.runPostInstallCmds(ctx, native); err != nil {
			log.Printf("Warning: post-install commands failed for %s: %v", native.ServiceName, err)
		}
	}

	return nil
}

func (nm *NativeManager) installApt(ctx context.Context, native *models.NativeConfig) error {
	// Add repo if needed
	if native.AptRepo != nil {
		if err := nm.addAptRepo(ctx, native.AptRepo); err != nil {
			return fmt.Errorf("failed to add apt repo: %w", err)
		}
	}

	// Fix any interrupted dpkg state and wait for apt lock
	fix := exec.CommandContext(ctx, "dpkg", "--configure", "-a")
	fix.Run()
	// Wait up to 60s for apt lock to be released
	for i := 0; i < 12; i++ {
		if _, err := os.Stat("/var/lib/dpkg/lock-frontend"); err != nil {
			break
		}
		cmd := exec.CommandContext(ctx, "fuser", "/var/lib/dpkg/lock-frontend")
		if err := cmd.Run(); err != nil {
			break // no process holding lock
		}
		log.Printf("Waiting for apt lock to be released... (%d/12)", i+1)
		time.Sleep(5 * time.Second)
	}

	// Install packages
	for _, pkg := range native.AptPackages {
		if err := nm.aptInstall(ctx, pkg); err != nil {
			return fmt.Errorf("failed to install %s: %w", pkg, err)
		}
	}

	// Install systemd unit if provided
	if native.ServiceUnit != "" {
		if err := nm.installSystemdUnit(native.ServiceName, native.ServiceUnit); err != nil {
			return err
		}
	}

	// Enable and start
	return nm.enableAndStart(ctx, native.ServiceName)
}

func (nm *NativeManager) installBinary(ctx context.Context, native *models.NativeConfig) error {
	url := nm.expandArch(native.BinaryURL)
	binPath := native.BinaryPath
	if binPath == "" {
		binPath = "/usr/local/bin/" + native.ServiceName
	}

	// If URL uses latest/download pattern, resolve via GitHub API to find versioned asset
	url, err := nm.resolveGitHubURL(url)
	if err != nil {
		return fmt.Errorf("failed to resolve download URL: %w", err)
	}

	// Download binary
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download binary: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Check if it's a tar archive
	if strings.HasSuffix(url, ".tar.gz") || strings.HasSuffix(url, ".tgz") || strings.HasSuffix(url, ".tar.xz") {
		tmpDir, err := os.MkdirTemp("", "velour-install-*")
		if err != nil {
			return err
		}
		defer os.RemoveAll(tmpDir)

		ext := "archive.tar.gz"
		tarFlag := "-xzf"
		if strings.HasSuffix(url, ".tar.xz") {
			ext = "archive.tar.xz"
			tarFlag = "-xJf"
		}
		tmpFile := filepath.Join(tmpDir, ext)
		f, err := os.Create(tmpFile)
		if err != nil {
			return err
		}
		if _, err := io.Copy(f, resp.Body); err != nil {
			f.Close()
			return err
		}
		f.Close()

		// Extract — ensure target dir exists
		extractDir := filepath.Dir(binPath)
		os.MkdirAll(extractDir, 0755)
		cmd := exec.CommandContext(ctx, "tar", tarFlag, tmpFile, "-C", extractDir)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("extract failed: %s: %w", string(out), err)
		}
	} else {
		// Direct binary
		f, err := os.OpenFile(binPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return fmt.Errorf("failed to create binary: %w", err)
		}
		if _, err := io.Copy(f, resp.Body); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}

	// Install systemd unit
	if native.ServiceUnit != "" {
		if err := nm.installSystemdUnit(native.ServiceName, native.ServiceUnit); err != nil {
			return err
		}
	}

	return nm.enableAndStart(ctx, native.ServiceName)
}

func (nm *NativeManager) installScript(ctx context.Context, native *models.NativeConfig) error {
	if native.InstallScript == "" {
		return fmt.Errorf("no install script provided")
	}

	cmd := exec.CommandContext(ctx, "bash", "-c", native.InstallScript)
	cmd.Env = append(os.Environ(),
		"DATA_DIR="+nm.dataDir,
		"VELOUR_DATA_DIR="+nm.dataDir,
		"VELOUR_SERVICE="+native.ServiceName,
	)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("install script failed: %s: %w", string(out), err)
	}

	if native.ServiceUnit != "" {
		log.Printf("Writing systemd unit for %s", native.ServiceName)
		if err := nm.installSystemdUnit(native.ServiceName, native.ServiceUnit); err != nil {
			return fmt.Errorf("failed to write systemd unit for %s: %w", native.ServiceName, err)
		}
		log.Printf("Systemd unit written successfully for %s", native.ServiceName)
	} else {
		log.Printf("WARNING: No ServiceUnit defined for %s, skipping unit file creation", native.ServiceName)
	}

	return nm.enableAndStart(ctx, native.ServiceName)
}

func (nm *NativeManager) Start(ctx context.Context, serviceID string) error {
	def := FindByID(serviceID)
	if def == nil || def.Native == nil {
		return fmt.Errorf("native config not found for %s", serviceID)
	}
	return nm.systemctl(ctx, "start", def.Native.ServiceName)
}

func (nm *NativeManager) Stop(ctx context.Context, serviceID string) error {
	def := FindByID(serviceID)
	if def == nil || def.Native == nil {
		return fmt.Errorf("native config not found for %s", serviceID)
	}
	return nm.systemctl(ctx, "stop", def.Native.ServiceName)
}

func (nm *NativeManager) Restart(ctx context.Context, serviceID string) error {
	def := FindByID(serviceID)
	if def == nil || def.Native == nil {
		return fmt.Errorf("native config not found for %s", serviceID)
	}
	return nm.systemctl(ctx, "restart", def.Native.ServiceName)
}

func (nm *NativeManager) Remove(ctx context.Context, serviceID string) error {
	def := FindByID(serviceID)
	if def == nil || def.Native == nil {
		return fmt.Errorf("native config not found for %s", serviceID)
	}
	native := def.Native

	// Stop and disable
	nm.systemctl(ctx, "stop", native.ServiceName)
	nm.systemctl(ctx, "disable", native.ServiceName)

	// Remove systemd unit
	unitPath := filepath.Join("/etc/systemd/system", native.ServiceName+".service")
	os.Remove(unitPath)
	nm.systemctl(ctx, "daemon-reload", "")

	// Remove packages if apt
	if native.Method == "apt" {
		for _, pkg := range native.AptPackages {
			exec.CommandContext(ctx, "apt-get", "remove", "-y", pkg).Run()
		}
	}

	// Remove binary if binary install
	if native.Method == "binary" && native.BinaryPath != "" {
		os.Remove(native.BinaryPath)
	}

	return nil
}

func (nm *NativeManager) Status(ctx context.Context, serviceID string) (models.ServiceStatus, error) {
	def := FindByID(serviceID)
	if def == nil || def.Native == nil {
		return models.StatusUnknown, fmt.Errorf("native config not found for %s", serviceID)
	}

	cmd := exec.CommandContext(ctx, "systemctl", "is-active", def.Native.ServiceName)
	out, err := cmd.Output()
	status := strings.TrimSpace(string(out))

	if err != nil || status != "active" {
		// Check if installed (handle template units like syncthing@velour)
		unitFile := def.Native.ServiceName + ".service"
		templateFile := ""
		if idx := strings.Index(def.Native.ServiceName, "@"); idx >= 0 {
			templateFile = def.Native.ServiceName[:idx+1] + ".service"
		}
		_, e1 := os.Stat(filepath.Join("/etc/systemd/system", unitFile))
		_, e2 := os.Stat(filepath.Join("/lib/systemd/system", unitFile))
		installed := e1 == nil || e2 == nil
		if !installed && templateFile != "" {
			_, e1 = os.Stat(filepath.Join("/etc/systemd/system", templateFile))
			_, e2 = os.Stat(filepath.Join("/lib/systemd/system", templateFile))
			installed = e1 == nil || e2 == nil
		}
		if installed {
			return models.StatusStopped, nil
		}
		return models.StatusUnknown, nil
	}

	return models.StatusRunning, nil
}

func (nm *NativeManager) ListManaged(ctx context.Context) ([]models.Service, error) {
	var services []models.Service

	for _, def := range Registry {
		if def.Native == nil {
			continue
		}

		// Check if installed by looking for systemd unit (apt packages use /lib, custom use /etc)
		// For template units like "syncthing@velour", check both the instance and the template
		unitFile := def.Native.ServiceName + ".service"
		templateFile := ""
		if idx := strings.Index(def.Native.ServiceName, "@"); idx >= 0 {
			templateFile = def.Native.ServiceName[:idx+1] + ".service"
		}
		_, err1 := os.Stat(filepath.Join("/etc/systemd/system", unitFile))
		_, err2 := os.Stat(filepath.Join("/lib/systemd/system", unitFile))
		_, err3 := os.Stat(filepath.Join("/usr/lib/systemd/system", unitFile))
		found := err1 == nil || err2 == nil || err3 == nil
		if !found && templateFile != "" {
			_, err1 = os.Stat(filepath.Join("/etc/systemd/system", templateFile))
			_, err2 = os.Stat(filepath.Join("/lib/systemd/system", templateFile))
			_, err3 = os.Stat(filepath.Join("/usr/lib/systemd/system", templateFile))
			found = err1 == nil || err2 == nil || err3 == nil
		}
		if !found {
			continue
		}

		status, _ := nm.Status(ctx, def.ID)
		port := def.Native.Port
		webURL := ""
		if port > 0 {
			webURL = fmt.Sprintf("http://localhost:%d", port)
		}

		services = append(services, models.Service{
			ID:          def.ID,
			Name:        def.Name,
			Description: def.Description,
			Icon:        def.Icon,
			Category:    def.Category,
			Port:        port,
			WebURL:      webURL,
			Status:      status,
			Type:        "native",
			Installed:   true,
		})
	}

	return services, nil
}

// Exec runs a command directly on the host (native mode equivalent of docker exec)
func (nm *NativeManager) Exec(ctx context.Context, cmd []string) (string, error) {
	if len(cmd) == 0 {
		return "", fmt.Errorf("empty command")
	}
	c := exec.CommandContext(ctx, cmd[0], cmd[1:]...)
	out, err := c.CombinedOutput()
	return string(out), err
}

// Helper methods

// runPostInstallCmds waits for the service to generate its config, then stops it,
// runs the post-install commands (e.g. patching bind address), and restarts.
// If PostStartCmds are defined, it then waits for the service port and runs those.
func (nm *NativeManager) runPostInstallCmds(ctx context.Context, native *models.NativeConfig) error {
	// Give the service time to create its config files on first start
	time.Sleep(3 * time.Second)

	// Stop service before patching config
	nm.systemctl(ctx, "stop", native.ServiceName)

	for i, cmdStr := range native.PostInstallCmds {
		cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)
		cmd.Env = append(os.Environ(),
			"DATA_DIR="+nm.dataDir,
			"VELOUR_DATA_DIR="+nm.dataDir,
			"VELOUR_USER="+nm.appUsername,
			"VELOUR_PASS="+nm.appPassword,
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("post-install cmd [%d] failed: %s\noutput: %s\nerror: %v", i, cmdStr, string(out), err)
		} else {
			log.Printf("post-install cmd [%d] ok: %s", i, cmdStr[:min(len(cmdStr), 80)])
		}
	}

	// Restart with patched config
	if err := nm.systemctl(ctx, "start", native.ServiceName); err != nil {
		return err
	}

	// Run post-start commands (require running service, e.g. API calls)
	if len(native.PostStartCmds) > 0 && native.Port > 0 {
		// Wait for the service port to be ready (up to 30s)
		for i := 0; i < 30; i++ {
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", native.Port), time.Second)
			if err == nil {
				conn.Close()
				break
			}
			time.Sleep(time.Second)
		}
		// Small extra delay for the app to fully initialize
		time.Sleep(2 * time.Second)

		for i, cmdStr := range native.PostStartCmds {
			cmd := exec.CommandContext(ctx, "bash", "-c", cmdStr)
			cmd.Env = append(os.Environ(),
				"DATA_DIR="+nm.dataDir,
				"VELOUR_DATA_DIR="+nm.dataDir,
				"VELOUR_USER="+nm.appUsername,
				"VELOUR_PASS="+nm.appPassword,
			)
			out, err := cmd.CombinedOutput()
			if err != nil {
				log.Printf("post-start cmd [%d] failed: %s\noutput: %s\nerror: %v", i, cmdStr, string(out), err)
			} else {
				log.Printf("post-start cmd [%d] ok: %s", i, cmdStr[:min(len(cmdStr), 80)])
			}
		}
	}

	return nil
}

func (nm *NativeManager) expandPath(path string) string {
	return strings.ReplaceAll(path, "${DATA_DIR}", nm.dataDir)
}

// resolveGitHubURL resolves GitHub latest/download URLs to actual versioned asset URLs.
// Many projects use versioned filenames (e.g. autobrr_1.74.0_linux_x86_64.tar.gz),
// so latest/download/autobrr_linux_x86_64.tar.gz returns 404.
// This method uses the GitHub API to find the matching asset.
func (nm *NativeManager) resolveGitHubURL(url string) (string, error) {
	re := regexp.MustCompile(`^https://github\.com/([^/]+/[^/]+)/releases/latest/download/(.+)$`)
	matches := re.FindStringSubmatch(url)
	if matches == nil {
		return url, nil // not a GitHub latest URL, use as-is
	}

	repo := matches[1]
	filename := matches[2]

	// Try direct URL first (works for projects with stable filenames)
	resp, err := http.Head(url)
	if err == nil && resp.StatusCode == 200 {
		return url, nil
	}

	// Fetch latest release assets from GitHub API
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	resp, err = http.Get(apiURL)
	if err != nil {
		return url, nil // fallback to original URL
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return url, nil
	}

	var release struct {
		Assets []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return url, nil
	}

	// Strip version-agnostic parts from filename to build a match pattern
	// e.g. "autobrr_linux_x86_64.tar.gz" should match "autobrr_1.74.0_linux_x86_64.tar.gz"
	parts := strings.SplitN(filename, "_", 2) // ["autobrr", "linux_x86_64.tar.gz"]
	if len(parts) < 2 {
		return url, nil
	}
	prefix := parts[0] + "_"
	suffix := "_" + parts[1]

	for _, asset := range release.Assets {
		if strings.HasPrefix(asset.Name, prefix) && strings.HasSuffix(asset.Name, suffix) {
			return asset.BrowserDownloadURL, nil
		}
	}

	return url, nil // no match found, try original
}

func (nm *NativeManager) expandArch(url string) string {
	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		arch = "x86_64"
	case "arm64":
		arch = "aarch64"
	}
	return strings.ReplaceAll(url, "${ARCH}", arch)
}

func (nm *NativeManager) ensureUser(username string) error {
	if err := exec.Command("id", username).Run(); err == nil {
		return nil // user exists
	}
	cmd := exec.Command("useradd", "-r", "-s", "/usr/sbin/nologin", "-d", "/var/lib/"+username, "-m", username)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %w", string(out), err)
	}
	return nil
}

func (nm *NativeManager) aptInstall(ctx context.Context, pkg string) error {
	cmd := exec.CommandContext(ctx, "apt-get", "install", "-y", pkg)
	cmd.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s: %w", string(out), err)
	}
	return nil
}

func (nm *NativeManager) addAptRepo(ctx context.Context, repo *models.AptRepo) error {
	// Add GPG key via URL or keyserver
	if repo.KeyURL != "" {
		// Extract keyring path from signed-by in RepoLine, or generate one
		keyringPath := ""
		if idx := strings.Index(repo.RepoLine, "signed-by="); idx >= 0 {
			rest := repo.RepoLine[idx+len("signed-by="):]
			if end := strings.IndexByte(rest, ']'); end >= 0 {
				keyringPath = rest[:end]
			}
		}
		if keyringPath == "" {
			keyringPath = fmt.Sprintf("/usr/share/keyrings/velour-%s.gpg",
				strings.ReplaceAll(filepath.Base(repo.KeyURL), ".", "-"))
		}
		cmd := exec.CommandContext(ctx, "bash", "-c",
			fmt.Sprintf("curl -fsSL '%s' | gpg --batch --yes --dearmor -o '%s'",
				repo.KeyURL, keyringPath))
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to add key: %s: %w", string(out), err)
		}
	} else if repo.KeyServer != "" && repo.KeyID != "" {
		// Extract domain from repo line for keyring filename
		keyringName := "keyserver"
		for _, field := range strings.Fields(repo.RepoLine) {
			if strings.HasPrefix(field, "http") {
				keyringName = strings.ReplaceAll(strings.ReplaceAll(
					strings.Split(field, "/")[2], ".", "-"), ":", "")
				break
			}
		}
		keyringPath := fmt.Sprintf("/usr/share/keyrings/velour-%s.gpg", keyringName)
		cmd := exec.CommandContext(ctx, "gpg", "--batch", "--yes",
			"--no-default-keyring", "--keyring", keyringPath,
			"--keyserver", repo.KeyServer, "--recv-keys", repo.KeyID)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to add key from keyserver: %s: %w", string(out), err)
		}
	}

	// Add repo — extract domain from repo line for filename
	repoName := "repo"
	for _, field := range strings.Fields(repo.RepoLine) {
		if strings.HasPrefix(field, "http") {
			repoName = strings.ReplaceAll(strings.ReplaceAll(
				strings.Split(field, "/")[2], ".", "-"), ":", "")
			break
		}
	}
	repoFile := fmt.Sprintf("/etc/apt/sources.list.d/velour-%s.list", repoName)
	if err := os.WriteFile(repoFile, []byte(repo.RepoLine+"\n"), 0644); err != nil {
		return err
	}

	// Update
	cmd := exec.CommandContext(ctx, "apt-get", "update")
	cmd.Env = append(os.Environ(), "DEBIAN_FRONTEND=noninteractive")
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("apt update failed: %s: %w", string(out), err)
	}

	return nil
}

func (nm *NativeManager) installSystemdUnit(name, content string) error {
	unitPath := filepath.Join("/etc/systemd/system", name+".service")
	return os.WriteFile(unitPath, []byte(content), 0644)
}

func (nm *NativeManager) enableAndStart(ctx context.Context, serviceName string) error {
	if err := nm.systemctl(ctx, "daemon-reload", ""); err != nil {
		return err
	}
	if err := nm.systemctl(ctx, "enable", serviceName); err != nil {
		return err
	}
	return nm.systemctl(ctx, "start", serviceName)
}

func (nm *NativeManager) systemctl(ctx context.Context, action, service string) error {
	var cmd *exec.Cmd
	if service == "" {
		cmd = exec.CommandContext(ctx, "systemctl", action)
	} else {
		cmd = exec.CommandContext(ctx, "systemctl", action, service)
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("systemctl %s %s failed: %s: %w", action, service, string(out), err)
	}
	return nil
}
