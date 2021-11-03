THIS_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

build:
	go build $(THIS_DIR)

update-dependencies:
	go get -u
	go mod tidy -compat=1.17

build-all-platforms: update-dependencies
	for goos in "darwin" "windows"; do \
		for goarch in "amd64"; do \
			GOOS=$$goos GOARCH=$$goarch go build -o out/$$goos-$$goarch/ $(THIS_DIR); \
		done \
	done

run:
	go run $(THIS_DIR)
