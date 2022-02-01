BINARY := jshell

default: $(BINARY)

$(BINARY): *.go go.* apps/*
	CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -o $@
