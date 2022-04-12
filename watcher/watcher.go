package watcher

import (
	"com.young/agent/util"
	log "github.com/sirupsen/logrus"
	"runtime"
	"time"
)

type WatchDog struct {
	Enable bool      `ini:"enable"       default:"true"`
	Level int 		 `ini:"level"        default:"0"`
	MemLimit uint64   `ini:"memLimit"     default:"0"`
	CPULimit uint64   `ini:"cpuLimit"     default:"0"`
	LatencyLimit uint64 `ini:"latencyLimit" default:"12"`
	Interval uint64	 `ini:"interval"     default:"3"`
}

type PerformanceState struct {
	sustainedLatency uint64
	userTime uint64
	kernelTime uint64
	initialFootprint uint64
}
type PerformanceChange struct{
	sustainedLatency uint64
	footprint uint64
	iv uint64
}

var (
	watchDog = new(WatchDog)
)

func Initialize(config string) error {
	err := util.LoadConfig(config, watchDog)
	if err != nil {
		return err
	}
	if watchDog.CPULimit == 0 {
		switch watchDog.Level {
		case 0:
			watchDog.CPULimit = 10
		case 1:
			watchDog.CPULimit = 5
		default:
			watchDog.CPULimit = 100
		}
	}
	if watchDog.MemLimit == 0 {
		switch watchDog.Level {
		case 0:
			watchDog.MemLimit = 200
		case 1:
			watchDog.MemLimit = 100
		default:
			watchDog.MemLimit = 10000
		}
	}
	return nil
}

func IsEnable() bool {
	return watchDog.Enable
}

func Watcher() <-chan struct{} {
	watcherState := make(chan struct{})
	go isWatcherHealthy(watcherState)
	return watcherState
}

func isWatcherHealthy(watcherState chan<- struct{}) {
	state := new(PerformanceState)

	for {
		userTime, kernelTime, residentSize, totalSize := getProcessRow()
		//log.Infof("userTime = %d, kernelTime = %d, residentSize = %d, totalSize = %d", userTime, kernelTime, residentSize, totalSize)
		change := getChange(userTime, kernelTime, residentSize, totalSize, state)
		if exceedCyclesLimit(change) {
			log.Warn("exceed cycles cpu limit.")
			break
		}
		if exceedMemoryLimit(change) {
			log.Warn("exceed memory limit.")
			break
		}
		time.Sleep(time.Duration(watchDog.Interval)*time.Second)
	}
	close(watcherState)
}

// getChange
// getProcess with osquery:  select "parent", "user_time", "system_time", "resident_size", "total_size" from processes;
// RSS: resident_size, VMS: total_size
func getChange(userTime, kernelTime, residentSize, totalSize uint64, state *PerformanceState) *PerformanceChange {
	change := new(PerformanceChange)

	// handle cpu time
	cpuExpectedTime := watchDog.CPULimit *  watchDog.Interval*1000 * uint64(runtime.NumCPU())/ 100
	cpuTakeTime := userTime - state.userTime + kernelTime - state.kernelTime
	if cpuTakeTime > cpuExpectedTime {
		state.sustainedLatency++
	} else {
		state.sustainedLatency = 0
	}
	state.userTime = userTime
	state.kernelTime = kernelTime

	change.sustainedLatency = state.sustainedLatency
	// handle memory
	// set change.footprint to current memory size
	if runtime.GOOS == "windows" {
		change.footprint = totalSize
	} else {
		change.footprint = residentSize
	}
	if state.initialFootprint == 0 {
		state.initialFootprint = change.footprint
	}

	if change.footprint < state.initialFootprint {
		change.footprint = 0
	} else {
		change.footprint = change.footprint - state.initialFootprint
	}
	return change
}

func exceedCyclesLimit(change *PerformanceChange) bool {
	if change.sustainedLatency == 0 {
		return false
	}
	return change.sustainedLatency* watchDog.Interval >= watchDog.LatencyLimit
}

func exceedMemoryLimit(change *PerformanceChange) bool {
	if change.footprint == 0 {
		return false
	}
	return change.footprint > watchDog.MemLimit*1024*1024
}