# SPDX-License-Identifier: Apache-2.0
# Copyright 2019 Open Networking Foundation
# Copyright 2019 Rimedo Labs
# Copyright 2024 Intel Corporation

.PHONY: build
export CGO_ENABLED=1
export GO111MODULE=on

RIMEDO_TS_VERSION ?= latest

GOLANG_CI_VERSION := v1.52.2

all: build docker-build

build: # @HELP build the Go binaries and run all validations (default)
	GOPRIVATE="github.com/onosproject/*" go build -o build/_output/rimedo-ts ./cmd/rimedo-ts

test: # @HELP run the unit tests and source code validation
test: build lint license
	go test -race github.com/onosproject/rimedo-ts/pkg/...
	go test -race github.com/onosproject/rimedo-ts/cmd/...

docker-build-rimedo-ts: # @HELP build Docker image
	@go mod vendor
	docker build --network host . -f build/rimedo-ts/Dockerfile \
		-t onosproject/rimedo-ts:${RIMEDO_TS_VERSION}
	@rm -rf vendor

docker-build: # @HELP build all Docker images
docker-build: build docker-build-rimedo-ts

docker-push-rimedo-ts: # @HELP push Docker image
	docker push onosproject/rimedo-ts:${RIMEDO_TS_VERSION}

docker-push: # @HELP push docker images
docker-push: docker-push-rimedo-ts

lint: # @HELP examines Go source code and reports coding problems
	golangci-lint --version | grep $(GOLANG_CI_VERSION) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b `go env GOPATH`/bin $(GOLANG_CI_VERSION)
	golangci-lint run --timeout 15m

license: # @HELP run license checks
	rm -rf venv
	python3 -m venv venv
	. ./venv/bin/activate;\
	python3 -m pip install --upgrade pip;\
	python3 -m pip install reuse;\
	reuse lint

check-version: # @HELP check version is duplicated
	./build/bin/version_check.sh all

clean: # @HELP remove all the build artifacts
	rm -rf ./build/_output ./vendor ./cmd/rimedo-ts/rimedo-ts ./cmd/onos/onos venv
	go clean github.com/onosproject/rimedo-ts/...

help:
	@grep -E '^.*: *# *@HELP' $(MAKEFILE_LIST) \
    | sort \
    | awk ' \
        BEGIN {FS = ": *# *@HELP"}; \
        {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}; \
    '