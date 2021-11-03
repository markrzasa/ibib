THIS_DIR := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))

build:
	go build $(THIS_DIR)

update-dependencies:
	go get -u
	go mod tidy -compat=1.17

build-linux-windows:
	for goos in "linux" "windows"; do \
		for goarch in "amd64"; do \
			GOOS=$$goos GOARCH=$$goarch go build -o out/$$goos-$$goarch/ $(THIS_DIR); \
		done \
	done

build-darwin:
	for goos in "darwin"; do \
		for goarch in "amd64"; do \
			GOOS=$$goos GOARCH=$$goarch go build -o out/$$goos-$$goarch/ $(THIS_DIR); \
		done \
	done

run:
	go run $(THIS_DIR)

build-compose-up:
	docker-compose -f $(THIS_DIR)/compose/ubuntu-go/docker-compose.yml up -d
