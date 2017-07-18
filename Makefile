export ROOT=$(realpath $(dir $(firstword $(MAKEFILE_LIST))))
include $(ROOT)/bin/build/variables.mk
all:
	$(BUILD) ./...
include $(ROOT)/bin/build/linter.mk
include $(ROOT)/bin/build/test.mk
