.PHONY: build test clean

APP_NAME := go-template
BUILD_DIR := bin

build:
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd

test:
	go test ./... -json | sift

clean:
	rm -rf $(BUILD_DIR)
