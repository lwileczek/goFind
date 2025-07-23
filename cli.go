package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

// Config the programs configuration used to set behavior
type Config struct {
	QueueSize   int
	MaxResults  int
	Workers     int
	Dir         string
	Pattern     string
	Insensative bool
	Version     bool
}

// myUsage overwrite default usage message from ./<bin> --help
func myUsage() {
	bin := "gof"
	if len(os.Args) > 0 {
		bin = os.Args[0]
	}

	fmt.Printf("Usage: %s [OPTIONS] <pattern> <path>\n", bin)
	fmt.Println("  set pattern & path positionally or via flags")
	flag.PrintDefaults()
}

// parseArgs all arguments from the user to use with the program
func parseArgs() *Config {
	flag.Usage = myUsage

	//The number of workers. If there are more workers the system can read from
	//the work queue more often and a larger queue is not required.
	workers := flag.Int("w", -1, "Number of workers")
	//If the queue overflows we'll use a slice to store work which might slow the system
	queueSize := flag.Int("q", 512, "The max work queue size")
	maxResults := flag.Int("c", -1, "The maximum number of results to find")
	dir := flag.String("d", ".", "The starting directory to check for files")
	pattern := flag.String("p", "", "A pattern to check for within the file names")
	insensative := flag.Bool("i", false, "perform a case insensative search")
	v := flag.Bool("v", false, "print the version and build information")
	flag.Parse()

	p := *pattern
	if p == "" {
		p = flag.Arg(0)
	}

	w := *workers
	if *workers <= 0 {
		// magic 2, anecdotal evidence of better performance over NumCPU
		w = runtime.NumCPU() + 2
	}

	if *dir == "." && flag.Arg(1) != "" {
		*dir = flag.Arg(1)
	}

	//Only for OSX/Linux, sorry windows
	//Remove any trailing slashes in the path
	if (*dir)[len(*dir)-1:] == "/" {
		*dir = string((*dir)[0 : len(*dir)-1])
	}

	return &Config{
		QueueSize:   *queueSize,
		MaxResults:  *maxResults,
		Workers:     w,
		Dir:         *dir,
		Pattern:     p,
		Insensative: *insensative,
		Version:     *v,
	}
}
