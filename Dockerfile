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
RUN gzip -9 /build/ottlplayground.wasm
RUN mv /build/ottlplayground.wasm.gz /build/ottlplayground.wasm
RUN make build-server

# Static server
FROM scratch

COPY --from=web /web/public /ottlplayground/web/public
COPY --from=builder /build/ottlplayground.wasm /ottlplayground/web/public
COPY --from=builder /build/server /ottlplayground

ENTRYPOINT ["./ottlplayground/server"]