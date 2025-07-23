# Go Find
A file finder written in Go.

Uses a worker pool to search through directory trees quicker by taking advantage of multi-core machines.
Go concurrency allows it to run on single core machines. 

### Features
 - Concurrency search directories for desired pattern
 - Exit early with a max results option
 - Common paths are ignored to save time and checks, such as `node_modules`, `.git`, and python `venv`.
 - Searches for a pattern using basic sub-string
 - Use a worker pool to avoid spinning up too many go routines if there are _many_ directories

## Warning
> **warning**
> Does not work on windows and I'm not trying too support this

## Using
### Usage
Customise your search with the following flags

```bash
Usage: goFind [OPTIONS] <pattern> <path>
  set pattern & path positionally or via flags
  -c int
        The maximum number of results to find (default -1)
  -d string
        The starting directory to check for files (default ".")
  -i    perform a case insensative search
  -p string
        A pattern to check for within the file names
  -q int
        The max work queue size (default 512)
  -v    print the version and build information
  -w int
        Number of workers (default -1)
```

## Contributing
### Build
build with `go build -o gf main.go` and run with `./gf -d <starting-path> -p <pattern-to-match-on>`
#### Requirements
 - Go 1.21+

### Makefile
You can use the [makefile](./Makefile) to build a production release with `make build`
#### Requirements
 - Make
 - Go 1.21+

## Roadmap
 - [x] can we create a worker pool instead of the number of initial directories
 - [x] cap the number of returned responses, say 1 (quit after the first match!)
 - [x] Auto-generate releases
 - [x] cap the depth of the search
 - [x] use select statement with a fallback queue to prevent deadlock from happening
 - [x] use postional args or flags to set the pattern/path
 - [ ] Add flag to print which paths will be ignored
 - [ ] create case sensative/insensative searching
 - [ ] exclude additional directories or patters
 - [ ] include specific patterns in the search path
 - [ ] switch to ignore all hidden files/directories
 - [ ] look for `.gitignore` files and use contents to build ignored patterns
