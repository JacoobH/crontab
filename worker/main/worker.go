package main

import (
	"flag"
	"fmt"
	"github.com/JacoobH/crontab/worker"
	"runtime"
)

var (
	confFile string // Configuration file path
)

// Parses command line arguments
func initArgs() {
	// worker -config ./worker.json
	// worker -h
	flag.StringVar(&confFile, "config", "./worker.json", "Specify the worker.json")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	// Initialize command line arguments
	initArgs()

	//初始化线程
	initEnv()

	// Load the configuration
	if err = worker.InitConfig(confFile); err != nil {
		goto ERR
	}

	if err = worker.InitJobMgr(); err != nil {
		goto ERR
	}

	//正常退出
	return

ERR:
	fmt.Println(err)
}
