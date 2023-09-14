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
	jobs := flag.Int("j", 4, "Estimated number of workers?")
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

	printCh := make(chan string, 2**jobs)
	workQ := make(chan string, 4**jobs)
	dirCount := make(chan int, *jobs)
	//Track how many dirs are open and close the work queue when we hit zero
	go dirChecker(dirCount, workQ)
	defer close(dirCount)
	go createWorkerPool(pattern, workQ, printCh, dirCount, jobs)

    //Send first work request
    workQ <- *dir
    
	//Print all results
	for item := range printCh {
		fmt.Println(item)
	}
}

func dirChecker(in chan int, fin chan string) {
	n := 1
	for i := range in {
		n += i
		if n <= 0 {
			close(fin)
			return
		}
	}
}

func createWorkerPool(p *string, in chan string, out chan string, cnt chan int, w *int) {
	var wg sync.WaitGroup
	for i := 0; i < *w; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			search(p, in, out, cnt)
		}()
	}
	wg.Wait()
	close(out)
}
func search(pattern *string, in chan string, out chan string, cnt chan int) {
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
				cnt <- 1
				in <- fmt.Sprintf("%s/%s", path, item.Name())
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
