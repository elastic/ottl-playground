GOCMD?=go
GO_BUILD_LDFLAGS?="-s -w"
WEB_OUTPUT_DIR?=../web/public
WASM_OUTPUT_DIR?=../web/public/wasm

.PHONY: clean
clean:
	rm -rf web/node_modules
	rm -rf web/public/wasm
	rm -f web/public/bundle.js
	rm -f web/public/bundle.js.map

.PHONY: build-test-server
build-build-test-server:
	 CGO_ENABLED=0 $(GOCMD) build -ldflags $(GO_BUILD_LDFLAGS) -o server

.PHONY: validate-registered-versions
validate-registered-versions:
	$(GOCMD) run ci-tools/main.go validate-registered-versions

.PHONY: build-unregistered-versions
build-unregistered-versions:
	$(eval PROCESSORS_VERSIONS ?= $(shell $(GOCMD) run ci-tools/main.go get-unregistered-versions))
	for v in $(PROCESSORS_VERSIONS); do \
		export PROCESSORS_VERSION=$$v ; \
		$(MAKE) update-processor-version && $(MAKE) build-wasm && $(MAKE) register-version ; \
	done

.PHONY: update-processor-version
update-processor-version:
	$(eval PARAMS = $(shell $(GOCMD) run ci-tools/main.go generate-processors-update -version=$(PROCESSORS_VERSION)))
	$(GOCMD) get $(PARAMS)
	@FIRST_PROCESSOR=$$(echo "$(PARAMS)" | awk '{print $$1}'); \
	COLLECTOR_DEPENDENCIES=$$($(GOCMD) mod graph | grep $$FIRST_PROCESSOR | grep "go.opentelemetry.io/collector/" | awk '{print $$2}' | sort -u); \
	COLLECTOR_PARAMS=""; \
	for DEP in $$COLLECTOR_DEPENDENCIES; do \
    	DEP_NAME=$$(echo $$DEP | cut -d'@' -f1); \
    	DEP_VERSION=$$(echo $$DEP | cut -d'@' -f2); \
    	if ! $(GOCMD) list -m -json $$DEP_NAME | grep -q "\"Indirect\":true"; then \
    		COLLECTOR_PARAMS="$$DEP_NAME@$$DEP_VERSION $$COLLECTOR_PARAMS"; \
    	fi; \
    done; \
    $(GOCMD) get $$COLLECTOR_PARAMS;
	$(MAKE) tidy

.PHONY: build-wasm
build-wasm:
	$(eval PROCESSORS_VERSION ?= $(shell $(GOCMD) run ci-tools/main.go get-version))
	$(GOCMD) run ci-tools/main.go generate-constants -version=$(PROCESSORS_VERSION)
	cd wasm; GOARCH=wasm GOOS=js $(GOCMD) build -ldflags $(GO_BUILD_LDFLAGS) -o $(WASM_OUTPUT_DIR)/ottlplayground-$(PROCESSORS_VERSION).wasm

.PHONY: register-version
register-version:
	$(eval PROCESSORS_VERSION ?= $(shell go run ci-tools/main.go get-version))
	$(GOCMD) run ci-tools/main.go register-wasm -version=$(PROCESSORS_VERSION)

.PHONY: build-web
build-web:
	cd web; npm install; npm run build

.PHONY: update-wasm-exec
update-wasm-exec:
	cp "$(shell go env GOROOT)/lib/wasm/wasm_exec.js" web/src/wasm_exec.js

.PHONY: build
build: update-wasm-exec build-web build-wasm register-version

.PHONY: fmt
fmt:
	gofmt  -w -s ./

.PHONY: tidy
tidy:
	rm -fr go.sum
	$(GOCMD) mod tidy -compat=1.22.0