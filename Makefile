THIS_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

build:
	go build $(THIS_DIR)

run:
	go run $(THIS_DIR)

update-dependencies:
	go get -u
	go mod tidy
