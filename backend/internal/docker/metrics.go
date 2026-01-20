package docker

import (
	"time"

	"github.com/docker/docker/api/types/container"
	"go.uber.org/zap"
)

// ContainerStats represents processed container statistics
type ContainerStats struct {
	ContainerID   string    `json:"container_id"`
	Timestamp     time.Time `json:"timestamp"`
	CPUPercent    float64   `json:"cpu_percent"`
	MemoryUsage   uint64    `json:"memory_usage"`
	MemoryLimit   uint64    `json:"memory_limit"`
	MemoryPercent float64   `json:"memory_percent"`
	NetworkRx     uint64    `json:"network_rx"`
	NetworkTx     uint64    `json:"network_tx"`
	BlockRead     uint64    `json:"block_read"`
	BlockWrite    uint64    `json:"block_write"`
	PIDs          uint64    `json:"pids"`
}

// StatsCalculator handles container statistics calculation
type StatsCalculator struct {
	previousStats map[string]*container.StatsResponse
	logger        *zap.Logger
}

// NewStatsCalculator creates a new stats calculator
func NewStatsCalculator(logger *zap.Logger) *StatsCalculator {
	return &StatsCalculator{
		previousStats: make(map[string]*container.StatsResponse),
		logger:        logger,
	}
}

// CalculateStats processes raw Docker stats and calculates percentages
func (sc *StatsCalculator) CalculateStats(containerID string, stats *container.StatsResponse) (*ContainerStats, error) {
	prevStats, hasPrevious := sc.previousStats[containerID]

	// Calculate CPU percentage
	var cpuPercent float64
	if hasPrevious {
		cpuPercent = sc.calculateCPUPercent(prevStats, stats)
	} else {
		cpuPercent = 0.0
	}

	// Calculate memory
	memoryUsage := stats.MemoryStats.Usage
	memoryLimit := stats.MemoryStats.Limit
	memoryPercent := 0.0
	if memoryLimit > 0 {
		// Exclude cache for accurate RSS
		cache := uint64(0)
		if stats.MemoryStats.Stats != nil {
			// Stats is a map[string]uint64, access cache by key
			if cacheVal, ok := stats.MemoryStats.Stats["cache"]; ok {
				cache = cacheVal
			}
		}
		rss := memoryUsage - cache
		memoryPercent = float64(rss) / float64(memoryLimit) * 100.0
	}

	// Network stats
	var networkRx, networkTx uint64
	if stats.Networks != nil {
		for _, network := range stats.Networks {
			networkRx += network.RxBytes
			networkTx += network.TxBytes
		}
	}

	// Block I/O stats
	var blockRead, blockWrite uint64
	if stats.BlkioStats.IoServiceBytesRecursive != nil {
		for _, entry := range stats.BlkioStats.IoServiceBytesRecursive {
			switch entry.Op {
			case "Read":
				blockRead += entry.Value
			case "Write":
				blockWrite += entry.Value
			}
		}
	}

	result := &ContainerStats{
		ContainerID:   containerID,
		Timestamp:     time.Now(),
		CPUPercent:    cpuPercent,
		MemoryUsage:   memoryUsage,
		MemoryLimit:   memoryLimit,
		MemoryPercent: memoryPercent,
		NetworkRx:     networkRx,
		NetworkTx:     networkTx,
		BlockRead:     blockRead,
		BlockWrite:    blockWrite,
		PIDs:          stats.PidsStats.Current,
	}

	// Store current stats for next calculation
	sc.previousStats[containerID] = stats

	return result, nil
}

// calculateCPUPercent calculates CPU usage percentage
// Formula: CPU% = (ΔTotalUsage / ΔSystemUsage) * OnlineCPUs * 100
func (sc *StatsCalculator) calculateCPUPercent(prevStats, currStats *container.StatsResponse) float64 {
	var (
		prevCPU    = prevStats.CPUStats.CPUUsage.TotalUsage
		prevSystem = prevStats.CPUStats.SystemUsage
		currCPU    = currStats.CPUStats.CPUUsage.TotalUsage
		currSystem = currStats.CPUStats.SystemUsage
	)

	// Calculate deltas
	// Note: deltaCPU can be negative if counter wraps, but since it's uint64,
	// we check if currCPU < prevCPU instead (which indicates wrap-around)
	deltaCPU := currCPU - prevCPU
	deltaSystem := currSystem - prevSystem

	// Handle edge cases
	if deltaSystem == 0 || currCPU < prevCPU {
		return 0.0
	}

	// Get number of online CPUs
	onlineCPUs := float64(len(currStats.CPUStats.CPUUsage.PercpuUsage))
	if onlineCPUs == 0 {
		onlineCPUs = float64(currStats.CPUStats.OnlineCPUs)
	}
	if onlineCPUs == 0 {
		onlineCPUs = 1.0 // Fallback
	}

	// Calculate percentage
	cpuPercent := (float64(deltaCPU) / float64(deltaSystem)) * onlineCPUs * 100.0

	// Clamp to reasonable range (0-1000% to handle multi-core)
	if cpuPercent < 0 {
		cpuPercent = 0
	}
	if cpuPercent > 1000 {
		cpuPercent = 1000
	}

	return cpuPercent
}

// ResetStats clears previous stats for a container (useful after restart)
func (sc *StatsCalculator) ResetStats(containerID string) {
	delete(sc.previousStats, containerID)
	sc.logger.Debug("Reset stats for container", zap.String("container_id", containerID))
}

// ClearAllStats clears all stored previous stats
func (sc *StatsCalculator) ClearAllStats() {
	sc.previousStats = make(map[string]*container.StatsResponse)
}

