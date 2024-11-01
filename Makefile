build:
	GOOS=darwin GOARCH=amd64 go build -ldflags='-s -w' -o ./bin/gofile_darwin_amd64; \
	GOOS=darwin GOARCH=arm64 go build -ldflags='-s -w' -o ./bin/gofile_darwin_arm64; \
	GOOS=linux GOARCH=386 go build -ldflags='-s -w' -o ./bin/gofile_linux_386; \
	GOOS=linux GOARCH=amd64 go build -ldflags='-s -w' -o ./bin/gofile_linux_amd64; \
	GOOS=linux GOARCH=arm64 go build -ldflags='-s -w' -o ./bin/gofile_linux_arm64; \
	GOOS=windows GOARCH=386 go build -ldflags='-s -w' -o ./bin/gofile_windows_386.exe; \
	GOOS=windows GOARCH=amd64 go build -ldflags='-s -w' -o ./bin/gofile_windows_amd64.exe; \
	GOOS=windows GOARCH=arm64 go build -ldflags='-s -w' -o ./bin/gofile_windows_arm64.exe; \

check:
	gofmt -d .; \
	goimports -d .; \
	golint . ; \
	go vet . ; \
	golangci-lint run; \
	gosec ./...; \

release:
	goreleaser release --snapshot --clean

.PHONY: build check release
