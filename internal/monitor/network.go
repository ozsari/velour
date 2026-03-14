package monitor

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"github.com/ozsari/velour/internal/models"
	"github.com/shirou/gopsutil/v3/net"
)

// NetworkTracker persists network usage per month in SQLite.
// It snapshots cumulative OS counters every 60 seconds and accumulates
// deltas into the current month's row.
type NetworkTracker struct {
	db   *sql.DB
	mu   sync.Mutex
	prev struct {
		sent uint64
		recv uint64
		set  bool
	}
	stopCh chan struct{}
}

func NewNetworkTracker(db *sql.DB) *NetworkTracker {
	return &NetworkTracker{
		db:     db,
		stopCh: make(chan struct{}),
	}
}

// InitDB creates the monthly_network_stats table
func (nt *NetworkTracker) InitDB() error {
	_, err := nt.db.Exec(`
		CREATE TABLE IF NOT EXISTS monthly_network_stats (
			month TEXT PRIMARY KEY,
			bytes_sent INTEGER NOT NULL DEFAULT 0,
			bytes_recv INTEGER NOT NULL DEFAULT 0
		)
	`)
	return err
}

// Start begins the background snapshot loop (every 60s)
func (nt *NetworkTracker) Start() {
	go nt.loop()
}

// Stop halts the snapshot loop
func (nt *NetworkTracker) Stop() {
	close(nt.stopCh)
}

func (nt *NetworkTracker) loop() {
	// Take an initial reading so the first tick has a baseline
	nt.snapshot()

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-nt.stopCh:
			return
		case <-ticker.C:
			nt.snapshot()
		}
	}
}

func (nt *NetworkTracker) snapshot() {
	counters, err := net.IOCounters(false)
	if err != nil || len(counters) == 0 {
		return
	}

	curSent := counters[0].BytesSent
	curRecv := counters[0].BytesRecv

	nt.mu.Lock()
	defer nt.mu.Unlock()

	if !nt.prev.set {
		// First reading - just store baseline, no delta yet
		nt.prev.sent = curSent
		nt.prev.recv = curRecv
		nt.prev.set = true
		return
	}

	// Calculate delta since last snapshot
	deltaSent := uint64(0)
	deltaRecv := uint64(0)
	if curSent >= nt.prev.sent {
		deltaSent = curSent - nt.prev.sent
	}
	if curRecv >= nt.prev.recv {
		deltaRecv = curRecv - nt.prev.recv
	}

	nt.prev.sent = curSent
	nt.prev.recv = curRecv

	if deltaSent == 0 && deltaRecv == 0 {
		return
	}

	month := time.Now().Format("2006-01") // e.g. "2026-03"

	_, err = nt.db.Exec(`
		INSERT INTO monthly_network_stats (month, bytes_sent, bytes_recv)
		VALUES (?, ?, ?)
		ON CONFLICT(month) DO UPDATE SET
			bytes_sent = bytes_sent + excluded.bytes_sent,
			bytes_recv = bytes_recv + excluded.bytes_recv
	`, month, deltaSent, deltaRecv)
	if err != nil {
		log.Printf("[network-tracker] Failed to persist stats: %v", err)
	}
}

// GetMonth returns stats for a specific month (e.g. "2026-03")
func (nt *NetworkTracker) GetMonth(month string) (*models.MonthlyNetStats, error) {
	var stats models.MonthlyNetStats
	err := nt.db.QueryRow(
		`SELECT month, bytes_sent, bytes_recv FROM monthly_network_stats WHERE month = ?`,
		month,
	).Scan(&stats.Month, &stats.BytesSent, &stats.BytesRecv)
	if err == sql.ErrNoRows {
		return &models.MonthlyNetStats{Month: month, BytesSent: 0, BytesRecv: 0}, nil
	}
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

// ListMonths returns all months that have recorded data, newest first
func (nt *NetworkTracker) ListMonths() ([]models.MonthlyNetStats, error) {
	rows, err := nt.db.Query(
		`SELECT month, bytes_sent, bytes_recv FROM monthly_network_stats ORDER BY month DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.MonthlyNetStats
	for rows.Next() {
		var s models.MonthlyNetStats
		if err := rows.Scan(&s.Month, &s.BytesSent, &s.BytesRecv); err != nil {
			return nil, err
		}
		result = append(result, s)
	}
	return result, nil
}
