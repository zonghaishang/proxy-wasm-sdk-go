.DEFAULT_GOAL := build

.PHONY: build build-image lint test

build:
	mkdir -p examples/${name}/build
	@rm -rf examples/${name}/build/${name}-go.wasm
	tinygo build -o ./examples/${name}/build/${name}-go.wasm -scheduler=none -target=wasi ./examples/${name}/main/main.go

build-image:
	@rm -rf ./examples/${name}/build
	mkdir -p examples/${name}/build
	docker run -it -w /tmp/build-proxy-wasm-go -v $(shell pwd):/tmp/build-proxy-wasm-go tinygo/tinygo:0.16.0 \
		tinygo build -o /tmp/build-proxy-wasm-go/examples/${name}/build/${name}-go.wasm -scheduler=none -target=wasi \
		/tmp/build-proxy-wasm-go/examples/${name}/main/main.go

lint:
	golangci-lint run --build-tags proxytest

test:
	go test -tags=proxytest $(shell go list ./... | grep -v e2e | sed 's/github.com\/zonghaishang\/proxy-wasm-sdk-go/./g')

#run:
#mosn -c ./examples/${name}/mosn_config.json
