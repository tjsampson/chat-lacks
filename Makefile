SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=chat-lacks

VERSION=0.0.1
BUILD_TIME=`date +%FT%T%z`

LDFLAGS=-ldflags "-X github.com/tjsampson/chat-lacks/core.Version=${VERSION} -X github.com/tjsampson/chat-lacks/core.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go build ${LDFLAGS} -o ${BINARY} main.go

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi


.PHONY: test
test:
	go test ./... -v

.PHONY: cover
cover:
	go test ./... -cover
	
