GOCMD?=go
GO_BUILD_LDFLAGS?="-s -w"
WEB_OUTPUT_DIR?=../web/public
WASM_OUTPUT_DIR?=../web/public/wasm

# Versions to build (latest first)
# Note: v0.125.0 not supported due to missing ProfileStatements in transformprocessor
WASM_VERSIONS?=v0.142.0 v0.138.0

.PHONY: clean
clean:
	rm -rf web/node_modules
	rm -rf web/public/wasm
	rm -f web/public/bundle.js
	rm -f web/public/bundle.js.map

.PHONY: build-web
build-web:
	cd web; npm install; npm run build

.PHONY: update-wasm-exec
update-wasm-exec:
	cp "$(shell go env GOROOT)/misc/wasm/wasm_exec.js" web/src/wasm_exec.js || cp "$(shell go env GOROOT)/lib/wasm/wasm_exec.js" web/src/wasm_exec.js || true

.PHONY: build-wasm
build-wasm:
	$(eval PROCESSORS_VERSION ?= $(shell $(GOCMD) run ci-tools/main.go get-version))
	$(eval BUILD_TAGS := $(shell echo "$(PROCESSORS_VERSION)" | sed 's/v//' | awk -F. '{if ($$1 > 0 || ($$1 == 0 && $$2 >= 142)) print "-tags ottl_ptr"}'))
	$(GOCMD) run ci-tools/main.go generate-constants -version=$(PROCESSORS_VERSION)
	cd wasm; GOARCH=wasm GOOS=js $(GOCMD) build $(BUILD_TAGS) -trimpath -ldflags $(GO_BUILD_LDFLAGS) -o $(WASM_OUTPUT_DIR)/ottlplayground-$(PROCESSORS_VERSION).wasm

.PHONY: build-all-wasm
build-all-wasm:
	@for v in $(WASM_VERSIONS); do \
		echo "Building WASM for $$v..."; \
		PROCESSORS_VERSION=$$v $(MAKE) update-processor-version && \
		PROCESSORS_VERSION=$$v $(MAKE) build-wasm && \
		PROCESSORS_VERSION=$$v $(MAKE) register-version || exit 1; \
	done
	@echo "Restoring to latest version..."
	PROCESSORS_VERSION=$(firstword $(WASM_VERSIONS)) $(MAKE) update-processor-version

.PHONY: update-processor-version
update-processor-version:
	$(eval PARAMS = $(shell $(GOCMD) run ci-tools/main.go generate-executors-update -version=$(PROCESSORS_VERSION)))
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

.PHONY: register-version
register-version:
	$(eval PROCESSORS_VERSION ?= $(shell $(GOCMD) run ci-tools/main.go get-version))
	$(GOCMD) run ci-tools/main.go register-wasm -version=$(PROCESSORS_VERSION)

.PHONY: build
build: update-wasm-exec build-web build-wasm

.PHONY: build-all
build-all: update-wasm-exec build-web build-all-wasm

.PHONY: fmt
fmt:
	gofmt -w -s ./

.PHONY: tidy
tidy:
	rm -fr go.sum
	$(GOCMD) mod tidy -compat=1.22.0
