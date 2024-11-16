package Helpers

import (
	"github.com/shirou/gopsutil/cpu"
	"runtime"
)

type Win32_PerfFormattedData_PerfOS_Process struct {
	Name                 string
	PercentProcessorTime uint64
}

func CalculateBatchSize(totalWork int, numCores int) int {
	if numCores <= 0 {
		numCores = 1
	}

	batchSize := totalWork / numCores

	if batchSize < 1 {
		batchSize = 1
	}

	return batchSize * 3
}

func determineStrategy() string {
	highCPULoadThreshold := 70.0
	mediumCPULoadThreshold := 50.0

	numCPU := runtime.NumCPU()

	percent, err := cpu.Percent(0, false)
	if err != nil {
		return "Guardian"
	}

	load := percent[0]

	if load > highCPULoadThreshold {
		return "Guardian"
	}

	if load > mediumCPULoadThreshold && numCPU > 2 {
		return "Vortex"
	}

	if numCPU <= 2 {
		return "Sentry"
	}

	return "Lympia"
}
