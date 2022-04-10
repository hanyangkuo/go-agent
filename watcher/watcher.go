package watcher

import (
	"com.young/agent/util"
	"fmt"
	"github.com/shirou/gopsutil/v3/process"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type WatchDog struct {
	Enable bool      `ini:"enable"       default:"true"`
	Level int 		 `ini:"level"        default:"0"`
	MemLimit int64   `ini:"memLimit"     default:"0"`
	CPULimit int     `ini:"cpuLimit"     default:"0"`
	LatencyLimit int `ini:"latencyLimit" default:"12"`
	Interval int	 `ini:"interval"     default:"3"`
}

type PerformanceState struct {
	sustainedLatency uint64
	user_time uint64
	system_time uint64
	initial_footprint uint64
}
type PerformanceChange struct{
	sustained_latency uint64
	footprint uint64
	iv uint64
}


func Watch() error {
	_, err := loadWatchDogConfig()
	if err != nil {
		log.Fatal(err)
	}
	//getChange()
	return nil
	//return isAgentSane
}

func loadWatchDogConfig() (*WatchDog, error){
	watchDog := new(WatchDog)
	exepath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	err = util.LoadConfig(strings.TrimSuffix(exepath, ".exe")+".ini", watchDog)
	if err != nil {
		return nil, err
	}
	if watchDog.MemLimit == 0 {
		switch watchDog.Level {
		case 0:
			watchDog.MemLimit = 200
		case 1:
			watchDog.MemLimit = 100
		}
	}
	if watchDog.CPULimit == 0 {
		switch watchDog.Level {
		case 0:
			watchDog.CPULimit = 20
		case 1:
			watchDog.CPULimit = 10
		}
	}
	return watchDog, nil
}

// getChange
// getProcess with osquery:  select "parent", "user_time", "system_time", "resident_size", "total_size" from processes;
// RSS: resident_size, VMS: total_size
func getChange(state PerformanceState) {
	userTime, kernelTime, residentSize, totalSize := getProcessRow()
	//p := &process.Process{Pid: int32(os.Getpid())}
	//m, err := p.MemoryInfo()
	//if err != nil {
	//	log.Println(err)
	//}
	//fmt.Printf("RSS: %f\tHWM: %f\tVMS: %d\n", m.RSS, m.HWM, m.VMS)
}