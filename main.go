package main

import (
	"fmt"
)

var IGNORE_PATHS = [7]string{
	".git",
	"build",
	"dist",
	"env",
	"node_modules",
	"target",
	"venv",
}

var (
	version = "development"
	commit  = "unknown"
	date    = "2024-01-01"
	builtBy = "lwileczek"
)

func printBuildInfo() {
	fmt.Printf("Version: %s\nCommit: %s\nBuild Date: %s | by: %s\n", version, commit, date, builtBy)
}

func main() {
	cfg := parseCLI()
	if cfg.Version {
		printBuildInfo()
		return
	}

	if cfg.Pattern == "" {
		fmt.Println("No pattern found, please provide a search pattern")
		return
	}

	printCh := make(chan string, cfg.Workers)
	//The system will reach deadlock if the work queue reaches capacity
	workQ := make(chan string, cfg.QueueSize)
	//To avoid deadlock, send tasks here which will have a non-blocky retry
	//func to add tasks back to workQ
	failover := make(chan string)
	dirCount := make(chan int)
	//Track how many dirs are open and close the work queue when we hit zero
	go dirChecker(dirCount, workQ)
	//Not closing as goroutines will continue to try and write if we exit early
	//with the -c flag but these should be fine and killed when the program exits
	//defer close(dirCount)
	//defer close(failover)
	go handleFailover(workQ, failover)
	go createWorkerPool(cfg.Pattern, workQ, failover, printCh, dirCount, cfg.Workers)

	//Send first work request
	workQ <- cfg.Dir

	//Print all results
	showResults(printCh, &cfg.MaxResults)
}
