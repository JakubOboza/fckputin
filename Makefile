.PHONE: test deps

build:
	go build -o bin/fckputin main.go

test:
	go test-v ./... -count=1

deps:
	go get

release: build-release package-release
	@echo "Release build and package"

build-release:
	GOOS=darwin GOARCH=amd64 go build -o release/osx-amd64/fckputin main.go
	GOOS=darwin GOARCH=arm64 go build -o release/osx-arm64/fckputin main.go
	GOOS=linux GOARCH=amd64 go build -o release/linux-amd64/fckputin main.go
	GOOS=windows GOARCH=amd64 go build -o release/win-amd64/fckputin.exe main.go

package-release:
	tar -czvf release/fckputin.osx-amd64.tar.gz --directory=release/osx-amd64 fckputin
	tar -czvf release/fckputin.osx-arm64.tar.gz --directory=release/osx-arm64 fckputin
	tar -czvf release/fckputin.linux-amd64.tar.gz --directory=release/linux-amd64 fckputin
	zip -j release/fckputin.win-amd64.zip release/win-amd64/fckputin.exe