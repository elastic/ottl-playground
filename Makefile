-include ../../Makefile.Common

GOCMD?=go
GO_BUILD_LDFLAGS?="-s -w"
WEB_OUTPUT_DIR?=../web/public
WASM_OUTPUT_DIR?=../web/public

.PHONY: clean
clean:
	rm -rf web/node_modules
	rm -f web/public/ottlplayground.wasm
	rm -f web/public/bundle.js
	rm -f web/public/bundle.js.map

.PHONY: build-server
build-server:
	 CGO_ENABLED=0 $(GOCMD) build -ldflags $(GO_BUILD_LDFLAGS) -o server

.PHONY: build-wasm
build-wasm:
	cd wasm; GOARCH=wasm GOOS=js $(GOCMD) build -ldflags $(GO_BUILD_LDFLAGS) -o $(WASM_OUTPUT_DIR)/ottlplayground.wasm

.PHONY: build-web
build-web:
	cd web; npm install; npm run build

.PHONY: build
build: build-web build-wasm