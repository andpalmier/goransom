all: goransom

goransom: build
	@go build -o build/goransom cmd/goransom/*.go

install: goransom
	@cp build/goransom /usr/bin/
	@chmod a+x /usr/bin/goransom

build:
	@mkdir -p build

clean:
	@rm -rf build
