# Go and compilation related variables
BUILD_DIR ?= out

ORG := github.com/code-ready
REPOPATH ?= $(ORG)/goodhosts
BINARY_NAME := crc-goodhosts

vendor:
	go mod vendor

$(BUILD_DIR):
	mkdir -p $(BUILD_DIR)

.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -rf vendor

$(BUILD_DIR)/macos-amd64/$(BINARY_NAME):
	GOARCH=amd64 GOOS=darwin go build -o $(BUILD_DIR)/macos-amd64/$(BINARY_NAME) ./cmd/main.go

$(BUILD_DIR)/linux-amd64/$(BINARY_NAME):
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/linux-amd64/$(BINARY_NAME) ./cmd/main.go

$(BUILD_DIR)/windows-amd64/$(BINARY_NAME).exe:
	GOARCH=amd64 GOOS=windows go build -o $(BUILD_DIR)/windows-amd64/$(BINARY_NAME).exe ./cmd/main.go

.PHONY: cross ## Cross compiles all binaries
cross: $(BUILD_DIR)/macos-amd64/$(BINARY_NAME) $(BUILD_DIR)/linux-amd64/$(BINARY_NAME) $(BUILD_DIR)/windows-amd64/$(BINARY_NAME).exe

.PHONY: install
	go install -o $(BINARY_NAME) ./cmd/main.go
