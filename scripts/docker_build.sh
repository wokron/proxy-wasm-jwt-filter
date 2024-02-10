#!/bin/bash

docker build \
    -t wokron/envoy:demo \
    -f ./examples/Dockerfile \
    ./
