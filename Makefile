# These are the values we want to pass for Version and BuildTime
VERSION=0.0.1
BUILD_TIME=$(shell date +%s)

# Setup the -ldflags option for go build here, interpolate the variable values

LDFLAGS += -X \"main.Version=$(VERSION)\"
LDFLAGS += -X \"main.BuildTime=$(BUILD_TIME)\"
LDFLAGS += -X \"main.BotToken=$(BOT_TOKEN)\"
LDFLAGS += -X \"main.BotToken=$(HTTP_PORT)\"
LDFLAGS += -X \"main.BotToken=$(IP)\"
LDFLAGS += -X \"main.BotToken=$(WEBHOOK_PORT)\"
LDFLAGS += -X \"main.BotToken=$(WEBHOOK_CERT)\"
LDFLAGS += -X \"main.BotToken=$(WEBHOOK_KEY)\"
LDFLAGS += -X \"main.BotToken=$(ENABLE_WEBHOOK)\"

.PHONY: build clean

setup:
ifndef BOT_TOKEN
	$(error BOT_TOKEN is not set.)
endif

build: setup
	go build -ldflags "$(LDFLAGS)"

install: setup
	go get -d ./...
	go install -ldflags "$(LDFLAGS)"

clean:
	go clean -i ./...

clean-mac: clean
	find . -name ".DS_Store" -print0 | xargs -0 rm
