package main

import (
	"fmt"
	"runtime"

	"github.com/shirou/gopsutil/v3/cpu"
)

func getCPUUsage() (float64, error) {
	// Get the system CPU usage
	systemCPUPercent, err := cpu.Percent(0, false)
	if err != nil {
		return 0, fmt.Errorf("cannot get system CPU usage: %v", err)
	}
	if len(systemCPUPercent) < 1 {
		return 0, fmt.Errorf("cannot get system CPU usage")
	}

	// Calculate the total system CPU usage by multiplying by the number of CPUs
	totalSystemCPUPercent := systemCPUPercent[0] * float64(runtime.NumCPU())

	return totalSystemCPUPercent, nil
}
