BINARY := ipapk-server
VERSION ?= v1.0.0
PLATFORMS := windows linux darwin

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	CGO_ENABLED=1 GO15VENDOREXPERIMENT=1 GOOS=$@ GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY)-$(VERSION)-$@-amd64

.PHONY: release
release: windows linux darwin