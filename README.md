# OTTL Playground

### Building

By default, the built resources are placed into the `ottlplayground/web/public` directory.
After successfully compiling the WebAssembly and Frontend, this directory is ready to be deployed as a static site.
Considering the WebAssembly's size, it's highly recommended to serve it using a compression method, such as 
`gzip` or `brotli`.

#### WebAssembly

```shell
make build-wasm
```

#### Frontend

Requirements:
- Node.js [installation](https://nodejs.org/en/download/package-manager)

```shell
make build-web
```

##### Developing 

The `ottlplayground/web` contains the frontend source code, and uses `npm` as package manager.
To install the project dependencies:

```shell
npm install
```

Start local development server:

```shell
npm run serve public
```

Automatic reload the code changes:

```shell
npm run watch
```

### Running

#### Local

After building the project resources, it can be run - for testing purpose - by using the `main.go` implementation. 
To improve the load performance and saving bandwidth in real deployments, please confider hosting it using
a server with compressing capabilities, such as `gzip` or `brotli`.

```
go run main.go
```

#### Docker

```shell
docker build . -t ottlplayground
```

```shell
docker run -d -p 8080:8080 ottlplayground
```

The listening address can be changed by setting the environment variable `ADDR`,
in the form "host:port", If empty, ":8080" is used.

```shell
docker run -d -p 80:80 -e ADDR=":80" ottlplayground
```

The Docker image server does not support traffic compressing, and its use is not recommended
for hosting production environments.