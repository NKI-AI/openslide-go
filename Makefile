all: deps build test

deps: FORCE
	CGO_CFLAGS_ALLOW=-Xpreprocessor go get ./...

build: FORCE
	CGO_CFLAGS_ALLOW=-Xpreprocessor go build ./openslide

test: FORCE
	CGO_CFLAGS_ALLOW=-Xpreprocessor go test -v ./...

FORCE: