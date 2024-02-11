# Proxy JWT Filter in WebAssembly
使用 [Proxy-WASM 的 Go SDK](https://github.com/tetratelabs/proxy-wasm-go-sdk) 编写的 JWT 过滤器插件。使用 Tinygo 进行编译。适用于 [Envoy](https://github.com/envoyproxy/envoy)、[Istio](https://github.com/istio/proxy) 等环境。

## Features
- 提供 JWT 验证
- 支持使用 JSON 进行配置
- 支持路径匹配
- 自定义验证规则
- 集成统计功能
- 高度可扩展

## Build
本插件使用 [Tinygo](https://tinygo.org/) 将 Go 语言源代码编译为 WASM。请遵循官方文档安装 Tinygo。

在本项目根路径下，执行脚本进行编译：
```console
$ bash ./script/build.sh
```

编译产生的目标文件为 `jwt-filter.wasm`。

## Run Example
本项目在 `./examples` 目录下提供了样例 `envoy.yaml` 配置文件和 `Dockerfile`，可以构建一个简单的使用本插件的 Envoy 镜像。请预先在环境中安装 [Docker](https://www.docker.com/)。

执行如下脚本构建 Docker 镜像：

```console
$ bash ./script/docker_build.sh
```

这将构建名为 `wokron/envoy:demo` 的镜像。使用如下命令运行该镜像：

```console
$ docker run \
    -it \
    --rm \
    -p 18000:18000 \
    -p 8001:8001 \
    wokron/envoy:demo
```

这将运行一个具有 JWT Filter 的 Envoy 容器。会对路径为 `/api/v1` 下的请求进行 JWT 验证，除了 `/api/v1/abc`；此外禁止对其他路径的请求。按照配置，JWT 所使用的密钥为 `your-secure-key`。

```console
$ curl localhost:18000/api/v1/abc
hello from the server
$ curl localhost:18000/api/v1/xxx
Forbidden
$ curl localhost:18000/api/v2/abc
Forbidden
```

此外本插件还提供了统计功能，可以在 Envoy 的 Prometheus 接口中得知本插件的处理情况。
```console
$ curl localhost:8001/stats/prometheus | grep "envoy_wasm_jwt_filter"
```

## Configuration
本插件支持使用 JSON 进行配置。具体配置可以参考[配置文档](./docs/CONFIGURATION.zh.md)。
