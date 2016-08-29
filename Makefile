# These are the values we want to pass for Version and BuildTime
VERSION=0.0.1
BUILD_TIME=$(shell date +%s)

# Setup the -ldflags option for go build here, interpolate the variable values

LDFLAGS += -X \"main.Version=$(VERSION)\"
LDFLAGS += -X \"main.BuildTime=$(BUILD_TIME)\"

.PHONY: build clean

build:
	go build -ldflags "$(LDFLAGS)"

deps:
	go get -d ./...

install: setup deps
	go install -ldflags "$(LDFLAGS)"

clean:
	go clean -i ./...

clean-mac: clean
	find . -name ".DS_Store" -print0 | xargs -0 rm
