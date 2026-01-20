package docker

import (
	"testing"

	"github.com/docker/docker/api/types/container"
	"go.uber.org/zap"
)

func TestStatsCalculator_CalculateStats(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	calculator := NewStatsCalculator(logger)

	containerID := "test-container-1"

	// First stats (no previous stats)
	stats1 := &container.StatsResponse{
		CPUStats: container.CPUStats{
			CPUUsage: container.CPUUsage{
				TotalUsage:  1000000000, // 1 second in nanoseconds
				PercpuUsage: []uint64{100000000, 100000000, 100000000, 100000000},
			},
			SystemUsage: 2000000000, // 2 seconds
			OnlineCPUs:  4,
		},
		MemoryStats: container.MemoryStats{
			Usage: 1000000000, // 1GB
			Limit: 2000000000, // 2GB
			Stats: map[string]uint64{
				"cache": 200000000, // 200MB cache
			},
		},
		Networks: map[string]container.NetworkStats{
			"eth0": {
				RxBytes: 1000,
				TxBytes: 2000,
			},
		},
		BlkioStats: container.BlkioStats{
			IoServiceBytesRecursive: []container.BlkioStatEntry{
				{Op: "Read", Value: 5000},
				{Op: "Write", Value: 10000},
			},
		},
		PidsStats: container.PidsStats{
			Current: 10,
		},
	}

	result1, err := calculator.CalculateStats(containerID, stats1)
	if err != nil {
		t.Fatalf("CalculateStats failed: %v", err)
	}

	if result1.CPUPercent != 0.0 {
		t.Errorf("Expected CPU percent to be 0.0 for first stats, got %f", result1.CPUPercent)
	}

	if result1.MemoryUsage != 1000000000 {
		t.Errorf("Expected memory usage to be 1000000000, got %d", result1.MemoryUsage)
	}

	if result1.MemoryLimit != 2000000000 {
		t.Errorf("Expected memory limit to be 2000000000, got %d", result1.MemoryLimit)
	}

	// Second stats (with previous stats for CPU calculation)
	stats2 := &container.StatsResponse{
		CPUStats: container.CPUStats{
			CPUUsage: container.CPUUsage{
				TotalUsage:  2000000000, // 2 seconds
				PercpuUsage: []uint64{200000000, 200000000, 200000000, 200000000},
			},
			SystemUsage: 4000000000, // 4 seconds
			OnlineCPUs:  4,
		},
		MemoryStats: container.MemoryStats{
			Usage: 1200000000, // 1.2GB
			Limit: 2000000000, // 2GB
			Stats: map[string]uint64{
				"cache": 200000000, // 200MB cache
			},
		},
		Networks: map[string]container.NetworkStats{
			"eth0": {
				RxBytes: 2000,
				TxBytes: 4000,
			},
		},
		BlkioStats: container.BlkioStats{
			IoServiceBytesRecursive: []container.BlkioStatEntry{
				{Op: "Read", Value: 10000},
				{Op: "Write", Value: 20000},
			},
		},
		PidsStats: container.PidsStats{
			Current: 12,
		},
	}

	result2, err := calculator.CalculateStats(containerID, stats2)
	if err != nil {
		t.Fatalf("CalculateStats failed: %v", err)
	}

	// CPU should be calculated now (delta CPU / delta System * CPUs * 100)
	// delta CPU = 1 second, delta System = 2 seconds, CPUs = 4
	// CPU% = (1/2) * 4 * 100 = 200%
	expectedCPU := 200.0
	if result2.CPUPercent < expectedCPU-1.0 || result2.CPUPercent > expectedCPU+1.0 {
		t.Errorf("Expected CPU percent around %f, got %f", expectedCPU, result2.CPUPercent)
	}

	if result2.MemoryUsage != 1200000000 {
		t.Errorf("Expected memory usage to be 1200000000, got %d", result2.MemoryUsage)
	}

	if result2.NetworkRx != 2000 {
		t.Errorf("Expected network RX to be 2000, got %d", result2.NetworkRx)
	}

	if result2.NetworkTx != 4000 {
		t.Errorf("Expected network TX to be 4000, got %d", result2.NetworkTx)
	}
}

func TestStatsCalculator_ResetStats(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	calculator := NewStatsCalculator(logger)

	containerID := "test-container-1"

	stats := &container.StatsResponse{
		CPUStats: container.CPUStats{
			CPUUsage: container.CPUUsage{
				TotalUsage:  1000000000,
				PercpuUsage: []uint64{100000000},
			},
			SystemUsage: 2000000000,
			OnlineCPUs:  1,
		},
		MemoryStats: container.MemoryStats{
			Usage: 1000000000,
			Limit: 2000000000,
		},
	}

	// Calculate stats to store previous stats
	_, err := calculator.CalculateStats(containerID, stats)
	if err != nil {
		t.Fatalf("CalculateStats failed: %v", err)
	}

	// Reset stats
	calculator.ResetStats(containerID)

	// Calculate again - should have no previous stats
	result, err := calculator.CalculateStats(containerID, stats)
	if err != nil {
		t.Fatalf("CalculateStats failed: %v", err)
	}

	if result.CPUPercent != 0.0 {
		t.Errorf("Expected CPU percent to be 0.0 after reset, got %f", result.CPUPercent)
	}
}

func TestStatsCalculator_ClearAllStats(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	calculator := NewStatsCalculator(logger)

	stats := &container.StatsResponse{
		CPUStats: container.CPUStats{
			CPUUsage: container.CPUUsage{
				TotalUsage:  1000000000,
				PercpuUsage: []uint64{100000000},
			},
			SystemUsage: 2000000000,
			OnlineCPUs:  1,
		},
		MemoryStats: container.MemoryStats{
			Usage: 1000000000,
			Limit: 2000000000,
		},
	}

	// Calculate stats for multiple containers
	_, err = calculator.CalculateStats("container-1", stats)
	if err != nil {
		t.Fatalf("CalculateStats failed: %v", err)
	}
	_, err = calculator.CalculateStats("container-2", stats)
	if err != nil {
		t.Fatalf("CalculateStats failed: %v", err)
	}

	// Clear all stats
	calculator.ClearAllStats()

	// Calculate again - should have no previous stats
	result, err := calculator.CalculateStats("container-1", stats)
	if err != nil {
		t.Fatalf("CalculateStats failed: %v", err)
	}

	if result.CPUPercent != 0.0 {
		t.Errorf("Expected CPU percent to be 0.0 after clear, got %f", result.CPUPercent)
	}
}


