THIS_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

build:
	go build $(THIS_DIR)

build-all-platforms:
	for goos in "darwin" "linux" "windows"; do \
		for goarch in "amd64"; do \
			GOOS=$$goos GOARCH=$$goarch go build -o out/$$goos-$$goarch/ $(THIS_DIR); \
		done \
	done

run:
	go run $(THIS_DIR)

update-dependencies:
	go get -u
	go mod tidy
