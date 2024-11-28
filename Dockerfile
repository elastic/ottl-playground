# Frontend
FROM node:lts AS web
WORKDIR /web
COPY web .
COPY Makefile .
RUN make build-web

# Web-assembly and server
FROM golang:1.22 AS builder
WORKDIR /build
COPY ./ .

ENV WASM_OUTPUT_DIR=..
RUN make build-wasm
RUN make build-test-server

# Static server (TODO: Replace by nginx or any other server that supports gzip or brotli)
FROM scratch

COPY --from=web /web/public /ottlplayground/web/public
COPY --from=builder /build/web/public/wasm /ottlplayground/web/public/wasm
COPY --from=builder /build/server /ottlplayground

ENTRYPOINT ["./ottlplayground/server"]