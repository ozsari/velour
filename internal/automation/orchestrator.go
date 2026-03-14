package automation

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ozsari/velour/internal/models"
	"github.com/ozsari/velour/internal/services"
)

// Orchestrator manages automation rules and watches for trigger events.
// In Docker mode it uses `docker exec` to run commands across containers.
// In Native mode it runs commands directly on the host.
type Orchestrator struct {
	db          *sql.DB
	docker      *services.DockerManager
	native      *services.NativeManager
	installMode string
	rules       []models.AutomationRule
	mu          sync.RWMutex
	stopCh      chan struct{}
}

func New(db *sql.DB, docker *services.DockerManager, native *services.NativeManager, installMode string) *Orchestrator {
	return &Orchestrator{
		db:          db,
		docker:      docker,
		native:      native,
		installMode: installMode,
		stopCh:      make(chan struct{}),
	}
}

// InitDB creates the automation_rules table
func (o *Orchestrator) InitDB() error {
	_, err := o.db.Exec(`
		CREATE TABLE IF NOT EXISTS automation_rules (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			enabled INTEGER NOT NULL DEFAULT 1,
			trigger_json TEXT NOT NULL,
			action_json TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			last_run_at DATETIME,
			last_run_ok INTEGER,
			last_run_log TEXT DEFAULT '',
			run_count INTEGER NOT NULL DEFAULT 0
		)
	`)
	return err
}

// Start begins the polling loop for trigger events
func (o *Orchestrator) Start() {
	go o.pollLoop()
}

// Stop halts the polling loop
func (o *Orchestrator) Stop() {
	close(o.stopCh)
}

// pollLoop checks torrent clients for completed downloads every 15 seconds
func (o *Orchestrator) pollLoop() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	// Track already-seen completed torrents to avoid re-firing
	seen := make(map[string]map[string]bool) // serviceID -> torrentHash -> seen

	for {
		select {
		case <-o.stopCh:
			return
		case <-ticker.C:
			rules, err := o.ListRules()
			if err != nil {
				log.Printf("[automation] Failed to load rules: %v", err)
				continue
			}

			for _, rule := range rules {
				if !rule.Enabled {
					continue
				}
				if rule.Trigger.Type != models.TriggerTorrentDone {
					continue
				}

				serviceID := rule.Trigger.ServiceID
				if seen[serviceID] == nil {
					seen[serviceID] = make(map[string]bool)
				}

				completed, err := o.getCompletedTorrents(serviceID)
				if err != nil {
					log.Printf("[automation] Failed to check %s: %v", serviceID, err)
					continue
				}

				for _, torrent := range completed {
					if seen[serviceID][torrent.Hash] {
						continue
					}
					seen[serviceID][torrent.Hash] = true

					log.Printf("[automation] Trigger fired: %s completed in %s", torrent.Name, serviceID)
					go o.executeAction(rule, torrent)
				}
			}
		}
	}
}

// CompletedTorrent holds info about a finished download
type CompletedTorrent struct {
	Hash     string
	Name     string
	SavePath string
	Category string
}

// getCompletedTorrents queries a torrent client's API for finished downloads
func (o *Orchestrator) getCompletedTorrents(serviceID string) ([]CompletedTorrent, error) {
	// In a real implementation, this would query the torrent client's API:
	// - qBittorrent: GET /api/v2/torrents/info?filter=completed
	// - Deluge: JSON-RPC core.get_torrents_status
	// - Transmission: RPC method "torrent-get"
	// - rTorrent: XMLRPC d.multicall2
	//
	// For now, return empty - the actual API polling will be implemented
	// when we add the torrent client API integrations.
	return nil, nil
}

// executeAction runs the action defined in a rule
func (o *Orchestrator) executeAction(rule models.AutomationRule, torrent CompletedTorrent) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	var output string
	var err error

	switch rule.Action.Type {
	case models.ActionExecInService:
		cmd := o.buildCommand(rule.Action, torrent)
		output, err = o.execInService(ctx, rule.Action.ServiceID, cmd)
	case models.ActionRestartService:
		err = o.restartService(ctx, rule.Action.ServiceID)
		output = "service restarted"
	}

	success := err == nil
	logMsg := output
	if err != nil {
		logMsg = fmt.Sprintf("ERROR: %v\n%s", err, output)
	}

	now := time.Now()
	o.db.Exec(`
		UPDATE automation_rules
		SET last_run_at = ?, last_run_ok = ?, last_run_log = ?, run_count = run_count + 1, updated_at = ?
		WHERE id = ?
	`, now, success, logMsg, now, rule.ID)
}

// buildCommand constructs the command to execute, substituting torrent variables
func (o *Orchestrator) buildCommand(action models.Action, torrent CompletedTorrent) []string {
	// If using a template, expand it
	if action.Template == "filebot_amc" {
		return []string{
			"filebot", "-script", "fn:amc",
			"--output", "/data/media",
			"--action", "hardlink",
			"-non-strict",
			"--def",
			fmt.Sprintf("ut_dir=%s", torrent.SavePath),
			fmt.Sprintf("ut_title=%s", torrent.Name),
			fmt.Sprintf("ut_label=%s", torrent.Category),
		}
	}

	// Custom command - substitute variables
	cmd := make([]string, 0, len(action.Args)+1)
	cmd = append(cmd, action.Command)
	for _, arg := range action.Args {
		expanded := arg
		// Replace template variables
		for _, pair := range []struct{ k, v string }{
			{"{torrent_path}", torrent.SavePath},
			{"{torrent_name}", torrent.Name},
			{"{torrent_category}", torrent.Category},
			{"{torrent_hash}", torrent.Hash},
		} {
			if expanded == pair.k {
				expanded = pair.v
			}
		}
		cmd = append(cmd, expanded)
	}
	return cmd
}

// execInService runs a command inside a service (docker exec or native shell)
func (o *Orchestrator) execInService(ctx context.Context, serviceID string, cmd []string) (string, error) {
	if o.installMode == "docker" {
		if o.docker == nil {
			return "", fmt.Errorf("docker not available")
		}
		containerName := "velour-" + serviceID
		return o.docker.Exec(ctx, containerName, cmd)
	}
	// Native mode: run command directly
	if o.native == nil {
		return "", fmt.Errorf("native manager not available")
	}
	return o.native.Exec(ctx, cmd)
}

// restartService restarts a managed service
func (o *Orchestrator) restartService(ctx context.Context, serviceID string) error {
	if o.installMode == "docker" {
		if o.docker != nil {
			return o.docker.Restart(ctx, serviceID)
		}
	} else {
		if o.native != nil {
			return o.native.Restart(ctx, serviceID)
		}
	}
	return fmt.Errorf("no manager available")
}

// ── CRUD Operations ──

func (o *Orchestrator) CreateRule(rule models.AutomationRule) error {
	triggerJSON, _ := json.Marshal(rule.Trigger)
	actionJSON, _ := json.Marshal(rule.Action)

	_, err := o.db.Exec(`
		INSERT INTO automation_rules (id, name, enabled, trigger_json, action_json, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, rule.ID, rule.Name, rule.Enabled, string(triggerJSON), string(actionJSON), rule.CreatedAt, rule.UpdatedAt)
	return err
}

func (o *Orchestrator) ListRules() ([]models.AutomationRule, error) {
	rows, err := o.db.Query(`
		SELECT id, name, enabled, trigger_json, action_json, created_at, updated_at,
		       last_run_at, last_run_ok, last_run_log, run_count
		FROM automation_rules ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []models.AutomationRule
	for rows.Next() {
		var r models.AutomationRule
		var triggerJSON, actionJSON string
		var lastRunAt sql.NullTime
		var lastRunOK sql.NullBool

		err := rows.Scan(&r.ID, &r.Name, &r.Enabled, &triggerJSON, &actionJSON,
			&r.CreatedAt, &r.UpdatedAt, &lastRunAt, &lastRunOK, &r.LastRunLog, &r.RunCount)
		if err != nil {
			return nil, err
		}

		json.Unmarshal([]byte(triggerJSON), &r.Trigger)
		json.Unmarshal([]byte(actionJSON), &r.Action)

		if lastRunAt.Valid {
			r.LastRunAt = &lastRunAt.Time
		}
		if lastRunOK.Valid {
			v := lastRunOK.Bool
			r.LastRunOK = &v
		}

		rules = append(rules, r)
	}
	return rules, nil
}

func (o *Orchestrator) GetRule(id string) (*models.AutomationRule, error) {
	rules, err := o.ListRules()
	if err != nil {
		return nil, err
	}
	for _, r := range rules {
		if r.ID == id {
			return &r, nil
		}
	}
	return nil, fmt.Errorf("rule not found")
}

func (o *Orchestrator) UpdateRule(rule models.AutomationRule) error {
	triggerJSON, _ := json.Marshal(rule.Trigger)
	actionJSON, _ := json.Marshal(rule.Action)

	_, err := o.db.Exec(`
		UPDATE automation_rules SET name=?, enabled=?, trigger_json=?, action_json=?, updated_at=?
		WHERE id=?
	`, rule.Name, rule.Enabled, string(triggerJSON), string(actionJSON), time.Now(), rule.ID)
	return err
}

func (o *Orchestrator) DeleteRule(id string) error {
	_, err := o.db.Exec(`DELETE FROM automation_rules WHERE id=?`, id)
	return err
}

func (o *Orchestrator) ToggleRule(id string, enabled bool) error {
	_, err := o.db.Exec(`UPDATE automation_rules SET enabled=?, updated_at=? WHERE id=?`, enabled, time.Now(), id)
	return err
}

// GetTemplates returns pre-configured action templates
func (o *Orchestrator) GetTemplates() []models.ActionTemplate {
	return []models.ActionTemplate{
		{
			ID:          "filebot_amc",
			Name:        "FileBot AMC",
			Description: "Automatically rename & organize media using FileBot's AMC script",
			ServiceID:   "filebot",
			Command:     "filebot",
			Args:        []string{"-script", "fn:amc", "--output", "/data/media", "--action", "hardlink", "-non-strict", "--def", "ut_dir={torrent_path}", "ut_title={torrent_name}", "ut_label={torrent_category}"},
			Icon:        "🎬",
		},
		{
			ID:          "unpackerr",
			Name:        "Unpackerr Extract",
			Description: "Extract compressed downloads automatically",
			ServiceID:   "unpackerr",
			Command:     "unpackerr",
			Args:        []string{"--path", "{torrent_path}"},
			Icon:        "📦",
		},
		{
			ID:          "rclone_sync",
			Name:        "Rclone Sync",
			Description: "Sync a local folder to a cloud remote (Google Drive, OneDrive, etc.)",
			ServiceID:   "rclone",
			Command:     "rclone",
			Args:        []string{"sync", "/data/media", "remote:", "--progress", "--transfers", "4"},
			Icon:        "☁️",
		},
		{
			ID:          "rclone_move",
			Name:        "Rclone Move",
			Description: "Move files to cloud storage and delete local copies after upload",
			ServiceID:   "rclone",
			Command:     "rclone",
			Args:        []string{"move", "{torrent_path}", "remote:uploads/", "--progress"},
			Icon:        "📤",
		},
		{
			ID:          "custom_script",
			Name:        "Custom Script",
			Description: "Run a custom command in any service",
			ServiceID:   "",
			Command:     "",
			Args:        []string{},
			Icon:        "⚡",
		},
	}
}
