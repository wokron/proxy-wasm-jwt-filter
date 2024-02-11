# Proxy JWT Filter in WebAssembly
JWT filter plugin written in Go using [Go SDK for Proxy-WASM](https://github.com/tetratelabs/proxy-wasm-go-sdk). Compiled using Tinygo. Suitable for environments such as [Envoy](https://github.com/envoyproxy/envoy), [Istio](https://github.com/istio/proxy), etc.

## Features
- Provides JWT validation
- Supports configuration using JSON
- Supports path matching
- Custom validation rules
- Integrated statistics functionality
- Highly scalable

## Build
This plugin compiles Go language source code to WASM using [Tinygo](https://tinygo.org/). Follow the official documentation to install Tinygo.

In the root directory of this project, execute the script to build:

```console
$ bash ./script/build.sh
```

The target file is `jwt-filter.wasm`。

## Run Example
This project provides an example `envoy.yaml` configuration file and `Dockerfile` in the `./examples` directory, which can be used to build a simple Envoy image using this plugin. [Docker](https://www.docker.com/) should be installed in the environment beforehand.

Execute the following script to build the Docker image:

```console
$ bash ./script/docker_build.sh
```

This will build an image named `wokron/envoy:demo`. Use the following command to run the image:

```console
$ docker run \
    -it \
    --rm \
    -p 18000:18000 \
    -p 8001:8001 \
    wokron/envoy:demo
```

This will run an Envoy container with the JWT Filter. Requests under the path `/api/v1` will undergo JWT validation, except for `/api/v1/abc`; other paths will be denied. According to the configuration, the JWT key used is `your-secure-key`.

```console
$ curl localhost:18000/api/v1/abc
hello from the server
$ curl localhost:18000/api/v1/xxx
Forbidden
$ curl localhost:18000/api/v2/abc
Forbidden
```

Additionally, this plugin provides statistics functionality, which can be obtained through the Prometheus interface of Envoy.

```console
$ curl localhost:8001/stats/prometheus | grep "envoy_wasm_jwt_filter"
```

## Configuration
This plugin supports configuration using JSON. For specific configurations, refer to the [configuration documentation](./docs/CONFIGURATION.md).
