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
ARG SKIP_BUILD_UNREGISTERED_VERSIONS
RUN if [ "$SKIP_BUILD_UNREGISTERED_VERSIONS" = "" ]; then make build-unregistered-versions ;  fi

# NGINX with brotli
FROM alpine
RUN apk add brotli nginx nginx-mod-http-brotli
COPY docker/nginx/default.conf /etc/nginx/http.d/default.conf
COPY --from=web /web/public /usr/share/nginx/html
COPY --from=builder /build/web/public/wasm /usr/share/nginx/html/wasm
RUN for file in /usr/share/nginx/html/wasm/ottlplayground-*.wasm ; do  brotli -9j "$file" ; done
CMD ["nginx", "-g", "daemon off;"]
EXPOSE 8080