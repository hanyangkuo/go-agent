package main

import (
	"com.young/agent/util"
	"flag"
	"fmt"
	_ "github.com/kardianos/service"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const version = "1.0.0"

var (
	rootDir = ""
)

func init() {
	flagVerbose := flag.Bool("verbose", false, "")
	flagVer := flag.Bool("version", false, "print")
	flag.Parse()
	if *flagVer {
		fmt.Print(version)
		os.Exit(0)
	}

	exepath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	rootDir = filepath.Dir(exepath)
	os.Mkdir(filepath.Join(rootDir, "log"), os.ModePerm)
	file, err := rotatelogs.New(
		filepath.Join(rootDir, "log", "agent.log.%Y%m%d"),
		rotatelogs.WithMaxAge(time.Duration(72)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.SetFormatter(&log.JSONFormatter{})
	if *flagVerbose {
		log.SetOutput(io.MultiWriter(file, os.Stdout))
	} else {
		log.SetOutput(file)
	}
}


func main() {

}


func testFunc(){
	defer util.TimeTrack(time.Now(), "testFunc")
	time.Sleep(1*time.Second)
	runtime.GC()
}
