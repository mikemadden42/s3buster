build:
	goreleaser build --snapshot --clean

check:
	golangci-lint run

release:
	goreleaser release --snapshot --clean

.PHONY: build check release
