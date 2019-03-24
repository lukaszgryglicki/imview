GO_BIN_FILES=cmd/imview/main.go cmd/iview/main.go
GO_LIB_FILES=run.go	texture.go util.go window.go
GO_BIN_CMDS=github.com/lukaszgryglicki/imview/cmd/imview github.com/lukaszgryglicki/iview/cmd/iview
GO_ENV=CGO_ENABLED=1
GO_BUILD=go build -ldflags '-s -w'
GO_INSTALL=go install -ldflags '-s'
GO_FMT=gofmt -s -w
GO_LINT=golint -set_exit_status
GO_VET=go vet
GO_CONST=goconst
GO_IMPORTS=goimports -w
GO_USEDEXPORTS=usedexports
GO_ERRCHECK=errcheck -asserts -ignore '[FS]?[Pp]rint*'
BINARIES=imview iview
STRIP=strip

all: check ${BINARIES}

imview: cmd/imview/main.go ${GO_LIB_FILES}
	 ${GO_ENV} ${GO_BUILD} -o imview cmd/imview/main.go

iview: cmd/iview/main.go ${GO_LIB_FILES}
	 ${GO_ENV} ${GO_BUILD} -o iview cmd/iview/main.go

fmt: ${GO_BIN_FILES} ${GO_LIB_FILES}
	./for_each_go_file.sh "${GO_FMT}"

lint: ${GO_BIN_FILES} ${GO_LIB_FILES}
	./for_each_go_file.sh "${GO_LINT}"

vet: ${GO_BIN_FILES} ${GO_LIB_FILES}
	go vet *.go
	go vet cmd/imview/main.go
	go vet cmd/iview/main.go

imports: ${GO_BIN_FILES} ${GO_LIB_FILES}
	./for_each_go_file.sh "${GO_IMPORTS}"

const: ${GO_BIN_FILES} ${GO_LIB_FILES}
	${GO_CONST} ./...

usedexports: ${GO_BIN_FILES} ${GO_LIB_FILES}
	${GO_USEDEXPORTS} ./...

errcheck: ${GO_BIN_FILES} ${GO_LIB_FILES}
	${GO_ERRCHECK} ./...

check: fmt lint imports vet const usedexports errcheck

install: check ${BINARIES}
	${GO_INSTALL} ${GO_BIN_CMDS}

strip: ${BINARIES}
	${STRIP} ${BINARIES}

clean:
	-rm -f ${BINARIES}

.PHONY: all
