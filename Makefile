.PHONY: install
install: build
	cp ./bin/crdflow $(GOPATH)/bin/crdflow

.PHONY: build
build:
	CGO_ENABLED=0 go build -o bin/crdflow main.go