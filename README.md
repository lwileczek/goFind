# Go Find
A file finder written in Go for my needs.

Uses a worker pool to search through directory trees quicker and take advantage of multi-core machines.
Go concurrency allows it to run on single core machines.

If you're looking for something specific you can limit the number of responses

### Features
 - Concurrency which may speed up or slow down the search depending on the tree structure. 
 - Common paths are ignored to save time and checks, such as `node_modules`, `.git`, and python `venv`.
 - Searches for a pattern using basic sub-string

## Warning
> **warning**
> Does not work on windows and I'm not trying too support this

## TODO
 - sub command to print paths that are ignored
 - create case sensative/insensative searching
 - add the ability exclude additional directories or patters
 - add the ability to include specific patterns in the search path
 - switch to ignore hidden files/directories
 - [x] can we create a worker pool instead of the number of initial directories
 - [x] cap the number of returned responses, say 1 (quit after the first match!)
 - cap the depth of the search

## Using
### Build
build with `go build -o gf main.go` and run with `./gf -d <starting-path> -p <pattern-to-match-on>`

### Flags
Customise your search with the following flags
```bash
Usage of ./gf:
  -c int
    	The maximum number of results to find (default -1)
  -d string
    	The starting directory to check for files (default ".")
  -p string
    	A pattern to check for within the file names
  -q int
    	The max work queue size (default 2048)
  -w int
    	Number of workers (default 8)
```
