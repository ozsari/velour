package models

import "time"

// AutomationRule defines a post-processing rule
// e.g. "When qBittorrent finishes a download → run FileBot AMC"
type AutomationRule struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Enabled     bool      `json:"enabled"`
	Trigger     Trigger   `json:"trigger"`
	Action      Action    `json:"action"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	LastRunAt   *time.Time `json:"last_run_at,omitempty"`
	LastRunOK   *bool     `json:"last_run_ok,omitempty"`
	LastRunLog  string    `json:"last_run_log,omitempty"`
	RunCount    int       `json:"run_count"`
}

// TriggerType defines what kind of event starts the rule
type TriggerType string

const (
	TriggerTorrentDone   TriggerType = "torrent_done"     // torrent finished downloading
	TriggerServiceStart  TriggerType = "service_start"    // a service started
	TriggerServiceStop   TriggerType = "service_stop"     // a service stopped
	TriggerSchedule      TriggerType = "schedule"         // cron-like schedule
)

type Trigger struct {
	Type      TriggerType `json:"type"`
	ServiceID string      `json:"service_id"`           // which service triggers this
	// For schedule triggers
	Cron      string      `json:"cron,omitempty"`       // cron expression e.g. "*/5 * * * *"
}

// ActionType defines what to do when triggered
type ActionType string

const (
	ActionExecInService ActionType = "exec_in_service"  // run command inside another service
	ActionWebhook       ActionType = "webhook"          // call a URL
	ActionRestartService ActionType = "restart_service" // restart a service
)

type Action struct {
	Type      ActionType `json:"type"`
	ServiceID string     `json:"service_id,omitempty"`  // target service
	Command   string     `json:"command,omitempty"`     // command to run
	Args      []string   `json:"args,omitempty"`        // command arguments
	// For webhook actions
	URL       string     `json:"url,omitempty"`
	// Pre-built templates
	Template  string     `json:"template,omitempty"`    // e.g. "filebot_amc"
}

// ActionTemplate is a pre-configured action users can pick
type ActionTemplate struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	ServiceID   string   `json:"service_id"`   // which service this runs in
	Command     string   `json:"command"`
	Args        []string `json:"args"`
	Icon        string   `json:"icon"`
}

// CreateRuleRequest is the API request to create a rule
type CreateRuleRequest struct {
	Name    string  `json:"name"`
	Trigger Trigger `json:"trigger"`
	Action  Action  `json:"action"`
}
