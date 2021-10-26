export $(shell sed 's/=.*//' .env)

GOPATH=$(shell go env GOPATH)

K = kubectl

.PHONY: help fmt vet

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

help:	# The following lines will print the available commands when entering just 'make'. ⚠️ This needs to be the first target, ever
ifeq ($(UNAME), Linux)
	@grep -P '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
else
	@awk -F ':.*###' '$$0 ~ FS {printf "%15s%s\n", $$1 ":", $$2}' \
		$(MAKEFILE_LIST) | grep -v '@awk' | sort
endif

fmt:	### Run go fmt against code
	go fmt ./...

vet:	### Run go vet against code
	go vet ./...

build: ### Run go test to build and test the code
	go test ./...
