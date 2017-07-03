PACKAGE=github.com/gkats/scraper
BINARIES=scraper server

BUILD = `git rev-parse HEAD`
BUILD_DIR = ${GOPATH}/src/${PACKAGE}

LDFLAGS = -ldflags "-X main.Build=${BUILD}"

# Default target
build:  vet fmt
	go build ${LDFLAGS} ${PACKAGE}/...

install: vet fmt
	go install ${LDFLAGS} ${PACKAGE}/...

clean:
	go clean ${PACKAGE}
	for binary in $(BINARIES); do \
		if [ -f ${GOPATH}/bin/$$binary ]; then rm ${GOPATH}/bin/$$binary; fi \
	done

fmt:
	cd ${BUILD_DIR}; \
	go fmt $$(go list ./... | grep -v /vendor/) ; \
	cd - >/dev/null

vet:
	cd ${BUILD_DIR}; \
	go vet ./... ; \
	cd - >/dev/null

.PHONY: clean fmt vet install
