package watcher

import (
	"com.young/agent/util"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type WatchDog struct {
	EnableWatchdog bool      `ini:"enableWatchdog"       default:"true"`
	WatchdogLevel int 		 `ini:"watchdogLevel"        default:"0"`
	WatchdogMemLimit int64   `ini:"watchdogMemLimit"     default:"0"`
	WatchdogCPULimit int     `ini:"watchdogCPULimit"     default:"0"`
	WatchdogLatencyLimit int `ini:"watchdogLatencyLimit" default:"12"`
	WatchdogInterval int	 `ini:"watchdogInterval"     default:"3"`
}

type state struct {
	sustainedLatency int

}


func Watch() chan<- struct{} {
	exepath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	err = util.LoadConfig(strings.TrimSuffix(exepath, ".exe")+".ini", WatchDog{})
	if err != nil {
		log.Fatal(err)
	}
	var isAgentSane chan struct{}


	return isAgentSane
}
