# OTTL Playground

The OTTL Playground is a powerful and user-friendly tool designed to allow users to experiment with OTTL effortlessly. 
The playground provides a rich interface for users to create, modify, and test statements in real-time, making it easier 
to understand how different configurations impact the OTLP data transformation.

### Building

**Requirements:**
- Go 1.22 (https://go.dev/doc/install)
- Node.js (https://nodejs.org/en/download/prebuilt-installer)

By default, the built resources are placed into the `web/public` directory. After successfully 
compiling the WebAssembly and Frontend, this directory is ready to be deployed as a static website.
Given that the WebAssembly size is relatively big, it's highly recommended to serve it using a compression 
method, such as `gzip` or `brotli`.

```shell
make build
```

##### Developing

The `web` directory contains the frontend source code, and uses `npm` as package manager.
To start the local development server:

```shell
npm run serve public
```

Automatic reload the code changes:

```shell
npm run watch
```

### Running

#### Local

For **testing** purpose only, after successfully building the project resources, it can be run by 
using the `main.go` server implementation. 

To improve the load performance and saving bandwidth in real deployments, 
please confider hosting it using a server with compressing capabilities, such as `gzip` or `brotli`.

```
go run main.go
```

#### Docker

The Docker image delivered with this project serves the static website using Nginx (port 8080), and 
applying static `brotli` compression to the WebAssemblies files.

```shell
docker build . -t ottlplayground
```

```shell
docker run -d -p 8080:8080 ottlplayground
```
