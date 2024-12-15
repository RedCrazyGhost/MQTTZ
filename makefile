AUTHOR := "RedCrazyGhost"

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

all: build run

build:
	go fmt ./... && go build -o ./bin/MQTTZ -ldflags "-w -s -X main.VERSION=$(VERSION) -w -s -X  main.AUTHOR=$(AUTHOR)" ./main.go

run:
	./bin/MQTTZ