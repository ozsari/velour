package monitor

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ozsari/velour/internal/models"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type Monitor struct{}

func New() *Monitor {
	return &Monitor{}
}

func (m *Monitor) GetSystemInfo() (*models.SystemInfo, error) {
	hostInfo, err := host.Info()
	if err != nil {
		return nil, err
	}

	cpuInfo, err := m.getCPUInfo()
	if err != nil {
		return nil, err
	}

	memInfo, err := m.getMemInfo()
	if err != nil {
		return nil, err
	}

	diskInfo, err := m.getDiskInfo()
	if err != nil {
		return nil, err
	}

	netInfo, err := m.getNetInfo()
	if err != nil {
		return nil, err
	}

	hostname, _ := os.Hostname()

	return &models.SystemInfo{
		Hostname:    hostname,
		OS:          fmt.Sprintf("%s %s", strings.Title(hostInfo.Platform), hostInfo.PlatformVersion),
		Platform:    runtime.GOARCH,
		Kernel:      hostInfo.KernelVersion,
		Uptime:      hostInfo.Uptime,
		UptimeHuman: formatUptime(hostInfo.Uptime),
		CPU:         *cpuInfo,
		Memory:      *memInfo,
		Disk:        *diskInfo,
		Network:     *netInfo,
	}, nil
}

func (m *Monitor) getCPUInfo() (*models.CPUInfo, error) {
	info, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	percentages, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, err
	}

	modelName := "Unknown"
	if len(info) > 0 {
		modelName = info[0].ModelName
	}

	usage := 0.0
	if len(percentages) > 0 {
		usage = percentages[0]
	}

	return &models.CPUInfo{
		Model:   modelName,
		Cores:   runtime.NumCPU(),
		Threads: runtime.NumCPU(),
		Usage:   usage,
	}, nil
}

func (m *Monitor) getMemInfo() (*models.MemInfo, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &models.MemInfo{
		Total:     v.Total,
		Used:      v.Used,
		Free:      v.Free,
		UsagePerc: v.UsedPercent,
	}, nil
}

func (m *Monitor) getDiskInfo() (*models.DiskInfo, error) {
	d, err := disk.Usage("/")
	if err != nil {
		return nil, err
	}

	return &models.DiskInfo{
		Total:     d.Total,
		Used:      d.Used,
		Free:      d.Free,
		UsagePerc: d.UsedPercent,
	}, nil
}

func (m *Monitor) getNetInfo() (*models.NetInfo, error) {
	counters, err := net.IOCounters(false)
	if err != nil {
		return nil, err
	}

	if len(counters) > 0 {
		return &models.NetInfo{
			BytesSent: counters[0].BytesSent,
			BytesRecv: counters[0].BytesRecv,
		}, nil
	}

	return &models.NetInfo{}, nil
}

func formatUptime(seconds uint64) string {
	days := seconds / 86400
	hours := (seconds % 86400) / 3600
	minutes := (seconds % 3600) / 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, minutes)
	}
	return fmt.Sprintf("%dm", minutes)
}
