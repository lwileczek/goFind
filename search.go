package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
)

func showResults(ch chan string, limit *int) {
	if *limit > 0 {
		n := 0
		for item := range ch {
			fmt.Println(item)
			n++
			if n >= *limit {
				return
			}
		}
		return
	}

	for item := range ch {
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

func createWorkerPool(p string, in chan string, failover chan string, results chan string, cnt chan int, w int) {
	var wg sync.WaitGroup
	for i := 0; i < w; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			search(p, in, failover, results, cnt)
		}()
	}
	wg.Wait()
	close(results)
}

func search(pattern string, in chan string, failover chan string, out chan string, cnt chan int) {
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
			}
			//Always check if the name of the thing matches pattern, including directory names
			if strings.Contains(item.Name(), pattern) {
				//subPath is repeated but no point in creating an allocation if not required
				subPath := fmt.Sprintf("%s/%s", path, item.Name())
				out <- subPath
			}

		}
		//We finished reading everything in the dir, tell the accounted we finished
		cnt <- -1
	}
}

// handleFailover if the work queue is backed up, append the task to a slice
// This helps to prevent deadlock since our works both read and write to the task
// queue
func handleFailover(work, fail chan string) {
	q := make([]string, 0, 64)
	for {
		task := <-fail
		q = append(q, task)
		slog.Debug("failover task added to queue", "queue length", len(q))
		for {
			select {
			case work <- q[0]:
				q = q[1:]
			case task := <-fail:
				q = append(q, task)
				slog.Debug("failover task added to queue", "queue length", len(q))
			default:
			}

			// avoid indexing q in `work <- q[0]` if no elements to avoid panic
			if len(q) == 0 {
				break
			}
		}
	}
}
