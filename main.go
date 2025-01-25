package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"sync"
	"syscall"
	"time"
)

var (
	cpuLoad     = flag.Float64("c", 0, "CPU load to apply (*100% CPU usage)")
	memGB       = flag.Float64("m", 0, "Memory to allocate in GiB (0 to disable)")
	absolute    = flag.Bool("a", false, "Enable absolute CPU usage mode (target system CPU usage with -c)")
	memStore    [][]byte
	loadMutex   sync.Mutex
	currentLoad float64
)

func main() {
	flag.Parse()

	if flag.NFlag() > 0 {
		handleFlagsMode()
	} else {
		handleLegacyMode()
	}

	waitForInterrupt()
}

func handleFlagsMode() {
	if len(flag.Args()) > 0 {
		fmt.Fprintln(os.Stderr, "Error: Positional arguments not allowed with flags")
		os.Exit(1)
	}

	if *cpuLoad < 0 || *memGB < 0 {
		fmt.Fprintln(os.Stderr, "Error: Negative values not allowed")
		os.Exit(1)
	}

	if *absolute {
		if *cpuLoad <= 0 {
			fmt.Fprintln(os.Stderr, "Error: -c must be specified and positive when using -a")
			os.Exit(1)
		}
		go absoluteCPUMode(*cpuLoad)
	} else {
		if *cpuLoad > 0 {
			startCPUWorkers(*cpuLoad)
		}
	}

	if *memGB > 0 {
		allocateMemory(*memGB)
	}

	if !*absolute && *cpuLoad == 0 && *memGB == 0 {
		fmt.Fprintln(os.Stderr, "Error: No stressors specified")
		os.Exit(1)
	}
}

func handleLegacyMode() {
	if len(flag.Args()) > 1 {
		fmt.Fprintln(os.Stderr, "Error: Too many arguments")
		os.Exit(1)
	}

	cpuLoad := float64(runtime.NumCPU())
	if len(flag.Args()) == 1 {
		n, err := strconv.ParseFloat(flag.Arg(0), 64)
		if err != nil || n <= 0 {
			fmt.Fprintf(os.Stderr, "Invalid CPU load value: %s\n", flag.Arg(0))
			os.Exit(1)
		}
		cpuLoad = n
	}

	startCPUWorkers(cpuLoad)
}

func startCPUWorkers(totalLoad float64) {
	fullCores := int(totalLoad)
	fractional := totalLoad - float64(fullCores)

	for i := 0; i < fullCores; i++ {
		go func() {
			for {
				_ = rand.Int()
				runtime.KeepAlive(struct{}{})
			}
		}()
	}

	if fractional > 0 {
		go cpuWorker(fractional)
	}

	fmt.Printf("Applying %d%% CPU load\n", int(totalLoad*100))
}

func cpuWorker(load float64) {
	const cycle = time.Second / 10
	if load <= 0 {
		return
	}
	if load > 1 {
		load = 1
	}

	busyDuration := time.Duration(float64(cycle) * load)
	idleDuration := cycle - busyDuration

	for {
		start := time.Now()
		end := start.Add(busyDuration)

		for time.Now().Before(end) {
			_ = rand.ExpFloat64()
		}

		time.Sleep(idleDuration)
	}
}

func absoluteCPUMode(target float64) {
	numCPUs := runtime.NumCPU()
	target *= 100
	fmt.Printf("Absolute mode: Targeting %.2f%% system CPU usage with %d cores\n", target, numCPUs)

	// PID controller parameters
	const (
		Kp = 0.9  // Proportional gain
		Ki = 0.05 // Integral gain
		Kd = 0.3  // Derivative gain
	)

	var (
		prevError float64
		integral  float64
	)

	for i := 0; i < numCPUs; i++ {
		go cpuWorkerDynamic()
	}

	ticker := time.NewTicker(time.Second / 3)
	defer ticker.Stop()

	for range ticker.C {
		currentSystem, err := getCPUUsage()
		if currentSystem == 0 {
			// avoid empty
			continue
		}
		if err != nil {
			fmt.Printf("Error getting CPU usage: %v\n", err)
			continue
		}

		// Calculate error
		error := target - currentSystem

		// Calculate integral term
		integral += error

		// Calculate derivative term
		derivative := error - prevError

		// Calculate PID output
		output := Kp*error + Ki*integral + Kd*derivative

		// Update previous error
		prevError = error

		// Calculate per-core load
		perCoreLoad := output / float64(numCPUs) / 100.0

		// Clamp the load between 0 and 1
		if perCoreLoad < 0 {
			perCoreLoad = 0
		} else if perCoreLoad > 1 {
			perCoreLoad = 1
		}

		loadMutex.Lock()
		currentLoad = perCoreLoad
		loadMutex.Unlock()

		fmt.Printf("Adjusting: System=%.2f%%, Target=%.2f%%, New Load=%.2f%% per core\n",
			currentSystem, target, perCoreLoad*100)
	}
}

func cpuWorkerDynamic() {
	const cycle = time.Second / 10

	for {
		loadMutex.Lock()
		load := currentLoad
		loadMutex.Unlock()

		if load <= 0 {
			time.Sleep(cycle)
			continue
		}
		if load > 1 {
			load = 1
		}

		busyDuration := time.Duration(float64(cycle) * load)
		idleDuration := cycle - busyDuration

		start := time.Now()
		end := start.Add(busyDuration)

		for time.Now().Before(end) {
			_ = rand.ExpFloat64()
		}

		time.Sleep(idleDuration)
	}
}

func allocateMemory(gib float64) {
	bytes := int64(gib * 1024 * 1024 * 1024)
	if bytes <= 0 {
		return
	}

	mem := make([]byte, bytes)
	for i := range mem {
		mem[i] = byte(i % 256)
	}

	memStore = append(memStore, mem)
	fmt.Printf("Allocated %.2f GiB\n", gib)
}

func waitForInterrupt() {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
	fmt.Println("\nExiting...")
	os.Exit(0)
}
