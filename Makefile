SHELL=/bin/bash
PACKAGE_NAME=cosmart-backend-test
BUILD_DIR=./build

clean/cache:
	@echo "clean cache and test cache"
	go clean -cache -testcache -modcache

run/download:
	@echo "download all package from go mod"
	go mod download

clean/package:
	@echo "remove unused package"
	go mod tidy

test/coverage:
	@echo "running all unit testing include coverage"
	go test -cover ./...

run/service:
	@echo "running service golang in $(PACKAGE_NAME)"
	go run main.go