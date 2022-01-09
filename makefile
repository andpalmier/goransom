all: goransom

goransom: build
	@go build -o build/goransom cmd/goransom/*.go

install:
	@go install ./cmd/goransom

build:
	@mkdir -p build

clean:
	@rm -rf build
