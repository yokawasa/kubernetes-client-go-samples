.PHONY: clean all list-svc list-nodes list-pods pods-in-svc istio-get-vs istio-update-vs

.DEFAULT_GOAL := all

TARGETS=\
	list-svc \
	list-nodes \
	list-pods \
	pods-in-svc \
	istio-get-vs \
	istio-update-vs

cur   := $(shell pwd)

${TARGETS}:
	cd ${cur}/$@ && GO111MODULE=on go build -o ${cur}/dist/$@ main.go

all: $(TARGETS)

clean:
	rm -rf dist

lint:
	golint -set_exit_status $$(go list ./...)
	go vet ./...
