#!/bin/bash

tinygo build \
    -o jwt-filter.wasm \
    -scheduler=none \
    -target=wasi \
    main.go
