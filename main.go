package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
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

func main() {
	//The number of workers. If there are more workers the system can read from
	//the work queue more often and a larger queue is not required. It's a blance
	jobs := flag.Int("j", 4, "Estimated number of workers?")
	//If this is reached the system could end up in deadlock
	//The bigger the queue size the more memory is used
	//Smaller could be faster but you coud have deadlock
	queueSize := flag.Int("q", 128, "The max work queue size")
	dir := flag.String("d", ".", "The starting directory to check for files")
	pattern := flag.String("p", "", "A pattern to check for within the file names")
	flag.Parse()

	if *pattern == "" {
		fmt.Println("No pattern provided")
		return
	}

	//Only for OSX/Linux, sorry windows
	//Remove any trailing slashes in the path
	if (*dir)[len(*dir)-1:] == "/" {
		*dir = string((*dir)[0 : len(*dir)-1])
	}

	printCh := make(chan string, *jobs)
	//The system will reach deadlock if the work queue reaches capacity
	workQ := make(chan string, *queueSize)
	//To avoid deadlock, send tasks here which will have a non-blocky retry
	//func to add tasks back to workQ
	failover := make(chan string)
	dirCount := make(chan int)
	//Track how many dirs are open and close the work queue when we hit zero
	go dirChecker(dirCount, workQ)
	defer close(dirCount)
	defer close(failover)
	go handleFailover(workQ, failover)
	go createWorkerPool(pattern, workQ, failover, printCh, dirCount, jobs)

	//Send first work request
	workQ <- *dir

	//Print all results
	for item := range printCh {
		fmt.Println(item)
	}
}

func dirChecker(in chan int, work chan string) {
	n := 1
	for i := range in {
		n += i
		if n <= 0 {
			close(work)
			return
		}
	}
}

func createWorkerPool(p *string, in chan string, failover chan string, results chan string, cnt chan int, w *int) {
	var wg sync.WaitGroup
	for i := 0; i < *w; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			search(p, in, failover, results, cnt)
		}()
	}
	wg.Wait()
	close(results)
}
func search(pattern *string, in chan string, failover chan string, out chan string, cnt chan int) {
	for path := range in {
		items, err := os.ReadDir(path)
		if err != nil {
			fmt.Println("Error reading the directory", path)
			fmt.Println(err)
			cnt <- -1
		}
	ItemSearch:
		for _, item := range items {
			if item.IsDir() {
				//Don't dive into directories I don't care about
				for _, p := range IGNORE_PATHS {
					if p == item.Name() {
						continue ItemSearch
					}
				}
				subPath := fmt.Sprintf("%s/%s", path, item.Name())
				cnt <- 1
				select {
				case in <- subPath:
				case failover <- subPath:
				}
			} else {
				if strings.Index(item.Name(), *pattern) >= 0 {
					//subPath is repeated but no point in creating an allocation if not required
					subPath := fmt.Sprintf("%s/%s", path, item.Name())
					out <- subPath
				}
			}
		}
		//We finished reading everything in the dir, tell the accounted we finished
		cnt <- -1
	}
}

func handleFailover(work, fail chan string) {
	var q []string
	for {
		task := <-fail
		q = append(q, task)
		//TODO: Add verbose logging here so users can check if the failover was used
		for {
			select {
			case work <- q[0]:
				q = q[1:]
			case task := <-fail:
				q = append(q, task)
			default:
			}
			//I don't know if we'll get an issue with `work <- q[0]` unless we have this
			if len(q) == 0 {
				break
			}
		}
	}
}
