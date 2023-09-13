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

	var wg sync.WaitGroup

	done := make(chan bool)
	printCh := make(chan string, 2**jobs)
	go printResults(printCh, done)

	//TODO: This is repeated in the function below except I allow the creation
	//of go routines. Repeated code is bad!
	items, err := os.ReadDir(*dir)
	if err != nil {
		panic(err)
	}
InitialItemSearch:
	for _, item := range items {
		subPath := fmt.Sprintf("%s/%s", *dir, item.Name())
		//Starting from the initial directory, create a new goroutine for each subdirectory
		//but no more, if there is one directory, use one goroutine
		if item.IsDir() {
			for _, p := range IGNORE_PATHS {
				if p == item.Name() {
					continue InitialItemSearch
				}
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				recursiveSearch(subPath, pattern, printCh)
			}()
		} else {
			if strings.Index(item.Name(), *pattern) >=0 {
				printCh<- subPath
			}
		}
	}

	wg.Wait()
	close(printCh)
	<-done
}

func recursiveSearch(path string, pattern *string, out chan string) {
	items, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("Error reading the directory", path)
		fmt.Println(err)
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
			recursiveSearch(subPath, pattern, out)
		} else {
			if strings.Index(item.Name(), *pattern) >=0 {
				//subPath is repeated but no point in creating an allocation if not required
				subPath := fmt.Sprintf("%s/%s", path, item.Name())
				out <- subPath
			}
		}
	}
}

func printResults(in chan string, done chan bool) {
	for item := range in {
		fmt.Println(item)
	}
	done <- true
}
