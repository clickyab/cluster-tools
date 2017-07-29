convey:
	$(GO) get -v github.com/smartystreets/goconvey
	$(GO) install -v github.com/smartystreets/goconvey

test-gui: convey
	cd $(ROOT) && goconvey -host=0.0.0.0

test: convey
	cd $(ROOT) && $(GO) test -v -race ./...
