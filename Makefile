.PHONY: build clean

build:
	mkdir -p dist
	go build -o dist/sftpx ./cmd/sftpx

build-windows:
	mkdir -p dist
	GOOS=windows GOARCH=amd64 go build -o dist/sftpx.exe ./cmd/sftpx

build-linux:
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build -o dist/sftpx ./cmd/sftpx

clean:
	rm -rf dist
