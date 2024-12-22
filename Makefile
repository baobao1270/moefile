#!/usr/bin/make -f
APP_NAME := $(or $(APP_NAME),MoeFile)
APP_VERSION := $(or $(APP_VERSION),$(shell ./scripts/getver))
BUILD_TIMESTAMP := $(or $(BUILD_TIMESTAMP),$(shell date -Iseconds))
BUILD_MODE := $(or $(BUILD_MODE),production)
BUILD_OUTPUT := $(or $(BUILD_OUTPUT),bin/)

METADATA_PACKAGE := moefile/internal/meta
CROSS_BUILD_TRIPLES = 	darwin/amd64 \
						darwin/arm64 \
						windows/386 \
						windows/amd64 \
						windows/arm64 \
						linux/386 \
						linux/amd64 \
						linux/arm64 \
						linux/loong64 \
						linux/mips \
						linux/mips64 \
						linux/mips64le \
						linux/mipsle \
						linux/riscv64 \
						freebsd/386 \
						freebsd/amd64 \
						freebsd/arm64 \
						freebsd/riscv64
LDFLAGS = -s -w
LDFLAGS += -X "$(METADATA_PACKAGE).AppName=$(APP_NAME)"
LDFLAGS += -X "$(METADATA_PACKAGE).AppVersion=$(APP_VERSION)"
LDFLAGS += -X "$(METADATA_PACKAGE).BuildTimestamp=$(BUILD_TIMESTAMP)"
LDFLAGS += -X "$(METADATA_PACKAGE).BuildMode=$(BUILD_MODE)"


.PHONY: clean build
build:
	CGO_ENABLED=0 go build -ldflags '$(LDFLAGS)' -o $(BUILD_OUTPUT) ./cmd/moefile


.PHONY: all
all: clean $(CROSS_BUILD_TRIPLES)
$(CROSS_BUILD_TRIPLES): GOOS = $(word 1,$(subst /, ,$@))
$(CROSS_BUILD_TRIPLES): GOARCH = $(word 2,$(subst /, ,$@))
$(CROSS_BUILD_TRIPLES): CROSS_BUILD_OUTPUT_DIR = $(BUILD_OUTPUT)build/$@
$(CROSS_BUILD_TRIPLES): CROSS_BUILD_OUTPUT_FILE = $(shell echo $(APP_NAME) | tr A-Z a-z)$($(if $(filter $(GOOS),windows),.exe,)
$(CROSS_BUILD_TRIPLES): CROSS_BUILD_OUTPUT = $(CROSS_BUILD_OUTPUT_DIR)/$(CROSS_BUILD_OUTPUT_FILE)
$(CROSS_BUILD_TRIPLES): ARCHIVE_DIR = $(BUILD_OUTPUT)archives
$(CROSS_BUILD_TRIPLES): ARCHIVE_NAME = $(CROSS_BUILD_OUTPUT_FILE)-$(subst /,.,$(APP_VERSION)-$(GOOS)-$(GOARCH)).$(if $(filter $(GOOS),windows),zip,tar.gz)
$(CROSS_BUILD_TRIPLES):
	$(MAKE) build GOOS=$(GOOS) GOARCH=$(GOARCH) BUILD_OUTPUT=$(CROSS_BUILD_OUTPUT) && \
	mkdir -p $(ARCHIVE_DIR) && \
	if [ "$(GOOS)" = "windows" ]; then \
		zip -j    $(ARCHIVE_DIR)/$(ARCHIVE_NAME)    $(CROSS_BUILD_OUTPUT_DIR)/$(CROSS_BUILD_OUTPUT_FILE); \
	else \
		tar -cvzf $(ARCHIVE_DIR)/$(ARCHIVE_NAME) -C $(CROSS_BUILD_OUTPUT_DIR) $(CROSS_BUILD_OUTPUT_FILE); \
	fi && \
	sh -c "cd $(ARCHIVE_DIR) && sha256sum $(ARCHIVE_NAME) > $(ARCHIVE_NAME).sha256sum";


.PHONY: clean
clean:
	rm -rf $(BUILD_OUTPUT)
