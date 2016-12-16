
BUILD_DIR := builds
BIN_NAME := mqttappbridge

.PHONY: all clean setup

all: $(BUILD_DIR)/$(BIN_NAME) $(BUILD_DIR)/$(BIN_NAME).osx $(BUILD_DIR)/$(BIN_NAME).exe

setup:
	go get -u github.com/tidwall/gjson
	mkdir -p $(BUILD_DIR)

$(BUILD_DIR)/$(BIN_NAME): setup main.go
	GOARCH=amd64 GOOS=linux go build -o $@

$(BUILD_DIR)/$(BIN_NAME).osx: setup main.go
	GOARCH=amd64 GOOS=darwin go build -o $@

$(BUILD_DIR)/$(BIN_NAME).exe: setup main.go
	GOARCH=386 GOOS=windows go build -o $@

clean:
	$(RM) -r $(BUILD_DIR)
