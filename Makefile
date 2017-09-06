BINARY := ipapk-server
VERSION ?= v1.0.0
PLATFORMS := windows linux darwin

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	CGO_ENABLED=1 GOOS=$@ GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o release/$(BINARY)-$(VERSION)-$@-amd64

.PHONY: release
release: windows linux darwin