# Go Find
A file finder written in Go for my needs.

Uses a worker pool to search through directory trees quicker and take advantage of multi-core machines.
Go concurrencly allows it to run on single core machines.

## Goal
The goal is to write this to use a worker pool but currently it spins up a go routine
for each directory within the starting directory and then linearly searches each sub-directory 

### Features
Concurrency which may speed up or slow down the search depending on the tree structure. 
Common paths are ignored to save time and checks, such as `node_modules`, `.git`, and python `venv`.
Searches for a pattern using basic sub-string

## Warning
> **warning**
> Does not work on windows and I'm not trying to

## TODO
 - sub command to print paths that are ignored
 - create case sensative/insensative searching
 - can we create a worker pool instead of the number of initial directories

## Using
build with `go build -o find` and run with `./find -d <starting-path> -p <pattern-to-match-on>`

```bash
  -d string
     The starting directory to check for files (default ".")
  -j int
     Estimated number of workers? (default 4)
  -p string
     A pattern to check for within the file names
```
