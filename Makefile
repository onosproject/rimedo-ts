.PHONY: build
#GO111MODULE=on 

XAPPNAME=rimedo-ts
VERSION=v0.0.1

build:
	GOPRIVATE="github.com/onosproject/*" go build -o build/_output/$(XAPPNAME) ./cmd/$(XAPPNAME)

build-tools:=$(shell if [ ! -d "./build/build-tools" ]; then cd build && git clone https://github.com/onosproject/build-tools.git; fi)
include ./build/build-tools/make/onf-common.mk

docker:
	@go mod vendor
	sudo docker build --network host -f build/Dockerfile -t onosproject/$(XAPPNAME):$(VERSION) . 
	@rm -rf vendor

install-xapp:
	helm install -n riab $(XAPPNAME) ./helm-chart/$(XAPPNAME) --values ./helm-chart/$(XAPPNAME)/values.yaml

delete-xapp:
	-helm uninstall -n riab $(XAPPNAME)

dev: delete-xapp docker install-xapp 
