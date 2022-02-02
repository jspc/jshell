BINARY := jshell

default: $(BINARY)

$(BINARY): *.go go.* apps/* apps/hex/wordlist.generated.go
	CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o $@

apps/hex/wordlist.generated.go: apps/hex/generator/main.go
	go run $< > $@
	go fmt $@
