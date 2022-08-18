APPNAME  := radiomomo
export APPNAME

build: build-linux build-mac

build-linux:
	mkdir -p bin/linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/linux/radiomomo

build-mac:
	mkdir -p bin/macosx
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -v -o bin/macosx/radiomomo

dirs:
	mkdir -p music

run:
	go run ./cmd/radio/main.go

clean:
	rm -fr bin
