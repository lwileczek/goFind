# Generic build command removing CGO and debugging info for smaller release
build:
	CGO_ENABLED=0 go build -o gof -ldflags="-s -w" ./*.go

# Format code
fmt:
	go fmt ./...

# check the code for mistakes
lint:
	go vet ./...

#This will build the binary specifically for newer x86-64 machines. It assumes the CPU
#is newer like most modern desktop cpus, hence amd64 lvl v3
#	Docs: https://github.com/golang/go/wiki/MinimumRequirements#amd64
# With a small program like this we likely won't see a speedup but I like the opporunity
# to go fast
buildx86:
	GOARCH=amd64 GOAMD64="v3" CGO_ENABLED=0 go build -o gof -ldflags="-s -w" ./*.go
