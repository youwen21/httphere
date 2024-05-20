.PHONY:build-all

build-all: build-darwin build-linux build-windows
	echo "build-all done"

build-linux:
	GOOS=linux GOARCH=386 go build -o bin/httphere-linux-386
	GOOS=linux GOARCH=amd64 go build -o bin/httphere-linux-amd64
	GOOS=linux GOARCH=arm go build -o bin/httphere-linux-arm
	GOOS=linux GOARCH=arm64 go build -o bin/httphere-linux-arm64

build-windows:
	GOOS=windows GOARCH=386 go build -o bin/httphere-windows-386.exe
	GOOS=windows GOARCH=amd64 go build -o bin/httphere-windows-amd64.exe
	GOOS=windows GOARCH=arm64 go build -o bin/httphere-windows-arm64.exe

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o bin/httphere-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -o bin/httphere-darwin-arm64

