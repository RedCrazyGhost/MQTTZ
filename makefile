AUTHOR := "RedCrazyGhost"
REPO_URL := "https://github.com/RedCrazyGhost/MQTTZ"

CURRENT_TAG := $(strip $(shell git describe --tags --always))
CURRENT_COMMIT_HASH := $(strip $(shell git rev-parse --short HEAD))

ifeq ($(CURRENT_TAG),$(CURRENT_COMMIT_HASH))
  VERSION := $(CURRENT_COMMIT_HASH)
else
  ifeq ($(CURRENT_COMMIT_HASH),$(shell git show-ref --hash=7 $(CURRENT_TAG)))
    VERSION := $(shell git describe --tags --abbrev=0)
  else
    VERSION := $(shell git describe --tags --abbrev=0)-$(CURRENT_COMMIT_HASH)
  endif
endif

local: local_build local_run
build: local_build docker_build version

version:
  @echo build app version: $(VERSION)

local_build:
	wire ./cmd
	go fmt ./...
	go build -o ./bin/MQTTZ -ldflags "-w -s -X main.VERSION=$(VERSION) -X main.AUTHOR=$(AUTHOR) -X main.REPO_URL=$(REPO_URL)" ./cmd

local_run:
	./bin/MQTTZ -config ./bin/conf/config.yaml