package models

type SystemInfo struct {
	Hostname    string  `json:"hostname"`
	OS          string  `json:"os"`
	Platform    string  `json:"platform"`
	Kernel      string  `json:"kernel"`
	Uptime      uint64  `json:"uptime"`
	UptimeHuman string  `json:"uptime_human"`
	CPU         CPUInfo `json:"cpu"`
	Memory      MemInfo `json:"memory"`
	Disk        DiskInfo `json:"disk"`
	Network     NetInfo  `json:"network"`
}

type CPUInfo struct {
	Model   string  `json:"model"`
	Cores   int     `json:"cores"`
	Threads int     `json:"threads"`
	Usage   float64 `json:"usage"`
}

type MemInfo struct {
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Free      uint64  `json:"free"`
	UsagePerc float64 `json:"usage_percent"`
}

type DiskInfo struct {
	Total     uint64  `json:"total"`
	Used      uint64  `json:"used"`
	Free      uint64  `json:"free"`
	UsagePerc float64 `json:"usage_percent"`
}

type NetInfo struct {
	BytesSent uint64 `json:"bytes_sent"`
	BytesRecv uint64 `json:"bytes_recv"`
}

// MonthlyNetStats holds network usage for a specific month
type MonthlyNetStats struct {
	Month     string `json:"month"`      // "2026-03"
	BytesSent uint64 `json:"bytes_sent"`
	BytesRecv uint64 `json:"bytes_recv"`
}
