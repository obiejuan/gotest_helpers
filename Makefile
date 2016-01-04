BUILD_DATE := `date +%Y-%m-%d\ %H:%M`
VERSIONFILE := version.go

SOURCES=$(wildcard *.go)
PACKAGES=.

all:	build

build:	$(SOURCES)
	go build .

build_linux: $(SOURCES)
	GOOS=linux GOARCH=amd64 go build

clean:
	rm -f testutils

fmt:
	go fmt $(PACKAGES)

test: build
	go test -v $(PACKAGES)

install:
	go install

install_deps:
	go get -u