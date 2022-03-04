.PHONY: build
#GO111MODULE=on 

XAPPNAME=rimedo-ts
RIMEDO_TS_VERSION := latest

build:
	GOPRIVATE="github.com/onosproject/*" go build -o build/_output/$(XAPPNAME) ./cmd/$(XAPPNAME)

build-tools:=$(shell if [ ! -d "./build/build-tools" ]; then cd build && git clone https://github.com/onosproject/build-tools.git; fi)
include ./build/build-tools/make/onf-common.mk

docker:
	@go mod vendor
	sudo docker build --network host -f build/Dockerfile -t onosproject/$(XAPPNAME):$(RIMEDO_TS_VERSION) .
	@rm -rf vendor

images: build
	@go mod vendor
	docker build -f build/Dockerfile -t onosproject/$(XAPPNAME):$(RIMEDO_TS_VERSION) .
	@rm -rf vendor

kind: images
	@if [ "`kind get clusters`" = '' ]; then echo "no kind cluster found" && exit 1; fi
	kind load docker-image onosproject/$(XAPPNAME):${RIMEDO_TS_VERSION}

install-xapp:
	helm install -n riab $(XAPPNAME) ./helm-chart/$(XAPPNAME) --values ./helm-chart/$(XAPPNAME)/values.yaml

delete-xapp:
	-helm uninstall -n riab $(XAPPNAME)

dev: delete-xapp docker install-xapp

test: build
jenkins-test: build

docker-login:
ifdef DOCKER_USER
ifdef DOCKER_PASSWORD
	echo ${DOCKER_PASSWORD} | docker login -u ${DOCKER_USER} --password-stdin
else
	@echo "DOCKER_USER is specified but DOCKER_PASSWORD is missing"
	@exit 1
endif
endif

docker-push-latest: docker-login
	docker push onosproject/$(XAPPNAME):latest

publish: # @HELP publish version on github and dockerhub
	./build/build-tools/publish-version ${VERSION} onosproject/$(XAPPNAME)

jenkins-publish: jenkins-tools images docker-push-latest # @HELP Jenkins calls this to publish artifacts
	./build/build-tools/release-merge-commit


