FROM tinygo/tinygo:0.30.0 AS builder
WORKDIR /wasm
COPY . /wasm
# ENV GOPROXY=https://goproxy.cn,direct \
#     GO111MODULE=on \
#     CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=amd64
RUN bash ./scripts/build.sh

FROM envoyproxy/envoy:v1.29.0
COPY --from=builder /wasm/jwt-filter.wasm /lib/jwt-filter.wasm
COPY ./examples/envoy.yaml /etc/envoy.yaml
ENTRYPOINT /usr/local/bin/envoy -c /etc/envoy.yaml -l debug --service-cluster proxy --log-format '%v'
